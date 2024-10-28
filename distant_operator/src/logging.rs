use log::Level;
use nu_ansi_term::{AnsiGenericString, Color};

pub struct Logger;

impl log::Log for Logger {
    fn enabled(&self, metadata: &log::Metadata) -> bool {
        metadata.level() >= Level::Debug
    }

    fn log(&self, record: &log::Record) {
        let prefix = log_prefix(record.level());
        let time = chrono::Local::now().format("%H:%M:%S").to_string();
        let str = format!(
            "{} {} {}",
            prefix,
            Color::LightGray.dimmed().paint(time),
            record.args()
        );

        print_terminal(str).unwrap_or_default()
    }

    fn flush(&self) {}
}

pub fn log_prefix<'a>(level: Level) -> AnsiGenericString<'a, str> {
    match level {
        Level::Error => Color::Red.bold().paint("[-]"),
        Level::Warn => Color::Yellow.bold().paint("[!]"),
        Level::Info => Color::Green.bold().paint("[+]"),
        Level::Debug => Color::Blue.bold().paint("[*]"),
        Level::Trace => Color::Purple.bold().paint("[$]"),
    }
}

fn print_terminal(str: String) -> Result<(), crossbeam_channel::SendError<String>> {
    crate::PRINTER.sender().send(str)
}
