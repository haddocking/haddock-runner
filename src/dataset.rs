use crate::utils::{extract_root, extract_root_from_filename};
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
) -> Vec<Target> {
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
            if let Some(root) = extract_root_from_filename(&file_name) {
                let builder = targets
                    .entry(root.clone())
                    .or_insert_with(|| TargetBuilder::new(root));
                builder.misc.push(path.clone());
            }
        }
    }

    // Convert builders to targets
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
        .collect()
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
