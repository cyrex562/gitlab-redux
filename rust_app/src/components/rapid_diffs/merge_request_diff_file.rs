use serde::{Deserialize, Serialize};
use std::collections::HashMap;

use crate::components::rapid_diffs::diff_file::{DiffFileComponent, DiffFileData};

/// Represents a merge request
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct MergeRequest {
    /// The merge request ID
    pub id: i64,
    /// The merge request IID
    pub iid: i64,
    /// The merge request title
    pub title: String,
    /// The merge request description
    pub description: Option<String>,
    /// The merge request state
    pub state: String,
    /// The merge request source branch
    pub source_branch: String,
    /// The merge request target branch
    pub target_branch: String,
}

/// Renders a merge request diff file component
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct MergeRequestDiffFileComponent {
    /// The diff file data
    diff_file: DiffFileData,
    /// The merge request
    merge_request: MergeRequest,
    /// Whether to use parallel view
    parallel_view: bool,
}

impl MergeRequestDiffFileComponent {
    /// Create a new merge request diff file component
    pub fn new(diff_file: DiffFileData, merge_request: MergeRequest, parallel_view: bool) -> Self {
        Self {
            diff_file,
            merge_request,
            parallel_view,
        }
    }

    /// Render the merge request diff file component
    pub fn render(&self) -> String {
        // Create a DiffFileComponent and render it
        let diff_file_component =
            DiffFileComponent::new(self.diff_file.clone(), self.parallel_view);
        diff_file_component.render()
    }
}
