use actix_web::{web, HttpRequest};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

/// This trait provides functionality for handling sorting preferences in controllers
pub trait SortingPreference {
    /// Get the current sorting preference from the request
    fn sorting_preference(&self, req: &HttpRequest) -> Option<String>;
    
    /// Set the sorting preference in the session
    fn set_sorting_preference(&self, req: &HttpRequest, preference: &str);
}

#[derive(Debug, Serialize, Deserialize)]
pub struct SortingPreferenceHandler {
    default_sort: String,
    allowed_sorts: Vec<String>,
}

impl SortingPreferenceHandler {
    pub fn new(default_sort: String, allowed_sorts: Vec<String>) -> Self {
        SortingPreferenceHandler {
            default_sort,
            allowed_sorts,
        }
    }
    
    pub fn validate_sort(&self, sort: &str) -> bool {
        self.allowed_sorts.contains(&sort.to_string())
    }
}

impl SortingPreference for SortingPreferenceHandler {
    fn sorting_preference(&self, req: &HttpRequest) -> Option<String> {
        // Get sort parameter from query string
        let query = req.query_string();
        let params: HashMap<_, _> = url::form_urlencoded::parse(query.as_bytes())
            .collect();
            
        // Get sort from params or session
        let sort = params.get("sort")
            .map(|s| s.to_string())
            .or_else(|| {
                req.extensions()
                    .get::<Session>()
                    .and_then(|session| session.get("sort"))
            });
            
        // Validate and return sort preference
        sort.filter(|s| self.validate_sort(s))
            .or(Some(self.default_sort.clone()))
    }
    
    fn set_sorting_preference(&self, req: &HttpRequest, preference: &str) {
        if let Some(session) = req.extensions().get::<Session>() {
            if self.validate_sort(preference) {
                session.insert("sort", preference);
            }
        }
    }
}

// This would be implemented in a separate module
pub struct Session {
    data: HashMap<String, String>,
}

impl Session {
    pub fn get(&self, key: &str) -> Option<String> {
        self.data.get(key).cloned()
    }
    
    pub fn insert(&mut self, key: &str, value: &str) {
        self.data.insert(key.to_string(), value.to_string());
    }
} 