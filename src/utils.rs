use chrono::Local;
use serde_yaml::Value;
use std::process::Command;

pub fn generate_timestamp() -> String {
    Local::now().format("%Y-%m-%d_%H-%M-%S").to_string()
}

pub fn extract_root(filename: &str) -> Option<String> {
    // Pattern: ROOT_*.xxx
    let parts: Vec<&str> = filename.split('_').collect();
    if parts.len() >= 2 {
        Some(parts[0].to_string())
    } else {
        None
    }
}

pub fn extract_root_from_filename(filename: &str) -> Option<String> {
    // Try to extract root by removing common suffixes
    let filename = filename.replace(".pdb", "");
    let filename = filename.replace(".tbl", "");
    let filename = filename.replace(".top", "");
    let filename = filename.replace(".param", "");

    // Remove any remaining suffixes after underscores
    if let Some(pos) = filename.find('_') {
        Some(filename[..pos].to_string())
    } else {
        Some(filename)
    }
}

pub fn format_toml_value(value: &Value) -> String {
    match value {
        Value::Bool(b) => b.to_string(),
        Value::Number(n) => n.to_string(),
        Value::String(s) => format!("\"{}\"", s),
        Value::Sequence(seq) => {
            let items: Vec<String> = seq.iter().map(format_toml_value).collect();
            format!("[{}]", items.join(", "))
        }
        _ => "null".to_string(), // Fallback for other types
    }
}

/// Check if a command exists in the system PATH
pub fn command_exists(command: &str) -> bool {
    if let Ok(output) = Command::new("which").arg(command).output() {
        output.status.success()
    } else {
        false
    }
}
