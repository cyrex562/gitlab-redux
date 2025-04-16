use serde::{Deserialize, Serialize};

/// Renders a streaming error component
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct StreamingErrorComponent {
    /// The error message
    message: String,
}

impl StreamingErrorComponent {
    /// Create a new streaming error component
    pub fn new(message: String) -> Self {
        Self { message }
    }

    /// Render the streaming error component
    pub fn render(&self) -> String {
        format!(
            r#"<div class="streaming-error" data-message="{}"></div>"#,
            self.message
        )
    }
} 