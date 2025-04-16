pub mod releases;

use actix_web::{web, HttpResponse, Responder};
use serde::{Deserialize, Serialize};

use crate::controllers::activity_pub::ApplicationController as BaseApplicationController;

/// Base controller for all ActivityPub project controllers
pub struct ApplicationController {
    /// The project being accessed
    pub project: Option<Project>,
}

impl ApplicationController {
    /// Create a new project application controller
    pub fn new() -> Self {
        Self { project: None }
    }

    /// Get the project from the request parameters
    pub fn project(&mut self, params: &web::Path<(String, String)>) -> Result<(), impl Responder> {
        let (namespace_id, project_id) = params.into_inner();
        
        // TODO: Implement project finding logic
        // This is a placeholder implementation
        self.project = Some(Project {
            id: 1,
            name: format!("{}/{}", namespace_id, project_id),
            is_public: true,
            is_pending_delete: false,
        });
        
        Ok(())
    }

    /// Ensure the project feature flag is enabled
    pub fn ensure_project_feature_flag(&self) -> Result<(), impl Responder> {
        // TODO: Implement proper feature flag checking
        // This is a placeholder implementation
        Ok(())
    }
}

/// Represents a project
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Project {
    /// The project ID
    pub id: i64,
    /// The project name
    pub name: String,
    /// Whether the project is public
    pub is_public: bool,
    /// Whether the project is pending deletion
    pub is_pending_delete: bool,
} 