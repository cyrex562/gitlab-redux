// Ported from: /home/azrael/Projects/gitlab-redux/orig_app/app/controllers/oauth/token_info_controller.rb
// This controller provides token info for OAuth tokens, extending Doorkeeper's TokenInfoController.

use crate::auth::doorkeeper::AccessToken;
use crate::auth::EnforcesTwoFactorAuthentication;
use actix_web::{get, HttpRequest, HttpResponse, Responder};
use serde_json::json;

#[derive(Debug, Clone)]
pub struct TokenInfoController;

impl TokenInfoController {
    fn get_doorkeeper_token(req: &HttpRequest) -> Option<AccessToken> {
        // TODO: Extract token from Authorization header and validate it
        match req.headers().get("Authorization") {
            Some(auth_header) => {
                if let Ok(auth_str) = auth_header.to_str() {
                    if auth_str.starts_with("Bearer ") {
                        let token = auth_str[7..].to_string();
                        AccessToken::find_by_token(&token)
                    } else {
                        None
                    }
                } else {
                    None
                }
            }
            None => None,
        }
    }

    fn invalid_token_response() -> HttpResponse {
        let error = json!({
            "error": "invalid_token",
            "error_description": "The access token is invalid"
        });
        HttpResponse::Unauthorized()
            .append_header(("Cache-Control", "no-store"))
            .append_header(("Pragma", "no-cache"))
            .json(error)
    }
}

/// GET /oauth/token/info
/// Returns information about the current token
#[get("/oauth/token/info")]
pub async fn show(req: HttpRequest) -> impl Responder {
    if let Some(token) = TokenInfoController::get_doorkeeper_token(&req) {
        if token.accessible() {
            let mut token_json = token.as_json();

            // Maintain backwards compatibility
            token_json["scopes"] = token_json["scope"].clone();
            token_json["expires_in_seconds"] = token_json["expires_in"].clone();

            HttpResponse::Ok()
                .append_header(("Cache-Control", "no-store"))
                .append_header(("Pragma", "no-cache"))
                .json(token_json)
        } else {
            TokenInfoController::invalid_token_response()
        }
    } else {
        TokenInfoController::invalid_token_response()
    }
}
