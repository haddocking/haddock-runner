use crate::utils::command_exists;
use anyhow::{Context, Result};
use log::info;
use std::fs::{self, canonicalize};
use std::io::Write;
use std::path::Path;
use std::process::Command;

pub fn prepare_job_file(work_dir: &Path, executable: &str) -> Result<()> {
    // Create SLURM job script
    let job_script = work_dir.join("job.sh");
    let mut file = fs::File::create(&job_script)?;

    let absolute_wd = canonicalize(work_dir).unwrap_or_else(|_| work_dir.to_path_buf());

    // Write SLURM header
    let header = "#!/bin/bash\n".to_string()
        + "#SBATCH --job-name=haddock\n"
        + "#SBATCH --output=haddock-%j.out\n"
        + "#SBATCH --error=haddock-%j.err\n";

    // Write job body
    let body = format!(
        "cd {}\n/trinity/login/rodrigo/repos/haddock-runner/.venv/bin/haddock3 run.toml\n",
        absolute_wd.display(),
    );

    file.write_all(header.as_bytes())?;
    file.write_all(body.as_bytes())?;

    // // Make the script executable
    // #[cfg(unix)]
    // {
    //     use std::os::unix::fs::PermissionsExt;
    //     let mut perms = file.metadata()?.permissions();
    //     perms.set_mode(0o755);
    //     file.set_permissions(perms)?;
    // }

    Ok(())
}

pub fn submit(job_script: &Path) -> Result<String> {
    let output = Command::new("sbatch")
        .arg(job_script)
        .output()
        .context("Failed to execute sbatch command")?;

    if !output.status.success() {
        let error_msg = String::from_utf8_lossy(&output.stderr);
        anyhow::bail!("sbatch failed: {}", error_msg);
    }

    let output_str = String::from_utf8_lossy(&output.stdout).trim().to_string();
    let job_id = output_str
        .split_whitespace()
        .last()
        .unwrap_or("")
        .to_string();
    info!("{}", job_id);
    Ok(job_id)
}

pub fn wait(job_id: &str) -> Result<()> {
    let args = vec!["--format=JobID,State", "-n", "-j", job_id];
    info!("sacct {:?}", args.join(" "));
    let output = Command::new("sacct")
        .args(args)
        .output()
        .context("Failed to execute sacct command")?;

    if !output.status.success() {
        let error_msg = String::from_utf8_lossy(&output.stderr);
        anyhow::bail!("sacct failed: {}", error_msg);
    }

    // Check both stdout and stderr for the status
    let stdout_status = String::from_utf8_lossy(&output.stdout);
    let stderr_status = String::from_utf8_lossy(&output.stderr);
    let status = if stdout_status.is_empty() {
        stderr_status
    } else {
        stdout_status
    };

    info!("Job status: {}", status);

    if status.contains("COMPLETED") {
        Ok(())
    } else if status.contains("RUNNING") || status.contains("PENDING") {
        anyhow::bail!("Job {} is still running/pending", job_id)
    } else if status.contains("CANCELLED") {
        anyhow::bail!("Job {} was cancelled", job_id)
    } else {
        anyhow::bail!("Job {} failed or has unexpected status: {}", job_id, status)
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
