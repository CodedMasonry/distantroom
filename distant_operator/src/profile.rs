use std::{fs, path::PathBuf};

use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct Profile {
    server: String,
    port: Option<u16>,
}

impl Profile {
    pub fn try_parse(str: &str) -> Result<Self, toml::de::Error> {
        toml::from_str::<Profile>(str)
    }

    pub fn read_from(path: PathBuf) -> Result<Self, anyhow::Error> {
        let file = fs::read_to_string(path)?;
        Ok(toml::from_str::<Profile>(&file)?)
    }

    pub fn write_to(&self, path: PathBuf) -> Result<(), anyhow::Error> {
        let str = toml::to_string_pretty(self)?;
        fs::write(path, str)?;
        Ok(())
    }
}
