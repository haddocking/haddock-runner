use anyhow::{Context, Result};
use std::fs;
use std::io::Write;
use std::path::Path;
use std::process::Command;

pub fn prepare_job_file(work_dir: &Path, executable: &str) -> Result<()> {
    // Create SLURM job script
    let job_script = work_dir.join("job.sh");
    let mut file = fs::File::create(&job_script)?;

    // Write SLURM header
    let header = "#!/bin/bash\n".to_string()
        + "#SBATCH --job-name=haddock\n"
        + "#SBATCH --output=haddock-%j.out\n"
        + "#SBATCH --error=haddock-%j.err\n";

    // Write job body
    let body = format!("cd {}\n{}", work_dir.display(), executable);

    file.write_all(header.as_bytes())?;
    file.write_all(body.as_bytes())?;

    // Make the script executable
    #[cfg(unix)]
    {
        use std::os::unix::fs::PermissionsExt;
        let mut perms = file.metadata()?.permissions();
        perms.set_mode(0o755);
        file.set_permissions(perms)?;
    }

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

    let output_str = String::from_utf8_lossy(&output.stdout);
    Ok(output_str.to_string())
}

pub fn wait(job_id: &str) -> Result<()> {
    // Check job status using sacct
    let output = Command::new("sacct")
        .args(["--format=JobID,State", "-n", "-j", job_id])
        .output()
        .context("Failed to execute sacct command")?;

    if !output.status.success() {
        let error_msg = String::from_utf8_lossy(&output.stderr);
        anyhow::bail!("sacct failed: {}", error_msg);
    }

    let status = String::from_utf8_lossy(&output.stdout);

    // Parse the status - this is a simplified version
    // In a real implementation, you'd want to properly parse the sacct output
    if status.contains("COMPLETED") {
        Ok(())
    } else if status.contains("RUNNING") || status.contains("PENDING") {
        anyhow::bail!("Job {} is still running/pending", job_id)
    } else {
        anyhow::bail!("Job {} failed or has unexpected status: {}", job_id, status)
    }
}

pub fn validate_slurm() -> Result<()> {
    // TODO: Should check if we can access the needed commands from the PATH
    todo!()
}
