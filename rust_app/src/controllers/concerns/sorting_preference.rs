// Ported from: orig_app/app/controllers/concerns/sorting_preference.rb
// Ported on: 2025-04-29

use actix_web::{web, HttpRequest};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

/// This trait provides functionality for handling sorting preferences in controllers
pub trait SortingPreference {
    /// Get the current sorting preference from the request
    fn sorting_preference(&self, req: &HttpRequest) -> Option<String>;
    /// Set the sorting preference in the session
    fn set_sorting_preference(&self, req: &HttpRequest, preference: &str);
    /// Set the sort order based on user preference, cookies, or params
    fn set_sort_order(&self, req: &HttpRequest, field: &str, default_order: &str) -> String;
}

#[derive(Debug, Serialize, Deserialize)]
pub struct SortingPreferenceHandler {
    pub default_sort: String,
    pub allowed_sorts: Vec<String>,
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

    /// Mimics the Ruby remember_sorting_key logic
    pub fn remember_sorting_key(&self, field: &str) -> String {
        let mut parts: Vec<&str> = field.split('_').collect();
        if parts.len() < 2 {
            return String::new();
        }
        parts.pop(); // remove last part (usually 'sort')
        let base: String = parts
            .iter()
            .map(|s| {
                if s.ends_with("ies") {
                    s.trim_end_matches("ies").to_owned() + "y"
                } else if s.ends_with("es") {
                    s.trim_end_matches("es").to_owned()
                } else if s.ends_with('s') {
                    s.trim_end_matches('s').to_owned()
                } else {
                    s.to_string()
                }
            })
            .collect();
        format!("{}_sort", base)
    }

    /// Mimics the Ruby update_cookie_value logic
    pub fn update_cookie_value(&self, value: &str) -> String {
        match value {
            "id_asc" => "created_asc".to_string(),
            "id_desc" => "created_desc".to_string(),
            "downvotes_asc" | "downvotes_desc" => "popularity".to_string(),
            _ => value.to_string(),
        }
    }

    /// Mimics the Ruby valid_sort_order? logic
    pub fn valid_sort_order(&self, sort_order: &str) -> bool {
        if sort_order.is_empty() {
            return false;
        }
        // Add custom logic for weight/merged_at if needed
        self.validate_sort(sort_order)
    }

    /// Set sort order from user preference (stubbed, as user/session logic is app-specific)
    pub fn set_sort_order_from_user_preference(
        &self,
        _req: &HttpRequest,
        _field: &str,
    ) -> Option<String> {
        // In a real app, fetch from user session or DB
        None
    }

    /// Set sort order from cookie (stubbed, as cookie logic is app-specific)
    pub fn set_sort_order_from_cookie(&self, _req: &HttpRequest, _field: &str) -> Option<String> {
        // In a real app, fetch from cookies
        None
    }
}

impl SortingPreference for SortingPreferenceHandler {
    fn sorting_preference(&self, req: &HttpRequest) -> Option<String> {
        // Get sort parameter from query string
        let query = req.query_string();
        let params: HashMap<_, _> = url::form_urlencoded::parse(query.as_bytes()).collect();
        let sort = params.get("sort").map(|s| s.to_string());
        sort.filter(|s| self.validate_sort(s))
            .or(Some(self.default_sort.clone()))
    }

    fn set_sorting_preference(&self, req: &HttpRequest, preference: &str) {
        if self.validate_sort(preference) {
            // In a real app, set in session or cookies
        }
    }

    fn set_sort_order(&self, req: &HttpRequest, field: &str, default_order: &str) -> String {
        let sort_order = self
            .set_sort_order_from_user_preference(req, field)
            .or_else(|| self.set_sort_order_from_cookie(req, field))
            .or_else(|| {
                let query = req.query_string();
                let params: HashMap<_, _> = url::form_urlencoded::parse(query.as_bytes()).collect();
                params.get("sort").map(|s| s.to_string())
            })
            .unwrap_or_else(|| default_order.to_string());
        if !self.valid_sort_order(&sort_order) {
            default_order.to_string()
        } else {
            sort_order
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
