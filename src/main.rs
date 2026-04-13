pub mod dataset;
pub mod input;
pub mod job;
pub mod runner;
pub mod utils;

use anyhow::Result;
use input::Input;
use std::path::Path;

fn main() -> Result<()> {
    let yaml_path = Path::new("example/bm.yml");

    let input = Input::new(yaml_path)?;

    let targets = dataset::load_dataset(
        &input.general.input_list,
        &input.general.mol_suffixes,
        input.general.shape_suffix.as_deref(),
    );

    // println!("{:?}", dataset);

    // let targets = dataset::organize_dataset(raw_targets, &input.general.work_dir)?;

    // println!("Organized {} targets successfully", targets.len());

    // scenario + target = job
    input
        .scenarios
        .iter()
        .for_each(|s| println!("{:?}", s.name));

    let mut jobs = job::create_jobs(input, targets);
    jobs.iter_mut().for_each(|j| {
        println!("{:?}", j.wd);
        j.setup().unwrap();
        j.run().unwrap();
    });
    // println!("{:?}", input.scenarios);

    Ok(())
}
