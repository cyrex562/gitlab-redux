use crate::models::user::User;
use crate::services::upload_service::UploadService;
use actix_multipart::Multipart;
use actix_web::{web, HttpResponse};
use futures::{StreamExt, TryStreamExt};
use std::io::Write;
use std::path::Path;
use std::sync::Arc;
use tokio::sync::RwLock;

use crate::controllers::concerns::send_file_upload::SendFileUpload;
use crate::models::upload::Upload;
use crate::utils::path_traversal::PathTraversal;
use crate::utils::strong_memoize::StrongMemoize;

const UPLOAD_MOUNTS: &[&str] = &[
    "avatar",
    "attachment",
    "file",
    "logo",
    "pwa_icon",
    "header_logo",
    "favicon",
    "screenshot",
];
const ID_BASED_UPLOAD_PATH_VERSION: i32 = 2;

/// Module for handling file uploads
pub trait UploadsActions: SendFileUpload + StrongMemoize {
    /// Create a new file upload
    async fn create(&self, user: Option<&User>, mut payload: Multipart) -> HttpResponse {
        let uploader = UploadService::new(
            self.model(),
            payload,
            self.uploader_class(),
            user.map(|u| u.id),
        )
        .execute()
        .await;

        match uploader {
            Ok(uploader) => HttpResponse::Ok().json(json!({
                "link": uploader.to_string()
            })),
            Err(_) => HttpResponse::UnprocessableEntity().json("Invalid file."),
        }
    }

    /// Show/download an uploaded file
    fn show(&self, filename: &str) -> HttpResponse {
        // Check for path traversal attacks
        if !self.is_safe_path(filename) {
            return HttpResponse::BadRequest().finish();
        }

        let uploader = self.uploader();
        if !uploader.exists() {
            return HttpResponse::NotFound().finish();
        }

        // Set cache headers
        let (ttl, directives) = self.cache_settings();
        self.set_cache_headers(ttl, directives);

        // Find the correct file version
        let file_uploader = uploader
            .versions
            .values()
            .find(|version| version.filename == filename)
            .unwrap_or(&uploader);

        if !file_uploader.exists() {
            return HttpResponse::NotFound().finish();
        }

        // Set content type and send file
        self.set_content_type();
        self.send_upload(
            file_uploader,
            Some(filename.to_string()),
            self.content_disposition(),
        )
    }

    /// Authorize a file upload
    fn authorize(&self) -> HttpResponse {
        self.set_workhorse_internal_api_content_type();

        let authorized = self.uploader_class().workhorse_authorize(
            false, // has_length
            self.maximum_size(),
        );

        match authorized {
            Ok(auth) => HttpResponse::Ok().json(auth),
            Err(_) => HttpResponse::InternalServerError().json("Error uploading file"),
        }
    }

    /// Get the uploader class
    fn uploader_class(&self) -> &dyn UploaderClass;

    /// Get the model
    fn model(&self) -> &dyn Model;

    /// Get the uploader
    fn uploader(&self) -> &dyn Uploader;

    /// Get cache settings
    fn cache_settings(&self) -> (Option<i32>, Option<CacheDirectives>);

    /// Get maximum file size
    fn maximum_size(&self) -> i32;

    /// Check if path is safe (no traversal)
    fn is_safe_path(&self, path: &str) -> bool;

    /// Set content type
    fn set_content_type(&self);

    /// Set workhorse internal API content type
    fn set_workhorse_internal_api_content_type(&self);

    /// Send upload
    fn send_upload(
        &self,
        uploader: &dyn Uploader,
        filename: Option<String>,
        disposition: String,
    ) -> HttpResponse;

    /// Get content disposition
    fn content_disposition(&self) -> String {
        if self.uploader().embeddable() || self.uploader().is_pdf() {
            "inline".to_string()
        } else {
            "attachment".to_string()
        }
    }
}

/// Cache directives for uploads
pub struct CacheDirectives {
    pub private: bool,
    pub must_revalidate: bool,
}

/// Trait for uploader classes
pub trait UploaderClass {
    /// Authorize upload with workhorse
    fn workhorse_authorize(
        &self,
        has_length: bool,
        maximum_size: i32,
    ) -> Result<serde_json::Value, std::io::Error>;
}

/// Trait for models that can be uploaded to
pub trait Model {}

/// Trait for uploaders
pub trait Uploader {
    /// Check if file exists
    fn exists(&self) -> bool;

    /// Check if file is embeddable
    fn embeddable(&self) -> bool;

    /// Check if file is PDF
    fn is_pdf(&self) -> bool;

    /// Get file versions
    fn versions(&self) -> &std::collections::HashMap<String, Box<dyn Uploader>>;
}
