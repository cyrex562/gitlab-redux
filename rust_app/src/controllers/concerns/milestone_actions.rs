use actix_web::{web, HttpResponse, Responder};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

#[derive(Debug, Serialize, Deserialize)]
pub struct Issue {
    pub id: i32,
    pub title: String,
    pub project_id: i32,
    pub project_name: Option<String>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct MergeRequest {
    pub id: i32,
    pub title: String,
    pub project_id: i32,
    pub project_name: Option<String>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct User {
    pub id: i32,
    pub name: String,
    pub username: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct Label {
    pub id: i32,
    pub title: String,
    pub color: String,
    pub description: Option<String>,
}

pub trait MilestoneActions {
    fn issues(&self, milestone_id: i32, show_project_name: bool) -> impl Responder {
        let issues = self.get_sorted_issues(milestone_id);
        let data = HashMap::from([
            ("issues", issues),
            ("show_project_name", show_project_name),
        ]);
        self.render_tab("shared/milestones/_issues_tab", data)
    }

    fn merge_requests(&self, milestone_id: i32, show_project_name: bool) -> impl Responder {
        let merge_requests = self.get_sorted_merge_requests(milestone_id);
        let data = HashMap::from([
            ("merge_requests", merge_requests),
            ("show_project_name", show_project_name),
        ]);
        self.render_tab("shared/milestones/_merge_requests_tab", data)
    }

    fn participants(&self, milestone_id: i32) -> impl Responder {
        let users = self.get_issue_participants(milestone_id);
        let data = HashMap::from([("users", users)]);
        self.render_tab("shared/milestones/_participants_tab", data)
    }

    fn labels(&self, milestone_id: i32) -> impl Responder {
        let labels = self.get_issue_labels(milestone_id);
        let data = HashMap::from([("labels", labels)]);
        self.render_tab("shared/milestones/_labels_tab", data)
    }

    // Required methods to be implemented by concrete types
    fn get_sorted_issues(&self, milestone_id: i32) -> Vec<Issue>;
    fn get_sorted_merge_requests(&self, milestone_id: i32) -> Vec<MergeRequest>;
    fn get_issue_participants(&self, milestone_id: i32) -> Vec<User>;
    fn get_issue_labels(&self, milestone_id: i32) -> Vec<Label>;
    fn render_tab(&self, template: &str, data: HashMap<&str, serde_json::Value>) -> HttpResponse;
    fn milestone_redirect_path(&self) -> String;

    // Helper method for JSON responses
    fn tabs_json(&self, template: &str, data: HashMap<&str, serde_json::Value>) -> HttpResponse {
        let html = self.render_tab(template, data);
        HttpResponse::Ok().json(HashMap::from([("html", html)]))
    }
} 