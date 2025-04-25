// Ported from: orig_app/app/controllers/concerns/registry/connection_errors_handler.rb
// Handles registry connection errors and sets error flags for the view

use actix_web::{HttpResponse, Responder};

pub struct RegistryErrorState {
    pub invalid_path_error: bool,
    pub connection_error: bool,
}

impl Default for RegistryErrorState {
    fn default() -> Self {
        Self {
            invalid_path_error: false,
            connection_error: false,
        }
    }
}

pub async fn invalid_registry_path(state: &mut RegistryErrorState) -> impl Responder {
    state.invalid_path_error = true;
    render_index(state)
}

pub async fn connection_error(state: &mut RegistryErrorState) -> impl Responder {
    state.connection_error = true;
    render_index(state)
}

fn render_index(_state: &RegistryErrorState) -> HttpResponse {
    // In a real app, render the index template and pass error state to the view
    HttpResponse::Ok().body("Registry index page (error state set)")
}

pub fn ping_container_registry() {
    // Call out to the registry client to check registry info
    // ContainerRegistry::Client.registry_info equivalent
    // ...existing code...
}
