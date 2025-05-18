// This file was ported from: orig_app/app/controllers/projects/aws/configuration_controller.rb
// Ported on: 2025-05-07

use crate::controllers::projects::aws::base_controller::track_event;
use actix_web::{get, web, HttpResponse, Responder};

#[get("/projects/aws/configuration")]
pub async fn index() -> impl Responder {
    let js_data = serde_json::json!({});
    let _js_data = js_data.to_string();
    track_event("render_page", None, "project", "current_user"); // TODO: Replace with actual project/user
    HttpResponse::Ok().json(js_data)
}
