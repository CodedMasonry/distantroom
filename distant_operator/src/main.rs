#![deny(clippy::all)]
use std::{path::PathBuf, process, thread};

use clap::{Parser, Subcommand};
use distant_operator::{server::{self, Server}, Profile, LOGGER, PROFILE_DIR};
use log::{error, info};
use nu_ansi_term::Color;
use reedline::{DefaultHinter, DefaultPrompt, Reedline, Signal};

#[derive(Debug, Parser)]
#[command(version, about, long_about=None)]
struct RootArgs {
    #[command(subcommand)]
    command: Option<RootSubCommands>,
}

#[derive(Debug, Subcommand)]
enum RootSubCommands {
    /// Import a profile
    Import {
        /// The path of a `profile.toml`
        path: PathBuf,
    },
}

fn main() -> Result<(), anyhow::Error> {
    // Init
    log::set_logger(&LOGGER)?;
    log::set_max_level(log::LevelFilter::Debug);
    parse_root()?;

    // Init Readline
    let mut line_editor = Reedline::create()
        .with_external_printer(distant_operator::PRINTER.clone())
        .with_hinter(Box::new(
            DefaultHinter::default().with_style(Color::LightGray.italic().dimmed()),
        ));
    let prompt = DefaultPrompt::default();

    let profile = distant_operator::select_profile()?;
    let server = Server::connect(&profile)?;

    // Handle Readline
    loop {
        if let Ok(sig) = line_editor.read_line(&prompt) {
            match sig {
                Signal::Success(buffer) => {
                    println!("We processed: {buffer}");
                }
                Signal::CtrlD | Signal::CtrlC => {
                    println!("\nExitting...");
                    break;
                }
            }
            continue;
        }
        break;
    }

    Ok(())
}

fn parse_root() -> Result<(), anyhow::Error> {
    let args = RootArgs::parse();

    if let Some(cmd) = args.command {
        match cmd {
            RootSubCommands::Import { path } => {
                let profile: Profile = Profile::parse(path.clone())?;
                let mut save_path = PROFILE_DIR.clone();
                save_path.push(path.file_name().unwrap());

                profile.write_to(save_path.clone())?;
                info!("Profile saved to {:?}", save_path);
                // Command Finished
                process::exit(0)
            }
        }
    }
    Ok(())
}
