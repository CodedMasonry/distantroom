use std::{fs, path::PathBuf};

use clap::{Parser, Subcommand};
use distant_operator::{Profile, PROFILE_DIR};
use nu_ansi_term::Color;
use promkit::preset::listbox::Listbox;
use reedline::{DefaultHinter, DefaultPrompt, Reedline, Signal};

#[derive(Debug, Parser)]
#[command(version, about, long_about=None)]
struct RootArgs {
    #[command(subcommand)]
    command: Option<RootSubCommands>,
}

#[derive(Debug, Subcommand)]
enum RootSubCommands {
    Import {
        #[arg(short, long)]
        path: PathBuf,
    },
}

fn main() -> Result<(), anyhow::Error> {
    // Init
    distant_operator::init()?;
    parse_root()?;

    let _profile = select_profile()?;

    // Init Readline
    let mut line_editor = Reedline::create()
        .with_external_printer(distant_operator::PRINTER.clone())
        .with_hinter(Box::new(
            DefaultHinter::default().with_style(Color::LightGray.italic().dimmed()),
        ));
    let prompt = DefaultPrompt::default();

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

                profile.write_to(save_path)?;
            }
        }
    }
    Ok(())
}

fn select_profile() -> Result<Profile, anyhow::Error> {
    let path = PROFILE_DIR.clone();

    let profiles = fs::read_dir(path)?
        .filter(|v| v.is_ok())
        .map(|v| v.unwrap().file_name().into_string().unwrap());

    let path = Listbox::new(profiles)
        .title("Select A Profile")
        .listbox_lines(5)
        .prompt()?
        .run()?;

    Profile::parse(PROFILE_DIR.join(path))
}
