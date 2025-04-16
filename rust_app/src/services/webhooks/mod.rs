pub mod hook_actions;
pub mod hook_execution_notice;
pub mod hook_log_actions;

use crate::models::user::User;
use actix_web::{web, HttpResponse};
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct WebHookResult {
    pub success: bool,
    pub message: String,
    pub http_status: Option<i32>,
}

// TODO: Implement WebHook model
pub struct WebHook {
    // Add fields based on hook_param_names
    pub enable_ssl_verification: bool,
    pub name: String,
    pub description: Option<String>,
    pub token: Option<String>,
    pub url: String,
    pub push_events_branch_filter: Option<String>,
    pub branch_filter_strategy: Option<String>,
    pub custom_webhook_template: Option<String>,
    pub url_variables: std::collections::HashMap<String, String>,
    pub custom_headers: std::collections::HashMap<String, String>,
}

// TODO: Implement WebHookLog model
pub struct WebHookLog {
    // Add necessary fields
}
