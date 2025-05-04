// Ported from: orig_app/app/controllers/oauth/tokens_controller.rb
// Ported on: 2025-05-04
//
// Handles OAuth token endpoint logic, including client validation and error responses.
// Extends Doorkeeper's TokensController functionality.

use crate::auth::doorkeeper;
use crate::auth::two_factor::EnforcesTwoFactorAuthentication;
use crate::controllers::concerns::request_payload_logger::RequestPayloadLogger;
use crate::models::user::User;
use actix_web::{post, web, HttpRequest, HttpResponse, Responder};
use serde_json::json;
use std::sync::Arc;

pub struct TokensController {
    current_user: Option<Arc<User>>,
}

impl TokensController {
    pub fn new(current_user: Option<Arc<User>>) -> Self {
        Self { current_user }
    }

    /// POST /oauth/token
    #[post("/oauth/token")]
    pub async fn create(req: HttpRequest) -> impl Responder {
        if let Err(resp) = Self::validate_presence_of_client(&req).await {
            return resp;
        }

        // TODO: Implement full token creation logic
        HttpResponse::Ok().json(json!({"access_token": "dummy_token"}))
    }

    async fn validate_presence_of_client(req: &HttpRequest) -> Result<(), HttpResponse> {
        // Check Doorkeeper config for password grant client authentication skip
        if doorkeeper::config::skip_client_authentication_for_password_grant() {
            return Ok(());
        }

        // Check if client credentials are present
        // See RFC 6749 Section 2.1 Client Authentication
        if doorkeeper::server::has_valid_client(req) {
            return Ok(());
        }

        // If validation fails, return error response conforming to RFC 6749 Section 5.2
        Err(Self::revocation_error_response())
    }

    fn revocation_error_response() -> HttpResponse {
        HttpResponse::Forbidden()
            .append_header(("Cache-Control", "no-store"))
            .append_header(("Pragma", "no-cache"))
            .json(json!({
                "error": "invalid_client",
                "error_description": "Client authentication failed"
            }))
    }

    // Alias for current_user, maintaining API compatibility
    pub fn auth_user(&self) -> Option<Arc<User>> {
        self.current_user.clone()
    }
}

// Register routes
pub fn configure(cfg: &mut web::ServiceConfig) {
    cfg.service(TokensController::create);
}
