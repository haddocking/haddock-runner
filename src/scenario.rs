use crate::{dataset::Target, input::Workflow, runner::status::JobStatus};
use anyhow::Result;
use std::path::PathBuf;

pub struct Scenario {
    id: String,
    path: PathBuf,
    status: JobStatus,
    target: Target,
    workflow: Workflow,
}

impl Scenario {
    fn setup() {
        todo!()
    }

    fn execute() {
        todo!()
    }
}
