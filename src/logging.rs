use log::LevelFilter;
use std::io::Write;

pub fn init_logging(log_level: LevelFilter) {
    env_logger::Builder::from_default_env()
        .filter_level(log_level)
        .format(|buf, record| {
            let level_style = match record.level() {
                log::Level::Error => "\x1b[31m", // Red
                log::Level::Warn => "\x1b[33m",  // Yellow
                log::Level::Info => "\x1b[32m",  // Green
                log::Level::Debug => "\x1b[36m", // Cyan
                log::Level::Trace => "\x1b[35m", // Magenta
            };
            let reset = "\x1b[0m";
            
            writeln!(
                buf,
                "{}[{}]{} {} - {}",
                level_style,
                record.level(),
                reset,
                record.target(),
                record.args()
            )
        })
        .init();
}
