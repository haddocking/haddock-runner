use anyhow::{Result, bail};

pub fn validate_slurm() -> Result<()> {
    bail!("SLURM not configured")
}
