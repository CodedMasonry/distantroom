use std::{fs::{self, DirEntry}, path::PathBuf, sync::LazyLock};

use anyhow::bail;
use logging::Logger;
use nu_ansi_term::Style;
use promkit::preset::{confirm::Confirm, listbox::Listbox};
use reedline::ExternalPrinter;
use serde::{Deserialize, Serialize};

pub mod logging;

// internal
pub static LOGGER: Logger = Logger;
pub static PRINTER: LazyLock<ExternalPrinter<String>> = LazyLock::new(|| ExternalPrinter::new(8));

// directories
pub static CONFIG_DIR: LazyLock<PathBuf> = LazyLock::new(|| {
    let mut dir = dirs::config_dir().unwrap();
    dir.push("distant_operator");
    fs::create_dir_all(dir.clone()).unwrap();

    dir
});
pub static PROFILE_DIR: LazyLock<PathBuf> = LazyLock::new(|| {
    let mut dir = dirs::config_dir().unwrap();
    dir.push("distant_operator");
    dir.push("profiles");
    fs::create_dir_all(dir.clone()).unwrap();

    dir
});

#[derive(Debug)]
pub struct Profile {
    path: PathBuf,
    inner: ProfileInner,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ProfileInner {
    pub address: String,
    pub port: u16,
    public_key: String,
    private_key: String,
}

impl Profile {
    pub fn parse(path: PathBuf) -> Result<Self, anyhow::Error> {
        let contents = fs::read_to_string(path.clone())?;
        let inner: ProfileInner = toml::from_str(&contents)?;
        Ok(Profile { path, inner })
    }

    pub fn save(&self) -> Result<(), anyhow::Error> {
        let str = toml::to_string(&self.inner)?;
        fs::write(self.path.clone(), str)?;
        Ok(())
    }

    pub fn write_to(&self, path: PathBuf) -> Result<(), anyhow::Error> {
        if fs::exists(path.clone()).unwrap_or_default() {
            if !Confirm::new(format!(
                "The file {:?} already exists; Do you want to replace it?",
                path.file_name().unwrap()
            ))
            .prompt()?
            .run()?
            .to_lowercase()
            .contains('y')
            {
                return Ok(());
            };
        }

        let str = toml::to_string(&self.inner)?;
        fs::write(path, str)?;
        Ok(())
    }
}

pub fn select_profile() -> Result<Profile, anyhow::Error> {
    let path = PROFILE_DIR.clone();

    let profiles: Vec<String> = fs::read_dir(path)?
        .filter(|v| v.is_ok())
        .map(|v| format_profile_name(v.unwrap()))
        .collect();

    if profiles.len() == 0 {
        bail!(
            "No profiles found in {:?}\n {}: run the `import` command to add a profile",
            PROFILE_DIR.clone(),
            Style::new().bold().paint("Note")
        )
    }

    let selected = Listbox::new(profiles)
        .title("Select A Profile")
        .listbox_lines(5)
        .prompt()?
        .run()?;
    let file = selected.split(" | ").next().unwrap();

    Profile::parse(PROFILE_DIR.join(file))
}

fn format_profile_name(file: DirEntry) -> String {
    let file_name = file.file_name().into_string().unwrap();
    let str = fs::read_to_string(file.path()).unwrap();
    let parsed: ProfileInner = toml::from_str(&str).unwrap();

    format!("{} | {}", file_name, parsed.address)
}