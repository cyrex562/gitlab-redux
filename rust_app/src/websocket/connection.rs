use axum::{
    extract::ws::{Message, WebSocket, WebSocketUpgrade},
    http::Request,
};
use futures_util::{SinkExt, StreamExt};
use serde_json::json;
use std::sync::Arc;
use tokio::sync::RwLock;
use tower_http::trace::TraceLayer;

use crate::models::User;
use crate::services::auth::AuthService;

use super::logging::Logging;

pub struct Connection {
    pub current_user: Option<Arc<User>>,
    pub request: Request<()>,
    pub auth_service: Arc<AuthService>,
}

impl Connection {
    pub fn new(request: Request<()>, auth_service: Arc<AuthService>) -> Self {
        Self {
            current_user: None,
            request,
            auth_service,
        }
    }

    pub async fn connect(&mut self) -> Result<(), String> {
        // Find user from bearer token or session
        self.current_user = self
            .find_user_from_bearer_token()
            .await
            .or_else(|| self.find_user_from_session_store());

        if self.current_user.is_none() {
            return Err("Unauthorized connection".to_string());
        }

        Ok(())
    }

    async fn find_user_from_bearer_token(&self) -> Option<Arc<User>> {
        // TODO: Implement token validation and user lookup
        None
    }

    fn find_user_from_session_store(&self) -> Option<Arc<User>> {
        // TODO: Implement session lookup
        None
    }

    pub fn notification_payload(&self, _: &str) -> serde_json::Value {
        json!({
            "params": self.request.uri().query().unwrap_or("")
        })
    }
}

impl Logging for Connection {
    fn notification_payload(&self, event: &str) -> serde_json::Value {
        self.notification_payload(event)
    }

    fn get_request(&self) -> &Request<()> {
        &self.request
    }

    fn get_current_user(&self) -> Option<&User> {
        self.current_user.as_deref()
    }
}
