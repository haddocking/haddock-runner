use crate::job::JOB_FILENAME;
use crate::runner::status::SlurmJobState;
use crate::utils::command_exists;
use anyhow::{Context, Result};
use std::path::PathBuf;
use std::process::Command;

pub struct SlurmJob {
    id: u16,
    wd: PathBuf,
    state: SlurmJobState,
}

impl SlurmJob {
    pub fn new(wd: PathBuf) -> Self {
        SlurmJob {
            id: u16::MIN,
            wd,
            state: SlurmJobState::Unknown,
        }
    }

    pub fn run(&mut self) -> Result<()> {
        self.submit()?;
        self.wait()?;

        Ok(())
    }

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
