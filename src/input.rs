use crate::runner::local::validate_local;
use crate::runner::slurm::validate_slurm;
use anyhow::{Context, Result, bail};
use indexmap::IndexMap;
use serde::{Deserialize, Serialize};
use serde_yaml::Value;
use std::{
    fs,
    path::{Path, PathBuf},
    thread,
};

#[derive(Debug, Deserialize, Serialize, Clone)]
pub struct Input {
    pub general: General,
    pub scenarios: Vec<Scenario>,
}

impl Input {
    pub fn new(yaml_path: &Path) -> Result<Self> {
        // let yaml_path = Path::new("example/bm.yml");
        let yaml_content =
            std::fs::read_to_string(yaml_path).context("Failed to read input file")?;
        let mut input: Input =
            serde_yaml::from_str(&yaml_content).context("Failed to parse YAML")?;

        // Convert relative paths to absolute paths based on YAML location
        let base_dir = yaml_path.parent().unwrap();

        // Convert work_dir to absolute path if it's relative
        if input.general.work_dir.is_relative() {
            input.general.work_dir = base_dir.join(&input.general.work_dir);
        }

        // Convert input_list to absolute path if it's relative
        let input_list_path = Path::new(&input.general.input_list);
        if input_list_path.is_relative() {
            input.general.input_list = base_dir
                .join(input_list_path)
                .to_string_lossy()
                .into_owned();
        }

        // Create the work_dir
        fs::create_dir_all(&input.general.work_dir)?;

        Ok(input)
    }

    /// Validates the input configuration.
    /// Checks for required fields, valid paths, and logical consistency.
    pub fn validate(&self) -> Result<()> {
        // Validate patterns
        self.validate_patterns()?;

        // Validate general fields
        self.validate_general()?;

        Ok(())
    }

    /// Validates patterns (suffixes) for duplicates and conflicts.
    fn validate_patterns(&self) -> Result<()> {
        // Check mol_suffixes
        if self.general.mol_suffixes.is_empty() {
            anyhow::bail!("mol_suffixes not defined in general section");
        }

        // Check for at least 2 suffixes (receptor and ligand)
        if self.general.mol_suffixes.len() < 2 {
            anyhow::bail!("mol_suffixes should contain at least 2 suffixes (receptor and ligand)");
        }

        // Check for duplicates in mol_suffixes
        let mut unique_suffixes = std::collections::HashSet::new();
        for suffix in &self.general.mol_suffixes {
            if !unique_suffixes.insert(suffix) {
                anyhow::bail!("Duplicate suffix found in mol_suffixes: {}", suffix);
            }
        }

        // TODO: Add more pattern validation as needed
        // - Check for conflicting patterns
        // - Validate shape_suffix if present

        Ok(())
    }

    /// Validates general configuration fields.
    fn validate_general(&self) -> Result<()> {
        if self.general.work_dir == PathBuf::new() {
            anyhow::bail!("work_dir not defined in general section");
        }

        if self.general.input_list.is_empty() {
            anyhow::bail!("input_list not defined in general section");
        }

        if self.general.max_concurrent == 0 {
            anyhow::bail!("max_concurrent must be greater than 0");
        }

        match &self.general.execution {
            Execution::Local => {
                validate_local()?;
                let num_cpus = thread::available_parallelism().unwrap().get();
                let needed_cpus = self.general.max_concurrent * self.general.ncores;
                if needed_cpus > num_cpus as u16 {
                    bail!(
                        "execution = local max_concurrent = {} and ncores = {}, you need {} cores for this but this machine only has {}",
                        self.general.max_concurrent,
                        self.general.ncores,
                        needed_cpus,
                        num_cpus,
                    );
                }
            }
            Execution::Slurm => validate_slurm()?,
        }

        Ok(())
    }
}

#[derive(Debug, Deserialize, Serialize, Clone)]
pub struct General {
    pub mol_suffixes: Vec<String>,
    pub shape_suffix: Option<String>,
    pub input_list: String,
    pub work_dir: PathBuf,
    pub max_concurrent: u16,
    pub ncores: u16,
    pub execution: Execution,
}

#[derive(Debug, Deserialize, Serialize, Clone)]
#[serde(rename_all = "lowercase")]
pub enum Execution {
    Local,
    Slurm,
}

#[derive(Debug, Deserialize, Serialize, PartialEq, Clone)]
pub struct Scenario {
    pub name: String,
    pub workflow: Workflow,
}

#[derive(Debug, Deserialize, Serialize, PartialEq, Clone)]
pub struct Workflow {
    #[serde(flatten)]
    pub modules: IndexMap<String, Value>,
}

impl Workflow {
    /// Returns the value associated with the given key in the workflow.
    /// Returns an error if the key is not found.
    /// The returned Value can be converted to specific types using methods like
    /// as_bool(), as_i64(), as_str(), etc.
    pub fn get_value(&self, key: &str) -> Result<&Value> {
        self.modules
            .get(key)
            .with_context(|| format!("Key '{}' not found in workflow", key))
    }

    /// Returns the value associated with the given module and parameter.
    /// For example, for a module like:
    /// ```yaml
    /// topoaa:
    ///   autohis: true
    /// ```
    /// You would call `get_module_param("topoaa", "autohis")` to get the value `true`.
    pub fn get_module_param(&self, module: &str, param: &str) -> Result<&Value> {
        let module_value = self.get_value(module)?;
        let mapping = module_value
            .as_mapping()
            .with_context(|| format!("Module '{}' is not a mapping", module))?;
        mapping
            .get(Value::String(param.to_string()))
            .with_context(|| format!("Parameter '{}' not found in module '{}'", param, module))
    }
}
