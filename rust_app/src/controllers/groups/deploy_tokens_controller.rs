// Ported from: orig_app/app/controllers/groups/deploy_tokens_controller.rb
// Ported on: 2025-05-01
// This file implements the Groups::DeployTokensController from the Ruby codebase.

use actix_web::{post, web, HttpResponse, Responder};

// TODO: Replace with actual group/user/service types
use crate::controllers::groups::application_controller::GroupsApplicationController;

#[post("/groups/{group_id}/deploy_tokens/revoke")]
pub async fn revoke(
    group_id: web::Path<String>,
    // TODO: Add user extraction (e.g., from session or request)
    params: web::Json<serde_json::Value>,
) -> impl Responder {
    // TODO: Implement authorization check (authorize_destroy_deploy_token!)
    // TODO: Call Groups::DeployTokens::RevokeService equivalent
    // For now, just simulate success

    // Redirect to group settings repository page with anchor
    let redirect_url = format!(
        "/groups/{}/-/settings/repository#js-deploy-tokens",
        group_id
    );
    HttpResponse::Found()
        .header("Location", redirect_url)
        .finish()
}
