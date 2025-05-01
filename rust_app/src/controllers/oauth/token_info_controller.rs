// Ported from: orig_app/app/controllers/oauth/token_info_controller.rb
// This controller provides token info for OAuth tokens, similar to the Ruby Doorkeeper::TokenInfoController.

use actix_web::{get, web, HttpRequest, HttpResponse, Responder};
use serde_json::json;

// Dummy traits and types for illustration; replace with actual implementations.
pub trait EnforcesTwoFactorAuthentication {}

#[derive(Debug, Clone)]
pub struct OAuthToken {
    pub scope: Vec<String>,
    pub expires_in: Option<i64>,
    pub accessible: bool,
    // ... other fields ...
}

impl OAuthToken {
    pub fn as_json(&self) -> serde_json::Value {
        json!({
            "scope": self.scope,
            "expires_in": self.expires_in,
            // ... other fields ...
        })
    }
}

fn get_doorkeeper_token(_req: &HttpRequest) -> Option<OAuthToken> {
    // TODO: Implement token extraction and validation
    None
}

#[get("/oauth/token/info")]
pub async fn show(req: HttpRequest) -> impl Responder {
    if let Some(token) = get_doorkeeper_token(&req) {
        if token.accessible {
            let mut token_json = token.as_json();
            // maintain backwards compatibility
            if let Some(scope) = token_json.get("scope").cloned() {
                token_json["scopes"] = scope;
            }
            if let Some(expires_in) = token_json.get("expires_in").cloned() {
                token_json["expires_in_seconds"] = expires_in;
            }
            HttpResponse::Ok().json(token_json)
        } else {
            invalid_token_response()
        }
    } else {
        invalid_token_response()
    }
}

fn invalid_token_response() -> HttpResponse {
    // Simulate Doorkeeper::OAuth::InvalidTokenResponse
    let error_body = json!({
        "error": "invalid_token"
    });
    HttpResponse::Unauthorized().json(error_body)
}
