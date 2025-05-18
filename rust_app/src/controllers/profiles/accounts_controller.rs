// Ported from: orig_app/app/controllers/profiles/accounts_controller.rb
// Ported on: 2025-05-05
// This file implements the Profiles::AccountsController from the Ruby codebase.

use actix_web::{get, post, web, HttpResponse, Responder};
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct UnlinkParams {
    pub provider: String,
}

#[get("/profile/account")]
pub async fn show() -> impl Responder {
    // TODO: Implement push_frontend_feature_flag(:delay_delete_own_user)
    // TODO: Implement show_view_variables
    HttpResponse::Ok().body("Account page (stub)")
}

#[post("/profile/account/unlink")]
pub async fn unlink(params: web::Json<UnlinkParams>) -> impl Responder {
    let provider = &params.provider;
    // TODO: Implement find_identity and unlink_provider_allowed logic
    let identity_found = true; // Placeholder
    let unlink_allowed = true; // Placeholder
    if !identity_found {
        return HttpResponse::NotFound().body("Identity not found");
    }
    if unlink_allowed {
        // TODO: Destroy identity
    } else {
        // TODO: Set flash alert: "You are not allowed to unlink your primary login account"
    }
    HttpResponse::Found()
        .header("Location", "/profile/account")
        .finish()
}

#[post("/profile/account/generate_support_pin")]
pub async fn generate_support_pin() -> impl Responder {
    // TODO: Implement Users::SupportPin::UpdateService equivalent
    let success = true; // Placeholder
    if success {
        // TODO: Set flash notice: "New Support PIN generated successfully."
    } else {
        // TODO: Set flash alert: "Failed to generate new Support PIN."
    }
    HttpResponse::Found()
        .header("Location", "/profile/account")
        .finish()
}

// TODO: Add integration with app state, user session, flash messages, and feature flags as needed.
