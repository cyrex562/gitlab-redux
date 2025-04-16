use crate::components::Component;
use std::collections::HashMap;

/// Component for displaying a settings section
pub struct SettingsSectionComponent {
    heading: String,
    description: Option<String>,
    id: Option<String>,
    testid: Option<String>,
    options: HashMap<String, String>,
    heading_slot: Option<String>,
    description_slot: Option<String>,
    body_slot: Option<String>,
}

impl Component for SettingsSectionComponent {
    fn render(&self) -> String {
        let heading = self.heading_slot.as_ref().unwrap_or(&self.heading);
        let description = self.description_slot.as_ref().or(self.description.as_ref());
        let body = self.body_slot.as_ref().unwrap_or(&String::new());

        // Format options attributes
        let options_attrs = self.format_options_attrs();

        // Build the HTML
        let mut html = format!(
            r#"<section class="settings-section js-search-settings-section" {}>",
            options_attrs
        );
        
        // Header
        html.push_str(r#"<div class="settings-sticky-header">"#
        );
        html.push_str(r#"<div class="settings-sticky-header-inner">"#);

        // Heading
        html.push_str(&format!(
            r#"<h2 class="gl-heading-2 !gl-mb-3">{}</h2>"#,
            heading
        ));

        // Description
        if let Some(description) = description {
            html.push_str(&format!(
                r#"<p class="gl-text-subtle gl-mb-3">{}</p>"#,
                description
            ));
        }

        html.push_str("</div></div>");

        // Body
        html.push_str(&format!(
            r#"<div class="gl-mt-3" data-testid="settings-section-body">{}</div>"#,
            body
        ));

        html.push_str("</section>");

        html
    }
}

impl SettingsSectionComponent {
    /// Create a new settings section component
    ///
    /// # Arguments
    ///
    /// * `heading` - The heading text
    /// * `description` - Optional description text
    /// * `id` - Optional ID for the section
    /// * `testid` - Optional test ID for testing
    /// * `options` - Optional HTML attributes
    pub fn new(
        heading: String,
        description: Option<String>,
        id: Option<String>,
        testid: Option<String>,
        options: HashMap<String, String>,
    ) -> Self {
        Self {
            heading,
            description,
            id,
            testid,
            options,
            heading_slot: None,
            description_slot: None,
            body_slot: None,
        }
    }

    /// Set the heading slot content
    pub fn with_heading(mut self, heading: String) -> Self {
        self.heading_slot = Some(heading);
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

    /// Format options attributes
    fn format_options_attrs(&self) -> String {
        let mut attrs = self.options.clone();

        // Handle data attributes
        let mut data = HashMap::new();

        // Add testid to data if provided
        if let Some(testid) = &self.testid {
            data.insert("testid".to_string(), testid.clone());
        }

        // Add other data attributes from options
        if let Some(data_attrs) = self.options.get("data") {
            // This is a simplified approach - in a real implementation,
            // you would parse the data attributes from the options
            data.insert("data".to_string(), data_attrs.clone());
        }

        // Add data attributes to attrs
        if !data.is_empty() {
            attrs.insert(
                "data".to_string(),
                data.iter()
                    .map(|(k, v)| format!("{}={}", k, v))
                    .collect::<Vec<String>>()
                    .join(" "),
            );
        }

        // Add ID if provided
        if let Some(id) = &self.id {
            attrs.insert("id".to_string(), id.clone());
        }

        // Format attributes
        attrs
            .iter()
            .map(|(k, v)| format!("{}={}", k, v))
            .collect::<Vec<String>>()
            .join(" ")
    }
}
