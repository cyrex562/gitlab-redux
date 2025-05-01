// Ported from: orig_app/app/controllers/jira_connect/oauth_application_ids_controller.rb

use actix_web::{get, HttpResponse, Responder};
use chrono::Utc;
use serde_json::json;

/// Simulate Gitlab.com? check. In real code, this would check config/env.
fn is_gitlab_com() -> bool {
    // TODO: Implement actual check
    false
}

/// Simulate fetching the Jira Connect application key from settings.
fn jira_connect_application_key() -> Option<String> {
    // TODO: Replace with actual settings fetch
    std::env::var("JIRA_CONNECT_APPLICATION_KEY").ok()
}

/// GET /jira_connect/oauth_application_id
#[get("/jira_connect/oauth_application_id")]
pub async fn show() -> impl Responder {
    if show_application_id() {
        let key = jira_connect_application_key().unwrap_or_default();
        HttpResponse::Ok().json(json!({ "application_id": key }))
    } else {
        HttpResponse::NotFound().finish()
    }
}

fn show_application_id() -> bool {
    if is_gitlab_com() {
        return false;
    }
    jira_connect_application_key().is_some()
}
