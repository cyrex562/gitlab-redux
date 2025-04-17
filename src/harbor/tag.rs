use serde::{Deserialize, Serialize};
use std::error::Error;

use crate::harbor::common::{HarborContainer, HarborIntegration, HarborQuery, PaginatedResult};

/// Query parameters for Harbor tags
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct TagQueryParams {
    pub repository_id: Option<String>,
    pub artifact_id: Option<String>,
    pub sort: Option<String>,
    pub page: Option<i32>,
    pub limit: Option<i32>,
}

/// Harbor tag representation
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct HarborTag {
    pub id: String,
    pub repository_id: String,
    pub artifact_id: String,
    pub name: String,
    pub size: Option<i64>,
    pub created_at: Option<String>,
    pub updated_at: Option<String>,
}

/// Trait for Harbor tag operations
pub trait HarborTag {
    /// Get query parameters for tags
    fn query_params(&self) -> TagQueryParams;

    /// Get tags based on query parameters
    fn tags(&self) -> Result<PaginatedResult<HarborTag>, Box<dyn Error>>;

    /// Get the container for Harbor integration
    fn container(&self) -> Result<Box<dyn HarborContainer>, Box<dyn Error>>;
}

/// Default implementation for HarborTag
pub struct DefaultHarborTag {
    query_params: TagQueryParams,
}

impl DefaultHarborTag {
    pub fn new(query_params: TagQueryParams) -> Self {
        Self { query_params }
    }
}

impl HarborTag for DefaultHarborTag {
    fn query_params(&self) -> TagQueryParams {
        self.query_params.clone()
    }

    fn tags(&self) -> Result<PaginatedResult<HarborTag>, Box<dyn Error>> {
        // This would be implemented to actually fetch tags from Harbor
        Ok(PaginatedResult {
            items: Vec::new(),
            total: 0,
            page: self.query_params.page.unwrap_or(1),
            limit: self.query_params.limit.unwrap_or(20),
            total_pages: 0,
        })
    }

    fn container(&self) -> Result<Box<dyn HarborContainer>, Box<dyn Error>> {
        Err("Not implemented".into())
    }
}
