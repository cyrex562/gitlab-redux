// Ported from orig_app/app/controllers/google_api/authorizations_controller.rb
// Date ported: 2025-04-30

use actix_web::{get, post, web, HttpRequest, HttpResponse, Responder};
use serde::Deserialize;

// Placeholder for session and flash handling
type Session = std::collections::HashMap<String, String>;

// Placeholder for GoogleApi::CloudPlatform::Client logic
struct GoogleApiCloudPlatformClient;

impl GoogleApiCloudPlatformClient {
    fn new(_user: Option<&str>, _callback_url: &str) -> Self {
        Self {}
    }
    fn get_token(&self, _code: &str) -> (String, String) {
        // Simulate token and expiry
        ("token_value".to_string(), "expires_at_value".to_string())
    }
    fn session_key_for_token() -> &'static str {
        "google_api_token"
    }
    fn session_key_for_expires_at() -> &'static str {
        "google_api_expires_at"
    }
    fn session_key_for_redirect_uri(_state: &str) -> String {
        format!("google_api_redirect_uri_{}", _state)
    }
}

#[derive(Deserialize)]
pub struct CallbackParams {
    pub error: Option<String>,
    pub code: Option<String>,
    pub state: Option<String>,
}

#[post("/google_api/authorizations/callback")]
pub async fn callback(
    req: HttpRequest,
    params: web::Form<CallbackParams>,
    session: web::Data<Session>,
) -> impl Responder {
    // Simulate flash messages
    let mut flash_alert = None;
    let mut redirect_uri = redirect_uri_from_session(&params, &session);

    if params.error.is_some() {
        flash_alert = Some("Google Cloud authorizations required".to_string());
        redirect_uri = session.get("error_uri").cloned();
    } else if let Some(code) = &params.code {
        let client = GoogleApiCloudPlatformClient::new(None, "callback_google_api_auth_url");
        let (token, expires_at) = client.get_token(code);
        // Save to session (simulated)
        // session.insert(GoogleApiCloudPlatformClient::session_key_for_token().to_string(), token);
        // session.insert(GoogleApiCloudPlatformClient::session_key_for_expires_at().to_string(), expires_at);
        redirect_uri = redirect_uri_from_session(&params, &session);
    }
    // Simulate Faraday errors as a generic error branch
    // In real code, handle HTTP client errors here
    // Always redirect
    let location = redirect_uri.unwrap_or_else(|| "/".to_string());
    let mut response = HttpResponse::Found().append_header(("Location", location));
    if let Some(alert) = flash_alert {
        response.append_header(("X-Flash-Alert", alert));
    }
    response.finish()
}

fn redirect_uri_from_session(params: &CallbackParams, session: &Session) -> Option<String> {
    if let Some(state) = &params.state {
        session
            .get(&GoogleApiCloudPlatformClient::session_key_for_redirect_uri(
                state,
            ))
            .cloned()
    } else {
        None
    }
}
