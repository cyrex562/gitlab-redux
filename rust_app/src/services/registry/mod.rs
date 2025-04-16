pub mod connection_errors_handler;

use actix_web::{web, HttpResponse};
use serde::{Deserialize, Serialize};

// TODO: Implement ContainerRegistry::Client
pub struct ContainerRegistryClient;

impl ContainerRegistryClient {
    pub fn registry_info() -> Result<(), Box<dyn std::error::Error>> {
        // TODO: Implement actual registry info retrieval
        Ok(())
    }
}

// TODO: Implement ContainerRegistry::Path::InvalidRegistryPathError
#[derive(Debug)]
pub struct InvalidRegistryPathError;

impl std::fmt::Display for InvalidRegistryPathError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "Invalid registry path")
    }
}

impl std::error::Error for InvalidRegistryPathError {}
