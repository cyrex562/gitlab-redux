use serde::{Deserialize, Serialize};
use std::collections::HashMap;

use crate::harbor::artifact::Artifact;
use crate::harbor::repository::Repository;
use crate::harbor::tag::Tag;

pub struct Query {
    params: HashMap<String, String>,
    errors: Vec<String>,
}

impl Query {
    pub fn new(params: HashMap<String, String>) -> Self {
        Query {
            params,
            errors: Vec::new(),
        }
    }

    pub fn is_valid(&self) -> bool {
        self.errors.is_empty()
    }

    pub fn errors(&self) -> &Vec<String> {
        &self.errors
    }

    pub fn artifacts(&self) -> Vec<Artifact> {
        // TODO: Implement artifact querying
        Vec::new()
    }

    pub fn repositories(&self) -> Vec<Repository> {
        // TODO: Implement repository querying
        Vec::new()
    }

    pub fn tags(&self) -> Vec<Tag> {
        // TODO: Implement tag querying
        Vec::new()
    }

    pub fn validate(&mut self) {
        // TODO: Implement query validation
        if let Some(repository_id) = self.params.get("repository_id") {
            if repository_id.parse::<i32>().is_err() {
                self.errors.push("Invalid repository_id".to_string());
            }
        }

        if let Some(artifact_id) = self.params.get("artifact_id") {
            if artifact_id.parse::<i32>().is_err() {
                self.errors.push("Invalid artifact_id".to_string());
            }
        }

        if let Some(page) = self.params.get("page") {
            if page.parse::<i32>().is_err() {
                self.errors.push("Invalid page number".to_string());
            }
        }

        if let Some(limit) = self.params.get("limit") {
            if limit.parse::<i32>().is_err() {
                self.errors.push("Invalid limit".to_string());
            }
        }
    }
} 