use anyhow::{Context, Result};
use indexmap::IndexMap;
use serde::{Deserialize, Serialize};
use serde_yaml::Value;
use std::path::PathBuf;

#[derive(Debug, Deserialize, Serialize)]
pub struct Input {
    pub general: General,
    pub slurm: Option<Slurm>,
    pub scenarios: Vec<Scenario>,
}

#[derive(Debug, Deserialize, Serialize)]
pub struct General {
    pub executable: PathBuf,
    pub mol_suffixes: Vec<String>,
    pub shape_suffix: String,
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
}
