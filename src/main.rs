pub mod checksum;
pub mod dataset;
pub mod input;
pub mod job;
pub mod logging;
pub mod queue;
pub mod runner;
pub mod utils;

use anyhow::Result;
use clap::Parser;
use input::Input;
use log::{LevelFilter, info};
use std::path::Path;

use crate::queue::Queue;

/// Print the welcome message with version information
///
/// This function displays the haddock-runner welcome banner including
/// the current version from the cargo package information.
fn print_welcome_message() {
    info!("###########################################");
    info!(" Starting haddock-runner {}", env!("CARGO_PKG_VERSION"));
    info!("###########################################");
}

/// Run HADDOCK on a dataset of complexes
#[derive(Parser, Debug)]
#[command(author, version, about, long_about = None)]
struct Args {
    /// Only perform the setup, do not execute the benchmark
    #[arg(long, short)]
    setup: bool,

    /// Enable debug logging
    #[arg(long, short)]
    debug: bool,

    /// Input file path
    input_file: String,
}

/// Main entry point for the haddock-runner application
///
/// This function parses command line arguments, initializes logging,
/// loads input configuration, validates inputs, creates jobs, and executes them.
///
/// # Returns
///
/// * `Result<()>` - Ok if execution completes successfully, error otherwise
fn main() -> Result<()> {
    // Parse command line arguments using clap
    let args = Args::parse();

    // Initialize logging with appropriate level
    let log_level = if args.debug {
        LevelFilter::Debug
    } else {
        LevelFilter::Info
    };
    logging::init_logging(log_level);

    print_welcome_message();

    let yaml_path = Path::new(&args.input_file);
    info!("Loading input file: {}", yaml_path.display());

    if args.setup {
        info!("`setup` argument passed, the benchmark will not be executed");
    }

    let input = Input::new(yaml_path)?;

    input.validate()?;

    let targets = dataset::load_dataset(&input.general.input_list, &input.general.mol_suffixes)?;

    // Validate checksums for all input files
    let checksum_file = input.general.work_dir.join("checksum.json");
    checksum::validate_checksums(&targets, &checksum_file)?;

    let jobs = job::create_jobs(input.clone(), targets);

    let queue = Queue::new(input.general.max_concurrent, jobs);

    queue.setup()?;

    if args.setup {
        log::info!("Benchmark setup finished successfully, exiting");
        return Ok(());
    }

    queue.start()?;

    Ok(())
}
