use crate::config::settings::Settings;
use actix_web::{web, HttpResponse};
use base64;
use std::io;

/// Module for sending blob data
pub trait SendsBlob {
    /// Get the blob data
    fn blob_data(&self) -> &[u8];

    /// Get the content type
    fn content_type(&self) -> Option<String>;

    /// Get the filename
    fn filename(&self) -> Option<String>;

    /// Send the blob data
    fn send_blob(&self) -> HttpResponse {
        let blob_data = self.blob_data();
        let content_type = self
            .content_type()
            .unwrap_or_else(|| "application/octet-stream".to_string());

        let mut response = HttpResponse::Ok().content_type(content_type);

        // Add filename header if provided
        if let Some(filename) = self.filename() {
            response = response.header(
                "Content-Disposition",
                format!("attachment; filename=\"{}\"", filename),
            );
        }

        // Add content length header
        response = response.header("Content-Length", blob_data.len().to_string());

        // Set the body
        response.body(blob_data.to_vec())
    }

    /// Send the blob data as base64
    fn send_blob_base64(&self) -> HttpResponse {
        let blob_data = self.blob_data();
        let base64_data = base64::encode(blob_data);

        HttpResponse::Ok()
            .content_type("text/plain")
            .body(base64_data)
    }

    /// Check if blob sending is allowed
    fn is_blob_sending_allowed(&self) -> bool {
        let settings = Settings::current();

        // Check blob size
        if self.blob_data().len() > settings.max_blob_size {
            return false;
        }

        // Check content type if provided
        if let Some(content_type) = self.content_type() {
            if !settings.allowed_blob_types.contains(&content_type) {
                return false;
            }
        }

        true
    }

    /// Get blob settings
    fn get_blob_settings(&self) -> HashMap<String, String> {
        let mut settings = HashMap::new();
        let settings = Settings::current();

        settings.insert(
            "max_blob_size".to_string(),
            settings.max_blob_size.to_string(),
        );
        settings.insert(
            "allowed_blob_types".to_string(),
            settings.allowed_blob_types.join(","),
        );

        settings
    }
}
