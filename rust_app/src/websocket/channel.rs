use axum::{
    extract::ws::{Message, WebSocket},
    http::Request,
};
use futures_util::{SinkExt, StreamExt};
use serde_json::json;
use std::sync::Arc;
use tokio::sync::RwLock;
use tokio::time::{interval, Duration};

use crate::models::User;
use crate::services::auth::AuthService;

use super::connection::Connection;
use super::logging::Logging;

pub struct Channel {
    pub connection: Arc<RwLock<Connection>>,
    pub params: serde_json::Value,
    pub subscription_rejected: bool,
    pub subscription_confirmed: bool,
}

impl Channel {
    pub fn new(connection: Arc<RwLock<Connection>>, params: serde_json::Value) -> Self {
        Self {
            connection,
            params,
            subscription_rejected: false,
            subscription_confirmed: false,
        }
    }

    pub async fn subscribe(&mut self) -> Result<(), String> {
        // Validate token scope before subscribing
        self.validate_token_scope().await?;

        // Set up periodic validation
        self.setup_periodic_validation();

        self.subscription_confirmed = true;
        Ok(())
    }

    pub async fn validate_token_scope(&self) -> Result<(), String> {
        // TODO: Implement token validation with proper scopes
        let scopes = self.authorization_scopes();

        // Validate token with scopes
        let connection = self.connection.read().await;
        if let Some(auth_service) = &connection.auth_service {
            // Validate token with scopes
            // auth_service.validate_token(scopes).await?;
        }

        Ok(())
    }

    fn setup_periodic_validation(&self) {
        let connection = Arc::clone(&self.connection);
        let params = self.params.clone();

        tokio::spawn(async move {
            let mut interval = interval(Duration::from_secs(600)); // 10 minutes

            loop {
                interval.tick().await;

                // Validate token scope
                let mut conn = connection.write().await;
                // TODO: Implement periodic validation
            }
        });
    }

    pub fn authorization_scopes(&self) -> Vec<String> {
        vec!["api".to_string(), "read_api".to_string()]
    }

    pub fn client_subscribed(&self) -> bool {
        !self.subscription_rejected && self.subscription_confirmed
    }

    pub async fn handle_authentication_error(&mut self) {
        if self.client_subscribed() {
            self.unsubscribe_from_channel().await;
        } else {
            self.reject();
        }
    }

    pub async fn unsubscribe_from_channel(&mut self) {
        // TODO: Implement unsubscribe logic
    }

    pub fn reject(&mut self) {
        self.subscription_rejected = true;
    }

    pub fn notification_payload(&self, event: &str) -> serde_json::Value {
        let mut payload = json!({});

        if let serde_json::Value::Object(ref mut map) = payload {
            // Add params except channel
            if let Some(params) = self.params.as_object() {
                for (key, value) in params {
                    if key != "channel" {
                        map.insert(key.clone(), value.clone());
                    }
                }
            }
        }

        payload
    }
}

impl Logging for Channel {
    fn notification_payload(&self, event: &str) -> serde_json::Value {
        self.notification_payload(event)
    }

    fn get_request(&self) -> &Request<()> {
        // This is a bit tricky in Rust - we need to access the connection's request
        // For now, we'll return a dummy request
        // TODO: Implement proper request access
        &Request::new(())
    }

    fn get_current_user(&self) -> Option<&User> {
        // This is also tricky in Rust - we need to access the connection's current_user
        // For now, we'll return None
        // TODO: Implement proper user access
        None
    }
}
