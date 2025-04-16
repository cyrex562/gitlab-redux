pub mod projects;

use actix_web::{web, HttpResponse, Responder};
use serde::{Deserialize, Serialize};

/// Base controller for all ActivityPub controllers
pub struct ApplicationController;

impl ApplicationController {
    /// Check if an object can perform an action on a subject
    pub fn can(object: &dyn std::any::Any, action: &str, subject: Option<&dyn std::any::Any>) -> bool {
        // TODO: Implement proper ability checking
        // This is a placeholder implementation
        true
    }

    /// Handle route not found
    pub fn route_not_found() -> impl Responder {
        HttpResponse::NotFound()
    }

    /// Set content type for ActivityPub responses
    pub fn set_content_type(response: &mut HttpResponse) {
        response.headers_mut().insert(
            "Content-Type",
            "application/activity+json".parse().unwrap(),
        );
    }

    /// Ensure the ActivityPub feature flag is enabled
    pub fn ensure_feature_flag() -> Result<(), impl Responder> {
        // TODO: Implement proper feature flag checking
        // This is a placeholder implementation
        Ok(())
    }
} 