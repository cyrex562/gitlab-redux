use crate::components::Component;
use std::collections::HashMap;

/// Component for displaying a horizontal section with title, description, and body
pub struct HorizontalSectionComponent {
    border: bool,
    options: HashMap<String, String>,
    title_slot: Option<String>,
    description_slot: Option<String>,
    body_slot: Option<String>,
}

impl Component for HorizontalSectionComponent {
    fn render(&self) -> String {
        let title = self.title_slot.as_ref();
        let description = self.description_slot.as_ref();
        let body = self.body_slot.as_ref();

        // Format options with border class if needed
        let mut formatted_options = self.options.clone();

        // Handle class attribute specially
        let mut classes = Vec::new();
        if self.border {
            classes.push("gl-border-b");
        }

        if let Some(class) = self.options.get("class") {
            classes.push(class);
        }

        if !classes.is_empty() {
            formatted_options.insert("class".to_string(), classes.join(" "));
        }

        // Convert options to HTML attributes
        let options_str = formatted_options
            .iter()
            .map(|(k, v)| format!("{}={}", k, v))
            .collect::<Vec<String>>()
            .join(" ");

        // Build the HTML
        let mut html = format!(
            r#"<div {}>",
            options_str
        );
        
        html.push_str(r#"<div class="row">"#
        );

        // Left column with title and description
        html.push_str(r#"<div class="col-lg-4">"#);

        if let Some(title) = title {
            html.push_str(&format!(r#"<h4 class="gl-mt-0">{}</h4>"#, title));
        }

        if let Some(description) = description {
            html.push_str(&format!(r#"<p>{}</p>"#, description));
        }

        html.push_str("</div>");

        // Right column with body
        html.push_str(r#"<div class="col-lg-8">"#);

        if let Some(body) = body {
            html.push_str(body);
        }

        html.push_str("</div></div></div>");

        html
    }
}

impl HorizontalSectionComponent {
    /// Create a new horizontal section component
    ///
    /// # Arguments
    ///
    /// * `border` - Whether to show a border at the bottom
    /// * `options` - Optional HTML attributes
    pub fn new(border: bool, options: HashMap<String, String>) -> Self {
        Self {
            border,
            options,
            title_slot: None,
            description_slot: None,
            body_slot: None,
        }
    }

    /// Set the title slot content
    pub fn with_title(mut self, title: String) -> Self {
        self.title_slot = Some(title);
        self
    }

    /// Set the description slot content
    pub fn with_description(mut self, description: String) -> Self {
        self.description_slot = Some(description);
        self
    }

    /// Set the body slot content
    pub fn with_body(mut self, body: String) -> Self {
        self.body_slot = Some(body);
        self
    }

    /// Check if description slot is set
    pub fn has_description(&self) -> bool {
        self.description_slot.is_some()
    }
}
