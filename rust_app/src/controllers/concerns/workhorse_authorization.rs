use actix_web::{web, HttpResponse, Responder};
use serde::{Deserialize, Serialize};
use std::path::Path;

use crate::controllers::concerns::workhorse_request::WorkhorseRequest;
use crate::models::uploaded_file::UploadedFile;
use crate::uploaders::import_export_uploader::ImportExportUploader;

/// Module for handling workhorse authorization
pub trait WorkhorseAuthorization: WorkhorseRequest {
    /// Authorize a workhorse request
    async fn authorize(&self) -> impl Responder {
        self.set_workhorse_internal_api_content_type();

        match self
            .uploader_class()
            .workhorse_authorize(false, self.maximum_size() as i64)
            .await
        {
            Ok(authorized) => HttpResponse::Ok().json(authorized),
            Err(_) => HttpResponse::InternalServerError().json("Error uploading file"),
        }
    }

    /// Check if a file is valid
    fn file_is_valid(&self, file: &UploadedFile) -> bool {
        if let Some(original_filename) = file.original_filename() {
            if let Some(extension) = Path::new(&original_filename).extension() {
                if let Some(extension_str) = extension.to_str() {
                    let extension_lower = extension_str.to_lowercase();
                    return self.file_extension_allowlist().contains(&extension_lower);
                }
            }
        }
        false
    }

    /// Get the uploader class
    fn uploader_class(&self) -> Box<dyn Uploader>;

    /// Get the maximum file size
    fn maximum_size(&self) -> usize {
        // Default implementation, should be overridden
        0
    }

    /// Get the file extension allowlist
    fn file_extension_allowlist(&self) -> Vec<String> {
        ImportExportUploader::extension_allowlist()
    }
}

/// Trait for uploaders
pub trait Uploader {
    /// Authorize a workhorse request
    async fn workhorse_authorize(
        &self,
        has_length: bool,
        maximum_size: i64,
    ) -> Result<serde_json::Value, Box<dyn std::error::Error>>;
}
