// Ported from: orig_app/app/controllers/import/url_controller.rb
// Ported on: 2025-05-01
// This file implements the Import::UrlController from the Ruby codebase.
// See porting_log.txt for details.

use actix_web::{post, web, HttpResponse, Responder};
use serde::Deserialize;
use serde_json::json;

#[derive(Debug, Deserialize)]
pub struct ValidateParams {
    pub user: Option<String>,
    pub password: Option<String>,
    pub url: Option<String>,
}

pub struct UrlController;

impl UrlController {
    /// POST /import/url/validate
    #[post("/import/url/validate")]
    pub async fn validate(params: web::Json<ValidateParams>) -> impl Responder {
        // Placeholder for Import::ValidateRemoteGitEndpointService logic
        // In production, call a real service here
        let valid = params.url.as_ref().map_or(false, |u| u.starts_with("http"));
        if valid {
            HttpResponse::Ok().json(json!({ "success": true }))
        } else {
            HttpResponse::Ok()
                .json(json!({ "success": false, "message": "Invalid or missing URL" }))
        }
    }
}
