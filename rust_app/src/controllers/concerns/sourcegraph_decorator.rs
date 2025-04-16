use actix_web::{web, HttpResponse};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

/// Module for decorating content with Sourcegraph information
pub trait SourcegraphDecorator {
    /// Get the content to decorate
    fn content(&self) -> String;

    /// Get the file path
    fn file_path(&self) -> String;

    /// Get the project ID
    fn project_id(&self) -> i32;

    /// Get the branch name
    fn branch_name(&self) -> String;

    /// Get the commit SHA
    fn commit_sha(&self) -> String;

    /// Get the Sourcegraph URL
    fn sourcegraph_url(&self) -> String {
        "https://sourcegraph.com".to_string()
    }

    /// Get the decoration type
    fn decoration_type(&self) -> String {
        "code".to_string()
    }

    /// Generate Sourcegraph URL
    fn generate_sourcegraph_url(&self) -> String {
        format!(
            "{}/{}/-/blob/{}/{}",
            self.sourcegraph_url(),
            self.project_id(),
            self.branch_name(),
            self.file_path()
        )
    }

    /// Add Sourcegraph decorations to content
    fn decorate_content(&self) -> Result<String, HttpResponse> {
        let mut decorated_content = self.content();
        let sourcegraph_url = self.generate_sourcegraph_url();

        // Add Sourcegraph link
        decorated_content = format!(
            "<!-- Sourcegraph: {} -->\n{}",
            sourcegraph_url, decorated_content
        );

        // Add decoration metadata
        let mut metadata = HashMap::new();
        metadata.insert("type".to_string(), self.decoration_type());
        metadata.insert("url".to_string(), sourcegraph_url);
        metadata.insert("project_id".to_string(), self.project_id().to_string());
        metadata.insert("branch".to_string(), self.branch_name());
        metadata.insert("commit".to_string(), self.commit_sha());

        let metadata_json = serde_json::to_string(&metadata).map_err(|e| {
            HttpResponse::InternalServerError().json(serde_json::json!({
                "error": format!("Failed to serialize metadata: {}", e)
            }))
        })?;

        decorated_content = format!(
            "<!-- Sourcegraph Metadata: {} -->\n{}",
            metadata_json, decorated_content
        );

        Ok(decorated_content)
    }

    /// Get decoration metadata
    fn get_decoration_metadata(&self) -> Result<HashMap<String, String>, HttpResponse> {
        let mut metadata = HashMap::new();

        metadata.insert("type".to_string(), self.decoration_type());
        metadata.insert("url".to_string(), self.generate_sourcegraph_url());
        metadata.insert("project_id".to_string(), self.project_id().to_string());
        metadata.insert("branch".to_string(), self.branch_name());
        metadata.insert("commit".to_string(), self.commit_sha());

        Ok(metadata)
    }
}
