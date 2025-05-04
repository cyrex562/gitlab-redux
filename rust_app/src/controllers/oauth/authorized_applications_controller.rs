// Ported from: orig_app/app/controllers/oauth/authorized_applications_controller.rb
// Ported on: 2025-05-04
//
// Handles listing and revoking authorized OAuth applications for a user.
// Actions:
// - index: Returns 404 (HTML or JSON) since this is handled client-side
// - destroy: Revokes tokens for an application and redirects back to settings

use crate::auth::doorkeeper;
use crate::models::user::User;
use actix_web::{get, post, web, HttpRequest, HttpResponse, Responder};

/// Controller for managing authorized OAuth applications
pub struct AuthorizedApplicationsController;

impl AuthorizedApplicationsController {
    /// GET /oauth/authorized_applications
    /// Returns 404 since this is handled client-side
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
            // Return 404 page with "errors" layout
            HttpResponse::NotFound()
                .content_type("text/html")
                .body("<!DOCTYPE html><html><head><title>Not Found</title></head><body><h1>Not Found</h1></body></html>")
        }
    }

    /// POST /oauth/authorized_applications/{id}/destroy
    /// Revokes tokens for an application and redirects back to settings
    #[post("/oauth/authorized_applications/{id}/destroy")]
    pub async fn destroy(
        path: web::Path<i32>,
        params: web::Query<std::collections::HashMap<String, String>>,
        user: Option<web::ReqData<User>>, // Current user will be extracted from request
    ) -> impl Responder {
        // Get current user from request data
        let current_user = match user {
            Some(u) => u,
            None => return HttpResponse::Unauthorized().finish(),
        };

        // Check if revoking a specific token or all tokens
        if let Some(token_id) = params.get("token_id") {
            // Revoke specific token
            if let Err(_) = doorkeeper::revoke_token(token_id, &current_user).await {
                return HttpResponse::NotFound().finish();
            }
        } else {
            // Revoke all tokens for this application
            let app_id = path.into_inner();
            if let Err(_) = doorkeeper::revoke_application_tokens(app_id, &current_user).await {
                return HttpResponse::NotFound().finish();
            }
        }

        // Redirect back to applications page with success notice
        HttpResponse::Found()
            .append_header(("Location", "/user/settings/applications"))
            .append_header(("X-Notice", "Application has been revoked."))
            .finish()
    }
}

// Register the controller routes
pub fn configure(cfg: &mut web::ServiceConfig) {
    cfg.service(AuthorizedApplicationsController::index)
        .service(AuthorizedApplicationsController::destroy);
}
