use serde::{Deserialize, Serialize};
use std::collections::HashMap;

/// Represents a line in a diff
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DiffLine {
    /// The text content of the line
    pub text_content: String,
    /// Whether the line was added
    pub added: bool,
    /// Whether the line was removed
    pub removed: bool,
}

impl DiffLine {
    /// Create a new diff line
    pub fn new(text_content: String, added: bool, removed: bool) -> Self {
        Self {
            text_content,
            added,
            removed,
        }
    }
}

/// Renders a line content component for a diff
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct LineContentComponent {
    /// The line to render
    line: Option<DiffLine>,
    /// The position of the line
    position: String,
}

impl LineContentComponent {
    /// Create a new line content component
    pub fn new(line: Option<DiffLine>, position: String) -> Self {
        Self { line, position }
    }

    /// Get the change type of the line
    fn change_type(&self) -> Option<String> {
        if let Some(line) = &self.line {
            if line.added {
                return Some("added".to_string());
            }
            if line.removed {
                return Some("removed".to_string());
            }
        }
        None
    }

    /// Render the line content component
    pub fn render(&self) -> String {
        if let Some(line) = &self.line {
            let change_type = self.change_type();
            let change_attr = change_type
                .map(|t| format!(" data-change=\"{}\"", t))
                .unwrap_or_default();

            format!(
                r#"<td class="rd-line-content" data-position="{}"{} tabindex="-1">{}</td>"#,
                self.position, change_attr, line.text_content
            )
        } else {
            String::new()
        }
    }
}
