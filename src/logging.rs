use log::LevelFilter;
use std::io::Write;

pub fn init_logging(log_level: LevelFilter) {
    env_logger::Builder::from_default_env()
        .filter_level(log_level)
        .format(|buf, record| {
            writeln!(
                buf,
                "[{}] {} - {}",
                record.level(),
                record.target(),
                record.args()
            )
        })
        .init();
}

