use async_trait::async_trait;
use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SlackInstallationResult {
    pub success: bool,
    pub message: Option<String>,
}

#[async_trait]
pub trait SlackControllerSettings {
    async fn slack_auth(&self) -> Result<SlackInstallationResult, Box<dyn std::error::Error>>;
    async fn destroy(&self) -> Result<(), Box<dyn std::error::Error>>;
}

pub struct SlackController {
    integration_id: i64,
    oauth_state: String,
}

impl SlackController {
    pub fn new(integration_id: i64, oauth_state: String) -> Self {
        Self {
            integration_id,
            oauth_state,
        }
    }

    pub async fn handle_oauth_error(&self, error: &str) -> Result<(), Box<dyn std::error::Error>> {
        if error == "access_denied" {
            return Err("Access request canceled".into());
        }
        Ok(())
    }

    pub fn check_oauth_state(&self, state: &str) -> Result<(), Box<dyn std::error::Error>> {
        if state != self.oauth_state {
            return Err("Invalid OAuth state".into());
        }
        Ok(())
    }
}

#[async_trait]
impl SlackControllerSettings for SlackController {
    async fn slack_auth(&self) -> Result<SlackInstallationResult, Box<dyn std::error::Error>> {
        // Implementation would depend on your OAuth flow
        Ok(SlackInstallationResult {
            success: true,
            message: None,
        })
    }

    async fn destroy(&self) -> Result<(), Box<dyn std::error::Error>> {
        // Implementation would depend on your integration deletion logic
        Ok(())
    }
} 