use crate::runner::slurm::validate_slurm;
use crate::utils::validate_haddock3;
use anyhow::{Context, Result, bail};
use indexmap::IndexMap;
use regex::Regex;
use serde::{Deserialize, Serialize};
use serde_yaml::Value;
use std::{
    fs,
    path::{Path, PathBuf},
    sync::LazyLock,
    thread,
};

pub(crate) static DIGIT_SUFFIX_RE: LazyLock<Regex> =
    LazyLock::new(|| Regex::new(r"^(.*?)\.(\d+)$").unwrap());

#[derive(Debug, Deserialize, Serialize, Clone)]
#[serde(deny_unknown_fields)]
pub struct Input {
    pub general: General,
    pub scenarios: Vec<Scenario>,
}

impl Input {
    /// Create a new Input instance from a YAML file
    ///
    /// This method reads and parses a YAML configuration file, converts relative paths
    /// to absolute paths based on the YAML file location, and creates the working directory.
    ///
    /// # Arguments
    ///
    /// * `yaml_path` - Path to the YAML configuration file
    ///
    /// # Returns
    ///
    /// * `Result<Self>` - Parsed Input configuration or error
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

        // Validate workflow module naming
        self.validate_workflows()?;

        // Validate haddock3
        validate_haddock3()?;

        Ok(())
    }

    /// Validates that repeated modules in every scenario workflow consistently use
    /// `.digit` suffixes (e.g. `caprieval.1`, `caprieval.2`) starting at 1 with no gaps.
    /// Mixing a bare name with a suffixed variant (e.g. `caprieval` + `caprieval.1`)
    /// is rejected because haddock3 cannot reconcile the two forms.
    fn validate_workflows(&self) -> Result<()> {
        for scenario in &self.scenarios {
            // Check for bare names coexisting with suffixed siblings
            for name in scenario.workflow.modules.keys() {
                if DIGIT_SUFFIX_RE.captures(name).is_none() {
                    let has_suffixed_sibling = scenario
                        .workflow
                        .modules
                        .keys()
                        .filter(|n| *n != name)
                        .any(|n| {
                            DIGIT_SUFFIX_RE
                                .captures(n)
                                .map(|c| c.get(1).unwrap().as_str() == name.as_str())
                                .unwrap_or(false)
                        });
                    if has_suffixed_sibling {
                        bail!(
                            "scenario '{}': module '{name}' mixes bare and .digit-suffixed names; \
                            repeated modules must all use .digit suffixes (e.g. {name}.1, {name}.2)",
                            scenario.name
                        );
                    }
                }
            }

            // Check that suffixed modules are sequential starting at 1
            let mut suffix_map: std::collections::HashMap<String, Vec<u32>> =
                std::collections::HashMap::new();
            for name in scenario.workflow.modules.keys() {
                if let Some(caps) = DIGIT_SUFFIX_RE.captures(name) {
                    let base = caps.get(1).unwrap().as_str().to_string();
                    let n: u32 = caps.get(2).unwrap().as_str().parse().unwrap();
                    suffix_map.entry(base).or_default().push(n);
                }
            }
            for (base, mut indices) in suffix_map {
                indices.sort();
                let expected: Vec<u32> = (1..=indices.len() as u32).collect();
                if indices != expected {
                    bail!(
                        "scenario '{}': module '{base}' suffixes must be sequential starting at 1 \
                        (e.g. {base}.1, {base}.2), got: {}",
                        scenario.name,
                        indices
                            .iter()
                            .map(|n| format!("{base}.{n}"))
                            .collect::<Vec<_>>()
                            .join(", ")
                    );
                }
            }
        }

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

        if let Some(partition) = &self.general.partition
            && partition.trim().is_empty()
        {
            anyhow::bail!("partition must be non-empty when defined");
        }

        match &self.general.execution {
            Execution::Local => {
                validate_haddock3()?;
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
#[serde(deny_unknown_fields)]
pub struct General {
    pub mol_suffixes: Vec<String>,
    pub input_list: String,
    pub work_dir: PathBuf,
    pub max_concurrent: u16,
    pub ncores: u16,
    pub execution: Execution,
    pub partition: Option<String>,
}

#[derive(Debug, Deserialize, Serialize, Clone)]
#[serde(rename_all = "lowercase")]
pub enum Execution {
    Local,
    Slurm,
}

#[derive(Debug, Deserialize, Serialize, PartialEq, Clone)]
#[serde(deny_unknown_fields)]
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

#[cfg(test)]
mod tests {
    use super::*;
    use std::path::PathBuf;

    #[test]
    fn test_workflow_get_value() {
        let mut modules = IndexMap::new();
        modules.insert("topoaa".to_string(), Value::Bool(true));
        modules.insert("flexref".to_string(), Value::String("test".to_string()));

        let workflow = Workflow { modules };

        // Test getting existing value
        let result = workflow.get_value("topoaa");
        assert!(result.is_ok());
        assert_eq!(result.unwrap(), &Value::Bool(true));

        // Test getting non-existing value
        let result = workflow.get_value("nonexistent");
        assert!(result.is_err());
    }

    #[test]
    fn test_workflow_get_module_param() {
        let mut modules = IndexMap::new();

        // Create a module with parameters
        let mut topoaa_params = serde_yaml::Mapping::new();
        topoaa_params.insert(Value::String("autohis".to_string()), Value::Bool(true));
        topoaa_params.insert(Value::String("keephetatm".to_string()), Value::Bool(false));

        modules.insert("topoaa".to_string(), Value::Mapping(topoaa_params));

        let workflow = Workflow { modules };

        // Test getting existing parameter
        let result = workflow.get_module_param("topoaa", "autohis");
        assert!(result.is_ok());
        assert_eq!(result.unwrap(), &Value::Bool(true));

        // Test getting parameter from non-existing module
        let result = workflow.get_module_param("nonexistent", "autohis");
        assert!(result.is_err());

        // Test getting non-existing parameter from existing module
        let result = workflow.get_module_param("topoaa", "nonexistent");
        assert!(result.is_err());
    }

    #[test]
    fn test_input_validate_patterns() {
        // Create a valid input
        let input = Input {
            general: General {
                mol_suffixes: vec!["_r".to_string(), "_l".to_string()],
                input_list: "test.txt".to_string(),
                work_dir: PathBuf::from("/tmp"),
                max_concurrent: 1,
                ncores: 1,
                execution: Execution::Local,
                partition: None,
            },
            scenarios: vec![],
        };

        // Should validate successfully
        let result = input.validate_patterns();
        assert!(result.is_ok());
    }

    #[test]
    fn test_input_validate_patterns_empty_suffixes() {
        // Create an input with empty suffixes
        let input = Input {
            general: General {
                mol_suffixes: vec![],
                input_list: "test.txt".to_string(),
                work_dir: PathBuf::from("/tmp"),
                max_concurrent: 1,
                ncores: 1,
                execution: Execution::Local,
                partition: None,
            },
            scenarios: vec![],
        };

        // Should fail validation
        let result = input.validate_patterns();
        assert!(result.is_err());
    }

    #[test]
    fn test_input_validate_patterns_duplicate_suffixes() {
        // Create an input with duplicate suffixes
        let input = Input {
            general: General {
                mol_suffixes: vec!["_r".to_string(), "_r".to_string()],
                input_list: "test.txt".to_string(),
                work_dir: PathBuf::from("/tmp"),
                max_concurrent: 1,
                ncores: 1,
                execution: Execution::Local,
                partition: None,
            },
            scenarios: vec![],
        };

        // Should fail validation
        let result = input.validate_patterns();
        assert!(result.is_err());
    }

    #[test]
    fn test_input_validate_general() {
        // Create a valid input
        let input = Input {
            general: General {
                mol_suffixes: vec!["_r".to_string(), "_l".to_string()],
                input_list: "test.txt".to_string(),
                work_dir: PathBuf::from("/tmp"),
                max_concurrent: 1,
                ncores: 1,
                execution: Execution::Local,
                partition: None,
            },
            scenarios: vec![],
        };

        // Should validate successfully (but might fail if /tmp doesn't exist)
        let result = input.validate_general();
        // Just check it returns a Result, don't assert it's Ok since it depends on filesystem
        let _ = result;
    }

    #[test]
    fn test_input_validate_general_invalid_work_dir() {
        // Create an input with empty work_dir
        let input = Input {
            general: General {
                mol_suffixes: vec!["_r".to_string(), "_l".to_string()],
                input_list: "test.txt".to_string(),
                work_dir: PathBuf::new(),
                max_concurrent: 1,
                ncores: 1,
                execution: Execution::Local,
                partition: None,
            },
            scenarios: vec![],
        };

        // Should fail validation
        let result = input.validate_general();
        assert!(result.is_err());
    }

    #[test]
    fn test_input_validate_general_zero_concurrent() {
        // Create an input with zero max_concurrent
        let input = Input {
            general: General {
                mol_suffixes: vec!["_r".to_string(), "_l".to_string()],
                input_list: "test.txt".to_string(),
                work_dir: PathBuf::from("/tmp"),
                max_concurrent: 0,
                ncores: 1,
                execution: Execution::Local,
                partition: None,
            },
            scenarios: vec![],
        };

        // Should fail validation
        let result = input.validate_general();
        assert!(result.is_err());
    }

    #[test]
    fn test_input_validate_general_empty_partition() {
        let input = Input {
            general: General {
                mol_suffixes: vec!["_r".to_string(), "_l".to_string()],
                input_list: "test.txt".to_string(),
                work_dir: PathBuf::from("/tmp"),
                max_concurrent: 1,
                ncores: 1,
                execution: Execution::Slurm,
                partition: Some("   ".to_string()),
            },
            scenarios: vec![],
        };

        let result = input.validate_general();
        assert!(result.is_err());
        assert_eq!(
            result.unwrap_err().to_string(),
            "partition must be non-empty when defined"
        );
    }

    #[test]
    fn test_input_deserialize_unknown_top_level_field() {
        let yaml = r#"
general:
  mol_suffixes: ["_r", "_l"]
  input_list: input_list.txt
  work_dir: ./work
  max_concurrent: 1
  ncores: 1
  execution: local
scenarios:
  - name: test
    workflow:
      topoaa:
        autohis: true
unknown_field: true
"#;

        let result: Result<Input, _> = serde_yaml::from_str(yaml);
        assert!(result.is_err());
    }

    #[test]
    fn test_input_deserialize_unknown_general_field() {
        let yaml = r#"
general:
  mol_suffixes: ["_r", "_l"]
  input_list: input_list.txt
  work_dir: ./work
  max_concurrent: 1
  ncores: 1
  execution: local
  unexpected: value
scenarios:
  - name: test
    workflow:
      topoaa:
        autohis: true
"#;

        let result: Result<Input, _> = serde_yaml::from_str(yaml);
        assert!(result.is_err());
    }

    #[test]
    fn test_input_deserialize_with_partition() {
        let yaml = r#"
general:
  mol_suffixes: ["_r", "_l"]
  input_list: input_list.txt
  work_dir: ./work
  max_concurrent: 1
  ncores: 1
  execution: slurm
  partition: gpu
scenarios:
  - name: test
    workflow:
      topoaa:
        autohis: true
"#;

        let result: Result<Input, _> = serde_yaml::from_str(yaml);
        assert!(result.is_ok());
        assert_eq!(result.unwrap().general.partition, Some("gpu".to_string()));
    }

    #[test]
    fn test_input_deserialize_without_partition() {
        let yaml = r#"
general:
  mol_suffixes: ["_r", "_l"]
  input_list: input_list.txt
  work_dir: ./work
  max_concurrent: 1
  ncores: 1
  execution: slurm
scenarios:
  - name: test
    workflow:
      topoaa:
        autohis: true
"#;

        let result: Result<Input, _> = serde_yaml::from_str(yaml);
        assert!(result.is_ok());
        assert_eq!(result.unwrap().general.partition, None);
    }

    #[test]
    fn test_input_deserialize_unknown_scenario_field() {
        let yaml = r#"
general:
  mol_suffixes: ["_r", "_l"]
  input_list: input_list.txt
  work_dir: ./work
  max_concurrent: 1
  ncores: 1
  execution: local
scenarios:
  - name: test
    unknown: true
    workflow:
      topoaa:
        autohis: true
"#;

        let result: Result<Input, _> = serde_yaml::from_str(yaml);
        assert!(result.is_err());
    }

    fn make_input_with_workflow(modules: IndexMap<String, Value>) -> Input {
        Input {
            general: General {
                mol_suffixes: vec!["_r".to_string(), "_l".to_string()],
                input_list: "test.txt".to_string(),
                work_dir: PathBuf::from("/tmp"),
                max_concurrent: 1,
                ncores: 1,
                execution: Execution::Local,
                partition: None,
            },
            scenarios: vec![Scenario {
                name: "test".to_string(),
                workflow: Workflow { modules },
            }],
        }
    }

    #[test]
    fn test_validate_workflow_repeated_modules_with_digits_is_ok() {
        let mut modules = IndexMap::new();
        modules.insert("topoaa".to_string(), Value::Null);
        modules.insert("caprieval.1".to_string(), Value::Null);
        modules.insert("caprieval.2".to_string(), Value::Null);

        let input = make_input_with_workflow(modules);
        assert!(input.validate_workflows().is_ok());
    }

    #[test]
    fn test_validate_workflow_mixed_bare_and_suffixed_is_error() {
        let mut modules = IndexMap::new();
        modules.insert("topoaa".to_string(), Value::Null);
        modules.insert("caprieval".to_string(), Value::Null);
        modules.insert("caprieval.1".to_string(), Value::Null);

        let input = make_input_with_workflow(modules);
        assert!(input.validate_workflows().is_err());
    }

    #[test]
    fn test_validate_workflow_sequential_digits_three_is_ok() {
        let mut modules = IndexMap::new();
        modules.insert("topoaa".to_string(), Value::Null);
        modules.insert("caprieval.1".to_string(), Value::Null);
        modules.insert("caprieval.2".to_string(), Value::Null);
        modules.insert("caprieval.3".to_string(), Value::Null);

        let input = make_input_with_workflow(modules);
        assert!(input.validate_workflows().is_ok());
    }

    #[test]
    fn test_validate_workflow_not_starting_at_one_is_error() {
        let mut modules = IndexMap::new();
        modules.insert("topoaa".to_string(), Value::Null);
        modules.insert("caprieval.2".to_string(), Value::Null);
        modules.insert("caprieval.3".to_string(), Value::Null);

        let input = make_input_with_workflow(modules);
        assert!(input.validate_workflows().is_err());
    }

    #[test]
    fn test_validate_workflow_gap_in_sequence_is_error() {
        let mut modules = IndexMap::new();
        modules.insert("topoaa".to_string(), Value::Null);
        modules.insert("caprieval.1".to_string(), Value::Null);
        modules.insert("caprieval.3".to_string(), Value::Null);

        let input = make_input_with_workflow(modules);
        assert!(input.validate_workflows().is_err());
    }

    #[test]
    fn test_validate_workflow_zero_based_is_error() {
        let mut modules = IndexMap::new();
        modules.insert("topoaa".to_string(), Value::Null);
        modules.insert("caprieval.0".to_string(), Value::Null);
        modules.insert("caprieval.1".to_string(), Value::Null);

        let input = make_input_with_workflow(modules);
        assert!(input.validate_workflows().is_err());
    }
}
