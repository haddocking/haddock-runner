use crate::job::JOB_FILENAME;
use crate::runner::status::SlurmJobState;
use crate::utils::command_exists;
use anyhow::{Context, Result};
use std::path::PathBuf;
use std::process::Command;

pub struct SlurmJob {
    id: u64,
    wd: PathBuf,
    state: SlurmJobState,
}

impl SlurmJob {
    /// Create a new SlurmJob instance
    ///
    /// This method creates a new SLURM job with the specified working directory.
    ///
    /// # Arguments
    ///
    /// * `wd` - Working directory for the SLURM job
    ///
    /// # Returns
    ///
    /// * `Self` - Newly created SlurmJob instance
    pub fn new(wd: PathBuf) -> Self {
        SlurmJob {
            id: u64::MIN,
            wd,
            state: SlurmJobState::Unknown,
        }
    }

    /// Run the SLURM job
    ///
    /// This method submits the job to SLURM and waits for its completion.
    ///
    /// # Returns
    ///
    /// * `Result<()>` - Ok if job completes successfully, error otherwise
    pub fn run(&mut self) -> Result<()> {
        self.submit()?;
        self.wait()?;

        Ok(())
    }

    /// Submit the job to SLURM
    ///
    /// This method submits the job script to SLURM using sbatch and captures the job ID.
    ///
    /// # Returns
    ///
    /// * `Result<()>` - Ok if submission successful, error otherwise
    pub fn submit(&mut self) -> Result<()> {
        let job_script = self.wd.join(JOB_FILENAME);
        let output = Command::new("sbatch")
            .arg(job_script.file_name().unwrap())
            .current_dir(job_script.parent().unwrap())
            .output()
            .context("Failed to execute sbatch command")?;

        if !output.status.success() {
            let error_msg = String::from_utf8_lossy(&output.stderr);
            anyhow::bail!("sbatch failed: {}", error_msg);
        }

        let stdout = String::from_utf8_lossy(&output.stdout);
        let job_id = stdout
            .split_whitespace()
            .last()
            .ok_or_else(|| anyhow::anyhow!("could not parse job_id"))?
            .parse()
            .map_err(|e| anyhow::anyhow!("could not parse job_id: {}", e))?;

        self.id = job_id;

        Ok(())
    }

    /// Wait for SLURM job completion
    ///
    /// This method monitors the job status using sacct and waits until the job
    /// reaches a terminal state (completed, failed, etc.).
    ///
    /// # Returns
    ///
    /// * `Result<()>` - Ok if job completes successfully, error otherwise
    pub fn wait(&mut self) -> Result<()> {
        loop {
            self.update_status()?;

            if self.state.is_terminal() {
                if self.state.is_success() {
                    return Ok(());
                } else {
                    anyhow::bail!("Job {} failed with status: {:?}", self.id, self.state);
                }
            }

            // If job is not terminal wait before checking again
            std::thread::sleep(std::time::Duration::from_secs(2));
        }
    }

    fn update_status(&mut self) -> Result<()> {
        let output = Command::new("sacct")
            .args(["-j", &self.id.to_string(), "-n", "--format=State"])
            .output()
            .context("Failed to execute sacct command")?;

        if !output.status.success() {
            let error_msg = String::from_utf8_lossy(&output.stderr);
            anyhow::bail!("sacct failed: {}", error_msg);
        }

        // Parse the output - sacct returns fixed-width format with leading spaces
        let stdout = String::from_utf8_lossy(&output.stdout);
        let status = stdout.split_whitespace().next().unwrap_or("");

        self.state = SlurmJobState::from_status_code(status);

        Ok(())
    }
}

/// Validate that SLURM commands are available
///
/// This function checks if the required SLURM commands (sbatch, sacct) are available
/// in the system PATH, which are needed for SLURM job submission and monitoring.
///
/// # Returns
///
/// * `Result<()>` - Ok if all required SLURM commands are available, error otherwise
pub fn validate_slurm() -> Result<()> {
    // Check if needed commands are available
    let needed_slurm_commands = vec!["sbatch", "sacct"];

    for command in needed_slurm_commands {
        // Check if sbatch command is available
        if !command_exists(command) {
            anyhow::bail!(format!(
                "{} command not found in PATH. Please ensure SLURM is installed and configured.",
                command
            ));
        }
    }

    Ok(())
}

#[cfg(test)]
mod tests {}
