use actix_web::{web, HttpResponse, Responder};
use std::sync::atomic::{AtomicBool, Ordering};

use crate::services::container_registry::client::ContainerRegistryClient;
use crate::services::container_registry::path::InvalidRegistryPathError;
use crate::utils::error::FaradayError;

pub trait ConnectionErrorsHandler {
    fn invalid_path_error(&self) -> &AtomicBool;
    fn set_invalid_path_error(&self, value: bool);
    fn connection_error(&self) -> &AtomicBool;
    fn set_connection_error(&self, value: bool);
    fn render_index(&self) -> HttpResponse;
    fn container_registry_client(&self) -> &ContainerRegistryClient;

    fn invalid_registry_path(&self) -> impl Responder {
        self.set_invalid_path_error(true);
        self.render_index()
    }

    fn handle_connection_error(&self) -> impl Responder {
        self.set_connection_error(true);
        self.render_index()
    }

    fn ping_container_registry(&self) -> Result<(), Box<dyn std::error::Error>> {
        self.container_registry_client().registry_info()
    }
}

pub struct ConnectionErrorsHandlerImpl {
    invalid_path_error: AtomicBool,
    connection_error: AtomicBool,
    container_registry_client: ContainerRegistryClient,
}

impl ConnectionErrorsHandlerImpl {
    pub fn new(container_registry_client: ContainerRegistryClient) -> Self {
        Self {
            invalid_path_error: AtomicBool::new(false),
            connection_error: AtomicBool::new(false),
            container_registry_client,
        }
    }
}

impl ConnectionErrorsHandler for ConnectionErrorsHandlerImpl {
    fn invalid_path_error(&self) -> &AtomicBool {
        &self.invalid_path_error
    }

    fn set_invalid_path_error(&self, value: bool) {
        self.invalid_path_error.store(value, Ordering::SeqCst);
    }

    fn connection_error(&self) -> &AtomicBool {
        &self.connection_error
    }

    fn set_connection_error(&self, value: bool) {
        self.connection_error.store(value, Ordering::SeqCst);
    }

    fn render_index(&self) -> HttpResponse {
        // In a real implementation, this would render the index view
        HttpResponse::Ok().finish()
    }

    fn container_registry_client(&self) -> &ContainerRegistryClient {
        &self.container_registry_client
    }
}
