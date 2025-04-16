use crate::config::settings::Settings;
use actix_web::{web, HttpResponse};
use serde::{Deserialize, Serialize};
use std::io;
use tokio::fs::File;
use tokio::io::{AsyncReadExt, BufReader};

/// Module for streaming diffs
pub trait StreamDiffs {
    /// Get the old content
    fn old_content(&self) -> &str;

    /// Get the new content
    fn new_content(&self) -> &str;

    /// Get the file path
    fn file_path(&self) -> Option<&str>;

    /// Stream the diff between old and new content
    async fn stream_diff(&self) -> Result<HttpResponse, io::Error> {
        let old_content = self.old_content();
        let new_content = self.new_content();

        // Generate diff
        let diff = self.generate_diff(old_content, new_content);

        // Create response with diff
        Ok(HttpResponse::Ok().content_type("text/plain").body(diff))
    }

    /// Generate a diff between old and new content
    fn generate_diff(&self, old_content: &str, new_content: &str) -> String {
        // TODO: Implement actual diff generation
        // This would typically involve:
        // 1. Using a diff library to generate the diff
        // 2. Formatting the diff output
        // 3. Returning the formatted diff

        format!(
            "--- a/{}\n+++ b/{}\n",
            self.file_path().unwrap_or("file"),
            self.file_path().unwrap_or("file")
        )
    }

    /// Stream a file diff
    async fn stream_file_diff(
        &self,
        old_path: &str,
        new_path: &str,
    ) -> Result<HttpResponse, io::Error> {
        // Read old file
        let old_file = File::open(old_path).await?;
        let mut old_content = String::new();
        BufReader::new(old_file)
            .read_to_string(&mut old_content)
            .await?;

        // Read new file
        let new_file = File::open(new_path).await?;
        let mut new_content = String::new();
        BufReader::new(new_file)
            .read_to_string(&mut new_content)
            .await?;

        // Generate diff
        let diff = self.generate_diff(&old_content, &new_content);

        // Create response with diff
        Ok(HttpResponse::Ok().content_type("text/plain").body(diff))
    }

    /// Get diff settings
    fn get_diff_settings(&self) -> HashMap<String, String> {
        let mut settings = HashMap::new();
        let settings = Settings::current();

        settings.insert(
            "diff_context_lines".to_string(),
            settings.diff_context_lines.to_string(),
        );
        settings.insert(
            "diff_max_files".to_string(),
            settings.diff_max_files.to_string(),
        );
        settings.insert(
            "diff_max_lines".to_string(),
            settings.diff_max_lines.to_string(),
        );

        settings
    }
}
