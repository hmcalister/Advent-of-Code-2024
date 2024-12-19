use std::fs::File;
use tracing_subscriber::{fmt, prelude::*, EnvFilter};

const LOG_FILEPATH: &str = "log";

fn get_log_env_filter() -> EnvFilter{
    EnvFilter::try_from_default_env().unwrap_or_else(|_| EnvFilter::new("info"))
}

pub fn set_logging() {
    let log_file_handle = File::create(LOG_FILEPATH).expect("failed to create log file");
    let log_env_filter = get_log_env_filter();
    let to_stdout = log_env_filter.to_string().contains("debug") || log_env_filter.to_string().contains("trace");

    let file_layer = fmt::layer()
        .with_writer(log_file_handle)
        .with_file(true)
        .with_line_number(true)
        .with_ansi(false)
        .with_filter(get_log_env_filter());

    let stdout_layer = if to_stdout {
        Some(
            fmt::layer()
                .with_writer(std::io::stdout)
                .with_filter(get_log_env_filter()),
        )
    } else {
        None
    };

    let subscriber = tracing_subscriber::registry()
        .with(file_layer)
        .with(stdout_layer);

    subscriber.init();
}
