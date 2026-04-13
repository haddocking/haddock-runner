#[derive(Debug, Clone, PartialEq)]
pub enum Status {
    Done,
    Failed,
    Queued,
    Incomplete,
    Unknown,
    Submitted,
    Prepared,
}
