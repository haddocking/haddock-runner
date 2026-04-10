use std::path::PathBuf;
use std::{fs, path::Path, process::Command};

use anyhow::{Context, Result, bail};

use crate::utils::generate_timestamp;

pub fn run(command: &str, arg: &str, path: &Path) -> Result<PathBuf> {
    if command.is_empty() {
        bail!("no command was passed")
    }

    if arg.is_empty() {
        bail!("no arg was passed")
    }

    let log_path = path.join(format!("log_{}.txt", generate_timestamp()));

    // Execute
    let output = Command::new(command).arg(arg).current_dir(path).output()?;

    // Get the contents of the stdout/err
    let mut contents = output.stdout;
    contents.extend_from_slice(&output.stderr);

    // Write
    fs::write(&log_path, contents).context(format!("failed to write log file: {:?}", log_path))?;

    // Check exit code
    if !output.status.success() {
        bail!("command failed with status: {}", output.status)
    }

    Ok(log_path)
}
