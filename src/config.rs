use std::{env::VarError, fs::File, io::{Error, Read}};

use serde::Deserialize;
use toml::from_str;

#[derive(Debug, Deserialize, Clone)]
pub struct APIConfig {
    pub database_url: String,
    pub address: String,
    pub port: u16
}

impl APIConfig {

    pub fn from_env() -> Result<APIConfig, VarError> {
        Ok(APIConfig{
            database_url: std::env::var("DATABASE_URL")?,
            address: String::from("0.0.0.0"),
            port: 8000
        })
    }

    #[allow(dead_code)]
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