use serde_json::{json, Value};
use actix_web::HttpRequest;
use std::sync::Arc;
use crate::models::user::User;

pub trait Logging {
    fn notification_payload(&self) -> Value;
    fn get_request(&self) -> &HttpRequest;
    fn get_current_user(&self) -> Option<Arc<User>>;
    fn log_connect(&self);
    fn log_disconnect(&self);
    fn log_subscribe(&self);
    fn log_unsubscribe(&self);
    fn log_event(&self, event: &str, data: &Value);
}

pub trait EnhancedLogging: Logging {
    fn enhanced_notification_payload(&self, event: &str) -> Value {
        let mut payload = self.notification_payload();

        if let Value::Object(ref mut map) = payload {
            // Add correlation ID
            map.insert(
                "correlation_id".to_string(),
                json!(self
                    .get_request()
                    .headers()
                    .get("x-request-id")
                    .and_then(|h| h.to_str().ok())
                    .unwrap_or("")),
            );

            // Add user information
            if let Some(user) = self.get_current_user() {
                map.insert("user_id".to_string(), json!(user.id));
                map.insert("username".to_string(), json!(user.username));
            }

            // Add request information
            map.insert(
                "remote_ip".to_string(),
                json!(self
                    .get_request()
                    .headers()
                    .get("x-forwarded-for")
                    .and_then(|h| h.to_str().ok())
                    .unwrap_or("")),
            );

            map.insert(
                "ua".to_string(),
                json!(self
                    .get_request()
                    .headers()
                    .get("user-agent")
                    .and_then(|h| h.to_str().ok())
                    .unwrap_or("")),
            );
        }

        payload
    }
}

pub struct Logger {
    request: HttpRequest,
    user: Option<Arc<User>>,
}

impl Logger {
    pub fn new(request: HttpRequest, user: Option<Arc<User>>) -> Self {
        Self { request, user }
    }
}

impl Logging for Logger {
    fn notification_payload(&self) -> Value {
        json!({
            "type": "log",
            "timestamp": chrono::Utc::now().to_rfc3339(),
        })
    }

    fn get_request(&self) -> &HttpRequest {
        &self.request
    }

    fn get_current_user(&self) -> Option<Arc<User>> {
        self.user.clone()
    }

    fn log_connect(&self) {
        // TODO: Implement actual logging
        println!("WebSocket connected");
    }

    fn log_disconnect(&self) {
        // TODO: Implement actual logging
        println!("WebSocket disconnected");
    }

    fn log_subscribe(&self) {
        // TODO: Implement actual logging
        println!("Client subscribed to channel");
    }

    fn log_unsubscribe(&self) {
        // TODO: Implement actual logging
        println!("Client unsubscribed from channel");
    }

    fn log_event(&self, event: &str, data: &Value) {
        // TODO: Implement actual logging
        println!("Event: {}, Data: {}", event, data);
    }
}

impl EnhancedLogging for Logger {}
