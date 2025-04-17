use actix_web::{web, HttpRequest, HttpResponse, Responder};
use serde::{Deserialize, Serialize};

use crate::utils::workhorse::Workhorse;

/// Module for handling workhorse requests
pub trait WorkhorseRequest {
    /// Verify the workhorse API request
    fn verify_workhorse_api(&self, req: &HttpRequest) -> Result<(), HttpResponse> {
        match Workhorse::verify_api_request(req.headers()) {
            Ok(_) => Ok(()),
            Err(_) => Err(HttpResponse::Unauthorized().json("Unauthorized workhorse request")),
        }
    }

    /// Set the workhorse internal API content type
    fn set_workhorse_internal_api_content_type(&self) {
        // This would typically set a response header
        // In Actix-web, this would be handled in the response builder
    }
}
