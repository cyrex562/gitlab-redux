use actix_web::{web, HttpResponse};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

/// Module for handling issuable collections actions
pub trait IssuableCollectionsAction {
    /// Get the collection type (issues, merge_requests, etc.)
    fn collection_type(&self) -> String;

    /// Get the collection scope
    fn collection_scope(&self) -> String {
        "all".to_string()
    }

    /// Get the collection state
    fn collection_state(&self) -> String {
        "opened".to_string()
    }

    /// Get the collection sort
    fn collection_sort(&self) -> String {
        "created_desc".to_string()
    }

    /// Get the collection filter params
    fn collection_filter_params(&self) -> HashMap<String, String> {
        HashMap::new()
    }

    /// Get the collection search params
    fn collection_search_params(&self) -> HashMap<String, String> {
        HashMap::new()
    }

    /// Get the collection page
    fn collection_page(&self) -> i32 {
        1
    }

    /// Get the collection per page
    fn collection_per_page(&self) -> i32 {
        20
    }

    /// Get the collection labels
    fn collection_labels(&self) -> Vec<String> {
        Vec::new()
    }

    /// Get the collection milestone
    fn collection_milestone(&self) -> Option<String> {
        None
    }

    /// Get the collection assignee
    fn collection_assignee(&self) -> Option<i32> {
        None
    }

    /// Get the collection author
    fn collection_author(&self) -> Option<i32> {
        None
    }

    /// Get the collection search
    fn collection_search(&self) -> Option<String> {
        None
    }

    /// Get the collection in
    fn collection_in(&self) -> Option<String> {
        None
    }

    /// Get the collection created after
    fn collection_created_after(&self) -> Option<String> {
        None
    }

    /// Get the collection created before
    fn collection_created_before(&self) -> Option<String> {
        None
    }

    /// Get the collection updated after
    fn collection_updated_after(&self) -> Option<String> {
        None
    }

    /// Get the collection updated before
    fn collection_updated_before(&self) -> Option<String> {
        None
    }

    /// Get the collection scope
    fn get_collection_scope(&self) -> Result<String, HttpResponse> {
        let scope = self.collection_scope();
        match scope.as_str() {
            "all" | "assigned_to_me" | "created_by_me" | "mentioned_to_me" => Ok(scope),
            _ => Err(HttpResponse::BadRequest().json(serde_json::json!({
                "error": format!("Invalid scope: {}", scope)
            })))
        }
    }

    /// Get the collection state
    fn get_collection_state(&self) -> Result<String, HttpResponse> {
        let state = self.collection_state();
        match state.as_str() {
            "opened" | "closed" | "all" => Ok(state),
            _ => Err(HttpResponse::BadRequest().json(serde_json::json!({
                "error": format!("Invalid state: {}", state)
            })))
        }
    }

    /// Get the collection sort
    fn get_collection_sort(&self) -> Result<String, HttpResponse> {
        let sort = self.collection_sort();
        match sort.as_str() {
            "created_desc" | "created_asc" | "updated_desc" | "updated_asc" => Ok(sort),
            _ => Err(HttpResponse::BadRequest().json(serde_json::json!({
                "error": format!("Invalid sort: {}", sort)
            })))
        }
    }

    /// Get the collection filter params
    fn get_collection_filter_params(&self) -> HashMap<String, String> {
        let mut params = self.collection_filter_params();
        
        if let Some(labels) = self.collection_labels().first() {
            params.insert("labels".to_string(), labels.clone());
        }
        
        if let Some(milestone) = self.collection_milestone() {
            params.insert("milestone".to_string(), milestone);
        }
        
        if let Some(assignee) = self.collection_assignee() {
            params.insert("assignee_id".to_string(), assignee.to_string());
        }
        
        if let Some(author) = self.collection_author() {
            params.insert("author_id".to_string(), author.to_string());
        }
        
        if let Some(search) = self.collection_search() {
            params.insert("search".to_string(), search);
        }
        
        if let Some(in_field) = self.collection_in() {
            params.insert("in".to_string(), in_field);
        }
        
        if let Some(created_after) = self.collection_created_after() {
            params.insert("created_after".to_string(), created_after);
        }
        
        if let Some(created_before) = self.collection_created_before() {
            params.insert("created_before".to_string(), created_before);
        }
        
        if let Some(updated_after) = self.collection_updated_after() {
            params.insert("updated_after".to_string(), updated_after);
        }
        
        if let Some(updated_before) = self.collection_updated_before() {
            params.insert("updated_before".to_string(), updated_before);
        }
        
        params
    }

    /// Get the collection search params
    fn get_collection_search_params(&self) -> HashMap<String, String> {
        let mut params = self.collection_search_params();
        params.insert("page".to_string(), self.collection_page().to_string());
        params.insert("per_page".to_string(), self.collection_per_page().to_string());
        params
    }

    /// Get the collection params
    fn get_collection_params(&self) -> Result<HashMap<String, String>, HttpResponse> {
        let mut params = HashMap::new();
        
        params.insert("scope".to_string(), self.get_collection_scope()?);
        params.insert("state".to_string(), self.get_collection_state()?);
        params.insert("sort".to_string(), self.get_collection_sort()?);
        
        params.extend(self.get_collection_filter_params());
        params.extend(self.get_collection_search_params());
        
        Ok(params)
    }
} 