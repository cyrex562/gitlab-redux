// Ported from: orig_app/app/controllers/groups/group_links_controller.rb
// Ported on: 2025-05-01
// This file implements the Groups::GroupLinksController from the Ruby codebase.

use actix_web::{delete, put, web, HttpResponse, Responder};
use serde::Deserialize;

// TODO: Replace with actual group/user/service types
use crate::controllers::groups::application_controller::GroupsApplicationController;

#[derive(Deserialize)]
pub struct GroupLinkParams {
    pub group_access: Option<String>,
    pub expires_at: Option<String>,
    pub member_role_id: Option<i64>,
}

#[put("/groups/{group_id}/group_links/{id}")]
pub async fn update(
    path: web::Path<(String, i64)>,
    params: web::Json<GroupLinkParams>,
    // TODO: Add user extraction (e.g., from session or request)
) -> impl Responder {
    // TODO: Authorization: authorize_admin_group!
    // TODO: Find group and group_link by IDs
    // TODO: Call Groups::GroupLinks::UpdateService equivalent
    // Simulate group_link with expires/expires_soon fields
    let expires = params.expires_at.is_some();
    let expires_soon = false; // TODO: Implement logic
    if expires {
        HttpResponse::Ok().json(serde_json::json!({
            "expires_in": "in 5 days", // TODO: Implement helpers.time_ago_with_tooltip
            "expires_soon": expires_soon
        }))
    } else {
        HttpResponse::Ok().json(serde_json::json!({}))
    }
}

#[delete("/groups/{group_id}/group_links/{id}")]
pub async fn destroy(
    path: web::Path<(String, i64)>,
    // TODO: Add user extraction (e.g., from session or request)
) -> impl Responder {
    // TODO: Authorization: authorize_admin_group!
    // TODO: Find group and group_link by IDs
    // TODO: Call Groups::GroupLinks::DestroyService equivalent
    // Simulate HTML and JS response
    // For now, always redirect to group members page
    let group_id = &path.0;
    let redirect_url = format!("/groups/{}/group_members", group_id);
    HttpResponse::Found()
        .header("Location", redirect_url)
        .finish()
}
