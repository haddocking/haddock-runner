use crate::utils::extract_root;
use anyhow::Context;
use regex::Regex;
use std::collections::HashMap;
use std::path::PathBuf;

#[derive(Debug, Clone)]
pub struct Target {
    pub id: String,
    pub molecules: Vec<PathBuf>,
    pub restraints: Vec<PathBuf>,
    pub toppar: Vec<PathBuf>,
    pub misc: Vec<PathBuf>,
    pub shape: Option<PathBuf>,
    pub size: u64,
}

impl Target {
    /// Create a new Target instance
    ///
    /// This method creates a new Target with the specified files and calculates
    /// the total size of all associated files.
    ///
    /// # Arguments
    ///
    /// * `id` - Target identifier
    /// * `molecules` - Vector of molecule file paths
    /// * `restraints` - Vector of restraint file paths
    /// * `toppar` - Vector of topology/parameter file paths
    /// * `misc` - Vector of miscellaneous file paths
    /// * `shape` - Optional shape file path
    ///
    /// # Returns
    ///
    /// * `Target` - Newly created Target instance
    fn new(
        id: String,
        molecules: Vec<PathBuf>,
        restraints: Vec<PathBuf>,
        toppar: Vec<PathBuf>,
        misc: Vec<PathBuf>,
        shape: Option<PathBuf>,
    ) -> Target {
        let size = calculate_target_size(&molecules, &restraints, &toppar, &misc, shape.as_ref());
        Target {
            id,
            molecules,
            restraints,
            toppar,
            misc,
            shape,
            size,
        }
    }
}

fn calculate_target_size(
    molecules: &[PathBuf],
    restraints: &[PathBuf],
    toppar: &[PathBuf],
    misc: &[PathBuf],
    shape: Option<&PathBuf>,
) -> u64 {
    let mut total_size = 0;

    // Sum sizes of all molecule files
    for path in molecules {
        if let Ok(metadata) = std::fs::metadata(path) {
            total_size += metadata.len();
        }
    }

    // Sum sizes of all restraint files
    for path in restraints {
        if let Ok(metadata) = std::fs::metadata(path) {
            total_size += metadata.len();
        }
    }

    // Sum sizes of all toppar files
    for path in toppar {
        if let Ok(metadata) = std::fs::metadata(path) {
            total_size += metadata.len();
        }
    }

    // Sum sizes of all misc files
    for path in misc {
        if let Ok(metadata) = std::fs::metadata(path) {
            total_size += metadata.len();
        }
    }

    // Add shape file size if present
    if let Some(shape_path) = shape
        && let Ok(metadata) = std::fs::metadata(shape_path)
    {
        total_size += metadata.len();
    }

    total_size
}

/// Load dataset from input list file
///
/// This function parses an input list file containing file paths and groups them
/// into targets based on their common root identifiers. It handles molecule files,
/// restraint files, topology/parameter files, miscellaneous files, and optional shape files.
///
/// # Arguments
///
/// * `input_list` - Path to the input list file
/// * `mol_suffixes` - Suffixes used to identify molecule files
/// * `shape_suffix` - Optional suffix for shape files
///
/// # Returns
///
/// * `Vec<Target>` - Vector of Target instances created from the input files
pub fn load_dataset(input_list: &str, mol_suffixes: &[String]) -> anyhow::Result<Vec<Target>> {
    // Parse the input_list contents and group the paths based on their common root
    let mut targets: HashMap<String, TargetBuilder> = HashMap::new();

    // Create regex patterns for molecule suffixes
    // Use more flexible pattern to match _l_u, _l_u_1, _l_u_2, etc.
    let mol_patterns: Vec<Regex> = mol_suffixes
        .iter()
        .map(|suffix| Regex::new(&format!(r"(.*){}(_\d+)?\.pdb$", regex::escape(suffix))).unwrap())
        .collect();

    let restraint_pattern = Regex::new(r".*_.*\.tbl$").unwrap();
    let toppar_pattern = Regex::new(r".*\.(top|param)$").unwrap();

    // Read input list file
    let content =
        std::fs::read_to_string(input_list).with_context(|| "could not read input list")?;

    for line in content.lines() {
        // Skip comments and empty lines
        if line.trim().is_empty() || line.trim().starts_with('#') {
            continue;
        }

        let path = PathBuf::from(line.trim());
        let file_name = path.file_name().unwrap().to_string_lossy();

        // Try to match molecule patterns
        let mut matched = false;
        for pattern in &mol_patterns {
            if let Some(captures) = pattern.captures(&file_name) {
                let root = captures.get(1).unwrap().as_str().to_string();

                // Add to target builder
                let builder = targets
                    .entry(root.clone())
                    .or_insert_with(|| TargetBuilder::new(root));
                builder.molecules.push(path.clone());
                matched = true;
                break;
            }
        }

        // Try to match restraint pattern
        if !matched && restraint_pattern.is_match(&file_name) {
            // Extract root from restraint file (pattern: ROOT_*.tbl)
            if let Some(root) = extract_root(&file_name) {
                let builder = targets
                    .entry(root.clone())
                    .or_insert_with(|| TargetBuilder::new(root));
                builder.restraints.push(path.clone());
                matched = true;
            }
        }

        // Try to match toppar pattern
        if !matched && toppar_pattern.is_match(&file_name) {
            // Extract root from toppar file (pattern: ROOT_*.top/param)
            if let Some(root) = extract_root(&file_name) {
                let builder = targets
                    .entry(root.clone())
                    .or_insert_with(|| TargetBuilder::new(root));
                builder.toppar.push(path.clone());
                matched = true;
            }
        }

        // If no pattern matched, add to misc
        if !matched {
            // Try to extract root from filename (remove suffixes)
            if let Some(root) = extract_root(&file_name) {
                let builder = targets
                    .entry(root.clone())
                    .or_insert_with(|| TargetBuilder::new(root));
                builder.misc.push(path.clone());
            }
        }
    }

    // Convert builders to targets
    let result = targets
        .into_values()
        .map(|builder| {
            Target::new(
                builder.id,
                builder.molecules,
                builder.restraints,
                builder.toppar,
                builder.misc,
                builder.shape,
            )
        })
        .collect();

    Ok(result)
}

// Helper struct for building targets
struct TargetBuilder {
    id: String,
    molecules: Vec<PathBuf>,
    restraints: Vec<PathBuf>,
    toppar: Vec<PathBuf>,
    misc: Vec<PathBuf>,
    shape: Option<PathBuf>,
}

impl TargetBuilder {
    /// Create a new TargetBuilder instance
    ///
    /// This method creates a new TargetBuilder with the specified target ID.
    ///
    /// # Arguments
    ///
    /// * `id` - Target identifier
    ///
    /// # Returns
    ///
    /// * `Self` - Newly created TargetBuilder instance
    fn new(id: String) -> Self {
        TargetBuilder {
            id,
            molecules: Vec::new(),
            restraints: Vec::new(),
            toppar: Vec::new(),
            misc: Vec::new(),
            shape: None,
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::io::Write;
    use tempfile::NamedTempFile;

    #[test]
    fn test_target_new() {
        let target = Target::new(
            "test_id".to_string(),
            vec![PathBuf::from("mol1.pdb")],
            vec![PathBuf::from("restraint1.tbl")],
            vec![PathBuf::from("toppar1.top")],
            vec![PathBuf::from("misc1.txt")],
            Some(PathBuf::from("shape1.pdb")),
        );

        assert_eq!(target.id, "test_id");
        assert_eq!(target.molecules.len(), 1);
        assert_eq!(target.restraints.len(), 1);
        assert_eq!(target.toppar.len(), 1);
        assert_eq!(target.misc.len(), 1);
        assert!(target.shape.is_some());
    }

    #[test]
    fn test_load_dataset_simple() {
        // Create a temporary input list file
        let mut temp_file = NamedTempFile::new().unwrap();
        writeln!(temp_file, "protein1_r.pdb").unwrap();
        writeln!(temp_file, "protein1_l.pdb").unwrap();
        writeln!(temp_file, "protein1_restraints.tbl").unwrap();
        let file_path = temp_file.path().to_str().unwrap().to_string();

        // Load dataset
        let targets = load_dataset(&file_path, &["_r".to_string(), "_l".to_string()]).unwrap();

        // Should have one target
        assert_eq!(targets.len(), 1);
        assert_eq!(targets[0].id, "protein1");
        assert_eq!(targets[0].molecules.len(), 2);
        assert_eq!(targets[0].restraints.len(), 1);
    }

    #[test]
    fn test_load_dataset_multiple_targets() {
        // Create a temporary input list file with multiple targets
        let mut temp_file = NamedTempFile::new().unwrap();
        writeln!(temp_file, "protein1_r.pdb").unwrap();
        writeln!(temp_file, "protein1_l.pdb").unwrap();
        writeln!(temp_file, "protein2_r.pdb").unwrap();
        writeln!(temp_file, "protein2_l.pdb").unwrap();
        let file_path = temp_file.path().to_str().unwrap().to_string();

        // Load dataset
        let targets = load_dataset(&file_path, &["_r".to_string(), "_l".to_string()]).unwrap();

        // Should have two targets
        assert_eq!(targets.len(), 2);
        // Check that both targets are present (order may vary)
        let target_ids: Vec<_> = targets.iter().map(|t| t.id.as_str()).collect();
        assert!(target_ids.contains(&"protein1"));
        assert!(target_ids.contains(&"protein2"));
    }

    #[test]
    fn test_load_dataset_with_comments_and_empty_lines() {
        // Create a temporary input list file with comments and empty lines
        let mut temp_file = NamedTempFile::new().unwrap();
        writeln!(temp_file, "# This is a comment").unwrap();
        writeln!(temp_file).unwrap(); // Empty line
        writeln!(temp_file, "protein1_r.pdb").unwrap();
        writeln!(temp_file, "protein1_l.pdb").unwrap();
        writeln!(temp_file, "# Another comment").unwrap();
        let file_path = temp_file.path().to_str().unwrap().to_string();

        // Load dataset
        let targets = load_dataset(&file_path, &["_r".to_string(), "_l".to_string()]).unwrap();

        // Should have one target (comments and empty lines should be ignored)
        assert_eq!(targets.len(), 1);
        assert_eq!(targets[0].id, "protein1");
    }

    #[test]
    fn test_load_dataset_with_misc_files() {
        // Create a temporary input list file with misc files
        let mut temp_file = NamedTempFile::new().unwrap();
        writeln!(temp_file, "protein1_r.pdb").unwrap();
        writeln!(temp_file, "protein1_l.pdb").unwrap();
        writeln!(temp_file, "protein1_other.txt").unwrap(); // This should be misc for protein1
        let file_path = temp_file.path().to_str().unwrap().to_string();

        // Load dataset
        let targets = load_dataset(&file_path, &["_r".to_string(), "_l".to_string()]).unwrap();

        // Should have one target with misc file
        assert_eq!(targets.len(), 1);
        assert_eq!(targets[0].misc.len(), 1);
    }
}
