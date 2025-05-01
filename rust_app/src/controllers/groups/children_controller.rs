// Ported from: orig_app/app/controllers/groups/children_controller.rb
// Ported on: 2025-05-01
// This controller handles listing children (subgroups and projects) for a group.

use actix_web::{get, web, HttpResponse, Responder};
use serde::Deserialize;

#[derive(Debug, Deserialize)]
pub struct ChildrenQuery {
    pub parent_id: Option<i32>,
    pub sort: Option<String>,
    pub filter: Option<String>,
    pub per_page: Option<u32>,
}

#[get("/groups/{group_id}/children")]
pub async fn index(
    group_id: web::Path<i32>,
    query: web::Query<ChildrenQuery>,
    // Add user/context extractors as needed
) -> impl Responder {
    // Validate per_page
    if let Some(per_page) = query.per_page {
        if per_page < 1 {
            return HttpResponse::BadRequest().json(serde_json::json!({
                "message": "per_page does not have a valid value"
            }));
        }
    }

    // Determine parent group
    let parent = if let Some(parent_id) = query.parent_id {
        // TODO: Replace with actual group finder logic
        Some(parent_id)
    } else {
        Some(group_id.into_inner())
    };

    if parent.is_none() {
        return HttpResponse::NotFound().finish();
    }

    // TODO: Setup children using a service/finder (stubbed here)
    let children = Vec::<serde_json::Value>::new();

    // TODO: Implement serializer and expand_hierarchy logic if needed
    HttpResponse::Ok().json(children)
}
