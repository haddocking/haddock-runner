use crate::utils::format_toml_value;
use anyhow::Context;
use log::{debug, info};
use std::fs;
use std::io::Write;
use std::path::PathBuf;

use itertools::Itertools;

use crate::input::{Execution, General};
use crate::runner::status::Status;
use crate::runner::{local, slurm};
use crate::{
    dataset::Target,
    input::{Input, Scenario},
};
use regex::Regex;

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
    pub status: Status,
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
            status: Status::Unknown,
        }
    }

    pub fn clean(&mut self) -> anyhow::Result<()> {
        fs::remove_dir_all(&self.wd)?;

        Ok(())
    }

    pub fn setup(&mut self) -> anyhow::Result<()> {
        // info!(
        //     "Setting up job {} in directory: {}",
        //     self.name,
        //     self.wd.display()
        // );

        // Create the working directory
        debug!("Creating working directory: {}", self.wd.display());
        fs::create_dir_all(&self.wd)?;

        // Copy the data
        debug!("Copying data files");
        self.copy_data()?;

        // Write the run.toml file
        debug!("Writing run.toml file");
        self.write_run_toml()?;

        if let Execution::Slurm = self.general.execution {
            debug!("Preparing SLURM job file");
            slurm::prepare_job_file(&self.wd, "haddock3")?;
        }

        // Mark it are ready for execution
        self.status = Status::Prepared;
        info!("> {}", self.name);

        Ok(())
    }

    pub fn run(&mut self) -> anyhow::Result<()> {
        info!("Starting {}", self.name);

        // TODO: Figure out if this job is incomplete and should be restarted

        match self.general.execution {
            Execution::Local => self.run_local(),
            Execution::Slurm => self.run_slurm(),
        }
    }

    pub fn run_slurm(&mut self) -> anyhow::Result<()> {
        // Prepare job file
        let job_script = self.wd.join("job.sh");

        // Submit job
        let output = slurm::submit(&job_script)?;

        // Parse job ID from output (simplified - would need proper parsing)
        let job_id = output.split_whitespace().last().unwrap_or("");

        if job_id.is_empty() {
            anyhow::bail!("Failed to parse job ID from sbatch output");
        }

        // Wait for job completion
        slurm::wait(job_id)?;

        // Update status to Done
        self.status = Status::Done;

        Ok(())
    }

    pub fn run_local(&mut self) -> anyhow::Result<()> {
        // info!("Running {} locally", self.name);

        // Execute haddock3 command in the working directory
        let _ = local::run(&self.wd)?;

        // Update status to Done
        self.status = Status::Done;

        // // Log the execution
        // info!(
        //     "Job '{}' executed successfully. Log: {}",
        //     self.name,
        //     log_path.display()
        // );
        //
        info!("{} done", self.name);

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
            size: self.target.size,
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

        // Add ncores section
        toml_content.push_str(&format!("ncores = {}\n", &self.general.ncores));

        // Add workflow modules from scenario
        for (module_name, module_params) in &self.scenario.workflow.modules {
            toml_content.push_str(&format!("[{}]\n", module_name));

            if let Some(params) = module_params.as_mapping() {
                for (key, value) in params.iter() {
                    if let Some(key_str) = key.as_str() {
                        if key_str.contains("_fname") 
                            // Handle _fname parameters by resolving the file pattern
                            && let Some(pattern) = value.as_str()
                        {
                            if let Some(resolved_path) =
                                self.resolve_fname_pattern(pattern, &all_files)
                            {
                                toml_content
                                    .push_str(&format!("{} = {}\n", key_str, resolved_path));
                            }
                        } else {
                            // Handle regular parameters
                            let value_str = format_toml_value(value);
                            toml_content.push_str(&format!("{} = {}\n", key_str, value_str));
                        }
                    }
                }
            }

            toml_content.push('\n');
        }

        Ok(toml_content)
    }

    /// Resolves _fname patterns to actual file paths
    fn resolve_fname_pattern(&self, pattern_str: &str, files: &[PathBuf]) -> Option<String> {
        let pattern = Regex::new(pattern_str).ok()?;
        let mut matching_file: Option<PathBuf> = None;

        for file in files {
            if let Some(file_name) = file.file_name().and_then(|n| n.to_str())
                && pattern.is_match(file_name)
            {
                if matching_file.is_some() {
                    // Multiple matches - return None to indicate ambiguity
                    return None;
                }
                matching_file = Some(file.clone());
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
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::input::{Execution, General, Scenario, Workflow};
    use crate::runner::status::Status;
    use indexmap::IndexMap;
    use std::path::PathBuf;
    use tempfile::tempdir;

    #[test]
    fn test_job_new() {
        let general = General {
            mol_suffixes: vec!["_r".to_string(), "_l".to_string()],
            shape_suffix: None,
            input_list: "test.txt".to_string(),
            work_dir: PathBuf::from("/tmp"),
            max_concurrent: 1,
            ncores: 1,
            execution: Execution::Local,
        };

        let scenario = Scenario {
            name: "test_scenario".to_string(),
            workflow: Workflow {
                modules: IndexMap::new(),
            },
        };

        let target = Target {
            id: "test_target".to_string(),
            molecules: vec![PathBuf::from("mol.pdb")],
            restraints: vec![PathBuf::from("restraint.tbl")],
            toppar: vec![PathBuf::from("toppar.top")],
            misc: vec![PathBuf::from("misc.txt")],
            shape: None,
            size: 0,
        };

        let job = Job::new(general, scenario, target);

        assert_eq!(job.name, "test_target-test_scenario");
        assert_eq!(job.wd, PathBuf::from("/tmp/test_scenario/test_target"));
        assert_eq!(job.target.id, "test_target");
        assert_eq!(job.scenario.name, "test_scenario");
    }

    #[test]
    fn test_create_jobs() {
        let input = Input {
            general: General {
                mol_suffixes: vec!["_r".to_string(), "_l".to_string()],
                shape_suffix: None,
                input_list: "test.txt".to_string(),
                work_dir: PathBuf::from("/tmp"),
                max_concurrent: 1,
                ncores: 1,
                execution: Execution::Local,
            },
            scenarios: vec![
                Scenario {
                    name: "scenario1".to_string(),
                    workflow: Workflow {
                        modules: IndexMap::new(),
                    },
                },
                Scenario {
                    name: "scenario2".to_string(),
                    workflow: Workflow {
                        modules: IndexMap::new(),
                    },
                },
            ],
        };

        let targets = vec![
            Target {
                id: "target1".to_string(),
                molecules: vec![PathBuf::from("mol1.pdb")],
                restraints: vec![],
                toppar: vec![],
                misc: vec![],
                shape: None,
                size: 0,
            },
            Target {
                id: "target2".to_string(),
                molecules: vec![PathBuf::from("mol2.pdb")],
                restraints: vec![],
                toppar: vec![],
                misc: vec![],
                shape: None,
                size: 0,
            },
        ];

        let jobs = create_jobs(input, targets);

        // Should create 2 scenarios × 2 targets = 4 jobs
        assert_eq!(jobs.len(), 4);
        assert_eq!(jobs[0].name, "target1-scenario1");
        assert_eq!(jobs[1].name, "target2-scenario1");
        assert_eq!(jobs[2].name, "target1-scenario2");
        assert_eq!(jobs[3].name, "target2-scenario2");
    }

    #[test]
    fn test_job_clean() {
        let temp_dir = tempdir().unwrap();
        let work_dir = temp_dir.path().join("test_work");

        // Create a job with a work directory
        let general = General {
            mol_suffixes: vec!["_r".to_string(), "_l".to_string()],
            shape_suffix: None,
            input_list: "test.txt".to_string(),
            work_dir: work_dir.clone(),
            max_concurrent: 1,
            ncores: 1,
            execution: Execution::Local,
        };

        let scenario = Scenario {
            name: "test_scenario".to_string(),
            workflow: Workflow {
                modules: IndexMap::new(),
            },
        };

        let target = Target {
            id: "test_target".to_string(),
            molecules: vec![],
            restraints: vec![],
            toppar: vec![],
            misc: vec![],
            shape: None,
            size: 0,
        };

        let mut job = Job::new(general, scenario, target);

        // Create the work directory
        std::fs::create_dir_all(&job.wd).unwrap();

        // Verify directory exists
        assert!(job.wd.exists());

        // Clean the job
        let result = job.clean();
        assert!(result.is_ok());

        // Verify directory was removed
        assert!(!job.wd.exists());
    }

    #[test]
    fn test_resolve_fname_pattern() {
        let temp_dir = tempdir().unwrap();
        let temp_path = temp_dir.path();

        // Create test files
        let file1 = temp_path.join("test_file1.pdb");
        let file2 = temp_path.join("test_file2.pdb");
        let file3 = temp_path.join("other_file.pdb");

        std::fs::write(&file1, "content1").unwrap();
        std::fs::write(&file2, "content2").unwrap();
        std::fs::write(&file3, "content3").unwrap();

        let files = vec![file1.clone(), file2.clone(), file3.clone()];

        let job = Job {
            name: "test".to_string(),
            status: Status::Unknown,
            wd: temp_path.to_path_buf(),
            target: Target {
                id: "test".to_string(),
                molecules: vec![],
                restraints: vec![],
                toppar: vec![],
                misc: vec![],
                shape: None,
                size: 0,
            },
            scenario: Scenario {
                name: "test".to_string(),
                workflow: Workflow {
                    modules: IndexMap::new(),
                },
            },
            general: General {
                mol_suffixes: vec!["_r".to_string(), "_l".to_string()],
                shape_suffix: None,
                input_list: "test.txt".to_string(),
                work_dir: temp_path.to_path_buf(),
                max_concurrent: 1,
                ncores: 1,
                execution: Execution::Local,
            },
        };

        // Test pattern matching
        let result = job.resolve_fname_pattern("test_file1\\.pdb", &files);
        assert!(result.is_some());
        assert!(result.unwrap().contains("test_file1.pdb"));

        // Test pattern with no match
        let result = job.resolve_fname_pattern("nonexistent.*\\.pdb", &files);
        assert!(result.is_none());
    }

    #[test]
    fn test_get_all_target_files() {
        let temp_dir = tempdir().unwrap();
        let temp_path = temp_dir.path();

        // Create test files
        let mol_file = temp_path.join("mol.pdb");
        let restraint_file = temp_path.join("restraint.tbl");
        let toppar_file = temp_path.join("toppar.top");
        let misc_file = temp_path.join("misc.txt");
        let shape_file = temp_path.join("shape.pdb");

        std::fs::write(&mol_file, "content").unwrap();
        std::fs::write(&restraint_file, "content").unwrap();
        std::fs::write(&toppar_file, "content").unwrap();
        std::fs::write(&misc_file, "content").unwrap();
        std::fs::write(&shape_file, "content").unwrap();

        let target = Target {
            id: "test".to_string(),
            molecules: vec![mol_file.clone()],
            restraints: vec![restraint_file.clone()],
            toppar: vec![toppar_file.clone()],
            misc: vec![misc_file.clone()],
            shape: Some(shape_file.clone()),
            size: 0,
        };

        let job = Job {
            name: "test".to_string(),
            status: Status::Unknown,
            wd: temp_path.to_path_buf(),
            target: target.clone(),
            scenario: Scenario {
                name: "test".to_string(),
                workflow: Workflow {
                    modules: IndexMap::new(),
                },
            },
            general: General {
                mol_suffixes: vec!["_r".to_string(), "_l".to_string()],
                shape_suffix: None,
                input_list: "test.txt".to_string(),
                work_dir: temp_path.to_path_buf(),
                max_concurrent: 1,
                ncores: 1,
                execution: Execution::Local,
            },
        };

        let all_files = job.get_all_target_files();

        // Should contain all files
        assert_eq!(all_files.len(), 5);
        assert!(all_files.contains(&mol_file));
        assert!(all_files.contains(&restraint_file));
        assert!(all_files.contains(&toppar_file));
        assert!(all_files.contains(&misc_file));
        assert!(all_files.contains(&shape_file));
    }
}
