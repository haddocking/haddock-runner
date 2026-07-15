use anyhow::{Context, bail};
use chrono::Local;
use serde_yaml::Value;
use std::process::Command;
use toml::Value as TomlValue;

/// Generate a timestamp string in the format YYYY-MM-DD_HH-MM-SS
///
/// This function creates a timestamp using the local system time, formatted
/// for use in filenames and log files.
///
/// # Returns
///
/// * `String` - Formatted timestamp string
pub fn generate_timestamp() -> String {
    Local::now().format("%Y-%m-%d_%H-%M-%S").to_string()
}

/// Extract the root part from a filename by splitting on underscores
///
/// This function splits a filename by underscores and returns the first part,
/// which is typically used as the root identifier in HADDOCK file naming conventions.
///
/// # Arguments
///
/// * `filename` - The filename to extract the root from
///
/// # Returns
///
/// * `Option<String>` - The root part if underscore exists, None otherwise
pub fn extract_root(filename: &str) -> Option<String> {
    // Pattern: ROOT_*.xxx
    let parts: Vec<&str> = filename.split('_').collect();
    if parts.len() >= 2 {
        Some(parts[0].to_string())
    } else {
        None
    }
}

/// Format a YAML value for TOML output
///
/// This function converts YAML values to their TOML string representation
/// using the toml crate for proper formatting.
///
/// # Arguments
///
/// * `value` - The YAML value to format
///
/// # Returns
///
/// * `String` - TOML-formatted string representation
pub fn format_toml_value(value: &Value) -> String {
    // Convert YAML Value to TOML Value and use its Display implementation
    let toml_value: TomlValue = match value {
        Value::Bool(b) => TomlValue::Boolean(*b),
        Value::Number(n) => {
            if n.is_i64() {
                TomlValue::Integer(n.as_i64().unwrap())
            } else if n.is_u64() {
                TomlValue::Integer(n.as_i64().unwrap_or(i64::MAX))
            } else if n.is_f64() {
                TomlValue::Float(n.as_f64().unwrap())
            } else {
                TomlValue::String(n.to_string())
            }
        }
        Value::String(s) => TomlValue::String(s.clone()),
        Value::Sequence(seq) => {
            let items: Vec<TomlValue> = seq
                .iter()
                .map(|v| {
                    // Recursively convert each item
                    match v {
                        Value::Bool(b) => TomlValue::Boolean(*b),
                        Value::Number(n) => {
                            if n.is_i64() {
                                TomlValue::Integer(n.as_i64().unwrap())
                            } else if n.is_f64() {
                                TomlValue::Float(n.as_f64().unwrap())
                            } else {
                                TomlValue::String(n.to_string())
                            }
                        }
                        Value::String(s) => TomlValue::String(s.clone()),
                        _ => TomlValue::String("null".to_string()),
                    }
                })
                .collect();
            TomlValue::Array(items)
        }
        _ => TomlValue::String("null".to_string()),
    };

    toml_value.to_string()
}

/// Check if a command exists in the system PATH
///
/// This function checks whether a specified command is available in the system's PATH
/// by attempting to locate it using the 'which' command.
///
/// # Arguments
///
/// * `command` - The command name to check
///
/// # Returns
///
/// * `bool` - True if command exists and is executable, false otherwise
pub fn command_exists(command: &str) -> bool {
    if let Ok(output) = Command::new("which").arg(command).output() {
        output.status.success()
    } else {
        false
    }
}

/// Find the haddock3 executable in the system PATH
///
/// This function attempts to locate the haddock3 executable using the 'which' command
/// and returns its full path if found.
///
/// # Returns
///
/// * `anyhow::Result<String>` - Path to haddock3 executable if found, error otherwise
pub fn find_haddock3_executable() -> anyhow::Result<String> {
    // Try to find haddock3 in PATH
    if let Ok(output) = Command::new("which").arg("haddock3").output()
        && output.status.success()
    {
        let path = String::from_utf8_lossy(&output.stdout);
        Ok(path.trim().to_string())
    } else {
        bail!("haddock3 executable not found in PATH")
    }
}

/// Validate that haddock3 is available in the system PATH
///
/// This function checks if the haddock3 executable can be found and is accessible.
///
/// # Returns
///
/// * `anyhow::Result<()>` - Ok if haddock3 is found, error otherwise
pub fn validate_haddock3() -> anyhow::Result<()> {
    find_haddock3_executable()?;
    Ok(())
}

/// Get the version of the haddock3 executable found in PATH
///
/// This function locates the haddock3 executable and runs it with `--version`.
/// haddock3 reports its version as `haddock3 - YYYY.MM.PATCH` (e.g.
/// `haddock3 - 2026.3.0`); this parses out just the version part.
///
/// # Returns
///
/// * `anyhow::Result<String>` - The haddock3 version string (e.g. `2026.3.0`), error otherwise
pub fn get_haddock3_version() -> anyhow::Result<String> {
    let exe = find_haddock3_executable()?;

    let output = Command::new(&exe)
        .arg("--version")
        .output()
        .with_context(|| format!("failed to run `{exe} --version`"))?;

    if !output.status.success() {
        bail!(
            "`{exe} --version` exited with status {}: {}",
            output.status,
            String::from_utf8_lossy(&output.stderr)
        );
    }

    let stdout = String::from_utf8_lossy(&output.stdout);
    parse_haddock3_version(stdout.trim()).with_context(|| {
        format!(
            "could not parse haddock3 version from output: `{}`",
            stdout.trim()
        )
    })
}

/// Extract the version from haddock3's `haddock3 - <version>` output
fn parse_haddock3_version(output: &str) -> Option<String> {
    let (_, version) = output.rsplit_once(" - ")?;
    let version = version.trim();
    if version.is_empty() {
        None
    } else {
        Some(version.to_string())
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use serde_yaml::Value;

    #[test]
    fn test_extract_root() {
        assert_eq!(extract_root("protein_1.pdb"), Some("protein".to_string()));
        assert_eq!(extract_root("ligand.pdb"), None);
        assert_eq!(extract_root("complex_1_2.pdb"), Some("complex".to_string()));
    }

    #[test]
    fn test_format_toml_value() {
        assert_eq!(format_toml_value(&Value::Bool(true)), "true");
        assert_eq!(format_toml_value(&Value::Bool(false)), "false");
        assert_eq!(format_toml_value(&Value::Number(42.into())), "42");
        assert_eq!(
            format_toml_value(&Value::Number(std::f64::consts::PI.into())),
            "3.141592653589793"
        );
        assert_eq!(
            format_toml_value(&Value::String("test".to_string())),
            "\"test\""
        );

        let seq = vec![
            Value::Number(1.into()),
            Value::Number(2.into()),
            Value::Number(3.into()),
        ];
        assert_eq!(format_toml_value(&Value::Sequence(seq)), "[1, 2, 3]");
        assert_eq!(format_toml_value(&Value::Null), "\"null\"");
    }

    #[test]
    fn test_command_exists() {
        // Test with a command that should exist
        assert!(command_exists("ls"));

        // Test with a command that likely doesn't exist
        assert!(!command_exists("nonexistent_command_12345"));
    }

    #[test]
    fn test_parse_haddock3_version() {
        assert_eq!(
            parse_haddock3_version("haddock3 - 2026.3.0"),
            Some("2026.3.0".to_string())
        );
    }

    #[test]
    fn test_parse_haddock3_version_no_match() {
        assert_eq!(parse_haddock3_version("some unrelated tool output"), None);
    }
}
