use serde::{Deserialize, Serialize};
use std::collections::HashMap;

/// Trait for components that need to render labels with help text
pub trait CheckboxRadioLabelWithHelpText {
    /// Get the label content
    fn label_content(&self) -> Option<String>;

    /// Get the help text content
    fn help_text_content(&self) -> Option<String>;

    /// Get the label options
    fn label_options(&self) -> &HashMap<String, String>;

    /// Get the help text options
    fn help_text_options(&self) -> &HashMap<String, String>;

    /// Format options with CSS classes
    fn format_options(
        &self,
        options: &HashMap<String, String>,
        css_classes: &[&str],
    ) -> HashMap<String, String>;

    /// Render a label with help text
    fn render_label_with_help_text(&self) -> String {
        let mut html = String::new();

        // Add label
        html.push_str(&self.render_label_tag_with_help_text());

        // Add help text if present
        if let Some(help_text) = self.help_text_content() {
            html.push_str(&format!(
                "<div class=\"form-text text-muted\">{}</div>",
                help_text
            ));
        }

        html
    }

    /// Render a label tag with help text
    fn render_label_tag_with_help_text(&self) -> String {
        let mut html = String::new();

        // Add label entry
        html.push_str(&self.render_label_entry());

        // Add help text if present
        if let Some(help_text) = self.help_text_content() {
            html.push_str(&format!(
                "<div class=\"form-text text-muted\">{}</div>",
                help_text
            ));
        }

        html
    }

    /// Render a label entry
    fn render_label_entry(&self) -> String {
        let mut html = String::new();

        // Add label
        if let Some(label) = self.label_content() {
            html.push_str(&format!("<span>{}</span>", label));
        }

        html
    }

    /// Format label options
    fn formatted_label_options(&self) -> HashMap<String, String> {
        self.format_options(self.label_options(), &["custom-control-label"])
    }
}
