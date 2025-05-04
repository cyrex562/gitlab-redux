use serde::{Deserialize, Serialize};
use serde_json::{json, Value};
use std::sync::Arc;

#[derive(Debug, Clone)]
pub struct PreAuth {
    pub client: Arc<Client>,
    pub error: Option<String>,
    pub scopes: Vec<String>,
}

impl PreAuth {
    pub fn new(client: Arc<Client>) -> Self {
        Self {
            client,
            error: None,
            scopes: Vec::new(),
        }
    }

    pub fn authorizable(&self) -> bool {
        self.error.is_none()
    }

    pub fn client(&self) -> Arc<Client> {
        self.client.clone()
    }

    pub fn set_error(&mut self, error: &str) {
        self.error = Some(error.to_string());
    }
}

#[derive(Debug, Clone)]
pub struct Client {
    pub application: Arc<Application>,
}

impl Client {
    pub fn application(&self) -> Option<Arc<Application>> {
        Some(self.application.clone())
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Application {
    pub id: i32,
    pub name: String,
    pub scopes: Vec<String>,
    pub confidential: bool,
    pub trusted: bool,
}

impl Application {
    pub fn includes_scope(&self, scope: &str) -> bool {
        self.scopes.contains(&scope.to_string())
    }

    pub fn confidential(&self) -> bool {
        self.confidential
    }

    pub fn trusted(&self) -> bool {
        self.trusted
    }

    pub fn update_scopes(&mut self, scopes: Vec<String>) -> Result<(), Box<dyn std::error::Error>> {
        self.scopes = scopes;
        Ok(())
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AccessToken {
    pub token: String,
    pub resource_owner_id: Option<i64>,
    pub application_id: Option<i64>,
    pub expires_in: Option<i64>,
    pub revoked_at: Option<chrono::DateTime<chrono::Utc>>,
    pub scopes: Vec<String>,
}

impl AccessToken {
    pub fn find_by_token(token: &str) -> Option<Self> {
        // TODO: Implement actual database lookup
        // For now return None to simulate token not found
        None
    }

    pub fn accessible(&self) -> bool {
        !self.revoked() && !self.expired()
    }

    pub fn revoked(&self) -> bool {
        self.revoked_at.is_some()
    }

    pub fn expired(&self) -> bool {
        if let Some(expires_in) = self.expires_in {
            // TODO: Implement proper expiration check against creation time
            expires_in <= 0
        } else {
            false
        }
    }

    pub fn as_json(&self) -> Value {
        json!({
            "resource_owner_id": self.resource_owner_id,
            "scope": self.scopes.join(" "),
            "expires_in": self.expires_in,
            "application": {
                "uid": self.application_id
            },
            "created_at": chrono::Utc::now().timestamp()
        })
    }
}

pub struct Authorization;

impl Authorization {
    pub fn authorize(&self) -> AuthorizationResponse {
        AuthorizationResponse {
            redirect_uri: String::new(),
        }
    }
}

#[derive(Debug)]
pub struct AuthorizationResponse {
    pub redirect_uri: String,
}

pub mod oauth {
    use super::*;

    pub struct Client;

    impl Client {
        pub fn find(client_id: &str) -> Option<Arc<super::Client>> {
            // TODO: Implement client lookup
            None
        }
    }
}
