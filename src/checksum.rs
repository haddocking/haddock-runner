use anyhow::{Context, Result, bail};
use md5::{Digest, Md5};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::fs;
use std::path::Path;

/// On-disk representation of `checksum.json`
///
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
struct ChecksumData {
    files: HashMap<String, String>,
    haddock3_version: String,
}

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
    // TODO: Use streaming instead of reading to avoid reading the whole file into memory
    let file_content = fs::read(&file_path)
        .with_context(|| format!("Failed to read file: {}", file_path.as_ref().display()))?;

    let mut hasher = Md5::new();
    hasher.update(&file_content);
    let result = hasher.finalize();

    Ok(format!("{:x}", result))
}

/// Calculate checksums for all files in a target
///
/// This function calculates MD5 checksums for all files associated with a target,
/// including molecules, restraints, toppar files, misc files, and optionally shape files.
///
/// # Arguments
///
/// * `target` - The target containing file paths to calculate checksums for
///
/// # Returns
///
/// * `Result<HashMap<String, String>>` - HashMap mapping file paths to their MD5 checksums
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
///
/// This function validates the integrity of input files by comparing their current MD5 checksums
/// against previously stored checksums. If the checksum file doesn't exist, it creates one.
/// If checksums don't match, it identifies which files have been modified, added, or removed.
///
/// # Arguments
///
/// * `targets` - Array of targets containing files to validate
/// * `yaml_file` - Path to the YAML input file
/// * `input_list_file` - Path to the input list file
/// * `checksum_file` - Path to the JSON file storing checksums
/// * `haddock3_version` - The haddock3 version currently in use
///
/// # Returns
///
/// * `Result<()>` - Ok if checksums match or new file created, Err if files have changed
pub fn validate_checksums(
    targets: &[crate::dataset::Target],
    yaml_file: &Path,
    input_list_file: &Path,
    checksum_file: &Path,
    haddock3_version: &str,
) -> Result<()> {
    // Check if checksum file exists
    let mut current_checksums = HashMap::new();

    // Calculate current checksums for all targets
    for target in targets {
        let target_checksums = calculate_target_checksums(target)?;
        current_checksums.extend(target_checksums);
    }

    // Add checksum for the YAML input file
    let yaml_checksum = calculate_checksum(yaml_file)?;
    current_checksums.insert(yaml_file.display().to_string(), yaml_checksum);

    // Add checksum for the input list file
    let input_list_checksum = calculate_checksum(input_list_file)?;
    current_checksums.insert(input_list_file.display().to_string(), input_list_checksum);

    if !checksum_file.exists() {
        // Create new checksum file with pretty printing
        let data = ChecksumData {
            files: current_checksums,
            haddock3_version: haddock3_version.to_string(),
        };
        let serialized =
            serde_json::to_string_pretty(&data).context("Failed to serialize checksums")?;
        fs::write(checksum_file, serialized).context("Failed to write checksum file")?;
        return Ok(()); // Fresh run
    }

    // Read stored checksums
    let stored_content =
        fs::read_to_string(checksum_file).context("Failed to read checksum file")?;
    let stored_data: ChecksumData =
        serde_json::from_str(&stored_content).context("Failed to parse checksum file")?;

    if stored_data.haddock3_version != haddock3_version {
        bail!(
            "\nHaddock3 version has changed since this benchmark's checksum file was created !!\n  - Expected: {}\n  - Found: {}\n\nIf you intended to switch haddock3 versions, start a new benchmark with a fresh work directory.\n",
            stored_data.haddock3_version,
            haddock3_version
        );
    }

    // Compare checksums
    if current_checksums == stored_data.files {
        Ok(()) // Checksums match, can resume
    } else {
        // Find which files have changed
        let error_msg = find_modified(current_checksums, stored_data.files);
        bail!(error_msg);
    }
}

/// Validate that the haddock3 executable currently on PATH matches the
/// version recorded in the checksum file
///
/// Unlike `validate_checksums`, which only runs once at startup, this is
/// meant to be called per-job so that a haddock3 version swap partway
/// through a long-running benchmark (e.g. an env-module or conda change)
/// is caught before the next job dispatches.
///
/// # Arguments
///
/// * `checksum_file` - Path to the JSON file storing checksums
///
/// # Returns
///
/// * `Result<()>` - Ok if the live version matches the stored one, Err otherwise
pub fn validate_haddock3_version_live(checksum_file: &Path) -> Result<()> {
    let stored_content =
        fs::read_to_string(checksum_file).context("Failed to read checksum file")?;
    let stored_data: ChecksumData =
        serde_json::from_str(&stored_content).context("Failed to parse checksum file")?;

    let live_version = crate::utils::get_haddock3_version()?;

    if live_version != stored_data.haddock3_version {
        bail!(
            "\nHaddock3 version changed mid-run !!\n  - Expected: {}\n  - Found: {}\n\nThe haddock3 executable on PATH no longer matches the version used to start this benchmark, making it non-reproducible. Restore the original version or start a new benchmark.\n",
            stored_data.haddock3_version,
            live_version
        );
    }

    Ok(())
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
    let _sep = "!! ========================================================================= !!";

    let mut error_msg = "\n".to_string() + _sep + "\n";
    error_msg.push_str("!! Input files have been modified during the benchmark run !!\n");
    error_msg.push_str("!! This makes the benchmark non-reproducible !!\n");
    error_msg.push_str(
        "!! If you change input or parameters mid execution, you need to start over !!\n",
    );
    error_msg
        .push_str("!! If this is an exploratory run, create a new configuration for testing !!\n");
    error_msg.push_str(&format!("{}\n\n", _sep));

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

#[cfg(test)]
mod tests {
    use super::*;
    use std::collections::HashMap;
    use std::io::Write;
    use tempfile::{NamedTempFile, tempdir};

    #[test]
    fn test_calculate_checksum() {
        // Create a temporary file with known content
        let mut temp_file = NamedTempFile::new().unwrap();
        writeln!(temp_file, "test content").unwrap();
        let file_path = temp_file.path().to_path_buf();

        // Calculate checksum
        let checksum = calculate_checksum(&file_path).unwrap();

        // Verify checksum is not empty and is a hex string
        assert!(!checksum.is_empty());
        assert!(checksum.len() == 32); // MD5 should be 32 characters
        assert!(checksum.chars().all(|c| c.is_ascii_hexdigit()));
    }

    #[test]
    fn test_calculate_checksum_different_content() {
        // Create two files with different content
        let mut temp_file1 = NamedTempFile::new().unwrap();
        writeln!(temp_file1, "content1").unwrap();
        let file_path1 = temp_file1.path().to_path_buf();

        let mut temp_file2 = NamedTempFile::new().unwrap();
        writeln!(temp_file2, "content2").unwrap();
        let file_path2 = temp_file2.path().to_path_buf();

        // Calculate checksums
        let checksum1 = calculate_checksum(&file_path1).unwrap();
        let checksum2 = calculate_checksum(&file_path2).unwrap();

        // Verify checksums are different
        assert_ne!(checksum1, checksum2);
    }

    #[test]
    fn test_calculate_target_checksums() {
        // Create a temporary target with some files
        let temp_dir = tempdir().unwrap();
        let temp_path = temp_dir.path();

        // Create test files
        let mol_file = temp_path.join("test_mol.pdb");
        let restraint_file = temp_path.join("test_restraint.tbl");
        let toppar_file = temp_path.join("test_toppar.top");
        let misc_file = temp_path.join("test_misc.txt");

        fs::write(&mol_file, "molecule content").unwrap();
        fs::write(&restraint_file, "restraint content").unwrap();
        fs::write(&toppar_file, "toppar content").unwrap();
        fs::write(&misc_file, "misc content").unwrap();

        // Create target with cloned paths
        let mol_file_clone = mol_file.clone();
        let restraint_file_clone = restraint_file.clone();
        let toppar_file_clone = toppar_file.clone();
        let misc_file_clone = misc_file.clone();

        let target = crate::dataset::Target {
            id: "test".to_string(),
            molecules: vec![mol_file],
            restraints: vec![restraint_file],
            toppar: vec![toppar_file],
            misc: vec![misc_file],
            shape: None,
            size: 0,
        };

        // Calculate target checksums
        let checksums = calculate_target_checksums(&target).unwrap();

        // Verify we have checksums for all files
        assert_eq!(checksums.len(), 4);
        assert!(checksums.contains_key(mol_file_clone.to_str().unwrap()));
        assert!(checksums.contains_key(restraint_file_clone.to_str().unwrap()));
        assert!(checksums.contains_key(toppar_file_clone.to_str().unwrap()));
        assert!(checksums.contains_key(misc_file_clone.to_str().unwrap()));
    }

    #[test]
    fn test_calculate_target_checksums_with_shape() {
        // Create a temporary target with shape file
        let temp_dir = tempdir().unwrap();
        let temp_path = temp_dir.path();

        // Create test files including shape
        let mol_file = temp_path.join("test_mol.pdb");
        let shape_file = temp_path.join("test_shape.pdb");

        fs::write(&mol_file, "molecule content").unwrap();
        fs::write(&shape_file, "shape content").unwrap();

        // Create target with shape using cloned paths
        let mol_file_clone = mol_file.clone();
        let shape_file_clone = shape_file.clone();

        let target = crate::dataset::Target {
            id: "test".to_string(),
            molecules: vec![mol_file],
            restraints: vec![],
            toppar: vec![],
            misc: vec![],
            shape: Some(shape_file),
            size: 0,
        };

        // Calculate target checksums
        let checksums = calculate_target_checksums(&target).unwrap();

        // Verify we have checksums for both files
        assert_eq!(checksums.len(), 2);
        assert!(checksums.contains_key(mol_file_clone.to_str().unwrap()));
        assert!(checksums.contains_key(shape_file_clone.to_str().unwrap()));
    }

    #[test]
    fn test_validate_checksums_new_file() {
        // Create a temporary directory for the checksum file
        let temp_dir = tempdir().unwrap();
        let checksum_file = temp_dir.path().join("checksum.json");

        // Create a test target
        let temp_target_dir = tempdir().unwrap();
        let mol_file = temp_target_dir.path().join("test_mol.pdb");
        fs::write(&mol_file, "molecule content").unwrap();

        let target = crate::dataset::Target {
            id: "test".to_string(),
            molecules: vec![mol_file],
            restraints: vec![],
            toppar: vec![],
            misc: vec![],
            shape: None,
            size: 0,
        };

        // Validate checksums (should create new file)
        let yaml_file = temp_dir.path().join("test.yml");
        fs::write(&yaml_file, "test yaml content").unwrap();
        let input_list_file = temp_dir.path().join("test_input_list.txt");
        fs::write(&input_list_file, "test input list content").unwrap();

        let result = validate_checksums(
            &[target],
            &yaml_file,
            &input_list_file,
            &checksum_file,
            "haddock3 3.0.0",
        );

        // Should succeed and create the file
        assert!(result.is_ok());
        assert!(checksum_file.exists());
    }

    #[test]
    fn test_find_modified() {
        // Create test checksums
        let mut current = HashMap::new();
        current.insert("file1.txt".to_string(), "checksum1".to_string());
        current.insert("file2.txt".to_string(), "checksum2".to_string());
        current.insert("file3.txt".to_string(), "checksum3".to_string());

        let mut stored = HashMap::new();
        stored.insert("file1.txt".to_string(), "checksum1_changed".to_string()); // Changed
        stored.insert("file4.txt".to_string(), "checksum4".to_string()); // Removed
        // file2.txt and file3.txt are new

        // Find modified files
        let error_msg = find_modified(current, stored);

        // Verify error message contains expected information
        assert!(error_msg.contains("Modified files:"));
        assert!(error_msg.contains("file1.txt"));
        assert!(error_msg.contains("New files added:"));
        assert!(error_msg.contains("file2.txt"));
        assert!(error_msg.contains("file3.txt"));
        assert!(error_msg.contains("Files removed:"));
        assert!(error_msg.contains("file4.txt"));
    }

    #[test]
    fn test_validate_checksums_yaml_file_changed() {
        // Create a temporary directory for the checksum file
        let temp_dir = tempdir().unwrap();
        let checksum_file = temp_dir.path().join("checksum.json");

        // Create test files
        let yaml_file = temp_dir.path().join("test.yml");
        let input_list_file = temp_dir.path().join("input_list.txt");
        let mol_file = temp_dir.path().join("test_mol.pdb");

        fs::write(&yaml_file, "original yaml content").unwrap();
        fs::write(&input_list_file, "input list content").unwrap();
        fs::write(&mol_file, "molecule content").unwrap();

        // Create target
        let target = crate::dataset::Target {
            id: "test".to_string(),
            molecules: vec![mol_file.clone()],
            restraints: vec![],
            toppar: vec![],
            misc: vec![],
            shape: None,
            size: 0,
        };

        // Create initial checksum file
        let result = validate_checksums(
            std::slice::from_ref(&target),
            &yaml_file,
            &input_list_file,
            &checksum_file,
            "haddock3 3.0.0",
        );
        assert!(result.is_ok());
        assert!(checksum_file.exists());

        // Now modify the YAML file
        fs::write(&yaml_file, "modified yaml content").unwrap();

        // Try to validate again - should fail because YAML file changed
        let result = validate_checksums(
            &[target],
            &yaml_file,
            &input_list_file,
            &checksum_file,
            "haddock3 3.0.0",
        );
        assert!(result.is_err());

        // Check that the error message mentions the YAML file
        let error_msg = result.unwrap_err().to_string();
        assert!(error_msg.contains("Modified files:"));
        assert!(error_msg.contains("test.yml"));
    }

    #[test]
    fn test_validate_checksums_haddock3_version_changed() {
        let temp_dir = tempdir().unwrap();
        let checksum_file = temp_dir.path().join("checksum.json");

        let yaml_file = temp_dir.path().join("test.yml");
        let input_list_file = temp_dir.path().join("input_list.txt");
        let mol_file = temp_dir.path().join("test_mol.pdb");

        fs::write(&yaml_file, "yaml content").unwrap();
        fs::write(&input_list_file, "input list content").unwrap();
        fs::write(&mol_file, "molecule content").unwrap();

        let target = crate::dataset::Target {
            id: "test".to_string(),
            molecules: vec![mol_file],
            restraints: vec![],
            toppar: vec![],
            misc: vec![],
            shape: None,
            size: 0,
        };

        // Create initial checksum file with one version
        let result = validate_checksums(
            std::slice::from_ref(&target),
            &yaml_file,
            &input_list_file,
            &checksum_file,
            "haddock3 3.0.0",
        );
        assert!(result.is_ok());

        // Validate again with a different version - should fail
        let result = validate_checksums(
            &[target],
            &yaml_file,
            &input_list_file,
            &checksum_file,
            "haddock3 3.1.0",
        );
        assert!(result.is_err());

        let error_msg = result.unwrap_err().to_string();
        assert!(error_msg.contains("Haddock3 version has changed"));
        assert!(error_msg.contains("3.0.0"));
        assert!(error_msg.contains("3.1.0"));
    }

    #[test]
    fn test_validate_haddock3_version_live_mismatch() {
        let temp_dir = tempdir().unwrap();
        let checksum_file = temp_dir.path().join("checksum.json");

        let data = ChecksumData {
            files: HashMap::new(),
            haddock3_version: "haddock3 9.9.9".to_string(),
        };
        fs::write(&checksum_file, serde_json::to_string_pretty(&data).unwrap()).unwrap();

        // The live haddock3 version (whatever is actually on PATH, or the
        // "not found" error if it isn't) will never equal the bogus stored
        // version above, so this must always fail.
        let result = validate_haddock3_version_live(&checksum_file);
        assert!(result.is_err());
    }
}
