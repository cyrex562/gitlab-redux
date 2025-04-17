use serde::{Deserialize, Serialize};
use std::collections::HashSet;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct IntegrationParams {
    pub id: Option<i64>,
    pub integration: Option<IntegrationConfig>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct IntegrationConfig {
    // Common fields
    pub active: Option<bool>,
    pub api_url: Option<String>,
    pub webhook: Option<String>,
    pub username: Option<String>,
    pub password: Option<String>,
    pub token: Option<String>,
    pub project_url: Option<String>,
    pub issues_url: Option<String>,
    pub new_issue_url: Option<String>,
    pub server_url: Option<String>,
    pub server_host: Option<String>,
    pub server_port: Option<i32>,
    pub room: Option<String>,
    pub recipients: Option<String>,
    pub channels: Option<Vec<String>>,
    pub color: Option<String>,
    pub notify_only_broken_pipelines: Option<bool>,
    pub branches_to_be_notified: Option<String>,
    pub labels_to_be_notified: Option<Vec<String>>,
    pub labels_to_be_notified_behavior: Option<String>,
    pub disable_diffs: Option<bool>,
    pub send_from_committer_email: Option<bool>,
    pub push_events: Option<bool>,
    pub tag_push_events: Option<bool>,
    pub note_events: Option<bool>,
    pub confidential_note_events: Option<bool>,
    pub issues_events: Option<bool>,
    pub confidential_issues_events: Option<bool>,
    pub merge_requests_events: Option<bool>,
    pub job_events: Option<bool>,
    pub pipeline_events: Option<bool>,
    pub wiki_page_events: Option<bool>,
    pub deployment_events: Option<bool>,
    pub alert_events: Option<bool>,
    pub comment_on_event_enabled: Option<bool>,
    pub comment_detail: Option<String>,
    pub group_mention_events: Option<bool>,
    pub group_confidential_mention_events: Option<bool>,
    pub incident_events: Option<bool>,
    pub archive_trace_events: Option<bool>,
    pub exclude_service_accounts: Option<bool>,
    pub inherit_from_id: Option<i64>,
    pub type: Option<String>,
}

impl IntegrationParams {
    pub fn new() -> Self {
        Self {
            id: None,
            integration: None,
        }
    }

    pub fn with_id(mut self, id: i64) -> Self {
        self.id = Some(id);
        self
    }

    pub fn with_integration(mut self, integration: IntegrationConfig) -> Self {
        self.integration = Some(integration);
        self
    }

    pub fn get_allowed_params() -> HashSet<&'static str> {
        let mut params = HashSet::new();
        params.extend([
            "active", "api_url", "webhook", "username", "password", "token",
            "project_url", "issues_url", "new_issue_url", "server_url",
            "server_host", "server_port", "room", "recipients", "channels",
            "color", "notify_only_broken_pipelines", "branches_to_be_notified",
            "labels_to_be_notified", "labels_to_be_notified_behavior",
            "disable_diffs", "send_from_committer_email", "push_events",
            "tag_push_events", "note_events", "confidential_note_events",
            "issues_events", "confidential_issues_events", "merge_requests_events",
            "job_events", "pipeline_events", "wiki_page_events", "deployment_events",
            "alert_events", "comment_on_event_enabled", "comment_detail",
            "group_mention_events", "group_confidential_mention_events",
            "incident_events", "archive_trace_events", "exclude_service_accounts",
            "inherit_from_id", "type"
        ]);
        params
    }
} 