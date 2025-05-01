// Ported from: orig_app/app/controllers/jira_connect/repositories_controller.rb
// This file implements the JiraConnect::RepositoriesController logic in Rust.

use actix_web::{get, post, web, HttpRequest, HttpResponse, Responder};
use serde::Deserialize;
use serde_json::json;
use uuid::Uuid;

use super::application_controller::JiraConnectApplicationController;

#[derive(Deserialize)]
pub struct QueryParams {
    pub id: Option<Uuid>,
    pub searchQuery: Option<String>,
    pub page: Option<u32>,
    pub limit: Option<u32>,
}

pub struct RepositoriesController {
    pub base: JiraConnectApplicationController,
}

impl RepositoriesController {
    pub fn new(base: JiraConnectApplicationController) -> Self {
        Self { base }
    }

    /// GET /jira_connect/repositories/search
    #[get("/jira_connect/repositories/search")]
    pub async fn search(req: HttpRequest, query: web::Query<QueryParams>) -> impl Responder {
        // TODO: Integrate with Project::with_jira_installation and RepositoryEntity
        // Placeholder: return empty containers array
        let containers = Vec::<serde_json::Value>::new();
        HttpResponse::Ok().json(json!({ "containers": containers }))
    }

    /// POST /jira_connect/repositories/associate
    #[post("/jira_connect/repositories/associate")]
    pub async fn associate(req: HttpRequest, query: web::Query<QueryParams>) -> impl Responder {
        // TODO: Integrate with Project::with_jira_installation and RepositoryEntity
        // Placeholder: always return not found
        HttpResponse::NotFound().json(json!({ "error": "Repository not found." }))
    }
}
