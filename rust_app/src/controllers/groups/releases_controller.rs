// Ported from: orig_app/app/controllers/groups/releases_controller.rb
// Ported on: 2025-05-01
// This file implements the Groups::ReleasesController from the Ruby codebase.

use actix_web::{web, HttpRequest, HttpResponse, Responder};
use serde::Serialize;

// Minimal stub for Release, to be replaced with real model
#[derive(Serialize)]
pub struct Release {
    pub id: i32,
    pub name: String,
    pub tag: String,
}

// Minimal stub for ReleaseSerializer
pub struct ReleaseSerializer;

impl ReleaseSerializer {
    pub fn represent(releases: &[Release]) -> Vec<serde_json::Value> {
        releases
            .iter()
            .map(|r| {
                serde_json::json!({
                    "id": r.id,
                    "name": r.name,
                    "tag": r.tag,
                })
            })
            .collect()
    }
}

pub struct GroupsReleasesController;

impl GroupsReleasesController {
    // GET /groups/{group_id}/releases (JSON only)
    pub async fn index(req: HttpRequest) -> impl Responder {
        // TODO: Extract group and current_user from request/context
        // TODO: Implement real finder logic and pagination
        let releases = vec![
            Release {
                id: 1,
                name: "Release 1".to_string(),
                tag: "v1.0.0".to_string(),
            },
            Release {
                id: 2,
                name: "Release 2".to_string(),
                tag: "v2.0.0".to_string(),
            },
        ];
        let json = ReleaseSerializer::represent(&releases);
        HttpResponse::Ok().json(json)
    }
}
