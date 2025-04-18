use crate::{
    config::Settings,
    models::User,
    utils::{
        application_context::ApplicationContext, cloudflare::CloudflareHelper,
        correlation::CorrelationId,
    },
};
use actix_web::{dev::ServiceRequest, web::Data};
use serde_json::Value;
use std::collections::HashMap;

pub trait RequestPayloadLogger {
    fn append_info_to_payload(&self, payload: &mut HashMap<String, Value>, req: &ServiceRequest);
}

pub struct RequestPayloadLoggerImpl {
    settings: Data<Settings>,
}

impl RequestPayloadLoggerImpl {
    pub fn new(settings: Data<Settings>) -> Self {
        Self { settings }
    }
}

impl RequestPayloadLogger for RequestPayloadLoggerImpl {
    fn append_info_to_payload(&self, payload: &mut HashMap<String, Value>, req: &ServiceRequest) {
        // Add user agent
        if let Some(user_agent) = req.headers().get("user-agent") {
            if let Ok(ua) = user_agent.to_str() {
                payload.insert("ua".to_string(), Value::String(ua.to_string()));
            }
        }

        // Add remote IP
        if let Some(remote_addr) = req.connection_info().peer_addr() {
            payload.insert(
                "remote_ip".to_string(),
                Value::String(remote_addr.to_string()),
            );
        }

        // Add correlation ID
        payload.insert(
            "correlation_id".to_string(),
            Value::String(CorrelationId::current().to_string()),
        );

        // Add application context
        if let Some(ctx) = ApplicationContext::current() {
            payload.insert(
                "metadata".to_string(),
                serde_json::to_value(ctx).unwrap_or(Value::Null),
            );
        }

        // Add user information if available
        if let Some(user) = req.extensions().get::<User>() {
            payload.insert("user_id".to_string(), Value::Number(user.id.into()));
            payload.insert("username".to_string(), Value::String(user.username.clone()));
        }

        // Add queue duration if available
        if let Some(duration) = req.extensions().get::<f64>() {
            payload.insert(
                "queue_duration_s".to_string(),
                Value::Number((*duration).into()),
            );
        }

        // Add Cloudflare headers
        CloudflareHelper::store_headers(payload, req);
    }
}
