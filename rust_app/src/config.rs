use config::{Config, ConfigError, Environment, File, FileFormat};
use serde::Deserialize;

#[derive(Debug, Deserialize)]
pub struct Settings {
    pub database_url: String,
    pub server_port: u16,
    pub redis_url: Option<String>,
    pub secret_key: String,
}

impl Settings {
    pub fn new() -> Result<Self, ConfigError> {
        let config = Config::builder()
            .add_source(File::new("config/default.yml", FileFormat::Yaml))
            .add_source(Environment::with_prefix("APP"))
            .build()?;

        config.try_deserialize()
    }
}