use crate::{
    models::{Namespace, Project, User},
    utils::{cookies::CookiesHelper, redis::HLLRedisCounter, tracking::Tracking},
};
use actix_web::{dev::ServiceRequest, error::Error};
use std::sync::Arc;
use uuid::Uuid;

pub trait ProductAnalyticsTracking {
    fn track_event(
        &self,
        name: &str,
        action: Option<&str>,
        label: Option<&str>,
        destinations: &[String],
        block: Option<Box<dyn Fn(&ServiceRequest) -> Option<String> + Send + Sync>>,
    ) -> Result<(), Error>;

    fn track_internal_event(
        &self,
        name: &str,
        user: Option<&User>,
        project: Option<&Project>,
        namespace: Option<&Namespace>,
        event_args: serde_json::Value,
    ) -> Result<(), Error>;
}

pub struct ProductAnalyticsTrackingImpl {
    tracking: Arc<Tracking>,
    cookies_helper: Box<dyn CookiesHelper>,
}

impl ProductAnalyticsTrackingImpl {
    pub fn new(tracking: Arc<Tracking>, cookies_helper: Box<dyn CookiesHelper>) -> Self {
        Self {
            tracking,
            cookies_helper,
        }
    }

    fn route_events_to(
        &self,
        req: &ServiceRequest,
        destinations: &[String],
        name: &str,
        action: Option<&str>,
        label: Option<&str>,
        block: Option<Box<dyn Fn(&ServiceRequest) -> Option<String> + Send + Sync>>,
    ) -> Result<(), Error> {
        if destinations.contains(&"redis_hll".to_string()) {
            self.track_unique_redis_hll_event(req, name, block)?;
        }

        if destinations.contains(&"snowplow".to_string()) {
            let action = action.ok_or_else(|| {
                actix_web::error::ErrorBadRequest("action is required for snowplow")
            })?;
            let label = label.ok_or_else(|| {
                actix_web::error::ErrorBadRequest("label is required for snowplow")
            })?;

            let mut optional_arguments = serde_json::json!({});
            if let Some(namespace) = req.extensions().get::<Namespace>() {
                optional_arguments["namespace"] = serde_json::to_value(namespace)?;
            }
            if let Some(project) = req.extensions().get::<Project>() {
                optional_arguments["project"] = serde_json::to_value(project)?;
            }

            self.tracking.event(
                req.path(),
                action,
                req.extensions().get::<User>().cloned(),
                name,
                label,
                serde_json::json!([{
                    "data_source": "redis_hll",
                    "event": name
                }]),
                optional_arguments,
            )?;
        }

        Ok(())
    }

    fn track_unique_redis_hll_event(
        &self,
        req: &ServiceRequest,
        event_name: &str,
        block: Option<Box<dyn Fn(&ServiceRequest) -> Option<String> + Send + Sync>>,
    ) -> Result<(), Error> {
        let custom_id = block.as_ref().and_then(|b| b(req));
        let unique_id = custom_id.or_else(|| self.visitor_id(req));

        if let Some(id) = unique_id {
            HLLRedisCounter::track_event(event_name, &[id])?;
        }

        Ok(())
    }

    fn visitor_id(&self, req: &ServiceRequest) -> Option<String> {
        if let Some(visitor_id) = self.cookies_helper.get(req, "visitor_id") {
            return Some(visitor_id);
        }

        if let Some(_user) = req.extensions().get::<User>() {
            let uuid = Uuid::new_v4().to_string();
            self.cookies_helper.add(req, "visitor_id", &uuid);
            Some(uuid)
        } else {
            None
        }
    }
}

impl ProductAnalyticsTracking for ProductAnalyticsTrackingImpl {
    fn track_event(
        &self,
        name: &str,
        action: Option<&str>,
        label: Option<&str>,
        destinations: &[String],
        block: Option<Box<dyn Fn(&ServiceRequest) -> Option<String> + Send + Sync>>,
    ) -> Result<(), Error> {
        self.route_events_to(req, destinations, name, action, label, block)
    }

    fn track_internal_event(
        &self,
        name: &str,
        user: Option<&User>,
        project: Option<&Project>,
        namespace: Option<&Namespace>,
        event_args: serde_json::Value,
    ) -> Result<(), Error> {
        self.tracking.track_internal_event(
            name,
            user.cloned(),
            project.cloned(),
            namespace.cloned(),
            event_args,
        )
    }
}
