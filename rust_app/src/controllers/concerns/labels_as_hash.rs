use actix_web::{web, HttpResponse};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

/// Module for handling labels as hash
pub trait LabelsAsHash {
    /// Get the labels
    fn labels(&self) -> Vec<String>;

    /// Get the project ID
    fn project_id(&self) -> Option<i32> {
        None
    }

    /// Get the group ID
    fn group_id(&self) -> Option<i32> {
        None
    }

    /// Get the label type
    fn label_type(&self) -> String {
        "project".to_string()
    }

    /// Get the label scope
    fn label_scope(&self) -> String {
        "all".to_string()
    }

    /// Get the label search
    fn label_search(&self) -> Option<String> {
        None
    }

    /// Get the label sort
    fn label_sort(&self) -> String {
        "name_asc".to_string()
    }

    /// Get the label page
    fn label_page(&self) -> i32 {
        1
    }

    /// Get the label per page
    fn label_per_page(&self) -> i32 {
        20
    }

    /// Convert labels to hash
    fn labels_as_hash(&self) -> Result<HashMap<String, serde_json::Value>, HttpResponse> {
        let mut hash = HashMap::new();

        for label in self.labels() {
            let mut label_data = HashMap::new();
            label_data.insert("name".to_string(), serde_json::json!(label));

            if let Some(project_id) = self.project_id() {
                label_data.insert("project_id".to_string(), serde_json::json!(project_id));
            }

            if let Some(group_id) = self.group_id() {
                label_data.insert("group_id".to_string(), serde_json::json!(group_id));
            }

            label_data.insert("type".to_string(), serde_json::json!(self.label_type()));
            label_data.insert("scope".to_string(), serde_json::json!(self.label_scope()));

            if let Some(search) = self.label_search() {
                label_data.insert("search".to_string(), serde_json::json!(search));
            }

            label_data.insert("sort".to_string(), serde_json::json!(self.label_sort()));
            label_data.insert("page".to_string(), serde_json::json!(self.label_page()));
            label_data.insert(
                "per_page".to_string(),
                serde_json::json!(self.label_per_page()),
            );

            hash.insert(
                label,
                serde_json::Value::Object(serde_json::Map::from_iter(
                    label_data.into_iter().map(|(k, v)| (k, v)),
                )),
            );
        }

        Ok(hash)
    }

    /// Get label metadata
    fn get_label_metadata(&self) -> Result<HashMap<String, serde_json::Value>, HttpResponse> {
        let mut metadata = HashMap::new();

        metadata.insert("type".to_string(), serde_json::json!(self.label_type()));
        metadata.insert("scope".to_string(), serde_json::json!(self.label_scope()));
        metadata.insert("sort".to_string(), serde_json::json!(self.label_sort()));
        metadata.insert("page".to_string(), serde_json::json!(self.label_page()));
        metadata.insert(
            "per_page".to_string(),
            serde_json::json!(self.label_per_page()),
        );

        if let Some(project_id) = self.project_id() {
            metadata.insert("project_id".to_string(), serde_json::json!(project_id));
        }

        if let Some(group_id) = self.group_id() {
            metadata.insert("group_id".to_string(), serde_json::json!(group_id));
        }

        if let Some(search) = self.label_search() {
            metadata.insert("search".to_string(), serde_json::json!(search));
        }

        Ok(metadata)
    }

    /// Validate label data
    fn validate_label_data(&self) -> Result<(), HttpResponse> {
        let label_type = self.label_type();
        if !["project", "group"].contains(&label_type.as_str()) {
            return Err(HttpResponse::BadRequest().json(serde_json::json!({
                "error": format!("Invalid label type: {}", label_type)
            })));
        }

        let label_scope = self.label_scope();
        if !["all", "title", "description"].contains(&label_scope.as_str()) {
            return Err(HttpResponse::BadRequest().json(serde_json::json!({
                "error": format!("Invalid label scope: {}", label_scope)
            })));
        }

        let label_sort = self.label_sort();
        if !["name_asc", "name_desc", "created_asc", "created_desc"].contains(&label_sort.as_str())
        {
            return Err(HttpResponse::BadRequest().json(serde_json::json!({
                "error": format!("Invalid label sort: {}", label_sort)
            })));
        }

        if self.labels().is_empty() {
            return Err(HttpResponse::BadRequest().json(serde_json::json!({
                "error": "Labels cannot be empty"
            })));
        }

        Ok(())
    }
}
