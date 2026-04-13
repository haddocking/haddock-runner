use anyhow::Context;
use std::fs;
use std::io::Write;
use std::path::PathBuf;

use itertools::Itertools;
use serde_yaml::Value;

use crate::input::General;
use crate::runner::status::JobStatus;
use crate::{
    dataset::Target,
    input::{Input, Scenario},
};
use regex::Regex;

// type Job struct {
// 	ID         string
// 	Path       string
// 	Params     map[string]interface{}
// 	Restraints input.Airs
// 	Toppar     input.TopologyParams
// 	Status     string
// }
//
//

pub fn create_jobs(input: Input, targets: Vec<Target>) -> Vec<Job> {
    input
        .scenarios
        .into_iter()
        .cartesian_product(targets.iter())
        .map(|(scenario, target)| Job::new(input.general.clone(), scenario, target.clone()))
        .collect::<Vec<Job>>()
}

#[derive(Debug, Clone)]
pub struct Job {
    pub name: String,
    pub status: JobStatus,
    pub wd: PathBuf,
    pub target: Target,
    pub scenario: Scenario,
    pub general: General,
}

impl Job {
    fn new(general: General, scenario: Scenario, target: Target) -> Self {
        let name = target.id.to_string() + "-" + &scenario.name;
        let wd = general.work_dir.join(&scenario.name).join(&target.id);
        Job {
            name,
            wd,
            target,
            scenario,
            general,
            status: JobStatus::Unknown,
        }
    }

    pub fn clean(&mut self) -> anyhow::Result<()> {
        fs::remove_dir_all(&self.wd)?;

        Ok(())
    }

    pub fn setup(&mut self) -> anyhow::Result<()> {
        // Create the working directory
        fs::create_dir_all(&self.wd)?;

        // Copy the data
        self.copy_data()?;

        // Write the run.toml file
        self.write_run_toml()?;

        // Mark it are ready for execution
        self.status = JobStatus::Prepared;

        Ok(())
    }

    pub fn run(&mut self) -> anyhow::Result<()> {
        // Execute haddock3 command in the working directory
        let log_path = crate::runner::execute::run("haddock3", "run.toml", &self.wd)?;

        // Update status to Done
        self.status = JobStatus::Done;

        // Log the execution
        println!(
            "Job '{}' executed successfully. Log: {}",
            self.name,
            log_path.display()
        );

        Ok(())
    }

    fn copy_data(&mut self) -> anyhow::Result<()> {
        // Helper function to copy file and update path
        let copy_and_update = |file: &PathBuf| -> anyhow::Result<PathBuf> {
            let dest = self.wd.join(file.file_name().unwrap());
            fs::copy(file, &dest).with_context(|| {
                format!("failed to copy {} to {}", file.display(), dest.display())
            })?;
            Ok(dest)
        };

        // Helper function to copy a collection of files
        let copy_collection = |files: &[PathBuf]| -> anyhow::Result<Vec<PathBuf>> {
            files
                .iter()
                .map(copy_and_update)
                .collect::<anyhow::Result<Vec<_>>>()
        };

        // Copy all file collections
        let molecules = copy_collection(&self.target.molecules)?;
        let restraints = copy_collection(&self.target.restraints)?;
        let toppar = copy_collection(&self.target.toppar)?;
        let misc = copy_collection(&self.target.misc)?;

        // Copy shape file if present
        let shape = &self
            .target
            .shape
            .as_ref()
            .map(copy_and_update)
            .transpose()?;

        // Create new organized target
        let organized_target = Target {
            id: self.target.id.clone(),
            molecules,
            restraints,
            toppar,
            misc,
            shape: shape.clone(),
        };

        self.target = organized_target;

        Ok(())
    }

    fn write_run_toml(&self) -> anyhow::Result<()> {
        let toml_path = self.wd.join("run.toml");
        let mut toml_file = fs::File::create(&toml_path)?;

        // Write the TOML content
        let toml_content = self.generate_run_toml()?;
        toml_file.write_all(toml_content.as_bytes())?;

        Ok(())
    }

    fn generate_run_toml(&self) -> anyhow::Result<String> {
        let mut toml_content = String::new();
        let all_files = self.get_all_target_files();

        // Add general HADDOCK3 configuration
        toml_content.push_str("run_dir = \"run1\"\n\n");

        // Add molecules section
        toml_content.push_str("molecules = [\n");
        for molecule in &self.target.molecules {
            if let Some(file_name) = molecule.file_name() {
                toml_content.push_str(&format!("    \"{}\",\n", file_name.to_string_lossy()));
            }
        }
        toml_content.push_str("]\n\n");

        // TODO: Add the general options

        // Add ncores section
        toml_content.push_str(&format!("ncores = {}\n", &self.general.ncores));

        // Add workflow modules from scenario
        for (module_name, module_params) in &self.scenario.workflow.modules {
            toml_content.push_str(&format!("[{}]\n", module_name));

            if let Some(params) = module_params.as_mapping() {
                for (key, value) in params.iter() {
                    if let Some(key_str) = key.as_str() {
                        if key_str.contains("_fname") {
                            // Handle _fname parameters by resolving the file pattern
                            if let Some(pattern) = value.as_str() {
                                if let Some(resolved_path) =
                                    self.resolve_fname_pattern(pattern, &all_files)
                                {
                                    toml_content
                                        .push_str(&format!("{} = {}\n", key_str, resolved_path));
                                }
                            }
                        } else {
                            // Handle regular parameters
                            let value_str = self.format_toml_value(value);
                            toml_content.push_str(&format!("{} = {}\n", key_str, value_str));
                        }
                    }
                }
            }

            toml_content.push('\n');
        }

        Ok(toml_content)
    }

    fn format_toml_value(&self, value: &Value) -> String {
        match value {
            Value::Bool(b) => b.to_string(),
            Value::Number(n) => n.to_string(),
            Value::String(s) => format!("\"{}\"", s),
            Value::Sequence(seq) => {
                let items: Vec<String> = seq.iter().map(|v| self.format_toml_value(v)).collect();
                format!("[{}]", items.join(", "))
            }
            _ => "null".to_string(), // Fallback for other types
        }
    }

    /// Resolves _fname patterns to actual file paths
    fn resolve_fname_pattern(&self, pattern_str: &str, files: &[PathBuf]) -> Option<String> {
        let pattern = Regex::new(pattern_str).ok()?;
        let mut matching_file: Option<PathBuf> = None;

        for file in files {
            if let Some(file_name) = file.file_name().and_then(|n| n.to_str()) {
                if pattern.is_match(file_name) {
                    if matching_file.is_some() {
                        // Multiple matches - return None to indicate ambiguity
                        return None;
                    }
                    matching_file = Some(file.clone());
                }
            }
        }

        matching_file.map(|path| format!("\"{}\"", path.file_name().unwrap().to_string_lossy()))
    }

    /// Collects all available files from the target
    fn get_all_target_files(&self) -> Vec<PathBuf> {
        let mut files = Vec::new();

        files.extend(self.target.molecules.clone());
        files.extend(self.target.restraints.clone());
        files.extend(self.target.toppar.clone());
        files.extend(self.target.misc.clone());

        if let Some(shape) = &self.target.shape {
            files.push(shape.clone());
        }

        files
    }
    //
    // fn execute() {
    //     todo!()
    // }
}
