use serde_json::json;
use std::sync::Arc;

use super::base::BaseComponent;
use super::overflow_warning::DiffFile;
use crate::components::Component;

/// Component for displaying diff statistics
pub struct StatsComponent {
    base: BaseComponent,
    diff_files: Arc<Vec<DiffFile>>,
    changed: usize,
    added: i32,
    removed: i32,
}

impl Component for StatsComponent {
    fn render(&self) -> String {
        let changed = self.changed;
        let added = self.added;
        let deleted = self.removed;
        let files = self.diff_files_data();

        // Using Tera template syntax for HTML with htmx
        format!(
            r#"<div class="js-diff-stats-dropdown" 
                data-changed="{}" 
                data-added="{}" 
                data-deleted="{}" 
                data-files='{}'
                hx-get="/api/diff-stats"
                hx-trigger="load"
                hx-swap="innerHTML">
            </div>"#,
            changed, added, deleted, files
        )
    }
}

impl StatsComponent {
    pub fn new(diff_files: Arc<Vec<DiffFile>>) -> Self {
        let changed = diff_files.len();
        let added = diff_files.iter().map(|f| f.added_lines).sum();
        let removed = diff_files.iter().map(|f| f.removed_lines).sum();

        Self {
            base: BaseComponent::new(),
            diff_files,
            changed,
            added,
            removed,
        }
    }

    fn diff_files_data(&self) -> String {
        let diffs_map: Vec<serde_json::Value> = self
            .diff_files
            .iter()
            .map(|f| {
                json!({
                    "href": format!("#{}", self.hexdigest(&f.file_path)),
                    "title": f.new_path,
                    "name": f.file_path,
                    "path": self.diff_file_path_text(f),
                    "icon": self.diff_file_changed_icon(f),
                    "iconColor": self.diff_file_changed_icon_color(f),
                    "added": f.added_lines,
                    "removed": f.removed_lines
                })
            })
            .collect();

        serde_json::to_string(&diffs_map).unwrap_or_default()
    }

    fn diff_file_path_text(&self, diff_file: &DiffFile, max: usize) -> String {
        let path = &diff_file.new_path;

        if path.len() <= max || max <= 3 {
            return path.clone();
        }

        format!("...{}", &path[path.len() - (max - 3)..])
    }

    fn diff_file_changed_icon(&self, diff_file: &DiffFile) -> &'static str {
        if diff_file.is_deleted_file() {
            "file-deletion"
        } else if diff_file.is_new_file() {
            "file-addition"
        } else {
            "file-modified"
        }
    }

    fn diff_file_changed_icon_color(&self, diff_file: &DiffFile) -> &'static str {
        if diff_file.is_deleted_file() {
            "danger"
        } else if diff_file.is_new_file() {
            "success"
        } else {
            "default"
        }
    }

    fn hexdigest(&self, input: &str) -> String {
        use sha2::{Digest, Sha256};
        let mut hasher = Sha256::new();
        hasher.update(input.as_bytes());
        format!("{:x}", hasher.finalize())
    }
}
