use std::path::PathBuf;

use clap::Parser;

use distant_operator::{self, logging::Logger, profile::Profile, State};
use log::LevelFilter;

static LOGGER: Logger = Logger;
#[derive(Parser, Debug)]
#[command(version, about, long_about = None)]
struct RootArgs {
    #[arg(short, long)]
    profile: Option<PathBuf>,
}

fn main() -> Result<(), anyhow::Error> {
    let _state = init()?;
    Ok(())
}

fn init() -> Result<State, anyhow::Error> {
    // init logging
    log::set_logger(&LOGGER).map(|()| log::set_max_level(LevelFilter::Debug))?;

    // Parse Args
    let root_args = RootArgs::parse();
    let mut config_dir = dirs::config_dir().expect("Failed to get Config Directory");
    config_dir.push("distantroom");
    config_dir.push("profiles");

    // Fetch Profile
    let profile = if root_args.profile.is_some() {
        Profile::read_from(root_args.profile.unwrap())?
    } else {
        anyhow::bail!(
            "No Profile Specified\n       Add `--profile /path/to/profile.json` to specify profile"
        )
    };

    return Ok(State { profile });
}
