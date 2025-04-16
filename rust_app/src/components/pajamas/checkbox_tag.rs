use serde::{Deserialize, Serialize};
use std::collections::HashMap;

/// Renders a Pajamas compliant checkbox element without a form builder
#[derive(Debug, Clone)]
pub struct CheckboxTagComponent {
    /// The name of the checkbox
    name: String,
    /// Options for the label
    label_options: HashMap<String, String>,
    /// Options for the checkbox input
    checkbox_options: HashMap<String, String>,
    /// The value of the checkbox
    value: String,
    /// Whether the checkbox is checked
    checked: bool,
    /// The label content slot
    label_slot: Option<String>,
    /// The help text content slot
    help_text_slot: Option<String>,
}

impl CheckboxTagComponent {
    /// Create a new checkbox tag component
    pub fn new(
        name: String,
        label_options: Option<HashMap<String, String>>,
        checkbox_options: Option<HashMap<String, String>>,
        value: Option<String>,
        checked: bool,
    ) -> Self {
        Self {
            name,
            label_options: label_options.unwrap_or_default(),
            checkbox_options: checkbox_options.unwrap_or_default(),
            value: value.unwrap_or_else(|| "1".to_string()),
            checked,
            label_slot: None,
            help_text_slot: None,
        }
    }

    /// Set the label content slot
    pub fn with_label_slot(mut self, label: String) -> Self {
        self.label_slot = Some(label);
        self
    }

    /// Set the help text content slot
    pub fn with_help_text_slot(mut self, help_text: String) -> Self {
        self.help_text_slot = Some(help_text);
        self
    }

    /// Get the label content
    fn label_content(&self) -> Option<String> {
        self.label_slot.clone()
    }

    /// Get the help text content
    fn help_text_content(&self) -> Option<String> {
        self.help_text_slot.clone()
    }

    /// Format the input options
    fn formatted_input_options(&self) -> String {
        let mut options = self.checkbox_options.clone();
        options.insert("class".to_string(), "custom-control-input".to_string());

        options
            .iter()
            .map(|(key, value)| format!(" {}=\"{}\"", key, value))
            .collect::<Vec<String>>()
            .join("")
    }

    /// Format the label options
    fn formatted_label_options(&self) -> String {
        let mut options = self.label_options.clone();
        options.insert("class".to_string(), "custom-control-label".to_string());

        options
            .iter()
            .map(|(key, value)| format!(" {}=\"{}\"", key, value))
            .collect::<Vec<String>>()
            .join("")
    }

    /// Render the label with help text
    fn render_label_tag_with_help_text(&self) -> String {
        let mut html = String::new();
        let label_options = self.formatted_label_options();

        html.push_str(&format!("<label for=\"{}\"{}>", self.name, label_options));

        if let Some(label) = self.label_content() {
            html.push_str(&format!("<span>{}</span>", label));
        }

        if let Some(help_text) = self.help_text_content() {
            html.push_str(&format!(
                r#"<p class="help-text" data-testid="pajamas-component-help-text">{}</p>"#,
                help_text
            ));
        }

        html.push_str("</label>");
        html
    }

    /// Render the checkbox tag component
    pub fn render(&self) -> String {
        let mut html = String::new();

        // Content wrapper
        html.push_str(r#"<div class="gl-form-checkbox custom-control custom-checkbox">"#);

        // Checkbox input
        let checked_attr = if self.checked { " checked" } else { "" };
        html.push_str(&format!(
            r#"<input type="checkbox" id="{}" name="{}" value="{}"{}>",
            self.name,
            self.name,
            self.value,
            self.formatted_input_options() + checked_attr
        ));

        // Label with help text
        html.push_str(&self.render_label_tag_with_help_text());

        html.push_str("</div>");
        html
    }
}

impl crate::components::pajamas::concerns::CheckboxRadioLabelWithHelpText for CheckboxTagComponent {
    fn render_label_with_help_text(&self) -> String {
        let label_options = self.formatted_label_options();
        let label_attrs = Self::format_attributes(&label_options);
        let label_entry = self.label_entry();

        format!(
            r#"<label for="{}"{}>{}</label>"#,
            self.method(),
            label_attrs,
            label_entry
        )
    }

    fn render_label_tag_with_help_text(&self) -> String {
        let label_options = self.formatted_label_options();
        let label_attrs = Self::format_attributes(&label_options);
        let label_entry = self.label_entry();

        format!(
            r#"<label for="{}"{}>{}</label>"#,
            self.name(),
            label_attrs,
            label_entry
        )
    }

    fn label_content(&self) -> String {
        self.label_slot.clone().unwrap_or_default()
    }

    fn help_text_content(&self) -> Option<String> {
        self.help_text_slot.clone()
    }

    fn label_options(&self) -> HashMap<String, String> {
        self.label_options.clone()
    }

    fn value(&self) -> String {
        self.value.clone()
    }

    fn method(&self) -> String {
        self.name.clone()
    }

    fn name(&self) -> String {
        self.name.clone()
    }
}

impl crate::components::pajamas::concerns::CheckboxRadioOptions for CheckboxTagComponent {
    fn input_options(&self) -> HashMap<String, String> {
        self.checkbox_options.clone()
    }
}
