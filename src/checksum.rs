use anyhow::{Context, Result, bail};
use md5::{Digest, Md5};
use std::collections::HashMap;
use std::fs;
use std::path::Path;

/// Calculate MD5 checksum for a file
///
/// An MD5 checksum is a 128-bit cryptographic hash value that serves as a
/// unique fingerprint for a file. It's commonly used to verify file integrity
/// and detect changes. Even a single byte change in the file will produce
/// a completely different MD5 hash.
///
/// This function reads the entire file content, computes its MD5 hash using
/// the MD5 algorithm, and returns the hash as a hexadecimal string.
///
/// # Arguments
///
/// * `file_path` - Path to the file to calculate checksum for
///
/// # Returns
///
/// * `Result<String>` - The MD5 checksum as a hexadecimal string, or an error if file reading fails
///
fn calculate_checksum<P: AsRef<Path>>(file_path: P) -> Result<String> {
    let file_content = fs::read(&file_path)
        .with_context(|| format!("Failed to read file: {}", file_path.as_ref().display()))?;

    let mut hasher = Md5::new();
    hasher.update(&file_content);
    let result = hasher.finalize();

    Ok(format!("{:x}", result))
}

/// Calculate checksums for all files in a target
fn calculate_target_checksums(target: &crate::dataset::Target) -> Result<HashMap<String, String>> {
    let mut checksums = HashMap::new();

    // Handle target collections
    let collections_to_check = vec![
        &target.molecules,
        &target.restraints,
        &target.toppar,
        &target.misc,
    ]
    .into_iter()
    .flatten();

    for element in collections_to_check {
        let checksum = calculate_checksum(element)?;
        checksums.insert(element.display().to_string(), checksum);
    }

    // Handle shape, its a special case because its an optional field
    if let Some(shape) = &target.shape {
        let checksum = calculate_checksum(shape)?;
        checksums.insert(shape.display().to_string(), checksum);
    }

    Ok(checksums)
}

/// Validate checksums against stored checksum file
pub fn validate_checksums(targets: &[crate::dataset::Target], checksum_file: &Path) -> Result<()> {
    // Calculate current checksums for all targets
    let mut current_checksums = HashMap::new();

    for target in targets {
        let target_checksums = calculate_target_checksums(target)?;
        current_checksums.extend(target_checksums);
    }

    // Create parent directory if it doesn't exist
    if let Some(parent) = checksum_file.parent() {
        fs::create_dir_all(parent).context("Failed to create checksum directory")?;
    }

    // Check if checksum file exists
    if !checksum_file.exists() {
        // Create new checksum file with pretty printing
        let serialized = serde_json::to_string_pretty(&current_checksums)
            .context("Failed to serialize checksums")?;
        fs::write(checksum_file, serialized).context("Failed to write checksum file")?;
        return Ok(()); // Fresh run
    }

    // Read stored checksums
    let stored_content =
        fs::read_to_string(checksum_file).context("Failed to read checksum file")?;
    let stored_checksums: HashMap<String, String> =
        serde_json::from_str(&stored_content).context("Failed to parse checksum file")?;

    // Compare checksums
    if current_checksums == stored_checksums {
        Ok(()) // Checksums match, can resume
    } else {
        // Find which files have changed
        let error_msg = find_modified(current_checksums, stored_checksums);
        bail!(error_msg);
    }
}

fn find_modified(
    current_checksums: HashMap<String, String>,
    stored_checksums: HashMap<String, String>,
) -> String {
    // Find which files have changed
    let mut changed_files = Vec::new();
    let mut new_files = Vec::new();
    let mut removed_files = Vec::new();

    // Files that exist in both but have different checksums
    for (file, current_checksum) in &current_checksums {
        if let Some(stored_checksum) = stored_checksums.get(file.as_str()) {
            if current_checksum != stored_checksum {
                changed_files.push(file);
            }
        } else {
            new_files.push(file);
        }
    }

    // Files that existed before but are now gone
    for file in stored_checksums.keys() {
        if !current_checksums.contains_key(file) {
            removed_files.push(file);
        }
    }

    // Build error message with details
    let mut error_msg = "Input files have changed since last run.\n".to_string();
    error_msg.push_str("Remove checksum.json to force fresh run.\n\n");

    if !changed_files.is_empty() {
        error_msg.push_str("Modified files:\n");
        for file in &changed_files {
            error_msg.push_str(&format!("  - {}\n", file,));
        }
        error_msg.push('\n');
    }

    if !new_files.is_empty() {
        error_msg.push_str("New files added:\n");
        for file in &new_files {
            error_msg.push_str(&format!("  - {}\n", file,));
        }
        error_msg.push('\n');
    }

    if !removed_files.is_empty() {
        error_msg.push_str("Files removed:\n");
        for file in &removed_files {
            error_msg.push_str(&format!("  - {}\n", file,));
        }
        error_msg.push('\n');
    }

    error_msg
}
