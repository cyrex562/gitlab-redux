// Ported from: orig_app/app/controllers/groups/registry/repositories_controller.rb
// Date ported: 2025-05-01
//
// Handles group-level container registry repositories endpoints.

use actix_web::{get, web, HttpResponse, Responder};

// Placeholder for feature flag and authorization logic
dfn push_frontend_feature_flag(_flag: &str, _group_id: i32) {
    // TODO: Implement feature flag logic
}

fn verify_container_registry_enabled() -> bool {
    // TODO: Replace with actual config check
    true
}

fn authorize_read_container_image(_user_id: i32, _group_id: i32) -> bool {
    // TODO: Replace with actual permission check
    true
}

#[get("/groups/{group_id}/-/container_registries")]
pub async fn index(path: web::Path<i32>, query: web::Query<serde_json::Value>) -> impl Responder {
    let group_id = path.into_inner();
    let user_id = 0; // TODO: Extract from session/auth

    if !verify_container_registry_enabled() {
        return HttpResponse::NotFound().finish();
    }
    if !authorize_read_container_image(user_id, group_id) {
        return HttpResponse::NotFound().finish();
    }

    push_frontend_feature_flag("show_container_registry_tag_signatures", group_id);

    // TODO: Implement actual finder and serializer logic
    // Placeholder response
    HttpResponse::Ok().json(serde_json::json!({
        "repositories": [],
        "pagination": {}
    }))
}

#[get("/groups/{group_id}/-/container_registries/{id}")]
pub async fn show(path: web::Path<(i32, i32)>) -> impl Responder {
    // The show action renders index to allow frontend routing to work on page refresh
    index(web::Path::from(path.0), web::Query(serde_json::Value::Null)).await
}
