use chrono::Local;

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
