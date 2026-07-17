use anyhow::{Context, Result};
use chrono::{DateTime, Local};
use log::info;
use serde::Serialize;
use std::collections::BTreeMap;
use std::fs;
use std::path::Path;
use std::time::Duration;

/// Outcome of a single job's execution, as reported by the queue
#[derive(Debug, Clone)]
pub enum JobOutcome {
    Completed,
    Skipped,
    Failed(String),
}

/// Result of running a single job, collected by the queue for reporting
#[derive(Debug, Clone)]
pub struct JobResult {
    pub name: String,
    pub scenario: String,
    pub target: String,
    pub duration: Duration,
    pub outcome: JobOutcome,
}

#[derive(Debug, Serialize)]
struct JobSummary {
    name: String,
    scenario: String,
    target: String,
    status: &'static str,
    duration_secs: f64,
    #[serde(skip_serializing_if = "Option::is_none")]
    error: Option<String>,
}

#[derive(Debug, Serialize)]
struct ScenarioSummary {
    name: String,
    total: usize,
    completed: usize,
    skipped: usize,
    failed: usize,
    duration_secs: f64,
}

#[derive(Debug, Serialize)]
struct RunSummary {
    started_at: String,
    finished_at: String,
    duration_secs: f64,
    total_jobs: usize,
    completed: usize,
    skipped: usize,
    failed: usize,
    generated_files: u64,
    generated_size_bytes: u64,
    scenarios: Vec<ScenarioSummary>,
    jobs: Vec<JobSummary>,
}

/// Recursively count files and sum their sizes under `dir`
fn walk_dir(dir: &Path) -> Result<(u64, u64)> {
    let mut file_count = 0u64;
    let mut total_bytes = 0u64;

    for entry in
        fs::read_dir(dir).with_context(|| format!("failed to read dir {}", dir.display()))?
    {
        let entry = entry?;
        let file_type = entry.file_type()?;

        if file_type.is_dir() {
            let (files, bytes) = walk_dir(&entry.path())?;
            file_count += files;
            total_bytes += bytes;
        } else if file_type.is_file() {
            file_count += 1;
            total_bytes += entry.metadata()?.len();
        }
    }

    Ok((file_count, total_bytes))
}

/// Build the run summary from job results and write it to `<work_dir>/summary.json`
///
/// This is best-effort reporting: it is called regardless of whether any jobs
/// failed, so the caller can still learn how much of the benchmark was
/// generated even on a partially failed run.
///
/// # Arguments
///
/// * `work_dir` - The benchmark's work directory, where `summary.json` is written
/// * `started_at` / `finished_at` - Wall-clock bounds of the run
/// * `results` - Per-job results collected by the queue
///
/// # Returns
///
/// * `Result<()>` - Ok if the summary was written successfully, error otherwise
pub fn build_and_write(
    work_dir: &Path,
    started_at: DateTime<Local>,
    finished_at: DateTime<Local>,
    results: &[JobResult],
) -> Result<()> {
    let mut completed = 0usize;
    let mut skipped = 0usize;
    let mut failed = 0usize;

    let mut jobs = Vec::with_capacity(results.len());
    let mut scenarios: BTreeMap<String, ScenarioSummary> = BTreeMap::new();

    for result in results {
        let (status, error) = match &result.outcome {
            JobOutcome::Completed => {
                completed += 1;
                ("completed", None)
            }
            JobOutcome::Skipped => {
                skipped += 1;
                ("skipped", None)
            }
            JobOutcome::Failed(msg) => {
                failed += 1;
                ("failed", Some(msg.clone()))
            }
        };

        let duration_secs = result.duration.as_secs_f64();

        let scenario_summary =
            scenarios
                .entry(result.scenario.clone())
                .or_insert_with(|| ScenarioSummary {
                    name: result.scenario.clone(),
                    total: 0,
                    completed: 0,
                    skipped: 0,
                    failed: 0,
                    duration_secs: 0.0,
                });
        scenario_summary.total += 1;
        scenario_summary.duration_secs += duration_secs;
        match status {
            "completed" => scenario_summary.completed += 1,
            "skipped" => scenario_summary.skipped += 1,
            "failed" => scenario_summary.failed += 1,
            _ => unreachable!(),
        }

        jobs.push(JobSummary {
            name: result.name.clone(),
            scenario: result.scenario.clone(),
            target: result.target.clone(),
            status,
            duration_secs,
            error,
        });
    }

    let (generated_files, generated_size_bytes) = walk_dir(work_dir)
        .with_context(|| format!("failed to collect file stats for {}", work_dir.display()))?;

    let duration_secs = (finished_at - started_at).num_milliseconds() as f64 / 1000.0;

    let summary = RunSummary {
        started_at: started_at.to_rfc3339(),
        finished_at: finished_at.to_rfc3339(),
        duration_secs,
        total_jobs: results.len(),
        completed,
        skipped,
        failed,
        generated_files,
        generated_size_bytes,
        scenarios: scenarios.into_values().collect(),
        jobs,
    };

    info!(
        "Run finished in {:.1}s: {} completed, {} skipped, {} failed, {} files generated ({} bytes)",
        summary.duration_secs,
        summary.completed,
        summary.skipped,
        summary.failed,
        summary.generated_files,
        summary.generated_size_bytes
    );

    let summary_path = work_dir.join("summary.json");
    let serialized =
        serde_json::to_string_pretty(&summary).context("failed to serialize run summary")?;
    fs::write(&summary_path, serialized)
        .with_context(|| format!("failed to write summary file: {}", summary_path.display()))?;

    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;
    use tempfile::tempdir;

    fn make_result(name: &str, scenario: &str, outcome: JobOutcome) -> JobResult {
        JobResult {
            name: name.to_string(),
            scenario: scenario.to_string(),
            target: "target1".to_string(),
            duration: Duration::from_secs(1),
            outcome,
        }
    }

    #[test]
    fn test_walk_dir_counts_files_and_bytes() {
        let temp_dir = tempdir().unwrap();
        let sub_dir = temp_dir.path().join("sub");
        fs::create_dir_all(&sub_dir).unwrap();

        fs::write(temp_dir.path().join("a.txt"), "1234").unwrap();
        fs::write(sub_dir.join("b.txt"), "12345678").unwrap();

        let (files, bytes) = walk_dir(temp_dir.path()).unwrap();

        assert_eq!(files, 2);
        assert_eq!(bytes, 12);
    }

    #[test]
    fn test_build_and_write_creates_summary_json() {
        let temp_dir = tempdir().unwrap();
        fs::write(temp_dir.path().join("run1.log"), "log content").unwrap();

        let results = vec![
            make_result("target1-scenario1", "scenario1", JobOutcome::Completed),
            make_result("target2-scenario1", "scenario1", JobOutcome::Skipped),
            make_result(
                "target1-scenario2",
                "scenario2",
                JobOutcome::Failed("boom".to_string()),
            ),
        ];

        let started_at = Local::now();
        let finished_at = started_at + chrono::Duration::seconds(3);

        build_and_write(temp_dir.path(), started_at, finished_at, &results).unwrap();

        let summary_path = temp_dir.path().join("summary.json");
        assert!(summary_path.exists());

        let content = fs::read_to_string(&summary_path).unwrap();
        let parsed: serde_json::Value = serde_json::from_str(&content).unwrap();

        assert_eq!(parsed["total_jobs"], 3);
        assert_eq!(parsed["completed"], 1);
        assert_eq!(parsed["skipped"], 1);
        assert_eq!(parsed["failed"], 1);
        assert_eq!(parsed["generated_files"], 1);
        assert_eq!(parsed["scenarios"].as_array().unwrap().len(), 2);

        let jobs = parsed["jobs"].as_array().unwrap();
        assert_eq!(jobs.len(), 3);
        let failed_job = jobs
            .iter()
            .find(|j| j["status"] == "failed")
            .expect("failed job present");
        assert_eq!(failed_job["error"], "boom");
    }

    #[test]
    fn test_build_and_write_empty_results() {
        let temp_dir = tempdir().unwrap();

        let started_at = Local::now();
        let finished_at = started_at;

        build_and_write(temp_dir.path(), started_at, finished_at, &[]).unwrap();

        let summary_path = temp_dir.path().join("summary.json");
        let content = fs::read_to_string(&summary_path).unwrap();
        let parsed: serde_json::Value = serde_json::from_str(&content).unwrap();

        assert_eq!(parsed["total_jobs"], 0);
        assert_eq!(parsed["scenarios"].as_array().unwrap().len(), 0);
    }
}
