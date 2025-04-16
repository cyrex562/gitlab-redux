use crate::registry::{ContainerRegistryClient, InvalidRegistryPathError};
use actix_web::{error::ResponseError, web, HttpResponse};

pub trait ConnectionErrorsHandler {
    // These fields are used to pass error state to the frontend
    fn set_invalid_path_error(&mut self, value: bool);
    fn set_connection_error(&mut self, value: bool);

    fn invalid_registry_path(&mut self) -> HttpResponse {
        self.set_invalid_path_error(true);
        // TODO: Implement proper index template rendering
        HttpResponse::Ok().json(serde_json::json!({
            "invalid_path_error": true
        }))
    }

    fn connection_error(&mut self) -> HttpResponse {
        self.set_connection_error(true);
        // TODO: Implement proper index template rendering
        HttpResponse::Ok().json(serde_json::json!({
            "connection_error": true
        }))
    }

    fn ping_container_registry(&self) -> Result<(), Box<dyn std::error::Error>> {
        ContainerRegistryClient::registry_info()
    }
}

// Implement error handling for actix-web
impl ResponseError for InvalidRegistryPathError {
    fn error_response(&self) -> HttpResponse {
        HttpResponse::BadRequest().json(serde_json::json!({
            "error": "Invalid registry path"
        }))
    }
}

// TODO: Implement Faraday::Error equivalent
#[derive(Debug)]
pub struct ConnectionError;

impl std::fmt::Display for ConnectionError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "Connection error")
    }
}

impl std::error::Error for ConnectionError {}

impl ResponseError for ConnectionError {
    fn error_response(&self) -> HttpResponse {
        HttpResponse::InternalServerError().json(serde_json::json!({
            "error": "Connection error"
        }))
    }
}
