pub mod dataset;
pub mod input;
pub mod runner;
pub mod scenario;
pub mod utils;

use anyhow::Result;
use input::Input;
use std::path::Path;

fn main() -> Result<()> {
    let yaml_path = Path::new("example/bm.yml");

    let input = Input::new(yaml_path)?;

    let dataset = dataset::load_dataset(
        &input.general.input_list,
        &input.general.mol_suffixes,
        input.general.shape_suffix.as_deref(),
    );

    // println!("{:?}", dataset);

    let targets = dataset.organize(&input.general.work_dir)?;

    println!("Organized {} targets successfully", targets.0.len());

    Ok(())
}
