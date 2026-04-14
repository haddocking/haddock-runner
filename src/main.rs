pub mod checksum;
pub mod dataset;
pub mod input;
pub mod job;
pub mod logging;
pub mod queue;
pub mod runner;
pub mod slurm;
pub mod utils;

use anyhow::Result;
use input::Input;
use log::LevelFilter;
use std::path::Path;

use crate::queue::Queue;

fn main() -> Result<()> {
    // Initialize logging
    logging::init_logging(LevelFilter::Info);

    let yaml_path = Path::new("example/bm.yml");

    let input = Input::new(yaml_path)?;

    input.validate()?;

    let targets = dataset::load_dataset(
        &input.general.input_list,
        &input.general.mol_suffixes,
        input.general.shape_suffix.as_deref(),
    );

    // Validate checksums for all input files
    let checksum_file = input.general.work_dir.join("checksum.json");
    checksum::validate_checksums(&targets, &checksum_file)?;

    // println!("{:?}", dataset);

    // let targets = dataset::organize_dataset(raw_targets, &input.general.work_dir)?;

    // println!("Organized {} targets successfully", targets.len());

    // scenario + target = job
    input
        .scenarios
        .iter()
        .for_each(|s| println!("{:?}", s.name));

    let jobs = job::create_jobs(input.clone(), targets);

    let queue = Queue::new(input.general.max_concurrent, jobs);

    queue.setup()?;

    queue.start()?;

    Ok(())
}
