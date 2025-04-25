// Ported from: orig_app/app/controllers/concerns/diff_for_path.rb
// This module provides the DiffForPath trait and handler for rendering diffs by file identifier.

use serde::{Deserialize, Serialize};
use std::collections::HashMap;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DiffFile {
    pub file_identifier: String,
    pub content: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DiffCollection {
    pub diff_files: Vec<DiffFile>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Response {
    pub status: u16,
    pub body: String,
    pub headers: HashMap<String, String>,
}

pub trait DiffForPath {
    fn render_diff_for_path(
        &self,
        diffs: &DiffCollection,
        file_identifier: &str,
    ) -> Result<Response, String>;
}

pub struct DiffForPathHandler;

impl DiffForPathHandler {
    pub fn new() -> Self {
        Self
    }

    fn view_to_html_string(&self, template: &str, diff_file: &DiffFile) -> String {
        // In a real implementation, this would render the template with the diff file
        format!("<div class='diff-file'>{}</div>", diff_file.content)
    }
}

impl DiffForPath for DiffForPathHandler {
    fn render_diff_for_path(
        &self,
        diffs: &DiffCollection,
        file_identifier: &str,
    ) -> Result<Response, String> {
        let diff_file = diffs
            .diff_files
            .iter()
            .find(|diff| diff.file_identifier == file_identifier);

        match diff_file {
            Some(file) => {
                let html = self.view_to_html_string("projects/diffs/_content", file);

                let mut headers = HashMap::new();
                headers.insert("Content-Type".to_string(), "application/json".to_string());

                let response_body = serde_json::json!({
                    "html": html
                })
                .to_string();

                Ok(Response {
                    status: 200,
                    body: response_body,
                    headers,
                })
            }
            None => {
                let mut headers = HashMap::new();
                headers.insert("Content-Type".to_string(), "text/plain".to_string());

                Ok(Response {
                    status: 404,
                    body: "Not Found".to_string(),
                    headers,
                })
            }
        }
    }
}
