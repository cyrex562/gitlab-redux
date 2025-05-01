// Ported from: orig_app/app/controllers/oauth/applications_controller.rb
// Ported: 2025-05-01
//
// This controller handles user OAuth applications (Doorkeeper::Application) management.
// Actions: index, show, create, renew

use crate::controllers::concerns::oauth_applications::{
    OauthApplications, OauthApplicationsHandler,
};
use actix_web::{get, post, web, HttpResponse, Responder};
use serde_json::json;

pub struct ApplicationsController {
    oauth_handler: OauthApplicationsHandler,
}

impl ApplicationsController {
    pub fn new() -> Self {
        Self {
            oauth_handler: OauthApplicationsHandler,
        }
    }

    #[get("/oauth/applications")]
    pub async fn index() -> impl Responder {
        // TODO: Set index vars (applications, tokens, etc.)
        HttpResponse::Ok().json(json!({"applications": [], "authorized_tokens": []}))
    }

    #[get("/oauth/applications/{id}")]
    pub async fn show(path: web::Path<i32>) -> impl Responder {
        // TODO: Show application details
        let id = path.into_inner();
        HttpResponse::Ok().json(json!({"application_id": id}))
    }

    #[post("/oauth/applications")]
    pub async fn create() -> impl Responder {
        // TODO: Create application logic
        // If persisted, return show; else, return index with errors
        HttpResponse::Ok().json(json!({"created": true}))
    }

    #[post("/oauth/applications/{id}/renew")]
    pub async fn renew(path: web::Path<i32>) -> impl Responder {
        // TODO: Renew secret logic
        let id = path.into_inner();
        HttpResponse::Ok().json(json!({"secret": "new_secret", "application_id": id}))
    }
}
