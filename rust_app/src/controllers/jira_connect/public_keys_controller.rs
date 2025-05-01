// Ported from: orig_app/app/controllers/jira_connect/public_keys_controller.rb
// Ported on: 2025-05-01
// This controller exposes public keys for Jira Connect integration.

use actix_web::{get, web, HttpResponse, Responder};

#[get("/jira_connect/public_keys/{id}")]
pub async fn show(path: web::Path<String>) -> impl Responder {
    // Stub: check if public key storage is enabled
    if !public_key_storage_enabled() {
        return HttpResponse::NotFound().finish();
    }

    // Stub: fetch the public key by id
    let id = path.into_inner();
    match find_public_key(&id) {
        Some(key) => HttpResponse::Ok().content_type("text/plain").body(key),
        None => HttpResponse::NotFound().finish(),
    }
}

// Stub for checking if public key storage is enabled
fn public_key_storage_enabled() -> bool {
    // TODO: Integrate with actual settings
    true
}

// Stub for finding a public key by id
fn find_public_key(id: &str) -> Option<String> {
    // TODO: Integrate with actual model/database
    if id == "example" {
        Some("EXAMPLE_PUBLIC_KEY_DATA".to_string())
    } else {
        None
    }
}
