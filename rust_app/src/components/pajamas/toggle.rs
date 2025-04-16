use serde::{Deserialize, Serialize};
use std::collections::HashMap;

/// Label position options for the toggle component
#[derive(Debug, Clone, Copy, PartialEq, Serialize, Deserialize)]
pub enum LabelPosition {
    Top,
    Left,
    Hidden,
}

impl LabelPosition {
    fn as_str(&self) -> &'static str {
        match self {
            LabelPosition::Top => "top",
            LabelPosition::Left => "left",
            LabelPosition::Hidden => "hidden",
        }
    }
}

/// Renders a GlToggle root element
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ToggleComponent {
    /// CSS classes
    classes: String,
    /// Label text
    label: Option<String>,
    /// Label position
    label_position: Option<LabelPosition>,
    /// HTML ID
    id: Option<String>,
    /// Input name
    name: Option<String>,
    /// Help text
    help: Option<String>,
    /// Data attributes
    data: HashMap<String, String>,
    /// Whether the toggle is disabled
    is_disabled: bool,
    /// Whether the toggle is checked
    is_checked: bool,
    /// Whether the toggle is loading
    is_loading: bool,
    /// Rich help text content
    content: Option<String>,
}

impl ToggleComponent {
    /// Create a new toggle component
    pub fn new(
        classes: String,
        label: Option<String>,
        label_position: Option<LabelPosition>,
        id: Option<String>,
        name: Option<String>,
        help: Option<String>,
        data: Option<HashMap<String, String>>,
        is_disabled: bool,
        is_checked: bool,
        is_loading: bool,
        content: Option<String>,
    ) -> Self {
        Self {
            classes,
            label,
            label_position,
            id,
            name,
            help,
            data: data.unwrap_or_default(),
            is_disabled,
            is_checked,
            is_loading,
            content,
        }
    }

    /// Render the toggle component
    pub fn render(&self) -> String {
        let mut html = format!(
            r#"<span class="{}"{}>"#,
            self.classes,
            self.render_data_attributes()
        );

        // Render rich help text content if present
        if let Some(content) = &self.content {
            html.push_str(&format!(r#"<div class="gl-help-label">{}</div>"#, content));
        }

        html.push_str("</span>");
        html
    }

    /// Render data attributes
    fn render_data_attributes(&self) -> String {
        let mut data_attrs = self.data.clone();

        // Add standard data attributes
        if let Some(name) = &self.name {
            data_attrs.insert("name".to_string(), name.clone());
        }
        if let Some(id) = &self.id {
            data_attrs.insert("id".to_string(), id.clone());
        }
        if let Some(label) = &self.label {
            data_attrs.insert("label".to_string(), label.clone());
        }
        if let Some(help) = &self.help {
            data_attrs.insert("help".to_string(), help.clone());
        }
        if let Some(label_position) = &self.label_position {
            data_attrs.insert(
                "label-position".to_string(),
                label_position.as_str().to_string(),
            );
        }

        // Add boolean attributes
        data_attrs.insert("is-checked".to_string(), self.is_checked.to_string());
        data_attrs.insert("disabled".to_string(), self.is_disabled.to_string());
        data_attrs.insert("is-loading".to_string(), self.is_loading.to_string());

        // Format data attributes
        data_attrs
            .iter()
            .map(|(key, value)| format!(" data-{}={}", key, value))
            .collect()
    }
}
