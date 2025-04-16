use actix_web::{web, HttpResponse};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

/// Module for handling strong pagination parameters
pub trait StrongPaginationParams {
    /// Get the current page number
    fn page(&self) -> Option<i32>;

    /// Get the number of items per page
    fn per_page(&self) -> Option<i32>;

    /// Get the total number of items
    fn total_items(&self) -> i32;

    /// Get the total number of pages
    fn total_pages(&self) -> i32 {
        let per_page = self.per_page().unwrap_or(20);
        if per_page <= 0 {
            return 1;
        }
        (self.total_items() as f64 / per_page as f64).ceil() as i32
    }

    /// Get the current page number with validation
    fn current_page(&self) -> i32 {
        let page = self.page().unwrap_or(1);
        if page < 1 {
            return 1;
        }
        if page > self.total_pages() {
            return self.total_pages();
        }
        page
    }

    /// Get the number of items per page with validation
    fn items_per_page(&self) -> i32 {
        let per_page = self.per_page().unwrap_or(20);
        if per_page < 1 {
            return 20;
        }
        if per_page > 100 {
            return 100;
        }
        per_page
    }

    /// Get the offset for pagination
    fn offset(&self) -> i32 {
        (self.current_page() - 1) * self.items_per_page()
    }

    /// Get pagination metadata
    fn pagination_metadata(&self) -> HashMap<String, i32> {
        let mut metadata = HashMap::new();
        metadata.insert("current_page".to_string(), self.current_page());
        metadata.insert("per_page".to_string(), self.items_per_page());
        metadata.insert("total_items".to_string(), self.total_items());
        metadata.insert("total_pages".to_string(), self.total_pages());
        metadata
    }

    /// Validate pagination parameters
    fn validate_pagination_params(&self) -> Result<(), HttpResponse> {
        if let Some(page) = self.page() {
            if page < 1 {
                return Err(HttpResponse::BadRequest().json(serde_json::json!({
                    "error": "Page number must be greater than 0"
                })));
            }
        }

        if let Some(per_page) = self.per_page() {
            if per_page < 1 {
                return Err(HttpResponse::BadRequest().json(serde_json::json!({
                    "error": "Items per page must be greater than 0"
                })));
            }
            if per_page > 100 {
                return Err(HttpResponse::BadRequest().json(serde_json::json!({
                    "error": "Items per page cannot exceed 100"
                })));
            }
        }

        Ok(())
    }
}
