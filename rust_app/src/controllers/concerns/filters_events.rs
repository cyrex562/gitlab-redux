use actix_web::web;
use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct EventFilter {
    pub filter: String,
}

impl EventFilter {
    pub fn new(filter: Option<String>) -> Self {
        EventFilter {
            filter: filter.unwrap_or_else(|| "all".to_string()),
        }
    }
}

pub trait FiltersEvents {
    fn event_filter(&self, req: &web::HttpRequest) -> EventFilter {
        let active_filter = req
            .query_string()
            .split('&')
            .find(|param| param.starts_with("event_filter="))
            .and_then(|param| param.split('=').nth(1))
            .or_else(|| self.get_cookie("event_filter"))
            .map(|s| s.to_string());

        EventFilter::new(active_filter)
    }

    // Required methods to be implemented by concrete types
    fn get_cookie(&self, name: &str) -> Option<String>;
}
