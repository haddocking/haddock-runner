use std::path::PathBuf;
use std::{fs, path::Path, process::Command};

use anyhow::{Context, Result, bail};
use log::{debug, error, info};

use crate::utils::generate_timestamp;

pub fn run(path: &Path) -> Result<PathBuf> {
    debug!("Running command in directory: {}", path.display());
    let command = "haddock3";
    let arg = "run.toml";

    let log_path = path.join(format!("log_{}.txt", generate_timestamp()));
    debug!("Log will be written to: {}", log_path.display());

    // Execute
    debug!(
        "Executing command: {} {} in directory: {}",
        command,
        arg,
        path.display()
    );
    let output = match Command::new(command).arg(arg).current_dir(path).output() {
        Ok(output) => output,
        Err(e) => {
            error!(
                "Failed to execute command '{} {}' in directory '{}': {}",
                command,
                arg,
                path.display(),
                e
            );
            return Err(e).context(format!(
                "Failed to execute command '{} {}' in directory '{}'",
                command,
                arg,
                path.display()
            ));
        }
    };

    // Get the contents of the stdout/err
    let mut contents = output.stdout;
    contents.extend_from_slice(&output.stderr);

    // Write
    debug!("Writing log file to: {}", log_path.display());
    if let Err(e) = fs::write(&log_path, contents) {
        error!("Failed to write log file '{}': {}", log_path.display(), e);
        return Err(e).context(format!("failed to write log file: {:?}", log_path));
    }

    // Check exit code
    if !output.status.success() {
        error!(
            "Command '{} {}' failed with status: {} log: {}",
            command,
            arg,
            output.status,
            log_path.display()
        );
        bail!(
            "command failed with status: {} log: {:?}",
            output.status,
            log_path
        )
    }

    info!("Command executed successfully");
    Ok(log_path)
}
