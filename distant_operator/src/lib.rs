use std::{fs, path::PathBuf, sync::LazyLock};

use logging::Logger;
use reedline::ExternalPrinter;
use serde::{Deserialize, Serialize};

pub mod logging;

// internal
static LOGGER: Logger = Logger;
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
    address: String,
    port: u16,
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
        let str = toml::to_string(&self.inner)?;
        fs::write(path, str)?;
        Ok(())
    }
}

pub fn init() -> Result<(), anyhow::Error> {
    log::set_logger(&LOGGER)?;
    log::set_max_level(log::LevelFilter::Debug);

    Ok(())
}
