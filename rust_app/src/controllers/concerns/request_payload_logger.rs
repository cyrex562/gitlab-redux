use actix_web::{web, HttpRequest};
use std::sync::Arc;

use crate::{
    models::user::User,
    services::{
        cloudflare::CloudflareService, correlation::CorrelationService, logging::LoggingService,
    },
    utils::{context::ApplicationContext, error::AppError},
};

/// Module for handling request payload logging
pub trait RequestPayloadLogger {
    /// Append information to payload
    fn append_info_to_payload(&self, payload: &mut serde_json::Value) -> Result<(), AppError> {
        // Get request
        let request = self.request();

        // Add user agent
        if let Some(user_agent) = request.headers().get("user-agent") {
            if let Ok(user_agent) = user_agent.to_str() {
                payload["ua"] = serde_json::Value::String(user_agent.to_string());
            }
        }

        // Add remote IP
        if let Some(remote_ip) = request.connection_info().peer_addr() {
            payload["remote_ip"] = serde_json::Value::String(remote_ip.to_string());
        }

        // Add correlation ID
        let correlation_id = CorrelationService::current_id();
        payload[CorrelationService::LOG_KEY] = serde_json::Value::String(correlation_id);

        // Add application context
        let context = ApplicationContext::current();
        payload["metadata"] = serde_json::to_value(context)?;

        // Add urgency if defined
        if let Some(urgency) = self.urgency() {
            payload["request_urgency"] = serde_json::Value::String(urgency.name());
            payload["target_duration_s"] = serde_json::Value::Number(urgency.duration().into());
        }

        // Add user information
        if let Some(user) = self.auth_user() {
            payload["user_id"] = serde_json::Value::Number(user.id().into());
            payload["username"] = serde_json::Value::String(user.username().to_string());
        }

        // Add queue duration
        if let Some(queue_duration) = request
            .headers()
            .get("x-gitlab-rails-queue-duration")
            .and_then(|v| v.to_str().ok())
            .and_then(|v| v.parse::<f64>().ok())
        {
            payload["queue_duration_s"] = serde_json::Value::Number(queue_duration.into());
        }

        // Store Cloudflare headers
        CloudflareService::store_headers(payload, request);

        Ok(())
    }

    // Required trait methods that need to be implemented by the controller
    fn request(&self) -> &HttpRequest;
    fn auth_user(&self) -> Option<&User>;
    fn urgency(&self) -> Option<&dyn RequestUrgency>;
}

/// Trait for request urgency
pub trait RequestUrgency {
    fn name(&self) -> String;
    fn duration(&self) -> f64;
}
