use actix_web::{web, HttpResponse, Responder};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

use crate::harbor::query::Query;

#[derive(Debug, Serialize, Deserialize)]
pub struct Tag {
    pub id: i32,
    pub repository_id: i32,
    pub artifact_id: i32,
    pub name: String,
    pub created_at: String,
    pub updated_at: String,
}

pub struct TagController;

impl TagController {
    pub async fn index(query: web::Query<HashMap<String, String>>) -> impl Responder {
        let query = Query::new(query.into_inner());

        if !query.is_valid() {
            return HttpResponse::UnprocessableEntity().json(serde_json::json!({
                "message": "Invalid parameters",
                "errors": query.errors()
            }));
        }

        let tags = query.tags();
        HttpResponse::Ok().json(tags)
    }

    pub async fn show(id: i32) -> impl Responder {
        // TODO: Implement tag retrieval
        HttpResponse::Ok().json(Tag {
            id,
            repository_id: 0,
            artifact_id: 0,
            name: String::new(),
            created_at: chrono::Utc::now().to_rfc3339(),
            updated_at: chrono::Utc::now().to_rfc3339(),
        })
    }
} 