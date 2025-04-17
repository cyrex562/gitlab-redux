use actix_web::{web, HttpRequest, HttpResponse};
use chrono::{DateTime, Duration, Utc};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use uuid::Uuid;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct TrackingEvent {
    pub name: String,
    pub action: Option<String>,
    pub label: Option<String>,
    pub user_id: Option<i32>,
    pub project_id: Option<i32>,
    pub namespace_id: Option<i32>,
    pub context: Option<Vec<TrackingContext>>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct TrackingContext {
    pub data_source: String,
    pub event: String,
}

pub trait ProductAnalyticsTracking {
    fn track_event(
        &self,
        name: &str,
        action: Option<&str>,
        label: Option<&str>,
        destinations: &[EventDestination],
        custom_id: Option<String>,
    ) -> Result<(), Box<dyn std::error::Error>> {
        if destinations.contains(&EventDestination::RedisHll) {
            self.track_unique_redis_hll_event(name, custom_id)?;
        }

        if destinations.contains(&EventDestination::Snowplow) {
            if action.is_none() {
                return Err("action is required when destination is snowplow".into());
            }
            if label.is_none() {
                return Err("label is required when destination is snowplow".into());
            }

            let event = TrackingEvent {
                name: name.to_string(),
                action: action.map(|s| s.to_string()),
                label: label.map(|s| s.to_string()),
                user_id: self.current_user_id(),
                project_id: self.tracking_project_source(),
                namespace_id: self.tracking_namespace_source(),
                context: Some(vec![TrackingContext {
                    data_source: "redis_hll".to_string(),
                    event: name.to_string(),
                }]),
            };

            self.track_snowplow_event(event)?;
        }

        Ok(())
    }

    fn track_internal_event(
        &self,
        name: &str,
        event_args: HashMap<String, serde_json::Value>,
    ) -> Result<(), Box<dyn std::error::Error>> {
        let mut args = event_args.clone();
        args.insert(
            "user_id".to_string(),
            serde_json::json!(self.current_user_id()),
        );
        args.insert(
            "project_id".to_string(),
            serde_json::json!(self.tracking_project_source()),
        );
        args.insert(
            "namespace_id".to_string(),
            serde_json::json!(self.tracking_namespace_source()),
        );

        self.track_internal_event_impl(name, args)
    }

    fn track_unique_redis_hll_event(
        &self,
        event_name: &str,
        custom_id: Option<String>,
    ) -> Result<(), Box<dyn std::error::Error>> {
        let unique_id = custom_id.unwrap_or_else(|| self.visitor_id().unwrap_or_default());
        if !unique_id.is_empty() {
            self.track_redis_hll_event(event_name, &unique_id)?;
        }
        Ok(())
    }

    fn visitor_id(&self) -> Option<String> {
        if let Some(id) = self.get_cookie("visitor_id") {
            return Some(id);
        }

        if self.current_user_id().is_some() {
            let uuid = Uuid::new_v4().to_string();
            self.set_cookie(
                "visitor_id",
                &uuid,
                Some(Utc::now() + Duration::days(730)), // 24 months
            );
            return Some(uuid);
        }

        None
    }

    // Required methods to be implemented by concrete types
    fn current_user_id(&self) -> Option<i32>;
    fn tracking_project_source(&self) -> Option<i32>;
    fn tracking_namespace_source(&self) -> Option<i32>;
    fn get_cookie(&self, name: &str) -> Option<String>;
    fn set_cookie(&self, name: &str, value: &str, expires: Option<DateTime<Utc>>);
    fn track_redis_hll_event(
        &self,
        event_name: &str,
        value: &str,
    ) -> Result<(), Box<dyn std::error::Error>>;
    fn track_snowplow_event(&self, event: TrackingEvent) -> Result<(), Box<dyn std::error::Error>>;
    fn track_internal_event_impl(
        &self,
        name: &str,
        args: HashMap<String, serde_json::Value>,
    ) -> Result<(), Box<dyn std::error::Error>>;
}

#[derive(Debug, Clone, Copy, PartialEq)]
pub enum EventDestination {
    RedisHll,
    Snowplow,
}
