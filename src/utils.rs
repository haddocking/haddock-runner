use anyhow::bail;
use chrono::Local;
use serde_yaml::Value;
use std::process::Command;

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

/// Extract root from filename by removing common file extensions and suffixes
///
/// This function removes common HADDOCK file extensions (.pdb, .tbl, .top, .param)
/// and then extracts the root part before the first underscore.
///
/// # Arguments
///
/// * `filename` - The filename to process
///
/// # Returns
///
/// * `Option<String>` - The extracted root part
pub fn extract_root_from_filename(filename: &str) -> Option<String> {
    // Try to extract root by removing common suffixes
    let filename = filename.replace(".pdb", "");
    let filename = filename.replace(".tbl", "");
    let filename = filename.replace(".top", "");
    let filename = filename.replace(".param", "");

    // Remove any remaining suffixes after underscores
    if let Some(pos) = filename.find('_') {
        Some(filename[..pos].to_string())
    } else {
        Some(filename)
    }
}

/// Format a YAML value for TOML output
///
/// This function converts YAML values to their TOML string representation,
/// handling booleans, numbers, strings, sequences, and providing a fallback for other types.
///
/// # Arguments
///
/// * `value` - The YAML value to format
///
/// # Returns
///
/// * `String` - TOML-formatted string representation
pub fn format_toml_value(value: &Value) -> String {
    match value {
        Value::Bool(b) => b.to_string(),
        Value::Number(n) => n.to_string(),
        Value::String(s) => format!("\"{}\"", s),
        Value::Sequence(seq) => {
            let items: Vec<String> = seq.iter().map(format_toml_value).collect();
            format!("[{}]", items.join(", "))
        }
        _ => "null".to_string(), // Fallback for other types
    }
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
        bail!("could not find `haddock3` executable in the PATH");
    }
}

/// Validate that haddock3 is available and executable
///
/// This function checks if the haddock3 executable can be found in the system PATH
/// and is ready for use.
///
/// # Returns
///
/// * `anyhow::Result<()>` - Ok if haddock3 is available, error otherwise
pub fn validate_haddock3() -> anyhow::Result<()> {
    find_haddock3_executable()?;
    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;
    use serde_yaml::Value;

    #[test]
    fn test_generate_timestamp() {
        let timestamp = generate_timestamp();

        // Should be in the format YYYY-MM-DD_HH-MM-SS
        assert_eq!(timestamp.len(), 19);
        assert!(timestamp.contains('-'));
        assert!(timestamp.contains('_'));
    }

    #[test]
    fn test_extract_root() {
        // Test normal case
        assert_eq!(extract_root("protein1_r.pdb"), Some("protein1".to_string()));

        // Test multiple underscores
        assert_eq!(extract_root("protein_1_r.pdb"), Some("protein".to_string()));

        // Test no underscore
        assert_eq!(extract_root("protein.pdb"), None);

        // Test single part
        assert_eq!(extract_root("protein"), None);
    }

    #[test]
    fn test_extract_root_from_filename() {
        // Test PDB file
        assert_eq!(
            extract_root_from_filename("protein1_r.pdb"),
            Some("protein1".to_string())
        );

        // Test TBL file
        assert_eq!(
            extract_root_from_filename("protein1_restraints.tbl"),
            Some("protein1".to_string())
        );

        // Test TOP file
        assert_eq!(
            extract_root_from_filename("protein1.top"),
            Some("protein1".to_string())
        );

        // Test PARAM file
        assert_eq!(
            extract_root_from_filename("protein1.param"),
            Some("protein1".to_string())
        );

        // Test file with multiple underscores
        assert_eq!(
            extract_root_from_filename("protein_1_r.pdb"),
            Some("protein".to_string())
        );

        // Test file without common suffixes - this should return the full filename without extension
        assert_eq!(
            extract_root_from_filename("protein.txt"),
            Some("protein.txt".to_string())
        );
    }

    #[test]
    fn test_format_toml_value() {
        // Test bool
        assert_eq!(format_toml_value(&Value::Bool(true)), "true");
        assert_eq!(format_toml_value(&Value::Bool(false)), "false");

        // Test number
        assert_eq!(format_toml_value(&Value::Number(42.into())), "42");

        // Test string
        assert_eq!(
            format_toml_value(&Value::String("test".to_string())),
            "\"test\""
        );

        // Test sequence
        let seq = vec![
            Value::String("item1".to_string()),
            Value::String("item2".to_string()),
        ];
        assert_eq!(
            format_toml_value(&Value::Sequence(seq)),
            "[\"item1\", \"item2\"]"
        );

        // Test null (fallback)
        assert_eq!(format_toml_value(&Value::Null), "null");
    }

    #[test]
    fn test_command_exists() {
        // Test with a command that should exist
        assert!(command_exists("ls"));

        // Test with a command that should not exist
        assert!(!command_exists("nonexistent_command_12345"));
    }
}
