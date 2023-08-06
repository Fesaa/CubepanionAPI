use std::{fs::File, io::{Error, Read}};

use serde::Deserialize;
use toml::from_str;

#[derive(Debug, Deserialize, Clone)]
pub struct APIConfig {
    pub database_url: String,
    pub address: String,
    pub port: u16
}

impl APIConfig {

    pub fn from_file(path: String) -> Result<APIConfig, Error> {
        let mut file = File::open(path)?;

        let mut content = String::new();
        file.read_to_string(&mut content)?;

        let config: APIConfig = match from_str(&content) {
            Ok(c) => c,
            Err(e) => panic!("Could not parse config: {}", e),
        };
        Ok(config)
    }

}