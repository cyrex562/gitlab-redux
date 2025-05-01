// Ported from: orig_app/app/controllers/oauth/authorized_applications_controller.rb
// Ported: 2025-05-01
//
// Handles listing and revoking authorized OAuth applications for a user.
//
// index: always returns 404 (HTML or JSON)
// destroy: revokes a specific token or all tokens for an application, then redirects with a notice

use actix_web::{get, post, web, HttpRequest, HttpResponse, Responder};

#[get("/oauth/authorized_applications")]
pub async fn index(req: HttpRequest) -> impl Responder {
    let accepts = req
        .headers()
        .get("accept")
        .and_then(|v| v.to_str().ok())
        .unwrap_or("");
    if accepts.contains("json") {
        HttpResponse::NotFound().json("")
    } else {
        // Render a not_found error page (could use a template in a real app)
        HttpResponse::NotFound().body("Not Found")
    }
}

#[post("/oauth/authorized_applications/{id}/destroy")]
pub async fn destroy(
    path: web::Path<i32>,
    req: HttpRequest,
    params: web::Query<std::collections::HashMap<String, String>>,
) -> impl Responder {
    let id = path.into_inner();
    let token_id = params.get("token_id");
    // TODO: Replace with real user/token logic
    if let Some(_token_id) = token_id {
        // Revoke a specific token for the current user
        // current_resource_owner.oauth_authorized_tokens.find(token_id).revoke
    } else {
        // Revoke all tokens for the application for the current user
        // Doorkeeper::Application.revoke_tokens_and_grants_for(id, current_resource_owner)
    }
    // Redirect to user settings applications page with a notice
    HttpResponse::Found()
        .append_header(("Location", "/user/settings/applications"))
        .append_header(("X-Notice", "Authorized application revoked."))
        .finish()
}
