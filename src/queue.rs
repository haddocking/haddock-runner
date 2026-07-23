use crate::job::{Job, RunOutcome};
use crate::summary::{self, JobOutcome, JobResult};
use anyhow::Context;
use chrono::Local;
use indicatif::{ProgressBar, ProgressStyle};
use log::info;
use std::sync::mpsc;
use std::thread;
use std::time::Instant;

pub struct Queue {
    concurrent: u16,
    workload: Vec<Job>,
}

impl Queue {
    /// Create a new Queue instance
    ///
    /// This method creates a job queue that will execute jobs with the specified concurrency.
    /// The workload is sorted by target size in descending order (largest first).
    ///
    /// # Arguments
    ///
    /// * `concurrent` - Maximum number of jobs to run concurrently
    /// * `workload` - Vector of jobs to be executed
    ///
    /// # Returns
    ///
    /// * `Self` - Configured Queue instance
    pub fn new(concurrent: u16, mut workload: Vec<Job>) -> Self {
        // Sort workload by target size in descending order (largest first)
        workload.sort_by_key(|job| std::cmp::Reverse(job.target.size));

        Queue {
            concurrent,
            workload,
        }
    }

    /// Set up all jobs in the queue
    ///
    /// This method runs the setup phase for all jobs sequentially, which includes
    /// creating working directories, copying input files, and generating configuration files.
    ///
    /// # Returns
    ///
    /// * `anyhow::Result<()>` - Ok if all jobs set up successfully, error otherwise
    pub fn setup(&self) -> anyhow::Result<()> {
        // Run setup for all jobs sequentially - this is I/O and doesn't benefit much from parallelism
        info!("Setting up {} jobs", self.workload.len());
        for job in &self.workload {
            let mut job_clone = job.clone();
            // info!("Setting up job: {}", job_clone.name);
            job_clone.setup()?;
        }
        Ok(())
    }

    /// Start executing jobs in the queue
    ///
    /// This method executes all jobs in the queue with the specified concurrency level,
    /// using a rolling window approach to maintain exactly N concurrent jobs.
    /// Progress is displayed with a progress bar.
    ///
    /// # Returns
    ///
    /// * `anyhow::Result<()>` - Ok if all jobs complete successfully, error otherwise
    pub fn start(&self) -> anyhow::Result<()> {
        // Nothing to run and nowhere to write a summary
        if self.workload.is_empty() {
            return Ok(());
        }

        info!("Start!");
        let started_at = Local::now();
        let concurrent = self.concurrent as usize;
        let total_jobs = self.workload.len();

        // Create progress bar with continuous spinner
        let pb = ProgressBar::new(total_jobs as u64);
        pb.set_style(
            ProgressStyle::with_template(
                "{spinner:.green} [{elapsed_precise}] {pos}/{len} [{wide_bar:.cyan/blue}] {percent}% (ETA: {eta})",
            )
            .unwrap()
            .progress_chars("##-"),
        );

        // Enable continuous spinner animation
        pb.enable_steady_tick(std::time::Duration::from_millis(100));

        // Rolling window: maintain exactly N concurrent jobs
        // - tx/rx: channel for communication between worker threads and main thread
        // - handles: keep track of all thread handles for later joining
        // - active_jobs: counter for currently running jobs
        // - job_index: tracks which job to dispatch next
        // - results: collects results from all jobs
        let (tx, rx) = mpsc::channel();
        let mut handles = vec![];
        let mut active_jobs = 0;
        let mut job_index = 0;
        let mut results: Vec<JobResult> = Vec::with_capacity(total_jobs);

        // Spawns a thread that runs `job_clone`, timing it and sending a
        // `JobResult` back through `tx` regardless of success or failure.
        fn dispatch(job_clone: Job, tx: mpsc::Sender<JobResult>) -> thread::JoinHandle<()> {
            let name = job_clone.name.clone();
            let scenario = job_clone.scenario.name.clone();
            let target = job_clone.target.id.clone();

            thread::spawn(move || {
                let mut job_clone = job_clone;
                let start = Instant::now();
                let outcome = match job_clone.run() {
                    Ok(RunOutcome::Executed) => JobOutcome::Completed,
                    Ok(RunOutcome::Skipped) => JobOutcome::Skipped,
                    Err(e) => JobOutcome::Failed(format!("{:#}", e)),
                };
                let duration = start.elapsed();

                tx.send(JobResult {
                    name,
                    scenario,
                    target,
                    duration,
                    outcome,
                })
                .unwrap();
            })
        }

        // PHASE 1: Initial dispatch - fill up to concurrency limit
        // Keep starting jobs until we either:
        // - Reach the concurrency limit (active_jobs == concurrent)
        // - Run out of jobs to dispatch (job_index == total_jobs)
        while active_jobs < concurrent && job_index < total_jobs {
            let job_clone = self.workload[job_index].clone();
            handles.push(dispatch(job_clone, tx.clone()));
            active_jobs += 1;
            job_index += 1;
        }

        // PHASE 2: Rolling window - maintain exactly N concurrent jobs
        // This loop continues until all jobs are processed
        while active_jobs > 0 {
            // Wait for ANY job to complete (blocks until a result arrives)
            let result = rx.recv().context("Failed to receive result")?;
            results.push(result); // Store the result
            active_jobs -= 1; // One job just finished
            pb.inc(1); // Update progress bar

            // If there are more jobs to process, start the next one immediately
            // This maintains the rolling window of exactly N concurrent jobs
            if job_index < total_jobs {
                let job_clone = self.workload[job_index].clone();
                handles.push(dispatch(job_clone, tx.clone()));
                active_jobs += 1; // New job started
                job_index += 1; // Move to next job
            }
        }

        // PHASE 3: Cleanup - wait for all threads to finish
        // This ensures all threads have properly terminated
        for handle in handles {
            // propagate error
            if let Err(e) = handle.join() {
                return Err(anyhow::anyhow!("Thread panicked: {:?}", e));
            }
        }

        pb.finish_with_message("All jobs completed!");

        // PHASE 4: Reporting - build and write the run summary before
        // checking for failures, so a summary is produced even on a
        // partially failed run.
        let finished_at = Local::now();
        let work_dir = self.workload[0].general.work_dir.clone();
        summary::build_and_write(&work_dir, started_at, finished_at, &results)?;

        // PHASE 5: Result checking - if any job failed, propagate an
        // aggregated error listing every failed job
        let failures: Vec<&JobResult> = results
            .iter()
            .filter(|r| matches!(r.outcome, JobOutcome::Failed(_)))
            .collect();

        if !failures.is_empty() {
            let mut error_msg = format!("{} job(s) failed:\n", failures.len());
            for failure in &failures {
                if let JobOutcome::Failed(msg) = &failure.outcome {
                    error_msg.push_str(&format!("  - {}: {}\n", failure.name, msg));
                }
            }
            anyhow::bail!(error_msg);
        }

        Ok(())
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::runner::status::Status;

    #[test]
    fn test_queue_new() {
        let workload = vec![
            Job {
                name: "job1".to_string(),
                status: Status::Unknown,
                wd: std::path::PathBuf::from("/tmp/job1"),
                target: crate::dataset::Target {
                    id: "target1".to_string(),
                    molecules: vec![],
                    restraints: vec![],
                    toppar: vec![],
                    misc: vec![],
                    shape: None,
                    size: 0,
                },
                scenario: crate::input::Scenario {
                    name: "scenario1".to_string(),
                    workflow: crate::input::Workflow {
                        modules: indexmap::IndexMap::new(),
                    },
                },
                general: crate::input::General {
                    mol_suffixes: vec!["_r".to_string(), "_l".to_string()],
                    input_list: "test.txt".to_string(),
                    work_dir: std::path::PathBuf::from("/tmp"),
                    max_concurrent: 1,
                    ncores: 1,
                    execution: crate::input::Execution::Local,
                    partition: None,
                    preprocess: None,
                    postprocess: None,
                    gen_archive: None,
                    slurm_header: None,
                    slurm_prologue: None,
                },
            },
            Job {
                name: "job2".to_string(),
                status: Status::Unknown,
                wd: std::path::PathBuf::from("/tmp/job2"),
                target: crate::dataset::Target {
                    id: "target2".to_string(),
                    molecules: vec![],
                    restraints: vec![],
                    toppar: vec![],
                    misc: vec![],
                    shape: None,
                    size: 0,
                },
                scenario: crate::input::Scenario {
                    name: "scenario2".to_string(),
                    workflow: crate::input::Workflow {
                        modules: indexmap::IndexMap::new(),
                    },
                },
                general: crate::input::General {
                    mol_suffixes: vec!["_r".to_string(), "_l".to_string()],
                    input_list: "test.txt".to_string(),
                    work_dir: std::path::PathBuf::from("/tmp"),
                    max_concurrent: 1,
                    ncores: 1,
                    execution: crate::input::Execution::Local,
                    partition: None,
                    preprocess: None,
                    postprocess: None,
                    gen_archive: None,
                    slurm_header: None,
                    slurm_prologue: None,
                },
            },
        ];

        let queue = Queue::new(2, workload);

        assert_eq!(queue.concurrent, 2);
        assert_eq!(queue.workload.len(), 2);
    }

    #[test]
    fn test_queue_setup() {
        // This test would normally test the setup method, but since it involves
        // file system operations and complex job setup, we'll just verify
        // that the method can be called without panicking
        let workload = vec![];
        let queue = Queue::new(1, workload);

        // This should not panic, though it won't do much with an empty workload
        let result = queue.setup();
        assert!(result.is_ok());
    }

    #[test]
    fn test_queue_sorting() {
        // Create test jobs with different target sizes
        let workload = vec![
            Job {
                name: "small_job".to_string(),
                status: Status::Unknown,
                wd: std::path::PathBuf::from("/tmp/small"),
                target: crate::dataset::Target {
                    id: "small".to_string(),
                    molecules: vec![],
                    restraints: vec![],
                    toppar: vec![],
                    misc: vec![],
                    shape: None,
                    size: 100,
                },
                scenario: crate::input::Scenario {
                    name: "scenario1".to_string(),
                    workflow: crate::input::Workflow {
                        modules: indexmap::IndexMap::new(),
                    },
                },
                general: crate::input::General {
                    mol_suffixes: vec!["_r".to_string(), "_l".to_string()],
                    input_list: "test.txt".to_string(),
                    work_dir: std::path::PathBuf::from("/tmp"),
                    max_concurrent: 1,
                    ncores: 1,
                    execution: crate::input::Execution::Local,
                    partition: None,
                    preprocess: None,
                    postprocess: None,
                    gen_archive: None,
                    slurm_header: None,
                    slurm_prologue: None,
                },
            },
            Job {
                name: "large_job".to_string(),
                status: Status::Unknown,
                wd: std::path::PathBuf::from("/tmp/large"),
                target: crate::dataset::Target {
                    id: "large".to_string(),
                    molecules: vec![],
                    restraints: vec![],
                    toppar: vec![],
                    misc: vec![],
                    shape: None,
                    size: 1000,
                },
                scenario: crate::input::Scenario {
                    name: "scenario2".to_string(),
                    workflow: crate::input::Workflow {
                        modules: indexmap::IndexMap::new(),
                    },
                },
                general: crate::input::General {
                    mol_suffixes: vec!["_r".to_string(), "_l".to_string()],
                    input_list: "test.txt".to_string(),
                    work_dir: std::path::PathBuf::from("/tmp"),
                    max_concurrent: 1,
                    ncores: 1,
                    execution: crate::input::Execution::Local,
                    partition: None,
                    preprocess: None,
                    postprocess: None,
                    gen_archive: None,
                    slurm_header: None,
                    slurm_prologue: None,
                },
            },
            Job {
                name: "medium_job".to_string(),
                status: Status::Unknown,
                wd: std::path::PathBuf::from("/tmp/medium"),
                target: crate::dataset::Target {
                    id: "medium".to_string(),
                    molecules: vec![],
                    restraints: vec![],
                    toppar: vec![],
                    misc: vec![],
                    shape: None,
                    size: 500,
                },
                scenario: crate::input::Scenario {
                    name: "scenario3".to_string(),
                    workflow: crate::input::Workflow {
                        modules: indexmap::IndexMap::new(),
                    },
                },
                general: crate::input::General {
                    mol_suffixes: vec!["_r".to_string(), "_l".to_string()],
                    input_list: "test.txt".to_string(),
                    work_dir: std::path::PathBuf::from("/tmp"),
                    max_concurrent: 1,
                    ncores: 1,
                    execution: crate::input::Execution::Local,
                    partition: None,
                    preprocess: None,
                    postprocess: None,
                    gen_archive: None,
                    slurm_header: None,
                    slurm_prologue: None,
                },
            },
        ];

        // Create queue (which should sort the workload)
        let queue = Queue::new(2, workload);

        // Verify that jobs are sorted by size in descending order
        assert_eq!(queue.workload[0].name, "large_job");
        assert_eq!(queue.workload[1].name, "medium_job");
        assert_eq!(queue.workload[2].name, "small_job");

        // Verify the sizes are in the correct order
        assert_eq!(queue.workload[0].target.size, 1000);
        assert_eq!(queue.workload[1].target.size, 500);
        assert_eq!(queue.workload[2].target.size, 100);
    }

    #[test]
    fn test_queue_start_writes_summary_even_on_failure() {
        use tempfile::tempdir;

        // haddock3 is not expected to be on PATH in the test environment, so
        // this job will fail (either resolving the checksum file or finding
        // the haddock3 executable) - either way it exercises the "write the
        // summary before propagating the failure" behavior of `start()`.
        let temp_dir = tempdir().unwrap();
        let work_dir = temp_dir.path().to_path_buf();

        let general = crate::input::General {
            mol_suffixes: vec!["_r".to_string(), "_l".to_string()],
            input_list: "test.txt".to_string(),
            work_dir: work_dir.clone(),
            max_concurrent: 1,
            ncores: 1,
            execution: crate::input::Execution::Local,
            partition: None,
            preprocess: None,
            postprocess: None,
            gen_archive: None,
            slurm_header: None,
            slurm_prologue: None,
        };

        let job = Job {
            name: "target1-scenario1".to_string(),
            status: Status::Unknown,
            wd: work_dir.join("scenario1").join("target1"),
            target: crate::dataset::Target {
                id: "target1".to_string(),
                molecules: vec![],
                restraints: vec![],
                toppar: vec![],
                misc: vec![],
                shape: None,
                size: 0,
            },
            scenario: crate::input::Scenario {
                name: "scenario1".to_string(),
                workflow: crate::input::Workflow {
                    modules: indexmap::IndexMap::new(),
                },
            },
            general,
        };

        let queue = Queue::new(1, vec![job]);
        let result = queue.start();

        assert!(result.is_err());

        let summary_path = work_dir.join("summary.json");
        assert!(summary_path.exists());

        let content = std::fs::read_to_string(&summary_path).unwrap();
        let parsed: serde_json::Value = serde_json::from_str(&content).unwrap();
        assert_eq!(parsed["total_jobs"], 1);
        assert_eq!(parsed["failed"], 1);
    }
}
