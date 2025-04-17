use actix_web::{web, HttpRequest, HttpResponse, Result};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::path::Path;

// Import the WorkhorseRequest trait
use super::workhorse_request::WorkhorseRequest;

// Define the UploadedFile struct
#[derive(Debug, Clone)]
pub struct UploadedFile {
    pub original_filename: String,
    pub content_type: String,
    pub size: usize,
}

// Define the authorization response
#[derive(Debug, Serialize, Deserialize)]
pub struct AuthorizationResponse {
    pub temp_path: String,
    pub max_size: usize,
    pub allowed_extensions: Vec<String>,
}

pub trait WorkhorseAuthorization {
    fn authorize(&self, req: &HttpRequest) -> Result<HttpResponse>;
    fn file_is_valid(&self, file: &UploadedFile) -> bool;
    fn uploader_class(&self) -> Box<dyn Uploader>;
    fn maximum_size(&self) -> usize;
    fn file_extension_allowlist(&self) -> Vec<String>;
}

pub struct WorkhorseAuthorizationHandler {
    workhorse_request: Box<dyn WorkhorseRequest>,
}

impl WorkhorseAuthorizationHandler {
    pub fn new(workhorse_request: Box<dyn WorkhorseRequest>) -> Self {
        WorkhorseAuthorizationHandler { workhorse_request }
    }

    fn set_workhorse_internal_api_content_type(&self, resp: &mut HttpResponse) {
        resp.content_type("application/json");
    }
}

impl WorkhorseAuthorization for WorkhorseAuthorizationHandler {
    fn authorize(&self, req: &HttpRequest) -> Result<HttpResponse> {
        // Verify the Workhorse API request
        self.workhorse_request.verify_workhorse_api(req)?;

        // Create a response with the appropriate content type
        let mut resp = HttpResponse::Ok();
        self.set_workhorse_internal_api_content_type(&mut resp);

        // Get the uploader class and authorize the upload
        let uploader = self.uploader_class();
        let max_size = self.maximum_size();

        // In a real implementation, this would call the uploader's workhorse_authorize method
        // For now, we'll create a simple response
        let authorized = AuthorizationResponse {
            temp_path: "/tmp/uploads".to_string(),
            max_size,
            allowed_extensions: self.file_extension_allowlist(),
        };

        // Return the JSON response
        Ok(resp.json(authorized))
    }

    fn file_is_valid(&self, file: &UploadedFile) -> bool {
        // Get the file extension
        let path = Path::new(&file.original_filename);
        let extension = path
            .extension()
            .and_then(|ext| ext.to_str())
            .map(|ext| ext.to_lowercase())
            .unwrap_or_default();

        // Check if the extension is in the allowlist
        self.file_extension_allowlist().contains(&extension)
    }

    fn uploader_class(&self) -> Box<dyn Uploader> {
        // This would be implemented by the concrete class
        unimplemented!("uploader_class must be implemented")
    }

    fn maximum_size(&self) -> usize {
        // This would be implemented by the concrete class
        unimplemented!("maximum_size must be implemented")
    }

    fn file_extension_allowlist(&self) -> Vec<String> {
        // Default implementation returns a standard allowlist
        vec!["gz".to_string(), "tar".to_string(), "zip".to_string()]
    }
}

// Define the Uploader trait
pub trait Uploader: Send + Sync {
    fn workhorse_authorize(
        &self,
        has_length: bool,
        maximum_size: usize,
    ) -> Result<AuthorizationResponse, Box<dyn std::error::Error>>;
}
