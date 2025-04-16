use serde::{Deserialize, Serialize};
use std::collections::HashMap;

/// Represents a slice of diffs
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DiffsSlice {
    /// The diffs in this slice
    pub diffs: Vec<DiffFile>,
}

/// Represents a diff file
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DiffFile {
    /// The file path
    pub path: String,
    /// The file content
    pub content: String,
    /// Whether the file is new
    pub new_file: bool,
    /// Whether the file is deleted
    pub deleted_file: bool,
    /// Whether the file is renamed
    pub renamed_file: bool,
    /// The old path if the file is renamed
    pub old_path: Option<String>,
}

/// Renders the main app component for rapid diffs
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AppComponent {
    /// The diffs slice to render
    diffs_slice: DiffsSlice,
    /// The URL to reload the stream
    reload_stream_url: String,
    /// The URL to stream diffs
    stream_url: Option<String>,
    /// Whether to show whitespace
    show_whitespace: bool,
    /// The diff view type
    diff_view: String,
    /// The endpoint to update user preferences
    update_user_endpoint: String,
    /// The endpoint to get metadata
    metadata_endpoint: String,
    /// Whether to preload the diffs
    preload: bool,
    /// The initial sidebar width
    initial_sidebar_width: Option<String>,
    /// The diffs list content
    diffs_list: Option<String>,
}

impl AppComponent {
    /// Create a new app component
    pub fn new(
        diffs_slice: DiffsSlice,
        reload_stream_url: String,
        stream_url: Option<String>,
        show_whitespace: bool,
        diff_view: String,
        update_user_endpoint: String,
        metadata_endpoint: String,
        preload: bool,
        initial_sidebar_width: Option<String>,
    ) -> Self {
        Self {
            diffs_slice,
            reload_stream_url,
            stream_url,
            show_whitespace,
            diff_view,
            update_user_endpoint,
            metadata_endpoint,
            preload,
            initial_sidebar_width,
            diffs_list: None,
        }
    }

    /// Set the diffs list content
    pub fn with_diffs_list(mut self, diffs_list: String) -> Self {
        self.diffs_list = Some(diffs_list);
        self
    }

    /// Render the app component
    pub fn render(&self) -> String {
        let mut html = String::new();

        // Add preload script if needed
        if self.preload {
            html.push_str(&format!(
                r#"<script nonce="content_security_policy_nonce">
                var controller = new AbortController();
                window.gl.rapidDiffsPreload = {{
                    controller: controller,
                    streamRequest: fetch('{}', {{ signal: controller.signal }})
                }}
                </script>"#,
                self.stream_url.as_ref().unwrap_or(&String::new())
            ));
        }

        // Add app container
        html.push_str(&format!(
            r#"<div class="rd-app" data-rapid-diffs="true" data-reload-stream-url="{}" data-metadata-endpoint="{}">"#,
            self.reload_stream_url, self.metadata_endpoint
        ));

        // Add app header
        html.push_str(r#"<div class="rd-app-header">"#);
        html.push_str(r#"<div class="rd-app-settings">"#);
        html.push_str(&format!(
            r#"<div data-view-settings="true" data-show-whitespace="{}" data-diff-view-type="{}" data-update-user-endpoint="{}"></div>"#,
            self.show_whitespace, self.diff_view, self.update_user_endpoint
        ));
        html.push_str("</div></div>");

        // Add app body
        html.push_str(r#"<div class="rd-app-body">"#);

        // Add sidebar
        let sidebar_style = self
            .initial_sidebar_width
            .as_ref()
            .map(|width| format!(" style=\"width: {}px\"", width))
            .unwrap_or_default();
        html.push_str(&format!(
            r#"<div class="rd-app-sidebar" data-file-browser="true"{}>"#,
            sidebar_style
        ));
        html.push_str(r#"<div class="rd-app-sidebar-loading">"#);
        html.push_str(r#"<span class="gl-spinner gl-spinner-sm"></span>"#);
        html.push_str("</div></div>");

        // Add content
        html.push_str(r#"<div class="rd-app-content" data-sidebar-visible="true">"#);
        html.push_str(r#"<div class="rd-app-content-header" data-hidden-files-warning="true"></div>"#);
        html.push_str(r#"<div class="code">"#);
        html.push_str(r#"<div data-diffs-list="true">"#);

        // Add performance mark
        html.push_str(r#"<script nonce="content_security_policy_nonce">"#);
        html.push_str(
            r#"requestAnimationFrame(() => { window.performance.mark('rapid-diffs-first-diff-file-shown') })"#,
        );
        html.push_str("</script>");

        // Add diffs list or render diffs
        if let Some(diffs_list) = &self.diffs_list {
            html.push_str(diffs_list);
        } else {
            // Render diffs slice
            for diff in &self.diffs_slice.diffs {
                // TODO: Implement DiffFileComponent rendering
                html.push_str(&format!(
                    r#"<div class="diff-file" data-path="{}"></div>"#,
                    diff.path
                ));
            }
        }

        // Add stream container or performance mark
        if let Some(stream_url) = &self.stream_url {
            html.push_str(&format!(
                r#"<div id="js-stream-container" data-diffs-stream-url="{}"></div>"#,
                stream_url
            ));
        } else {
            html.push_str(r#"<script nonce="content_security_policy_nonce">"#);
            html.push_str(
                r#"requestAnimationFrame(() => {
                    window.performance.mark('rapid-diffs-list-loaded');
                    window.performance.measure('rapid-diffs-list-loading', 'rapid-diffs-first-diff-file-shown', 'rapid-diffs-list-loaded');
                })"#,
            );
            html.push_str("</script>");
        }

        // Close all divs
        html.push_str("</div></div></div></div></div>");

        html
    }
} 