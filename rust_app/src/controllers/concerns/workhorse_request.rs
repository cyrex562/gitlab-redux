// Ported from: orig_app/app/controllers/concerns/workhorse_request.rb
// This module provides a trait and handler to verify GitLab Workhorse API requests.
// It checks for the 'Gitlab-Workhorse-Api-Request' header, as in the Ruby concern.
use actix_web::{web, HttpRequest, HttpResponse, Result};
use std::collections::HashMap;

pub trait WorkhorseRequest {
    fn verify_workhorse_api(&self, req: &HttpRequest) -> Result<()>;
}

pub struct WorkhorseRequestHandler;

impl WorkhorseRequestHandler {
    pub fn new() -> Self {
        WorkhorseRequestHandler
    }
}

impl WorkhorseRequest for WorkhorseRequestHandler {
    fn verify_workhorse_api(&self, req: &HttpRequest) -> Result<()> {
        // Check for the presence of the Workhorse API header (matching Ruby concern)
        if !req.headers().contains_key("Gitlab-Workhorse-Api-Request") {
            return Err(actix_web::error::ErrorForbidden(
                "Invalid Workhorse API request",
            ));
        }
        // TODO: Implement JWT verification as in Gitlab::Workhorse.verify_api_request!
        Ok(())
    }
}
