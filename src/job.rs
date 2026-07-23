use crate::job::fs::canonicalize;
use crate::runner::slurm::SlurmJob;
use crate::utils::{find_haddock3_executable, format_toml_value};
use anyhow::Context;
use log::debug;
use std::fs;
use std::io::Write;
use std::path::PathBuf;

use indexmap::IndexMap;
use itertools::Itertools;
use serde_yaml::Value;

use crate::input::{DIGIT_SUFFIX_RE, Execution, General};
use crate::runner::local;
use crate::runner::status::Status;
use crate::{
    dataset::Target,
    input::{Input, Scenario},
};
use regex::Regex;

pub const JOB_FILENAME: &str = "job.sh";
pub const WORKFLOW_FILENAME: &str = "run.toml";

/// Create jobs from input scenarios and targets
///
/// This function creates a Cartesian product of all scenarios and targets,
/// generating a separate job for each combination.
///
/// # Arguments
///
/// * `input` - The input configuration containing scenarios
/// * `targets` - Vector of targets to process
///
/// # Returns
///
/// * `Vec<Job>` - Vector of created jobs
pub fn create_jobs(input: Input, targets: Vec<Target>) -> Vec<Job> {
    input
        .scenarios
        .into_iter()
        .cartesian_product(targets.iter())
        .map(|(scenario, target)| Job::new(input.general.clone(), scenario, target.clone()))
        .collect::<Vec<Job>>()
}

/// Outcome of a completed `Job::run()` call
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum RunOutcome {
    /// The job actually executed haddock3
    Executed,
    /// The job was already done and was skipped
    Skipped,
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
    /// Create a new Job instance
    ///
    /// This method creates a new Job by combining general configuration,
    /// a specific scenario, and a target. The job name is formed by combining
    /// the target ID and scenario name.
    ///
    /// # Arguments
    ///
    /// * `general` - General configuration
    /// * `scenario` - Scenario to execute
    /// * `target` - Target to process
    ///
    /// # Returns
    ///
    /// * `Self` - Newly created Job instance
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

    /// Clean up job working directory
    ///
    /// This method removes the entire working directory and all its contents.
    ///
    /// # Returns
    ///
    /// * `anyhow::Result<()>` - Ok if cleanup successful, error otherwise
    pub fn clean(&mut self) -> anyhow::Result<()> {
        fs::remove_dir_all(&self.wd)?;

        Ok(())
    }

    /// Set up job working directory and files
    ///
    /// This method prepares the job for execution by creating the working directory,
    /// copying all required input files, generating the run.toml configuration file,
    /// and optionally preparing SLURM job files.
    ///
    /// # Returns
    ///
    /// * `anyhow::Result<()>` - Ok if setup successful, error otherwise
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
            self.prepare_job_file()?;
        }

        // Mark it are ready for execution
        self.status = Status::Prepared;
        debug!("> {}", self.name);

        Ok(())
    }

    /// Run the job
    ///
    /// This method executes the job using the configured execution method
    /// (local or SLURM) and updates the job status upon completion.
    ///
    /// # Returns
    ///
    /// * `anyhow::Result<RunOutcome>` - `Executed` if the job actually ran,
    ///   `Skipped` if it was already done, error otherwise
    pub fn run(&mut self) -> anyhow::Result<RunOutcome> {
        // info!("Starting {}", self.name);

        // Check the status of the job
        self.update_status()?;

        match self.status {
            // Job is incomplete, clean it up
            Status::Incomplete => {
                debug!("job {} is incomplete, cleaning and setup again", self.name);
                self.clean()?;
                self.setup()?;
            }
            // Job was already done, skip it by returning
            Status::Done => {
                debug!("job {} is complete, skipping", self.name);
                return Ok(RunOutcome::Skipped);
            }
            // All other, move ahead
            _ => {}
        }

        // Re-check the haddock3 version against the one recorded when
        // the benchmark started, this will catch a mid-run version swap
        let checksum_file = self.general.work_dir.join("checksum.json");
        crate::checksum::validate_haddock3_version_live(&checksum_file)?;

        match self.general.execution {
            Execution::Local => self.run_local(),
            Execution::Slurm => self.run_slurm(),
        }?;

        Ok(RunOutcome::Executed)
    }

    /// Check the status of the job
    ///
    /// This method checks the status of a job.
    ///
    /// # Returns
    ///
    /// * `anyhow::Result<()>` - Ok if check completes successfully, error otherwise
    fn update_status(&mut self) -> anyhow::Result<()> {
        // Check if working directory exists
        if !self.wd.exists() {
            // Directory doesn't exist, this is a new job
            self.status = Status::Unknown;
            return Ok(());
        }

        // Check if job is prepared
        let run_toml = self.wd.join(WORKFLOW_FILENAME);
        if run_toml.exists() {
            self.status = Status::Prepared
        }

        // Check if it has a log file
        let log_file = self.wd.join("run1/log");
        if log_file.exists() {
            // Read it and find some match words
            let contents = fs::read_to_string(log_file)?;
            if contents.contains("Finished at") {
                self.status = Status::Done
            } else {
                self.status = Status::Incomplete
            }
        }

        Ok(())
    }

    /// Run the job using SLURM
    ///
    /// This method submits and monitors a SLURM job for this HADDOCK run.
    ///
    /// # Returns
    ///
    /// * `anyhow::Result<()>` - Ok if SLURM job completes successfully, error otherwise
    pub fn run_slurm(&mut self) -> anyhow::Result<()> {
        // Create a SlurmJob
        let mut slurm_job = SlurmJob::new(self.wd.clone());

        // Run
        slurm_job.run()?;

        // Update status to Done
        self.status = Status::Done;

        Ok(())
    }

    /// Run the job locally
    ///
    /// This method executes the HADDOCK run directly on the local machine
    /// using the haddock3 command.
    ///
    /// # Returns
    ///
    /// * `anyhow::Result<()>` - Ok if local execution completes successfully, error otherwise
    pub fn run_local(&mut self) -> anyhow::Result<()> {
        // info!("Running {} locally", self.name);

        // Execute haddock3 command in the working directory
        let _ = local::run(&self.wd)?;

        // Update status to Done
        self.status = Status::Done;

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
        let toml_path = self.wd.join(WORKFLOW_FILENAME);
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
        toml_content.push_str("run_dir = \"run1\"\nmode = \"local\"\n");

        // Add molecules section
        toml_content.push_str("molecules = [\n");
        for molecule in &self.target.molecules {
            if let Some(file_name) = molecule.file_name() {
                toml_content.push_str(&format!("    \"{}\",\n", file_name.to_string_lossy()));
            }
        }
        toml_content.push_str("]\n\n");

        // Add ncores section
        toml_content.push_str(&format!("ncores = {}\n", self.general.ncores));

        if let Some(value) = self.general.preprocess {
            toml_content.push_str(&format!("preprocess = {}\n", value));
        }
        if let Some(value) = self.general.postprocess {
            toml_content.push_str(&format!("postprocess = {}\n", value));
        }
        if let Some(value) = self.general.gen_archive {
            toml_content.push_str(&format!("gen_archive = {}\n", value));
        }

        // Add workflow modules from scenario
        for (module_name, module_params) in &self.scenario.workflow.modules {
            // Modules with a .digit suffix are quoted so haddock3 treats the dot as a literal
            // character rather than a TOML nested-table separator.
            let header = if DIGIT_SUFFIX_RE.is_match(module_name) {
                format!("[\"{module_name}\"]\n")
            } else {
                format!("[{module_name}]\n")
            };
            toml_content.push_str(&header);

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
    /// Prepare SLURM job submission script
    ///
    /// This method creates a SLURM job script with the appropriate header
    /// and commands to execute the HADDOCK run.
    ///
    /// # Returns
    ///
    /// * `anyhow::Result<()>` - Ok if job script created successfully, error otherwise
    pub fn prepare_job_file(&self) -> anyhow::Result<()> {
        // Create SLURM job script
        let job_script = self.wd.join(JOB_FILENAME);
        let mut file = fs::File::create(&job_script)?;

        let absolute_wd = canonicalize(&self.wd).unwrap_or_else(|_| self.wd.to_path_buf());

        // Write SLURM header
        let header = self.generate_slurm_header();

        let prologue = self.format_slurm_prologue();

        // Write job body
        let body = format!(
            "cd {}\n{} {}\n",
            absolute_wd.display(),
            find_haddock3_executable()?,
            WORKFLOW_FILENAME
        );

        file.write_all(header.as_bytes())?;
        file.write_all(prologue.as_bytes())?;
        file.write_all(body.as_bytes())?;

        Ok(())
    }

    fn format_slurm_prologue(&self) -> String {
        let mut prologue = String::new();
        if let Some(commands) = &self.general.slurm_prologue {
            prologue.push_str(commands);
            if !commands.ends_with('\n') {
                prologue.push('\n');
            }
        }
        prologue
    }

    fn generate_slurm_header(&self) -> String {
        // `Some(value)` renders as `#SBATCH --key=value`, `None` as a bare `#SBATCH --key`
        let mut fields: IndexMap<String, Option<String>> = IndexMap::new();

        // cpus-per-task is always derived from ncores and cannot be overridden.
        fields.insert(
            "cpus-per-task".to_string(),
            Some(self.general.ncores.to_string()),
        );

        if let Some(partition) = &self.general.partition {
            fields.insert("partition".to_string(), Some(partition.clone()));
        }

        // Parse fields and assign its proper types
        if let Some(slurm_header) = &self.general.slurm_header {
            for (key, value) in slurm_header {
                // NOTE: `cpus-per-task` is a special field, its tied to haddock's `ncores`
                //  so here on purpose we do not let users override it
                if key == "cpus-per-task" {
                    log::warn!(
                        "slurm_header field 'cpus-per-task' is ignored: it is derived from ncores ({}) and cannot be overridden",
                        self.general.ncores
                    );
                    continue;
                }

                // Match values
                match value {
                    // Not provided / explicitly unset: drop the entry and let SLURM default apply.
                    Value::Null | Value::Bool(false) => {
                        fields.shift_remove(key);
                    }
                    // Empty
                    Value::String(s) if s.trim().is_empty() => {
                        fields.shift_remove(key);
                    }
                    // Bare no-argument flag, e.g. `exclusive: true` -> `#SBATCH --exclusive`.
                    Value::Bool(true) => {
                        fields.insert(key.clone(), None);
                    }
                    Value::String(s) => {
                        fields.insert(key.clone(), Some(s.clone()));
                    }
                    Value::Number(n) => {
                        fields.insert(key.clone(), Some(n.to_string()));
                    }
                    // Sequences/mappings are rejected by `Input::validate`; ignore defensively.
                    Value::Sequence(_) | Value::Mapping(_) | Value::Tagged(_) => {}
                }
            }
        }

        let mut header = "#!/bin/bash\n".to_string();
        for (key, value) in &fields {
            match value {
                // key=value
                Some(v) => header.push_str(&format!("#SBATCH --{key}={v}\n")),
                // bare
                None => header.push_str(&format!("#SBATCH --{key}\n")),
            }
        }

        header
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
    use serde_yaml::Value;
    use std::path::PathBuf;
    use tempfile::tempdir;

    #[test]
    fn test_job_new() {
        let general = General {
            mol_suffixes: vec!["_r".to_string(), "_l".to_string()],
            input_list: "test.txt".to_string(),
            work_dir: PathBuf::from("/tmp"),
            max_concurrent: 1,
            ncores: 1,
            execution: Execution::Local,
            partition: None,
            preprocess: None,
            postprocess: None,
            gen_archive: None,
            slurm_header: None,
            slurm_prologue: None,
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
                input_list: "test.txt".to_string(),
                work_dir: PathBuf::from("/tmp"),
                max_concurrent: 1,
                ncores: 1,
                execution: Execution::Local,
                partition: None,
                preprocess: None,
                postprocess: None,
                gen_archive: None,
                slurm_header: None,
                slurm_prologue: None,
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
            input_list: "test.txt".to_string(),
            work_dir: work_dir.clone(),
            max_concurrent: 1,
            ncores: 1,
            execution: Execution::Local,
            partition: None,
            preprocess: None,
            postprocess: None,
            gen_archive: None,
            slurm_header: None,
            slurm_prologue: None,
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
                input_list: "test.txt".to_string(),
                work_dir: temp_path.to_path_buf(),
                max_concurrent: 1,
                ncores: 1,
                execution: Execution::Local,
                partition: None,
                preprocess: None,
                postprocess: None,
                gen_archive: None,
                slurm_header: None,
                slurm_prologue: None,
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
                input_list: "test.txt".to_string(),
                work_dir: temp_path.to_path_buf(),
                max_concurrent: 1,
                ncores: 1,
                execution: Execution::Local,
                partition: None,
                preprocess: None,
                postprocess: None,
                gen_archive: None,
                slurm_header: None,
                slurm_prologue: None,
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

    #[test]
    fn test_generate_slurm_header_without_partition() {
        let job = Job {
            name: "test".to_string(),
            status: Status::Unknown,
            wd: PathBuf::from("/tmp"),
            target: Target {
                id: "target".to_string(),
                molecules: vec![],
                restraints: vec![],
                toppar: vec![],
                misc: vec![],
                shape: None,
                size: 0,
            },
            scenario: Scenario {
                name: "scenario".to_string(),
                workflow: Workflow {
                    modules: IndexMap::new(),
                },
            },
            general: General {
                mol_suffixes: vec!["_r".to_string(), "_l".to_string()],
                input_list: "test.txt".to_string(),
                work_dir: PathBuf::from("/tmp"),
                max_concurrent: 1,
                ncores: 4,
                execution: Execution::Slurm,
                partition: None,
                preprocess: None,
                postprocess: None,
                gen_archive: None,
                slurm_header: None,
                slurm_prologue: None,
            },
        };

        let header = job.generate_slurm_header();
        assert_eq!(header, "#!/bin/bash\n#SBATCH --cpus-per-task=4\n");
    }

    #[test]
    fn test_generate_slurm_header_with_partition() {
        let job = Job {
            name: "test".to_string(),
            status: Status::Unknown,
            wd: PathBuf::from("/tmp"),
            target: Target {
                id: "target".to_string(),
                molecules: vec![],
                restraints: vec![],
                toppar: vec![],
                misc: vec![],
                shape: None,
                size: 0,
            },
            scenario: Scenario {
                name: "scenario".to_string(),
                workflow: Workflow {
                    modules: IndexMap::new(),
                },
            },
            general: General {
                mol_suffixes: vec!["_r".to_string(), "_l".to_string()],
                input_list: "test.txt".to_string(),
                work_dir: PathBuf::from("/tmp"),
                max_concurrent: 1,
                ncores: 4,
                execution: Execution::Slurm,
                partition: Some("gpu".to_string()),
                preprocess: None,
                postprocess: None,
                gen_archive: None,
                slurm_header: None,
                slurm_prologue: None,
            },
        };

        let header = job.generate_slurm_header();
        assert_eq!(
            header,
            "#!/bin/bash\n#SBATCH --cpus-per-task=4\n#SBATCH --partition=gpu\n"
        );
    }

    fn make_slurm_job(
        partition: Option<String>,
        slurm_header: Option<IndexMap<String, Value>>,
        slurm_prologue: Option<String>,
    ) -> Job {
        Job {
            name: "test".to_string(),
            status: Status::Unknown,
            wd: PathBuf::from("/tmp"),
            target: Target {
                id: "target".to_string(),
                molecules: vec![],
                restraints: vec![],
                toppar: vec![],
                misc: vec![],
                shape: None,
                size: 0,
            },
            scenario: Scenario {
                name: "scenario".to_string(),
                workflow: Workflow {
                    modules: IndexMap::new(),
                },
            },
            general: General {
                mol_suffixes: vec!["_r".to_string(), "_l".to_string()],
                input_list: "test.txt".to_string(),
                work_dir: PathBuf::from("/tmp"),
                max_concurrent: 1,
                ncores: 4,
                execution: Execution::Slurm,
                partition,
                preprocess: None,
                postprocess: None,
                gen_archive: None,
                slurm_header,
                slurm_prologue,
            },
        }
    }

    #[test]
    fn test_generate_slurm_header_with_new_field() {
        let mut slurm_header = IndexMap::new();
        slurm_header.insert("nodes".to_string(), Value::Number(1.into()));
        slurm_header.insert(
            "account".to_string(),
            Value::String("project_XXXXXX".to_string()),
        );

        let job = make_slurm_job(None, Some(slurm_header), None);
        let header = job.generate_slurm_header();
        let expected = "#!/bin/bash\n\
#SBATCH --cpus-per-task=4\n\
#SBATCH --nodes=1\n\
#SBATCH --account=project_XXXXXX\n";
        assert_eq!(header, expected);
    }

    #[test]
    fn test_generate_slurm_header_cpus_per_task_cannot_be_overridden() {
        // cpus-per-task is always derived from general.ncores; a user-supplied
        // value in slurm_header is ignored (a warning is logged instead).
        let mut slurm_header = IndexMap::new();
        slurm_header.insert("cpus-per-task".to_string(), Value::Number(8.into()));

        let job = make_slurm_job(None, Some(slurm_header), None);
        let header = job.generate_slurm_header();
        assert_eq!(header, "#!/bin/bash\n#SBATCH --cpus-per-task=4\n");
    }

    #[test]
    fn test_generate_slurm_header_slurm_header_overrides_partition_field() {
        let mut slurm_header = IndexMap::new();
        slurm_header.insert("partition".to_string(), Value::String("gpu".to_string()));

        let job = make_slurm_job(Some("small".to_string()), Some(slurm_header), None);
        let header = job.generate_slurm_header();
        let expected = "#!/bin/bash\n\
#SBATCH --cpus-per-task=4\n\
#SBATCH --partition=gpu\n";
        assert_eq!(header, expected);
    }

    #[test]
    fn test_generate_slurm_header_null_value_is_skipped() {
        let mut slurm_header = IndexMap::new();
        slurm_header.insert("partition".to_string(), Value::Null);

        let job = make_slurm_job(Some("small".to_string()), Some(slurm_header), None);
        let header = job.generate_slurm_header();
        assert_eq!(header, "#!/bin/bash\n#SBATCH --cpus-per-task=4\n");
    }

    #[test]
    fn test_generate_slurm_header_false_value_is_skipped() {
        let mut slurm_header = IndexMap::new();
        slurm_header.insert("exclusive".to_string(), Value::Bool(false));

        let job = make_slurm_job(None, Some(slurm_header), None);
        let header = job.generate_slurm_header();
        assert!(!header.contains("exclusive"));
    }

    #[test]
    fn test_generate_slurm_header_empty_string_value_is_skipped() {
        let mut slurm_header = IndexMap::new();
        slurm_header.insert("qos".to_string(), Value::String("   ".to_string()));

        let job = make_slurm_job(None, Some(slurm_header), None);
        let header = job.generate_slurm_header();
        assert!(!header.contains("qos"));
    }

    #[test]
    fn test_generate_slurm_header_true_value_renders_bare_flag() {
        let mut slurm_header = IndexMap::new();
        slurm_header.insert("exclusive".to_string(), Value::Bool(true));

        let job = make_slurm_job(None, Some(slurm_header), None);
        let header = job.generate_slurm_header();
        assert!(
            header.contains("#SBATCH --exclusive\n"),
            "expected bare --exclusive flag but got:\n{header}"
        );
        assert!(!header.contains("--exclusive="));
    }

    fn make_job_with_modules(modules: IndexMap<String, Value>) -> Job {
        Job {
            name: "test".to_string(),
            status: Status::Unknown,
            wd: PathBuf::from("/tmp"),
            target: Target {
                id: "test".to_string(),
                molecules: vec![PathBuf::from("mol.pdb")],
                restraints: vec![],
                toppar: vec![],
                misc: vec![],
                shape: None,
                size: 0,
            },
            scenario: Scenario {
                name: "test".to_string(),
                workflow: Workflow { modules },
            },
            general: General {
                mol_suffixes: vec!["_r".to_string(), "_l".to_string()],
                input_list: "test.txt".to_string(),
                work_dir: PathBuf::from("/tmp"),
                max_concurrent: 1,
                ncores: 1,
                execution: Execution::Local,
                partition: None,
                preprocess: None,
                postprocess: None,
                gen_archive: None,
                slurm_header: None,
                slurm_prologue: None,
            },
        }
    }

    #[test]
    fn test_generate_run_toml_repeated_modules() {
        let mut modules = IndexMap::new();
        modules.insert("topoaa".to_string(), Value::Null);
        modules.insert("caprieval.1".to_string(), Value::Null);
        modules.insert("caprieval.2".to_string(), Value::Null);

        let job = make_job_with_modules(modules);
        let toml = job.generate_run_toml().unwrap();

        assert!(
            toml.contains("[\"caprieval.1\"]"),
            "expected [\"caprieval.1\"] but got:\n{toml}"
        );
        assert!(
            toml.contains("[\"caprieval.2\"]"),
            "expected [\"caprieval.2\"] but got:\n{toml}"
        );
        assert!(
            !toml.contains("[caprieval.1]"),
            "unquoted [caprieval.1] should not appear"
        );
    }

    #[test]
    fn test_generate_run_toml_single_suffixed_module() {
        // A module with a .digit suffix appearing alone should still be quoted
        let mut modules = IndexMap::new();
        modules.insert("topoaa".to_string(), Value::Null);
        modules.insert("emscoring.1".to_string(), Value::Null);

        let job = make_job_with_modules(modules);
        let toml = job.generate_run_toml().unwrap();

        assert!(
            toml.contains("[\"emscoring.1\"]"),
            "expected [\"emscoring.1\"] but got:\n{toml}"
        );
        assert!(
            !toml.contains("[emscoring.1]"),
            "unquoted [emscoring.1] should not appear"
        );
    }

    #[test]
    fn test_generate_run_toml_plain_modules_unaffected() {
        // Modules without a .digit suffix must keep plain [module] headers
        let mut modules = IndexMap::new();
        modules.insert("topoaa".to_string(), Value::Null);
        modules.insert("caprieval".to_string(), Value::Null);

        let job = make_job_with_modules(modules);
        let toml = job.generate_run_toml().unwrap();

        assert!(
            toml.contains("[topoaa]"),
            "expected [topoaa] but got:\n{toml}"
        );
        assert!(
            toml.contains("[caprieval]"),
            "expected [caprieval] but got:\n{toml}"
        );
    }

    #[test]
    fn test_generate_run_toml_gen_archive_some_true() {
        let mut job = make_job_with_modules(IndexMap::new());
        job.general.gen_archive = Some(true);
        let toml = job.generate_run_toml().unwrap();
        assert!(
            toml.contains("gen_archive = true"),
            "expected gen_archive = true but got:\n{toml}"
        );
    }

    #[test]
    fn test_generate_run_toml_gen_archive_some_false() {
        let mut job = make_job_with_modules(IndexMap::new());
        job.general.gen_archive = Some(false);
        let toml = job.generate_run_toml().unwrap();
        assert!(
            toml.contains("gen_archive = false"),
            "expected gen_archive = false but got:\n{toml}"
        );
    }

    #[test]
    fn test_generate_run_toml_optional_haddock_params_present() {
        let mut job = make_job_with_modules(IndexMap::new());
        job.general.preprocess = Some(true);
        job.general.postprocess = Some(true);
        job.general.gen_archive = Some(true);
        let toml = job.generate_run_toml().unwrap();
        assert!(
            toml.contains("preprocess = true"),
            "expected preprocess = true but got:\n{toml}"
        );
        assert!(
            toml.contains("postprocess = true"),
            "expected postprocess = true but got:\n{toml}"
        );
        assert!(
            toml.contains("gen_archive = true"),
            "expected gen_archive = true but got:\n{toml}"
        );
    }

    #[test]
    fn test_generate_run_toml_optional_haddock_params_absent() {
        let job = make_job_with_modules(IndexMap::new());
        let toml = job.generate_run_toml().unwrap();
        assert!(
            !toml.contains("preprocess"),
            "preprocess should not appear but got:\n{toml}"
        );
        assert!(
            !toml.contains("postprocess"),
            "postprocess should not appear but got:\n{toml}"
        );
        assert!(
            !toml.contains("gen_archive"),
            "gen_archive should not appear but got:\n{toml}"
        );
    }

    #[test]
    fn test_generate_run_toml_preprocess_some_true() {
        let mut job = make_job_with_modules(IndexMap::new());
        job.general.preprocess = Some(true);
        let toml = job.generate_run_toml().unwrap();
        assert!(
            toml.contains("preprocess = true"),
            "expected preprocess = true but got:\n{toml}"
        );
    }

    #[test]
    fn test_generate_run_toml_preprocess_some_false() {
        let mut job = make_job_with_modules(IndexMap::new());
        job.general.preprocess = Some(false);
        let toml = job.generate_run_toml().unwrap();
        assert!(
            toml.contains("preprocess = false"),
            "expected preprocess = false but got:\n{toml}"
        );
    }

    #[test]
    fn test_generate_run_toml_postprocess_some_true() {
        let mut job = make_job_with_modules(IndexMap::new());
        job.general.postprocess = Some(true);
        let toml = job.generate_run_toml().unwrap();
        assert!(
            toml.contains("postprocess = true"),
            "expected postprocess = true but got:\n{toml}"
        );
    }

    #[test]
    fn test_generate_run_toml_postprocess_some_false() {
        let mut job = make_job_with_modules(IndexMap::new());
        job.general.postprocess = Some(false);
        let toml = job.generate_run_toml().unwrap();
        assert!(
            toml.contains("postprocess = false"),
            "expected postprocess = false but got:\n{toml}"
        );
    }
}
