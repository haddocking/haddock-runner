pub mod dataset;
pub mod input;
pub mod runner;
pub mod scenario;
pub mod utils;

use input::Input;

use crate::dataset::load_dataset;

fn main() {
    let yaml_content = std::fs::read_to_string("example/bm.yml").unwrap();
    let input: Input = serde_yaml::from_str(&yaml_content).unwrap();

    // println!("{:?}", input);

    for scenario in input.scenarios {
        println!("{:?}", scenario.workflow)
    }

    let dataset = load_dataset(
        &input.general.input_list,
        &input.general.mol_suffixes,
        input.general.shape_suffix.as_deref(),
    );

    println!("{:?}", dataset)
}
