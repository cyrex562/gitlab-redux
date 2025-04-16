use actix_web::{web, HttpResponse, Responder};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

use crate::harbor::query::Query;

#[derive(Debug, Serialize, Deserialize)]
pub struct Repository {
    pub id: i32,
    pub name: String,
    pub project_id: i32,
    pub artifact_count: i32,
    pub created_at: String,
    pub updated_at: String,
}

pub struct RepositoryController;

impl RepositoryController {
    pub async fn index(query: web::Query<HashMap<String, String>>) -> impl Responder {
        let query = Query::new(query.into_inner());

        if !query.is_valid() {
            return HttpResponse::UnprocessableEntity().json(serde_json::json!({
                "message": "Invalid parameters",
                "errors": query.errors()
            }));
        }

        let repositories = query.repositories();
        HttpResponse::Ok().json(repositories)
    }

    pub async fn show(id: i32) -> impl Responder {
        // TODO: Implement repository retrieval
        HttpResponse::Ok().json(Repository {
            id,
            name: String::new(),
            project_id: 0,
            artifact_count: 0,
            created_at: chrono::Utc::now().to_rfc3339(),
            updated_at: chrono::Utc::now().to_rfc3339(),
        })
    }
} 