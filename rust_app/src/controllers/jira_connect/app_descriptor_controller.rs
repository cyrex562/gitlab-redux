// Ported from: orig_app/app/controllers/jira_connect/app_descriptor_controller.rb
// This controller returns an app descriptor for use with Jira in development mode.

use actix_web::{get, web, HttpResponse, Responder};
use serde_json::json;
use chrono::Utc;

#[get("/jira_connect/app_descriptor")]
pub async fn show() -> impl Responder {
    HttpResponse::Ok().json(json!({
        "name": "GitLab Jira Connect", // TODO: Replace with dynamic value if needed
        "description": "Integrate commits, branches and merge requests from GitLab into Jira",
        "key": "gitlab-jira-connect", // TODO: Replace with dynamic value if needed
        "baseUrl": "https://gitlab.com", // TODO: Replace with dynamic value if needed
        "lifecycle": {
            "installed": "/events/installed", // TODO: Implement path helpers
            "uninstalled": "/events/uninstalled"
        },
        "vendor": {
            "name": "GitLab",
            "url": "https://gitlab.com"
        },
        "links": {
            "documentation": "https://docs.gitlab.com/ee/integration/jira/development_panel.md"
        },
        "authentication": {
            "type": "jwt"
        },
        "modules": {}, // TODO: Implement modules logic
        "scopes": ["READ", "WRITE", "DELETE"],
        "apiVersion": 1,
        "apiMigrations": {
            "context-qsh": true,
            "signed-install": true,
            "gdpr": true
        }
    }))
}
