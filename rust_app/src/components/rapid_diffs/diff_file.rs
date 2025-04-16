use serde::{Deserialize, Serialize};
use std::collections::HashMap;

use crate::components::rapid_diffs::diff_file_header::{
    DiffFileHeaderComponent, DiffFileHeaderData,
};

/// Represents a diff file with additional metadata
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DiffFileData {
    /// The file hash
    pub file_hash: String,
    /// The file path
    pub file_path: String,
    /// The content SHA
    pub content_sha: String,
    /// Whether the file is collapsed
    pub collapsed: bool,
    /// Whether the file is modified
    pub modified_file: bool,
    /// Whether the file is diffable text
    pub diffable_text: bool,
    /// The repository project ID
    pub project_id: i64,
    /// The repository ID
    pub repository_id: i64,
    /// The old file path if renamed
    pub old_path: Option<String>,
    /// The new file path if renamed
    pub new_path: Option<String>,
    /// Whether the file is a submodule
    pub submodule: bool,
    /// Whether the file is renamed
    pub renamed_file: bool,
    /// Whether the file is deleted
    pub deleted_file: bool,
    /// Whether the file mode is changed
    pub mode_changed: bool,
    /// The old mode if changed
    pub a_mode: Option<String>,
    /// The new mode if changed
    pub b_mode: Option<String>,
    /// Whether the file is stored externally
    pub stored_externally: bool,
    /// The external storage type
    pub external_storage: Option<String>,
    /// The number of added lines
    pub added_lines: i32,
    /// The number of removed lines
    pub removed_lines: i32,
}

/// Renders a diff file component
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DiffFileComponent {
    /// The diff file data
    diff_file: DiffFileData,
    /// Whether to use parallel view
    parallel_view: bool,
}

impl DiffFileComponent {
    /// Create a new diff file component
    pub fn new(diff_file: DiffFileData, parallel_view: bool) -> Self {
        Self {
            diff_file,
            parallel_view,
        }
    }

    /// Get the ID of the diff file
    fn id(&self) -> String {
        self.diff_file.file_hash.clone()
    }

    /// Get the server data for the diff file
    fn server_data(&self) -> HashMap<String, String> {
        let mut data = HashMap::new();

        // Determine the viewer component
        let viewer = self.viewer_component();

        // Add viewer name
        data.insert("viewer".to_string(), viewer.to_string());

        // Add diff lines path
        let diff_lines_path = format!(
            "/projects/{}/repository/blobs/{}/diff_lines",
            self.diff_file.project_id, self.diff_file.content_sha
        );
        data.insert("diff_lines_path".to_string(), diff_lines_path);

        data
    }

    /// Get the viewer component name
    fn viewer_component(&self) -> &'static str {
        if self.diff_file.collapsed || !self.diff_file.modified_file {
            return "NoPreviewComponent";
        }

        if self.diff_file.diffable_text {
            if self.parallel_view {
                return "ParallelViewComponent";
            } else {
                return "InlineViewComponent";
            }
        }

        "NoPreviewComponent"
    }

    /// Convert DiffFileData to DiffFileHeaderData
    fn to_header_data(&self) -> DiffFileHeaderData {
        DiffFileHeaderData {
            file_path: self.diff_file.file_path.clone(),
            old_path: self.diff_file.old_path.clone(),
            new_path: self.diff_file.new_path.clone(),
            content_sha: self.diff_file.content_sha.clone(),
            submodule: self.diff_file.submodule,
            renamed_file: self.diff_file.renamed_file,
            deleted_file: self.diff_file.deleted_file,
            mode_changed: self.diff_file.mode_changed,
            a_mode: self.diff_file.a_mode.clone(),
            b_mode: self.diff_file.b_mode.clone(),
            stored_externally: self.diff_file.stored_externally,
            external_storage: self.diff_file.external_storage.clone(),
            added_lines: self.diff_file.added_lines,
            removed_lines: self.diff_file.removed_lines,
            project_id: self.diff_file.project_id,
        }
    }

    /// Render the diff file component
    pub fn render(&self) -> String {
        let mut html = String::new();

        // Convert server data to JSON
        let server_data_json = serde_json::to_string(&self.server_data()).unwrap_or_default();

        // Add diff file container
        html.push_str(&format!(
            r#"<div class="rd-diff-file-component" id="{}" data-server-data='{}'>"#,
            self.id(),
            server_data_json
        ));

        // Add diff file
        html.push_str(r#"<div class="rd-diff-file">"#);

        // Add diff file header
        let header_component = DiffFileHeaderComponent::new(self.to_header_data());
        html.push_str(&header_component.render());

        // Add file body wrapper
        html.push_str(r#"<div data-file-body="">"#);
        html.push_str(r#"<div class="rd-diff-file-body">"#);

        // Add viewer component
        // TODO: Implement viewer component rendering
        html.push_str(&format!(
            r#"<div class="rd-viewer-component" data-viewer="{}"></div>"#,
            self.viewer_component()
        ));

        // Close divs
        html.push_str("</div></div>");

        // Add diff file mounted marker
        html.push_str(r#"<div class="diff-file-mounted"></div>"#);

        // Close diff file container
        html.push_str("</div>");

        html
    }
}
