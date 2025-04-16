use crate::config::settings::Settings;
use actix_web::{web, HttpResponse};
use mime_guess::MimeGuess;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::io;
use std::path::PathBuf;
use tokio::fs;

/// Module for handling file upload sending
pub trait SendFileUpload {
    /// Get the file path
    fn file_path(&self) -> PathBuf;

    /// Get the file name
    fn file_name(&self) -> String;

    /// Get the content type
    fn content_type(&self) -> String;

    /// Get the file size
    fn file_size(&self) -> u64;

    /// Get the file extension
    fn file_extension(&self) -> String {
        self.file_name()
            .split('.')
            .last()
            .map(|s| s.to_string())
            .unwrap_or_default()
    }

    /// Check if the file exists
    async fn file_exists(&self) -> bool {
        self.file_path().exists()
    }

    /// Get the file metadata
    async fn get_file_metadata(&self) -> Result<HashMap<String, String>, HttpResponse> {
        let mut metadata = HashMap::new();

        metadata.insert("name".to_string(), self.file_name());
        metadata.insert("content_type".to_string(), self.content_type());
        metadata.insert("size".to_string(), self.file_size().to_string());
        metadata.insert("extension".to_string(), self.file_extension());

        Ok(metadata)
    }

    /// Send the file as a response
    async fn send_file(&self) -> Result<HttpResponse, HttpResponse> {
        if !self.file_exists().await {
            return Err(HttpResponse::NotFound().json(serde_json::json!({
                "error": "File not found"
            })));
        }

        let file = tokio::fs::File::open(&self.file_path())
            .await
            .map_err(|e| {
                HttpResponse::InternalServerError().json(serde_json::json!({
                    "error": format!("Failed to open file: {}", e)
                }))
            })?;

        Ok(HttpResponse::Ok()
            .content_type(self.content_type())
            .header(
                "Content-Disposition",
                format!("attachment; filename=\"{}\"", self.file_name()),
            )
            .streaming(file))
    }

    /// Send the file as a stream
    async fn send_file_stream(&self) -> Result<HttpResponse, HttpResponse> {
        if !self.file_exists().await {
            return Err(HttpResponse::NotFound().json(serde_json::json!({
                "error": "File not found"
            })));
        }

        let file = tokio::fs::File::open(&self.file_path())
            .await
            .map_err(|e| {
                HttpResponse::InternalServerError().json(serde_json::json!({
                    "error": format!("Failed to open file: {}", e)
                }))
            })?;

        Ok(HttpResponse::Ok()
            .content_type(self.content_type())
            .header(
                "Content-Disposition",
                format!("inline; filename=\"{}\"", self.file_name()),
            )
            .streaming(file))
    }

    /// Delete the file
    async fn delete_file(&self) -> Result<(), HttpResponse> {
        if !self.file_exists().await {
            return Err(HttpResponse::NotFound().json(serde_json::json!({
                "error": "File not found"
            })));
        }

        tokio::fs::remove_file(&self.file_path())
            .await
            .map_err(|e| {
                HttpResponse::InternalServerError().json(serde_json::json!({
                    "error": format!("Failed to delete file: {}", e)
                }))
            })?;

        // Get file metadata
        let metadata = fs::metadata(&file_path).await?;
        let file_size = metadata.len();

        // Determine content type
        let content_type = self.content_type().unwrap_or_else(|| {
            MimeGuess::from_path(&file_path)
                .first_or_octet_stream()
                .to_string()
        });

        // Create response with file
        let response = HttpResponse::Ok()
            .content_type(content_type)
            .header(
                "Content-Disposition",
                format!("attachment; filename=\"{}\"", file_name),
            )
            .header("Content-Length", file_size.to_string());

        // Stream the file
        let file = fs::File::open(&file_path).await?;
        Ok(response.streaming(file))
    }

    /// Check if file upload is allowed
    fn is_file_upload_allowed(&self) -> bool {
        let settings = Settings::current();

        // Check file size
        if let Ok(metadata) = std::fs::metadata(self.file_path()) {
            if metadata.len() > settings.max_file_upload_size {
                return false;
            }
        }

        // Check file type
        if let Some(content_type) = self.content_type() {
            if !settings.allowed_file_types.contains(&content_type) {
                return false;
            }
        }

        true
    }

    /// Get file upload settings
    fn get_file_upload_settings(&self) -> HashMap<String, String> {
        let mut settings = HashMap::new();
        let settings = Settings::current();

        settings.insert(
            "max_file_size".to_string(),
            settings.max_file_upload_size.to_string(),
        );
        settings.insert(
            "allowed_file_types".to_string(),
            settings.allowed_file_types.join(","),
        );
        settings.insert(
            "upload_directory".to_string(),
            settings.upload_directory.to_string(),
        );

        settings
    }
}
