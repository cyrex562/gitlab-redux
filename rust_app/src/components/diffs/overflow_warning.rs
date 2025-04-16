use serde_json::json;
use std::sync::Arc;

use super::base::BaseComponent;
use crate::components::Component;
use crate::models::{Commit, MergeRequest, Project};

/// Component for displaying overflow warnings in diffs
pub struct OverflowWarningComponent {
    base: BaseComponent,
    diffs: Arc<Vec<DiffFile>>,
    diff_files: Arc<Vec<DiffFile>>,
    project: Arc<Project>,
    commit: Option<Arc<Commit>>,
    merge_request: Option<Arc<MergeRequest>>,
    can_download_files: bool,
}

impl Component for OverflowWarningComponent {
    fn render(&self) -> String {
        let title = "Some changes are not shown.";
        let message = self.message();
        let diff_link = self.diff_link();
        let patch_link = self.patch_link();

        // Using Tera template syntax for HTML
        format!(
            r#"<div class="alert alert-warning gl-mb-5" data-testid="too-many-changes-alert">
                <div class="alert-body">
                    {}
                </div>
                <div class="alert-actions">
                    {}
                    {}
                </div>
            </div>"#,
            message, diff_link, patch_link
        )
    }
}

impl OverflowWarningComponent {
    pub fn new(
        diffs: Arc<Vec<DiffFile>>,
        diff_files: Arc<Vec<DiffFile>>,
        project: Arc<Project>,
        commit: Option<Arc<Commit>>,
        merge_request: Option<Arc<MergeRequest>>,
    ) -> Self {
        let mut component = Self {
            base: BaseComponent::new(),
            diffs,
            diff_files,
            project,
            commit,
            merge_request,
            can_download_files: false,
        };

        component.before_render();
        component
    }

    fn before_render(&mut self) {
        self.can_download_files = self.commit.is_some() || self.merge_request.is_some();
    }

    fn message(&self) -> String {
        let display_size = self.diff_files.len();
        let real_size = self.diffs.len();

        let base_message = format!(
            "For a faster browsing experience, only <strong>{} of {}</strong> files are shown.",
            display_size, real_size
        );

        let download_message = "Download one of the files below to see all changes.";

        if self.can_download_files {
            format!("{} {}", base_message, download_message)
        } else {
            base_message
        }
    }

    fn diff_link(&self) -> String {
        let text = "Plain diff";

        if let Some(commit) = &self.commit {
            format!(
                r#"<a href="/projects/{}/commits/{}/diff" class="btn btn-default gl-mr-3">{}</a>"#,
                self.project.id, commit.id, text
            )
        } else if let Some(merge_request) = &self.merge_request {
            format!(
                r#"<a href="/merge_requests/{}/diff" class="btn btn-default gl-mr-3">{}</a>"#,
                merge_request.id, text
            )
        } else {
            String::new()
        }
    }

    fn patch_link(&self) -> String {
        let text = "Email patch";

        if let Some(commit) = &self.commit {
            format!(
                r#"<a href="/projects/{}/commits/{}/patch" class="btn btn-default">{}</a>"#,
                self.project.id, commit.id, text
            )
        } else if let Some(merge_request) = &self.merge_request {
            format!(
                r#"<a href="/merge_requests/{}/patch" class="btn btn-default">{}</a>"#,
                merge_request.id, text
            )
        } else {
            String::new()
        }
    }
}

// Placeholder for DiffFile struct
pub struct DiffFile {
    pub file_path: String,
    pub new_path: String,
    pub added_lines: i32,
    pub removed_lines: i32,
}

impl DiffFile {
    pub fn is_deleted_file(&self) -> bool {
        // TODO: Implement logic to determine if file is deleted
        false
    }

    pub fn is_new_file(&self) -> bool {
        // TODO: Implement logic to determine if file is new
        false
    }
}
