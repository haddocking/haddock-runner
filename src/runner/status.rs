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

#[derive(Debug, PartialEq)]
pub enum SlurmJobState {
    Pending,     // PD
    Running,     // R
    Suspended,   // S
    Completing,  // CG
    Completed,   // COMPLETED
    Configuring, // CF
    Failed,      // FAILED
    Timeout,     // TIMEOUT
    Cancelled,   // CANCELLED
    NodeFail,    // NODE_FAIL
    Preempted,   // PREEMPTED
    BootFail,    // BOOT_FAIL
    StageOut,    // STAGE_OUT
    Stopped,     // STOPPED
    SpecialExit, // SPECIAL_EXIT
    Unknown,
}

impl SlurmJobState {
    pub fn from_status_code(code: &str) -> Self {
        match code {
            "PD" | "PENDING" => SlurmJobState::Pending,
            "R" | "RUNNING" => SlurmJobState::Running,
            "S" | "SUSPENDED" => SlurmJobState::Suspended,
            "CG" | "COMPLETING" => SlurmJobState::Completing,
            "COMPLETED" => SlurmJobState::Completed,
            "CF" | "CONFIGURING" => SlurmJobState::Configuring,
            "F" | "FAILED" => SlurmJobState::Failed,
            "TO" | "TIMEOUT" => SlurmJobState::Timeout,
            "CA" | "CANCELLED" => SlurmJobState::Cancelled,
            "NF" | "NODE_FAIL" => SlurmJobState::NodeFail,
            "PR" | "PREEMPTED" => SlurmJobState::Preempted,
            "BF" | "BOOT_FAIL" => SlurmJobState::BootFail,
            "SO" | "STAGE_OUT" => SlurmJobState::StageOut,
            "ST" | "STOPPED" => SlurmJobState::Stopped,
            "SE" | "SPECIAL_EXIT" => SlurmJobState::SpecialExit,
            _ => SlurmJobState::Unknown,
        }
    }

    /// Status that indicate job is no longer being executed
    pub fn is_terminal(&self) -> bool {
        matches!(
            self,
            SlurmJobState::Completed
                | SlurmJobState::Failed
                | SlurmJobState::Timeout
                | SlurmJobState::Cancelled
                | SlurmJobState::NodeFail
                | SlurmJobState::Preempted
                | SlurmJobState::BootFail
                | SlurmJobState::SpecialExit
        )
    }

    pub fn is_success(&self) -> bool {
        matches!(self, SlurmJobState::Completed)
    }
}
