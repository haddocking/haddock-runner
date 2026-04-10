use anyhow::{Context, Result};
use indexmap::IndexMap;
use serde::{Deserialize, Serialize};
use serde_yaml::Value;
use std::os::unix::fs::PermissionsExt;
use std::path::PathBuf;

#[derive(Debug, Deserialize, Serialize)]
pub struct Input {
    pub general: General,
    pub slurm: Option<Slurm>,
    pub scenarios: Vec<Scenario>,
}

impl Input {
    /// Validates the input configuration.
    /// Checks for required fields, valid paths, and logical consistency.
    pub fn validate(&self) -> Result<()> {
        // Validate executable
        self.validate_executable()?;

        // Validate patterns
        self.validate_patterns()?;

        // Validate general fields
        self.validate_general()?;

        Ok(())
    }

    /// Validates the executable path and permissions.
    fn validate_executable(&self) -> Result<()> {
        if self.general.executable.to_string_lossy().is_empty() {
            anyhow::bail!("executable not defined");
        }

        if !self.general.executable.is_absolute() {
            anyhow::bail!(
                "{} is not an absolute path",
                self.general.executable.display()
            );
        }

        // Check if file exists and is executable
        let metadata = std::fs::metadata(&self.general.executable).with_context(|| {
            format!(
                "Failed to access executable: {}",
                self.general.executable.display()
            )
        })?;

        if metadata.permissions().mode() & 0o111 == 0 {
            anyhow::bail!("executable is not executable");
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
        if self.general.work_dir.is_empty() {
            anyhow::bail!("work_dir not defined in general section");
        }

        if self.general.input_list.is_empty() {
            anyhow::bail!("input_list not defined in general section");
        }

        if self.general.max_concurrent == 0 {
            anyhow::bail!("max_concurrent must be greater than 0");
        }

        Ok(())
    }
}

#[derive(Debug, Deserialize, Serialize)]
pub struct General {
    pub executable: PathBuf,
    pub mol_suffixes: Vec<String>,
    pub shape_suffix: Option<String>,
    pub input_list: String,
    pub work_dir: String,
    pub max_concurrent: u16,
}

#[derive(Debug, Deserialize, Serialize)]
pub struct Slurm {
    pub cpus_per_task: u16,
}

#[derive(Debug, Deserialize, Serialize, PartialEq)]
pub struct Scenario {
    pub name: String,
    pub workflow: Workflow,
}

#[derive(Debug, Deserialize, Serialize, PartialEq)]
pub struct Workflow {
    #[serde(flatten)]
    modules: IndexMap<String, Value>,
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
    use serde_yaml;

    #[test]
    fn test_parse_bm_yaml() {
        let yaml_content = std::fs::read_to_string("example/bm.yml").unwrap();
        let input: Input = serde_yaml::from_str(&yaml_content).unwrap();

        assert_eq!(input.scenarios.len(), 1);
        assert_eq!(input.scenarios[0].name, "true-interface");

        let workflow = &input.scenarios[0].workflow;

        // Test if modules contains the expected
        assert!(workflow.modules.contains_key("topoaa"));
        assert!(workflow.modules.contains_key("rigidbody"));

        let topoaa_params = workflow
            .modules
            .get("topoaa")
            .unwrap()
            .as_mapping()
            .unwrap();
        assert!(
            topoaa_params
                .get("autohis".to_string())
                .unwrap()
                .as_bool()
                .unwrap(),
        );

        let rigidbody_params = workflow
            .modules
            .get("rigidbody")
            .unwrap()
            .as_mapping()
            .unwrap();
        assert_eq!(
            rigidbody_params
                .get("sampling".to_string())
                .unwrap()
                .as_i64()
                .unwrap(),
            10
        );
    }

    #[test]
    fn test_get_value() -> Result<()> {
        let yaml_content = std::fs::read_to_string("example/bm.yml").unwrap();
        let input: Input = serde_yaml::from_str(&yaml_content).unwrap();
        let workflow = &input.scenarios[0].workflow;

        // Test getting a module value
        let topoaa_value = workflow.get_value("topoaa")?;
        assert!(topoaa_value.as_mapping().is_some());

        // Test getting a non-existent key
        let non_existent_value = workflow.get_value("non_existent");
        assert!(non_existent_value.is_err());
        Ok(())
    }

    #[test]
    fn test_get_module_param() -> Result<()> {
        let yaml_content = std::fs::read_to_string("example/bm.yml").unwrap();
        let input: Input = serde_yaml::from_str(&yaml_content).unwrap();
        let workflow = &input.scenarios[0].workflow;

        // Test getting a module parameter
        let autohis_value = workflow.get_module_param("topoaa", "autohis")?;
        assert_eq!(autohis_value.as_bool(), Some(true));

        let sampling_value = workflow.get_module_param("rigidbody", "sampling")?;
        assert_eq!(sampling_value.as_i64(), Some(10));

        // Test getting a non-existent module
        let result = workflow.get_module_param("non_existent_module", "param");
        assert!(result.is_err());

        // Test getting a non-existent parameter
        let result = workflow.get_module_param("topoaa", "non_existent_param");
        assert!(result.is_err());

        Ok(())
    }

    #[test]
    fn test_module_order_preserved() {
        let yaml_content = std::fs::read_to_string("example/bm.yml").unwrap();
        let input: Input = serde_yaml::from_str(&yaml_content).unwrap();
        let workflow = &input.scenarios[0].workflow;

        // Check that modules are in the expected order
        let module_keys: Vec<_> = workflow.modules.keys().collect();
        assert_eq!(
            module_keys,
            vec![
                "topoaa",
                "rigidbody",
                "seletop",
                "flexref",
                "emref",
                "caprieval"
            ]
        );
    }

    #[test]
    fn test_validation() -> Result<()> {
        let yaml_content = std::fs::read_to_string("example/bm.yml").unwrap();
        let input: Input = serde_yaml::from_str(&yaml_content).unwrap();

        // This should pass validation
        input.validate()?;

        Ok(())
    }

    #[test]
    fn test_validation_errors() {
        let yaml_content = std::fs::read_to_string("example/bm.yml").unwrap();
        let mut input: Input = serde_yaml::from_str(&yaml_content).unwrap();

        // Test empty mol_suffixes
        input.general.mol_suffixes = vec![];
        let result = input.validate();
        assert!(result.is_err());

        // Test only one suffix
        input.general.mol_suffixes = vec!["_test".to_string()];
        let result = input.validate();
        assert!(result.is_err());

        // Test duplicate suffixes
        input.general.mol_suffixes = vec!["_test".to_string(), "_test".to_string()];
        let result = input.validate();
        assert!(result.is_err());
    }
}
