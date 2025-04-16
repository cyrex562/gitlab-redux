use axum::http::Request;
use serde_json::json;
use std::collections::HashMap;

pub trait Logging {
    fn notification_payload(&self, _: &str) -> serde_json::Value;
    fn get_request(&self) -> &Request<()>;
    fn get_current_user(&self) -> Option<&crate::models::User>;
}

impl<T: Logging> T {
    pub fn enhanced_notification_payload(&self, event: &str) -> serde_json::Value {
        let mut payload = self.notification_payload(event);

        if let serde_json::Value::Object(ref mut map) = payload {
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
