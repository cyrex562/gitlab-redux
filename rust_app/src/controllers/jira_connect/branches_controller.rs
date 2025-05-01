// Ported from: orig_app/app/controllers/jira_connect/branches_controller.rb
// Ported on: 2025-05-01
// This file implements the JiraConnect::BranchesController from the Ruby codebase.
//
// See porting_log.txt for details.

use actix_web::{get, post, web, HttpRequest, HttpResponse, Responder};
use serde_json::json;

use super::application_controller::JiraConnectApplicationController;

pub struct BranchesController {
    pub base: JiraConnectApplicationController,
}

impl BranchesController {
    pub fn new(base: JiraConnectApplicationController) -> Self {
        Self { base }
    }

    /// GET /jira_connect/branches/new
    #[get("/jira_connect/branches/new")]
    pub async fn new_branch(req: HttpRequest) -> impl Responder {
        // TODO: Implement user authentication check (authenticate_user!)
        let new_branch_data = Self::new_branch_data(&req);
        HttpResponse::Ok().json(new_branch_data)
    }

    /// GET /jira_connect/branches/route
    #[get("/jira_connect/branches/route")]
    pub async fn route(req: HttpRequest) -> impl Responder {
        // TODO: Implement current_jira_installation.proxy? logic
        // For now, just redirect to new_branch
        let query = req.query_string();
        let redirect_url = format!("/jira_connect/branches/new?{}", query);
        HttpResponse::Found()
            .header("Location", redirect_url)
            .finish()
    }

    fn initial_branch_name(req: &HttpRequest) -> Option<String> {
        let params = req.query_string();
        let query: Vec<_> = url::form_urlencoded::parse(params.as_bytes()).collect();
        let issue_key = query
            .iter()
            .find(|(k, _)| k == "issue_key")
            .map(|(_, v)| v.to_string());
        let issue_summary = query
            .iter()
            .find(|(k, _)| k == "issue_summary")
            .map(|(_, v)| v.to_string());
        if let Some(issue_key) = issue_key {
            // TODO: Implement Issue::to_branch_name equivalent
            Some(format!(
                "{}-{}",
                issue_key,
                issue_summary.unwrap_or_default()
            ))
        } else {
            None
        }
    }

    fn new_branch_data(req: &HttpRequest) -> serde_json::Value {
        json!({
            "initial_branch_name": Self::initial_branch_name(req),
            "success_state_svg_path": "/static/illustrations/empty-state/empty-merge-requests-md.svg"
        })
    }
}
