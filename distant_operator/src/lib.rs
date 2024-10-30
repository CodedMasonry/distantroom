use std::{
    fs::{self, DirEntry},
    path::PathBuf,
    process,
    sync::LazyLock,
};

use anyhow::bail;
use logging::{log_prefix, Logger};
use nu_ansi_term::Style;
use promkit::preset::{confirm::Confirm, listbox::Listbox};
use reedline::ExternalPrinter;
use serde::{Deserialize, Serialize};

pub mod logging;
pub mod server;

// internal
pub static LOGGER: Logger = Logger;
pub static PRINTER: LazyLock<ExternalPrinter<String>> = LazyLock::new(|| ExternalPrinter::new(8));

// directories
pub static CONFIG_DIR: LazyLock<PathBuf> = LazyLock::new(|| {
    let mut dir = dirs::config_dir().unwrap();
    dir.push("distantroom");
    fs::create_dir_all(dir.clone()).unwrap();

    dir
});
pub static PROFILE_DIR: LazyLock<PathBuf> = LazyLock::new(|| {
    let mut dir = dirs::config_dir().unwrap();
    dir.push("distantroom");
    dir.push("profiles");
    fs::create_dir_all(dir.clone()).unwrap();

    dir
});

#[derive(Debug)]
pub struct Profile {
    path: PathBuf,
    pub inner: ProfileInner,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ProfileInner {
    pub host: String,
    pub port: u16,
    certificate: String,
    private_key: String,
    server_certificate: String,
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

    pub fn get_client_cert(&self) -> String {
        self.inner.certificate.clone()
    }

    pub fn get_client_key(&self) -> String {
        self.inner.private_key.clone()
    }

    pub fn get_root_cert(&self) -> String {
        self.inner.server_certificate.clone()
    }

    pub fn write_to(&self, path: PathBuf) -> Result<(), anyhow::Error> {
        if fs::exists(path.clone()).unwrap_or_default()
            && !Confirm::new(format!(
                "The file {:?} already exists; Do you want to replace it?",
                path.file_name().unwrap()
            ))
            .prompt()?
            .run()?
            .to_lowercase()
            .contains('y')
        {
            return Ok(());
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

    if profiles.is_empty() {
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
    let parsed: ProfileInner = match toml::from_str(&str) {
        Ok(v) => v,
        Err(err) => {
            println!("{} Failed to parse {file_name}\nError: {}\n\nRead the `example.toml` in the repo to see the current format.", log_prefix(log::Level::Error), err.message());
            process::exit(1)
        }
    };

    format!("{} | {}", file_name, parsed.host)
}
