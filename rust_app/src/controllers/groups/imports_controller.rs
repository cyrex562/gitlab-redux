// Ported from: orig_app/app/controllers/groups/imports_controller.rb
// Ported on: 2025-05-01
// This file implements the Groups::ImportsController from the Ruby codebase.

use actix_web::{web, HttpResponse, Responder};
use serde::Deserialize;

// Feature category: importers
// Urgency: low

#[derive(Debug, Deserialize)]
pub struct ContinueParams {
    pub to: Option<String>,
    pub notice: Option<String>,
    pub notice_now: Option<String>,
}

pub struct ImportsController;

impl ImportsController {
    /// Show import status for a group
    pub async fn show(
        group_id: web::Path<String>,
        params: web::Query<ContinueParams>,
    ) -> impl Responder {
        // TODO: Replace with real group lookup and import state
        let import_state = get_group_import_state(&group_id);
        match import_state.as_deref() {
            Some("finished") | None => {
                if let Some(to) = &params.to {
                    // Simulate redirect with notice
                    HttpResponse::Found()
                        .append_header(("Location", to.clone()))
                        .body(params.notice.clone().unwrap_or_default())
                } else {
                    // Redirect to group page with success notice
                    let group_url = format!("/groups/{}", group_id);
                    HttpResponse::Found()
                        .append_header(("Location", group_url))
                        .body("The group was successfully imported.")
                }
            }
            Some("failed") => {
                // Redirect to new group page with error alert
                let new_group_url = format!("/groups/{}/new", group_id);
                let last_error = get_group_import_last_error(&group_id)
                    .unwrap_or_else(|| "Unknown error".to_string());
                HttpResponse::Found()
                    .append_header(("Location", new_group_url))
                    .body(format!("Failed to import group: {}", last_error))
            }
            _ => {
                // Show page with notice_now (simulate flash.now)
                HttpResponse::Ok().body(params.notice_now.clone().unwrap_or_default())
            }
        }
    }
}

// --- Helpers (mocked for now) ---
fn get_group_import_state(_group_id: &str) -> Option<String> {
    // TODO: Replace with real DB/model logic
    Some("finished".to_string())
}

fn get_group_import_last_error(_group_id: &str) -> Option<String> {
    // TODO: Replace with real error lookup
    Some("Example error message".to_string())
}
