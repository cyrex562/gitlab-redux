use actix_web::{web, HttpResponse, Responder};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

use crate::harbor::query::Query;

#[derive(Debug, Serialize, Deserialize)]
pub struct Artifact {
    pub id: i32,
    pub repository_id: i32,
    pub name: String,
    pub digest: String,
    pub size: i64,
    pub created_at: String,
    pub updated_at: String,
}

pub struct ArtifactController;

impl ArtifactController {
    pub async fn index(query: web::Query<HashMap<String, String>>) -> impl Responder {
        let query = Query::new(query.into_inner());

        if !query.is_valid() {
            return HttpResponse::UnprocessableEntity().json(serde_json::json!({
                "message": "Invalid parameters",
                "errors": query.errors()
            }));
        }

        let artifacts = query.artifacts();
        HttpResponse::Ok().json(artifacts)
    }

    pub async fn show(id: i32) -> impl Responder {
        // TODO: Implement artifact retrieval
        HttpResponse::Ok().json(Artifact {
            id,
            repository_id: 0,
            name: String::new(),
            digest: String::new(),
            size: 0,
            created_at: chrono::Utc::now().to_rfc3339(),
            updated_at: chrono::Utc::now().to_rfc3339(),
        })
    }
} 