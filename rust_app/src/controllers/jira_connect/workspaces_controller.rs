// Ported from: orig_app/app/controllers/jira_connect/workspaces_controller.rb
// Ported on: 2025-05-01
// This file implements the JiraConnect::WorkspacesController logic in Rust.

use actix_web::{get, web, HttpRequest, HttpResponse, Responder};
use serde::{Deserialize, Serialize};
use serde_json::json;

use super::application_controller::JiraConnectApplicationController;

#[derive(Serialize)]
pub struct WorkspaceJson {
    pub id: u64,
    pub name: String,
}

pub struct WorkspacesController {
    pub base: JiraConnectApplicationController,
}

impl WorkspacesController {
    pub fn new(base: JiraConnectApplicationController) -> Self {
        Self { base }
    }

    /// GET /jira_connect/workspaces/search
    #[get("/jira_connect/workspaces/search")]
    pub async fn search(req: HttpRequest, query: web::Query<SearchQuery>) -> impl Responder {
        // TODO: Integrate with Namespace::without_project_namespaces and with_jira_installation
        // and sanitize search_query
        let search_query = sanitize(&query.search_query);
        let workspaces = find_workspaces(&search_query);
        HttpResponse::Ok().json(json!({ "workspaces": workspaces }))
    }
}

#[derive(Deserialize)]
pub struct SearchQuery {
    #[serde(rename = "searchQuery")]
    pub search_query: String,
}

fn sanitize(input: &str) -> String {
    // TODO: Use a real HTML sanitizer if needed
    input.replace('<', "").replace('>', "")
}

fn find_workspaces(_search_query: &str) -> Vec<WorkspaceJson> {
    // TODO: Query namespaces with filters as in Ruby code
    // Placeholder: return empty list
    vec![]
}
