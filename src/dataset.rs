use anyhow::{Context, Result};
use regex::Regex;
use std::collections::HashMap;
use std::fs;
use std::path::{Path, PathBuf};

#[derive(Debug)]
pub struct Dataset(pub Vec<Target>);

impl Dataset {
    /// Organizes the dataset by copying files to a structured directory layout.
    pub fn organize(&self, work_dir: &Path) -> Result<Dataset> {
        let mut organized_targets = Vec::new();

        // Process each target
        for target in &self.0 {
            let target_dir = work_dir.join(&target.id);
            let data_dir = target_dir.join("data");

            // Create directories
            fs::create_dir_all(&data_dir)?;

            // Helper function to copy file and update path
            let copy_and_update = |file: &PathBuf| -> Result<PathBuf> {
                let dest = data_dir.join(file.file_name().unwrap());
                fs::copy(file, &dest).with_context(|| {
                    format!("failed to copy {} to {}", file.display(), dest.display())
                })?;
                Ok(dest)
            };

            // Helper function to copy a collection of files
            let copy_collection = |files: &[PathBuf]| -> Result<Vec<PathBuf>> {
                files
                    .iter()
                    .map(copy_and_update)
                    .collect::<Result<Vec<_>>>()
            };

            // Copy all file collections
            let molecules = copy_collection(&target.molecules)?;
            let restraints = copy_collection(&target.restraints)?;
            let toppar = copy_collection(&target.toppar)?;
            let misc = copy_collection(&target.misc)?;

            // Copy shape file if present
            let shape = target.shape.as_ref().map(copy_and_update).transpose()?;

            // Create new organized target
            let organized_target = Target {
                id: target.id.clone(),
                molecules,
                restraints,
                toppar,
                misc,
                shape,
            };

            organized_targets.push(organized_target);
        }

        Ok(Dataset(organized_targets))
    }
}

#[derive(Debug)]
pub struct Target {
    id: String,
    molecules: Vec<PathBuf>,
    restraints: Vec<PathBuf>,
    toppar: Vec<PathBuf>,
    misc: Vec<PathBuf>,
    shape: Option<PathBuf>,
}

impl Target {
    fn new(
        id: String,
        molecules: Vec<PathBuf>,
        restraints: Vec<PathBuf>,
        toppar: Vec<PathBuf>,
        misc: Vec<PathBuf>,
        shape: Option<PathBuf>,
    ) -> Target {
        Target {
            id,
            molecules,
            restraints,
            toppar,
            misc,
            shape,
        }
    }
}

pub fn load_dataset(
    input_list: &str,
    mol_suffixes: &[String],
    shape_suffix: Option<&str>,
) -> Dataset {
    // Parse the input_list contents and group the paths based on their common root
    let mut targets: HashMap<String, TargetBuilder> = HashMap::new();

    // Create regex patterns for molecule suffixes
    // Use more flexible pattern to match _l_u, _l_u_1, _l_u_2, etc.
    let mol_patterns: Vec<Regex> = mol_suffixes
        .iter()
        .map(|suffix| Regex::new(&format!(r"(.*){}(_\d+)?\.pdb$", regex::escape(suffix))).unwrap())
        .collect();

    let shape_pattern = shape_suffix
        .map(|suffix| Regex::new(&format!(r"(.*){}\.pdb$", regex::escape(suffix))).unwrap());

    let restraint_pattern = Regex::new(r".*_.*\.tbl$").unwrap();
    let toppar_pattern = Regex::new(r".*\.(top|param)$").unwrap();

    // Read input list file
    let content = std::fs::read_to_string(input_list).expect("Failed to read input list");

    for line in content.lines() {
        // Skip comments and empty lines
        if line.trim().is_empty() || line.trim().starts_with('#') {
            continue;
        }

        let path = PathBuf::from(line);
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

        // Try to match shape pattern
        if !matched {
            if let Some(pattern) = &shape_pattern {
                if let Some(captures) = pattern.captures(&file_name) {
                    let root = captures.get(1).unwrap().as_str().to_string();

                    // Add to target builder
                    let builder = targets
                        .entry(root.clone())
                        .or_insert_with(|| TargetBuilder::new(root));
                    builder.shape = Some(path.clone());
                    matched = true;
                }
            }
        }

        // Try to match restraint pattern
        if !matched && restraint_pattern.is_match(&file_name) {
            // Extract root from restraint file (pattern: ROOT_*.tbl)
            if let Some(root) = extract_root_from_restraint(&file_name) {
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
            if let Some(root) = extract_root_from_toppar(&file_name) {
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
            if let Some(root) = extract_root_from_filename(&file_name) {
                let builder = targets
                    .entry(root.clone())
                    .or_insert_with(|| TargetBuilder::new(root));
                builder.misc.push(path.clone());
            }
        }
    }

    // Convert builders to targets
    Dataset(
        targets
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
            .collect(),
    )
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

// Helper functions for extracting roots from filenames
fn extract_root_from_restraint(filename: &str) -> Option<String> {
    // Pattern: ROOT_*.tbl
    let parts: Vec<&str> = filename.split('_').collect();
    if parts.len() >= 2 {
        Some(parts[0].to_string())
    } else {
        None
    }
}

fn extract_root_from_toppar(filename: &str) -> Option<String> {
    // Pattern: ROOT_*.top or ROOT_*.param
    let parts: Vec<&str> = filename.split('_').collect();
    if parts.len() >= 2 {
        Some(parts[0].to_string())
    } else {
        None
    }
}

fn extract_root_from_filename(filename: &str) -> Option<String> {
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

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_load_dataset() {
        let mol_suffixes = vec!["_r_u".to_string(), "_l_u".to_string(), "_x_u".to_string()];
        let dataset = load_dataset(
            "example/docking/input_list.txt",
            &mol_suffixes,
            Some("_shape"),
        );

        let Dataset(targets) = dataset;

        // Should have 5 targets: 1A2K, 1GGR, 1PPE, 2OOB, 1QU9
        assert_eq!(targets.len(), 5);

        // Check 1A2K target
        let target_1a2k = targets.iter().find(|t| t.id == "1A2K").unwrap();
        assert_eq!(target_1a2k.molecules.len(), 2); // receptor and ligand
        assert!(!target_1a2k.restraints.is_empty()); // should have restraints
        assert_eq!(target_1a2k.toppar.len(), 2); // .top and .param files

        // Check 1GGR target (has multiple ligand files)
        let target_1ggr = targets.iter().find(|t| t.id == "1GGR").unwrap();
        assert_eq!(target_1ggr.molecules.len(), 4); // 1 receptor + 3 ligands
        assert!(!target_1ggr.restraints.is_empty()); // should have restraints

        // Check 1PPE target
        let target_1ppe = targets.iter().find(|t| t.id == "1PPE").unwrap();
        assert_eq!(target_1ppe.molecules.len(), 2); // receptor and ligand
        assert!(!target_1ppe.restraints.is_empty()); // should have restraints

        // Check 2OOB target
        let target_2oob = targets.iter().find(|t| t.id == "2OOB").unwrap();
        assert_eq!(target_2oob.molecules.len(), 2); // receptor and ligand
        assert!(!target_2oob.restraints.is_empty()); // should have restraints

        // Check 1QU9 target
        let target_1qu9 = targets.iter().find(|t| t.id == "1QU9").unwrap();
        assert_eq!(target_1qu9.molecules.len(), 3); // 3 molecules
        assert!(!target_1qu9.restraints.is_empty()); // should have restraints
    }

    #[test]
    fn test_organize_dataset() -> Result<()> {
        use tempfile::tempdir;

        let mol_suffixes = vec!["_r_u".to_string(), "_l_u".to_string()];
        let dataset = load_dataset(
            "example/docking/input_list.txt",
            &mol_suffixes,
            Some("_shape"),
        );

        // Create temporary directory for organization
        let temp_dir = tempdir()?;
        let work_dir = temp_dir.path();

        // Organize the dataset
        let organized_dataset = dataset.organize(work_dir)?;

        let Dataset(targets) = organized_dataset;

        // Should still have 5 targets
        assert_eq!(targets.len(), 5);

        // Check that files were copied to data directories
        for target in &targets {
            let data_dir = work_dir.join(&target.id).join("data");
            assert!(
                data_dir.exists(),
                "Data directory should exist for target {}",
                target.id
            );

            // Check that molecules were copied
            for molecule in &target.molecules {
                assert!(
                    molecule.exists(),
                    "Molecule file should exist: {}",
                    molecule.display()
                );
                assert!(
                    molecule.starts_with(&data_dir),
                    "Molecule should be in data directory"
                );
            }

            // Check that restraints were copied
            for restraint in &target.restraints {
                assert!(
                    restraint.exists(),
                    "Restraint file should exist: {}",
                    restraint.display()
                );
                assert!(
                    restraint.starts_with(&data_dir),
                    "Restraint should be in data directory"
                );
            }
        }

        Ok(())
    }
}
