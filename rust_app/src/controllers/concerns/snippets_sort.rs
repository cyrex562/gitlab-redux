use crate::models::snippet::Snippet;
use actix_web::{web, HttpResponse};
use serde::{Deserialize, Serialize};
use std::cmp::Ordering;
use std::collections::HashMap;

/// Module for handling snippet sorting
pub trait SnippetsSort {
    /// Get the sort field
    fn sort_field(&self) -> Option<String>;

    /// Get the sort direction
    fn sort_direction(&self) -> Option<String>;

    /// Get the default sort field
    fn default_sort_field(&self) -> String {
        "created_at".to_string()
    }

    /// Get the default sort direction
    fn default_sort_direction(&self) -> String {
        "desc".to_string()
    }

    /// Get the valid sort fields
    fn valid_sort_fields(&self) -> Vec<String> {
        vec![
            "created_at".to_string(),
            "updated_at".to_string(),
            "title".to_string(),
            "file_name".to_string(),
            "author".to_string(),
        ]
    }

    /// Get the valid sort directions
    fn valid_sort_directions(&self) -> Vec<String> {
        vec!["asc".to_string(), "desc".to_string()]
    }

    /// Validate sort parameters
    fn validate_sort_params(&self) -> Result<(), HttpResponse> {
        if let Some(field) = &self.sort_field() {
            if !self.valid_sort_fields().contains(field) {
                return Err(HttpResponse::BadRequest().json(serde_json::json!({
                    "error": format!("Invalid sort field: {}", field)
                })));
            }
        }

        if let Some(direction) = &self.sort_direction() {
            if !self.valid_sort_directions().contains(direction) {
                return Err(HttpResponse::BadRequest().json(serde_json::json!({
                    "error": format!("Invalid sort direction: {}", direction)
                })));
            }
        }

        Ok(())
    }

    /// Get sort parameters
    fn get_sort_params(&self) -> HashMap<String, String> {
        let mut params = HashMap::new();

        params.insert(
            "field".to_string(),
            self.sort_field()
                .unwrap_or_else(|| self.default_sort_field()),
        );

        params.insert(
            "direction".to_string(),
            self.sort_direction()
                .unwrap_or_else(|| self.default_sort_direction()),
        );

        params
    }

    /// Apply sorting to snippets
    async fn apply_sort(
        &self,
        snippets: Vec<HashMap<String, String>>,
    ) -> Result<Vec<HashMap<String, String>>, HttpResponse> {
        self.validate_sort_params()?;

        let field = self
            .sort_field()
            .unwrap_or_else(|| self.default_sort_field());
        let direction = self
            .sort_direction()
            .unwrap_or_else(|| self.default_sort_direction());

        let mut sorted_snippets = snippets;

        sorted_snippets.sort_by(|a, b| {
            let a_val = a.get(&field).unwrap_or(&String::new());
            let b_val = b.get(&field).unwrap_or(&String::new());

            match direction.as_str() {
                "asc" => a_val.cmp(b_val),
                "desc" => b_val.cmp(a_val),
                _ => a_val.cmp(b_val),
            }
        });

        Ok(sorted_snippets)
    }

    /// Sort snippets by the specified criteria
    fn sort_snippets(&self, snippets: &mut Vec<Snippet>) {
        let sort_by = self
            .sort_field()
            .unwrap_or_else(|| "created_at".to_string());
        let direction = self.sort_direction().unwrap_or_else(|| "desc".to_string());

        snippets.sort_by(|a, b| {
            let comparison = match sort_by.as_str() {
                "title" => a.title.cmp(&b.title),
                "created_at" => a.created_at.cmp(&b.created_at),
                "updated_at" => a.updated_at.cmp(&b.updated_at),
                "author" => a.author_name.cmp(&b.author_name),
                "language" => a.language.cmp(&b.language),
                "file_name" => a.file_name.cmp(&b.file_name),
                _ => a.created_at.cmp(&b.created_at),
            };

            match direction.as_str() {
                "asc" => comparison,
                "desc" => comparison.reverse(),
                _ => comparison,
            }
        });
    }

    /// Get the available sort options
    fn available_sort_options(&self) -> Vec<SortOption> {
        vec![
            SortOption {
                value: "created_at".to_string(),
                label: "Created date".to_string(),
            },
            SortOption {
                value: "updated_at".to_string(),
                label: "Last updated".to_string(),
            },
            SortOption {
                value: "title".to_string(),
                label: "Title".to_string(),
            },
            SortOption {
                value: "author".to_string(),
                label: "Author".to_string(),
            },
            SortOption {
                value: "language".to_string(),
                label: "Language".to_string(),
            },
            SortOption {
                value: "file_name".to_string(),
                label: "File name".to_string(),
            },
        ]
    }

    /// Get the current sort option
    fn current_sort_option(&self) -> SortOption {
        let sort_by = self
            .sort_field()
            .unwrap_or_else(|| "created_at".to_string());
        let direction = self.sort_direction().unwrap_or_else(|| "desc".to_string());

        SortOption {
            value: format!("{}_{}", sort_by, direction),
            label: self
                .available_sort_options()
                .into_iter()
                .find(|option| option.value == sort_by)
                .map(|option| option.label)
                .unwrap_or_else(|| "Created date".to_string()),
        }
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SortOption {
    pub value: String,
    pub label: String,
}
