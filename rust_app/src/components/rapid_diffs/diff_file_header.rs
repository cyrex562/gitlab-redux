use serde::{Deserialize, Serialize};
use std::collections::HashMap;

/// Represents a diff file header with additional metadata
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DiffFileHeaderData {
    /// The file path
    pub file_path: String,
    /// The old file path if renamed
    pub old_path: Option<String>,
    /// The new file path if renamed
    pub new_path: Option<String>,
    /// The content SHA
    pub content_sha: String,
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
    /// The project ID
    pub project_id: i64,
}

/// Renders a diff file header component
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DiffFileHeaderComponent {
    /// The diff file header data
    diff_file: DiffFileHeaderData,
}

impl DiffFileHeaderComponent {
    /// Create a new diff file header component
    pub fn new(diff_file: DiffFileHeaderData) -> Self {
        Self { diff_file }
    }

    /// Get the copy path button HTML
    fn copy_path_button(&self) -> String {
        format!(
            r#"<button class="btn btn-tertiary btn-sm" data-clipboard-text="{}" data-clipboard-gfm="`{}`" title="Copy file path" data-placement="top" data-boundary="viewport" data-testid="rd-diff-file-copy-clipboard">
                <i class="fa fa-copy"></i>
            </button>"#,
            self.diff_file.file_path, self.diff_file.file_path
        )
    }

    /// Get the options menu items JSON
    fn options_menu_items(&self) -> String {
        let view_title = format!("View file @ {}", self.diff_file.content_sha);
        let view_href = format!(
            "/projects/{}/repository/blobs/{}",
            self.diff_file.project_id, self.diff_file.content_sha
        );

        format!(r#"[{{"text":"{}","href":"{}"}}]"#, view_title, view_href)
    }

    /// Render the diff file header component
    pub fn render(&self) -> String {
        let mut html = String::new();

        // Add diff file header container
        html.push_str(r#"<div class="rd-diff-file-header" data-testid="rd-diff-file-header">"#);

        // Add diff file toggle
        html.push_str(r#"<div class="rd-diff-file-toggle">"#);
        html.push_str(r#"<button class="btn btn-tertiary btn-sm rd-diff-file-toggle-button" data-opened="" data-click="toggleFile" aria-label="Hide file contents"><i class="fa fa-chevron-down"></i></button>"#);
        html.push_str(r#"<button class="btn btn-tertiary btn-sm rd-diff-file-toggle-button" data-closed="" data-click="toggleFile" aria-label="Show file contents"><i class="fa fa-chevron-right"></i></button>"#);
        html.push_str("</div>");

        // Add diff file title
        html.push_str(r#"<div class="rd-diff-file-title">"#);

        if self.diff_file.submodule {
            // Submodule
            html.push_str(r#"<span data-testid="rd-diff-file-header-submodule">"#);
            html.push_str(r#"<i class="fa fa-folder"></i>"#);
            html.push_str(&format!(r#"<strong>{}</strong>"#, self.diff_file.file_path));
            html.push_str("</span>");
            html.push_str(&self.copy_path_button());
        } else {
            // Regular file
            if self.diff_file.renamed_file {
                // Renamed file
                if let (Some(old_path), Some(new_path)) =
                    (&self.diff_file.old_path, &self.diff_file.new_path)
                {
                    html.push_str(&format!(r#"<strong>{}</strong>"#, old_path));
                    html.push_str(" → ");
                    html.push_str(&format!(r#"<strong>{}</strong>"#, new_path));
                }
            } else {
                // Regular file
                html.push_str(&format!(r#"<strong>{}</strong>"#, self.diff_file.file_path));
                if self.diff_file.deleted_file {
                    html.push_str(" deleted");
                }
            }
            html.push_str(&self.copy_path_button());

            // Mode changed
            if self.diff_file.mode_changed {
                if let (Some(a_mode), Some(b_mode)) =
                    (&self.diff_file.a_mode, &self.diff_file.b_mode)
                {
                    html.push_str(&format!(r#"<small>{} → {}</small>"#, a_mode, b_mode));
                }
            }

            // External storage
            if self.diff_file.stored_externally {
                if let Some(storage) = &self.diff_file.external_storage {
                    if storage == "lfs" {
                        html.push_str(r#"<span class="badge badge-neutral">LFS</span>"#);
                    }
                }
            }
        }

        html.push_str("</div>");

        // Add diff file info
        html.push_str(r#"<div class="rd-diff-file-info">"#);

        // Add diff file stats
        html.push_str(r#"<div class="rd-diff-file-stats">"#);
        html.push_str(r#"<span class="rd-lines-added"><span>+</span><span data-testid="js-file-addition-line">"#);
        html.push_str(&format!("{}</span></span>", self.diff_file.added_lines));
        html.push_str(r#"<span class="rd-lines-removed"><span>−</span><span data-testid="js-file-deletion-line">"#);
        html.push_str(&format!("{}</span></span>", self.diff_file.removed_lines));
        html.push_str("</div>");

        // Add diff file options menu
        html.push_str(r#"<div class="rd-diff-file-options-menu">"#);
        html.push_str(r#"<div class="js-options-menu">"#);
        html.push_str(&format!(
            r#"<script type="application/json">{}</script>"#,
            self.options_menu_items()
        ));
        html.push_str(r#"<button class="btn btn-tertiary btn-sm js-options-button" data-click="toggleOptionsMenu" aria-label="Options"><i class="fa fa-ellipsis-v"></i></button>"#);
        html.push_str("</div></div>");

        // Close diff file info
        html.push_str("</div>");

        // Close diff file header
        html.push_str("</div>");

        html
    }
}
