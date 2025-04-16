use std::collections::HashMap;
use serde::{Serialize, Deserialize};

/// Renders a Pajamas compliant radio button element without a form builder
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct RadioTagComponent {
    /// The name of the radio button
    name: String,
    /// The value of the radio button
    value: String,
    /// Whether the radio button is checked
    checked: bool,
    /// The label text
    label: Option<String>,
    /// The help text
    help_text: Option<String>,
    /// Options for the label
    label_options: HashMap<String, String>,
    /// Options for the radio input
    radio_options: HashMap<String, String>,
    /// The label content slot
    label_slot: Option<String>,
    /// The help text content slot
    help_text_slot: Option<String>,
}

impl RadioTagComponent {
    /// Create a new radio tag component
    pub fn new(
        name: String,
        value: String,
        checked: bool,
        label: Option<String>,
        help_text: Option<String>,
        label_options: Option<HashMap<String, String>>,
        radio_options: Option<HashMap<String, String>>,
    ) -> Self {
        let mut label_opts = label_options.unwrap_or_default();
        
        // Set the 'for' attribute if not already set
        if !label_opts.contains_key("for") {
            label_opts.insert("for".to_string(), Self::label_for(&name, &value));
        }

        Self {
            name,
            value,
            checked,
            label,
            help_text,
            label_options: label_opts,
            radio_options: radio_options.unwrap_or_default(),
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
        self.help_text_slot.clone().or_else(|| self.help_text.clone())
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
        self.label_options
            .iter()
            .map(|(key, value)| format!(" {}=\"{}\"", key, value))
            .collect::<Vec<String>>()
            .join("")
    }

    /// Render the label with help text
    fn render_label_tag_with_help_text(&self) -> String {
        let mut html = String::new();
        let label_options = self.formatted_label_options();

        html.push_str(&format!("<label{}>", label_options));

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

    /// Generate a label for attribute
    fn label_for(name: &str, value: &str) -> String {
        format!("{}_{}", Self::sanitize_to_id(name), value)
    }

    /// Sanitize a string to be used as an ID
    fn sanitize_to_id(name: &str) -> String {
        name.replace(|c: char| !c.is_alphanumeric(), "_")
    }

    /// Render the radio tag component
    pub fn render(&self) -> String {
        let mut html = String::new();

        // Content wrapper
        html.push_str(r#"<div class="gl-form-radio custom-control custom-radio">"#);

        // Radio input
        let checked_attr = if self.checked { " checked" } else { "" };
        html.push_str(&format!(
            r#"<input type="radio" id="{}" name="{}" value="{}"{}>",
            Self::label_for(&self.name, &self.value),
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

impl crate::components::pajamas::concerns::CheckboxRadioLabelWithHelpText for RadioTagComponent {
    fn render_label_with_help_text(&self) -> String {
        let label_options = self.formatted_label_options();
        let label_attrs = Self::format_attributes(&label_options);
        let label_entry = self.label_entry();

        format!(
            r#"<label for="{}"{}>{}</label>"#,
            self.method(), label_attrs, label_entry
        )
    }

    fn render_label_tag_with_help_text(&self) -> String {
        let label_options = self.formatted_label_options();
        let label_attrs = Self::format_attributes(&label_options);
        let label_entry = self.label_entry();

        format!(
            r#"<label for="{}"{}>{}</label>"#,
            self.name(), label_attrs, label_entry
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
        self.value.clone()
    }

    fn method(&self) -> String {
        self.name.clone()
    }

    fn name(&self) -> String {
        self.name.clone()
    }
}

impl crate::components::pajamas::concerns::CheckboxRadioOptions for RadioTagComponent {
    fn input_options(&self) -> HashMap<String, String> {
        self.radio_options.clone()
    }
} 