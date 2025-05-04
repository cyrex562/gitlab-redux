// Ported from orig_app/app/controllers/concerns/oauth_applications.rb
// Handles OAuth application-related controller logic

use crate::models::oauth::Application;
use crate::models::user::User;
use std::collections::HashMap;
use std::sync::Arc;

pub const CREATED_SESSION_KEY: &str = "oauth_applications_created";

pub async fn load_scopes() -> Vec<String> {
    let mut scopes = vec![
        "api".to_string(),
        "read_user".to_string(),
        "read_repository".to_string(),
        "write_repository".to_string(),
        "admin_repository".to_string(),
    ];
    scopes.sort();
    scopes
}

pub async fn set_application(user: Arc<User>, id: i32) -> Result<Application, &'static str> {
    // TODO: Implement application lookup
    // For now, return error since we haven't implemented persistence yet
    Err("Application not found")
}

pub async fn set_index_vars(user: Arc<User>) -> Vec<Application> {
    // TODO: Implement loading user's applications
    vec![]
}

pub async fn verify_oauth_applications_enabled() -> bool {
    // TODO: Load from configuration
    true
}
