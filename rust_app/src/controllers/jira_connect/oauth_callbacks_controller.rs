// Ported from: orig_app/app/controllers/jira_connect/oauth_callbacks_controller.rb
// This controller serves as a landing page after installing/authenticating the GitLab.com for Jira App.

use actix_web::{get, HttpResponse, Responder};

#[get("/jira_connect/oauth/callback")]
pub async fn index() -> impl Responder {
    HttpResponse::Ok().body("Jira OAuth callback received. You may close this page.")
}
