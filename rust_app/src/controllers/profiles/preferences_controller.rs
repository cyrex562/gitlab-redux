// Ported from: orig_app/app/controllers/profiles/preferences_controller.rb
// Ported on: 2025-05-07
// This file implements the Profiles::PreferencesController from the Ruby codebase.

use actix_web::{get, post, web, HttpResponse, Responder};
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct PreferencesParams {
    pub color_scheme_id: Option<i32>,
    pub color_mode_id: Option<i32>,
    pub diffs_deletion_color: Option<String>,
    pub diffs_addition_color: Option<String>,
    pub home_organization_id: Option<i32>,
    pub layout: Option<String>,
    pub dashboard: Option<String>,
    pub project_view: Option<String>,
    pub theme_id: Option<i32>,
    pub first_day_of_week: Option<i32>,
    pub preferred_language: Option<String>,
    pub time_display_relative: Option<bool>,
    pub time_display_format: Option<String>,
    pub show_whitespace_in_diffs: Option<bool>,
    pub view_diffs_file_by_file: Option<bool>,
    pub tab_width: Option<i32>,
    pub sourcegraph_enabled: Option<bool>,
    pub gitpod_enabled: Option<bool>,
    pub extensions_marketplace_enabled: Option<bool>,
    pub render_whitespace_in_code: Option<bool>,
    pub project_shortcut_buttons: Option<bool>,
    pub keyboard_shortcuts_enabled: Option<bool>,
    pub markdown_surround_selection: Option<bool>,
    pub markdown_automatic_lists: Option<bool>,
    pub use_new_navigation: Option<bool>,
    pub enabled_following: Option<bool>,
    pub use_work_items_view: Option<bool>,
    pub text_editor: Option<String>,
}

#[get("/profiles/preferences")]
pub async fn show() -> impl Responder {
    // TODO: Implement fetching current user's preferences
    HttpResponse::Ok().json(serde_json::json!({
        "message": "Show user preferences (stub)"
    }))
}

#[post("/profiles/preferences")]
pub async fn update(params: web::Json<PreferencesParams>) -> impl Responder {
    // TODO: Implement updating user preferences
    // Simulate success
    let result_status = "success";
    if result_status == "success" {
        HttpResponse::Ok().json(serde_json::json!({
            "type": "notice",
            "message": "Preferences saved."
        }))
    } else {
        HttpResponse::BadRequest().json(serde_json::json!({
            "type": "alert",
            "message": "Failed to save preferences."
        }))
    }
}

// TODO: Add integration with user authentication and actual preferences storage.
