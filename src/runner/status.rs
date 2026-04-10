#[derive(Debug, Clone, PartialEq)]
pub enum JobStatus {
    Done,
    Failed,
    Queued,
    Incomplete,
    Unknown,
    Submitted,
}
