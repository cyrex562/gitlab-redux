use actix_web::{web, HttpRequest, HttpResponse};
use serde::{Deserialize, Serialize};
use std::sync::Arc;

/// This trait provides functionality for decorating content with Sourcegraph integration
pub trait SourcegraphDecorator {
    /// Get the Sourcegraph URL for the current request
    fn sourcegraph_url(&self, req: &HttpRequest) -> Option<String>;
    
    /// Check if Sourcegraph integration is enabled
    fn sourcegraph_enabled?(&self) -> bool;
    
    /// Get the Sourcegraph project URL
    fn sourcegraph_project_url(&self) -> Option<String>;
}

#[derive(Debug, Serialize, Deserialize)]
pub struct SourcegraphDecoratorHandler {
    enabled: bool,
    base_url: String,
    project_url: Option<String>,
}

impl SourcegraphDecoratorHandler {
    pub fn new(enabled: bool, base_url: String, project_url: Option<String>) -> Self {
        SourcegraphDecoratorHandler {
            enabled,
            base_url,
            project_url,
        }
    }
}

impl SourcegraphDecorator for SourcegraphDecoratorHandler {
    fn sourcegraph_url(&self, req: &HttpRequest) -> Option<String> {
        if !self.sourcegraph_enabled?() {
            return None;
        }
        
        // Get the current path and project URL
        let path = req.path();
        let project_url = self.sourcegraph_project_url()?;
        
        // Construct the Sourcegraph URL
        Some(format!("{}/{}", self.base_url, project_url))
    }
    
    fn sourcegraph_enabled?(&self) -> bool {
        self.enabled
    }
    
    fn sourcegraph_project_url(&self) -> Option<String> {
        self.project_url.clone()
    }
}

// Helper function to handle the Result type
pub fn sourcegraph_enabled(result: Result<bool, HttpResponse>) -> bool {
    result.unwrap_or(false)
} 