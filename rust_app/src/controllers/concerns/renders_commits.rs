use actix_web::{web, HttpRequest, HttpResponse};
use serde::{Deserialize, Serialize};
use std::sync::Arc;

/// This trait provides functionality for rendering commits in controllers
pub trait RendersCommits {
    /// Render commits for the current request
    fn render_commits(&self, req: &HttpRequest) -> HttpResponse;
}

#[derive(Debug, Serialize, Deserialize)]
pub struct Commit {
    id: String,
    title: String,
    message: String,
    author_name: String,
    author_email: String,
    authored_date: String,
    committed_date: String,
    created_at: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct RendersCommitsHandler {
    current_user: Option<Arc<User>>,
}

impl RendersCommitsHandler {
    pub fn new(current_user: Option<Arc<User>>) -> Self {
        RendersCommitsHandler { current_user }
    }

    fn fetch_commits(&self, ref_name: &str) -> Vec<Commit> {
        // This would be implemented to fetch commits from the database
        // For now, we'll return an empty vector
        Vec::new()
    }
}

impl RendersCommits for RendersCommitsHandler {
    fn render_commits(&self, req: &HttpRequest) -> HttpResponse {
        // Check if user is authenticated
        if self.current_user.is_none() {
            return HttpResponse::Unauthorized().finish();
        }

        // Get ref name from request
        let ref_name = req
            .match_info()
            .get("ref")
            .map(|s| s.to_string())
            .unwrap_or_else(|| "master".to_string());

        // Fetch commits
        let commits = self.fetch_commits(&ref_name);

        // Render commits as JSON
        HttpResponse::Ok()
            .content_type("application/json")
            .json(commits)
    }
}

// This would be implemented in a separate module
pub struct User {
    id: i32,
    // Add other fields as needed
}
