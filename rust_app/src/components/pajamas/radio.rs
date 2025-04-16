use serde::{Deserialize, Serialize};
use std::collections::HashMap;

/// Renders a Pajamas compliant radio button element
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct RadioComponent {
    /// The form builder
    form: String,
    /// The method name for the form field
    method: String,
    /// The label text
    label: Option<String>,
    /// The help text
    help_text: Option<String>,
    /// Options for the label
    label_options: HashMap<String, String>,
    /// Options for the radio button
    radio_options: HashMap<String, String>,
    /// The value of the radio button
    value: Option<String>,
    /// The label content slot
    label_slot: Option<String>,
    /// The help text content slot
    help_text_slot: Option<String>,
}

impl RadioComponent {
    /// Create a new radio component
    pub fn new(
        form: String,
        method: String,
        label: Option<String>,
        help_text: Option<String>,
        label_options: Option<HashMap<String, String>>,
        radio_options: Option<HashMap<String, String>>,
        value: Option<String>,
    ) -> Self {
        Self {
            form,
            method,
            label,
            help_text,
            label_options: label_options.unwrap_or_default(),
            radio_options: radio_options.unwrap_or_default(),
            value,
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
        self.label_slot.clone().or_else(|| self.label.clone())
    }

    /// Get the help text content
    fn help_text_content(&self) -> Option<String> {
        self.help_text_slot
            .clone()
            .or_else(|| self.help_text.clone())
    }

    /// Format the input options
    fn formatted_input_options(&self) -> String {
        let mut options = self.radio_options.clone();
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

    /// Render the radio component
    pub fn render(&self) -> String {
        let mut html = String::new();

        // Content wrapper
        html.push_str(r#"<div class="gl-form-radio custom-control custom-radio">"#);

        // Radio input
        let value = self.value.clone().unwrap_or_default();
        html.push_str(&format!(
            r#"<input type="radio" id="{}" name="{}" value="{}"{}>"#,
            self.method,
            self.method,
            value,
            self.formatted_input_options()
        ));

        // Label with help text
        html.push_str(&self.render_label_with_help_text());

        html.push_str("</div>");
        html
    }

    /// Render the label with help text
    fn render_label_with_help_text(&self) -> String {
        let mut html = String::new();
        let label_options = self.formatted_label_options();

        html.push_str(&format!("<label for=\"{}\"{}>", self.method, label_options));

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
}

impl crate::components::pajamas::concerns::CheckboxRadioLabelWithHelpText for RadioComponent {
    fn render_label_with_help_text(&self) -> String {
        let label_options = self.formatted_label_options();
        let label_attrs = Self::format_attributes(&label_options);
        let label_entry = self.label_entry();

        format!(
            r#"<label for="{}"{}>{}</label>"#,
            self.method, label_attrs, label_entry
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
        self.label.clone().unwrap_or_default()
    }

    fn help_text_content(&self) -> Option<String> {
        self.help_text.clone()
    }

    fn label_options(&self) -> HashMap<String, String> {
        self.label_options.clone()
    }

    fn value(&self) -> String {
        self.value.clone().unwrap_or_default()
    }

    fn method(&self) -> String {
        self.method.clone()
    }

    fn name(&self) -> String {
        self.method.clone()
    }
}

impl crate::components::pajamas::concerns::CheckboxRadioOptions for RadioComponent {
    fn input_options(&self) -> HashMap<String, String> {
        self.radio_options.clone()
    }
}
