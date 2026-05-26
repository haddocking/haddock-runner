use crate::job::Job;
use anyhow::Context;
use indicatif::{ProgressBar, ProgressStyle};
use log::info;
use std::sync::mpsc;
use std::thread;

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
        workload.sort_by(|a, b| b.target.size.cmp(&a.target.size));

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
        info!("Setting up {} jobs", &self.workload.len());
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
        info!("Start!");
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
        let mut results = Vec::with_capacity(total_jobs);

        // PHASE 1: Initial dispatch - fill up to concurrency limit
        // Keep starting jobs until we either:
        // - Reach the concurrency limit (active_jobs == concurrent)
        // - Run out of jobs to dispatch (job_index == total_jobs)
        while active_jobs < concurrent && job_index < total_jobs {
            let tx = tx.clone();
            let mut job_clone = self.workload[job_index].clone();

            // Spawn a new thread for this job
            let handle = thread::spawn(move || {
                // info!("Processing job: {}", job_clone.name);
                if let Err(e) = job_clone.run() {
                    // error!("Failed running job {}: {}", job_clone.name, e);
                    tx.send(Err(e)).unwrap(); // Send error result
                    return;
                }
                tx.send(Ok(())).unwrap(); // Send success result
            });

            // Track this thread and update counters
            handles.push(handle);
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
                let tx = tx.clone();
                let mut job_clone = self.workload[job_index].clone();

                // Spawn a new thread for the next job
                let handle = thread::spawn(move || {
                    // info!("Processing job: {}", job_clone.name);
                    if let Err(e) = job_clone.run() {
                        // error!("Failed processing job {}: {}", job_clone.name, e);
                        tx.send(Err(e)).unwrap();
                        return;
                    }
                    tx.send(Ok(())).unwrap();
                });

                // Track this new thread and update counters
                handles.push(handle);
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

        // PHASE 4: Result checking
        // If any job failed, propagate the error
        for result in results {
            result?; // The ? operator will return early if any result is Err
        }

        pb.finish_with_message("All jobs completed!");

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
}
