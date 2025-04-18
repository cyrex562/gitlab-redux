use serde::{Deserialize, Serialize};
use std::error::Error;

/// Common query parameters for Harbor
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct HarborQueryParams {
    pub search: Option<String>,
    pub sort: Option<String>,
    pub page: Option<i32>,
    pub limit: Option<i32>,
}

/// Query result with pagination
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PaginatedResult<T> {
    pub items: Vec<T>,
    pub total: i32,
    pub page: i32,
    pub limit: i32,
    pub total_pages: i32,
}

/// Trait for Harbor query
pub trait HarborQuery {
    /// Check if the query is valid
    fn is_valid(&self) -> bool;

    /// Get validation errors
    fn errors(&self) -> Vec<String>;
}

/// Trait for Harbor container
pub trait HarborContainer {
    /// Get the Harbor integration
    fn harbor_integration(&self) -> &dyn HarborIntegration;
}

/// Trait for Harbor integration
pub trait HarborIntegration {
    /// Get the URL of the Harbor integration
    fn url(&self) -> &str;

    /// Get the project name of the Harbor integration
    fn project_name(&self) -> &str;
}

/// Default implementation for HarborQuery
pub struct DefaultHarborQuery<T> {
    params: T,
    errors: Vec<String>,
}

impl<T> DefaultHarborQuery<T> {
    pub fn new(params: T) -> Self {
        Self {
            params,
            errors: Vec::new(),
        }
    }

    pub fn add_error(&mut self, error: String) {
        self.errors.push(error);
    }
}

impl<T> HarborQuery for DefaultHarborQuery<T> {
    fn is_valid(&self) -> bool {
        self.errors.is_empty()
    }

    fn errors(&self) -> Vec<String> {
        self.errors.clone()
    }
}
