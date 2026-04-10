use chrono::Local;

pub fn generate_timestamp() -> String {
    Local::now().format("%Y-%m-%d_%H-%M-%S").to_string()
}
