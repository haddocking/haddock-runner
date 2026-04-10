pub mod input;
pub mod runner;
pub mod utils;

use input::Input;

fn main() {
    let yaml_content = std::fs::read_to_string("example/bm.yml").unwrap();
    let input: Input = serde_yaml::from_str(&yaml_content).unwrap();

    // println!("{:?}", input);

    for scenario in input.scenarios {
        println!("{:?}", scenario.workflow)
    }
}
