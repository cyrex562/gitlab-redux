use actix_web::{web, HttpRequest, HttpResponse};
use lazy_static::lazy_static;
use std::sync::atomic::{AtomicBool, Ordering};
use std::sync::Arc;

lazy_static! {
    static ref PEEK_ENABLED: Arc<AtomicBool> = Arc::new(AtomicBool::new(false));
}

pub trait WithPerformanceBar {
    fn set_peek_enabled_for_current_request(&self, req: &HttpRequest);
    fn peek_enabled(&self) -> bool;
    fn cookie_or_default_value(&self, req: &HttpRequest) -> bool;
}

pub struct PerformanceBarHandler;

impl PerformanceBarHandler {
    pub fn new() -> Self {
        PerformanceBarHandler
    }
}

impl WithPerformanceBar for PerformanceBarHandler {
    fn set_peek_enabled_for_current_request(&self, req: &HttpRequest) {
        let enabled = self.cookie_or_default_value(req);
        PEEK_ENABLED.store(enabled, Ordering::SeqCst);
    }

    fn peek_enabled(&self) -> bool {
        PEEK_ENABLED.load(Ordering::SeqCst)
    }

    fn cookie_or_default_value(&self, req: &HttpRequest) -> bool {
        let is_development =
            std::env::var("RUST_ENV").unwrap_or_else(|_| "production".to_string()) == "development";

        // Get the performance bar cookie
        let cookie_enabled = req
            .cookies()
            .get("perf_bar_enabled")
            .and_then(|c| c.value().parse::<bool>().ok())
            .unwrap_or(false);

        // Set cookie to true in development if not set
        if is_development && !cookie_enabled {
            // Note: In a real implementation, you would set the cookie here
            // This is simplified for the example
            return true;
        }

        // Check if the user is allowed to see the performance bar
        let user_allowed = self.is_user_allowed(req);

        cookie_enabled && user_allowed
    }
}

impl PerformanceBarHandler {
    fn is_user_allowed(&self, req: &HttpRequest) -> bool {
        // This would be implemented based on your user permission system
        // For now, we'll return true as a placeholder
        true
    }
}
