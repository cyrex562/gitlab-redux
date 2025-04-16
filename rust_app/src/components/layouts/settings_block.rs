use crate::components::Component;
use std::collections::HashMap;

/// Component for displaying a settings block
pub struct SettingsBlockComponent {
    heading: String,
    description: Option<String>,
    id: Option<String>,
    testid: Option<String>,
    expanded: Option<bool>,
    button_options: HashMap<String, String>,
    css_class: Option<String>,
    heading_slot: Option<String>,
    description_slot: Option<String>,
    callout_slot: Option<String>,
    body_slot: Option<String>,
}

impl Component for SettingsBlockComponent {
    fn render(&self) -> String {
        let heading = self.heading_slot.as_ref().unwrap_or(&self.heading);
        let description = self.description_slot.as_ref().or(self.description.as_ref());
        let body = self.body_slot.as_ref().unwrap_or(&String::new());

        // Get section classes
        let section_classes = self.get_section_classes();

        // Get title classes
        let title_classes = self.get_title_classes();

        // Get button text and aria label
        let (button_text, aria_label) = self.get_button_text_and_aria_label();

        // Build the HTML
        let mut html = format!(
            r#"<section class="{}"{} data-testid="{}">"#,
            section_classes,
            if let Some(id) = &self.id {
                format!(" id=\"{}\"", id)
            } else {
                String::new()
            },
            self.testid.as_ref().unwrap_or(&String::new())
        );

        // Callout
        if let Some(callout) = &self.callout_slot {
            html.push_str(callout);
        }

        // Header
        html.push_str(
            r#"<div class="gl-flex gl-justify-between gl-items-start gl-gap-x-3 gl-pt-5">"#,
        );

        // Toggle button
        html.push_str(r#"<div class="gl-shrink-0 gl-px-2 gl-py-0 -gl-mr-3 sm:gl-p-2 sm:gl-mr-0">"#);
        html.push_str(&format!(
            r#"<button class="btn btn-tertiary btn-sm settings-toggle js-settings-toggle {}" aria-label="{}">"#,
            self.button_options.get("class").unwrap_or(&String::new()),
            aria_label
        ));
        html.push_str(r#"<span class="!-gl-mx-2"><img src="/assets/icons/chevron-lg-right.svg" alt="" /></span>"#);
        html.push_str(&format!(
            r#"<span class="gl-sr-only">{}</span>"#,
            button_text
        ));
        html.push_str("</button></div>");

        // Title and description
        html.push_str(r#"<div class="gl-grow">"#);
        html.push_str(&format!(
            r#"<h2 class="{}" data-settings-block-title="">{}</h2>"#,
            title_classes, heading
        ));

        if let Some(description) = description {
            html.push_str(&format!(
                r#"<p class="gl-text-subtle gl-m-0">{}</p>"#,
                description
            ));
        }

        html.push_str("</div></div>");

        // Content
        html.push_str(r#"<div class="settings-content">"#);
        html.push_str(r#"<div class="gl-pl-7 sm:gl-pl-8 gl-mt-5">"#);
        html.push_str(body);
        html.push_str("</div></div>");

        html.push_str("</section>");

        html
    }
}

impl SettingsBlockComponent {
    /// Create a new settings block component
    ///
    /// # Arguments
    ///
    /// * `heading` - The heading text
    /// * `description` - Optional description text
    /// * `id` - Optional ID for the section
    /// * `testid` - Optional test ID for testing
    /// * `expanded` - Whether the section is expanded
    /// * `button_options` - Optional HTML attributes for the button
    /// * `css_class` - Optional CSS class to add to the section
    pub fn new(
        heading: String,
        description: Option<String>,
        id: Option<String>,
        testid: Option<String>,
        expanded: Option<bool>,
        button_options: HashMap<String, String>,
        css_class: Option<String>,
    ) -> Self {
        Self {
            heading,
            description,
            id,
            testid,
            expanded,
            button_options,
            css_class,
            heading_slot: None,
            description_slot: None,
            callout_slot: None,
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

    /// Set the callout slot content
    pub fn with_callout(mut self, callout: String) -> Self {
        self.callout_slot = Some(callout);
        self
    }

    /// Set the body slot content
    pub fn with_body(mut self, body: String) -> Self {
        self.body_slot = Some(body);
        self
    }

    /// Check if callout slot is set
    pub fn has_callout(&self) -> bool {
        self.callout_slot.is_some()
    }

    /// Get section classes
    fn get_section_classes(&self) -> String {
        let mut classes = vec!["settings".to_string(), "no-animate".to_string()];

        if let Some(expanded) = self.expanded {
            if expanded {
                classes.push("expanded".to_string());
            }
        }

        if let Some(css_class) = &self.css_class {
            classes.push(css_class.clone());
        }

        classes.join(" ")
    }

    /// Get title classes
    fn get_title_classes(&self) -> String {
        vec![
            "gl-heading-2".to_string(),
            "gl-cursor-pointer".to_string(),
            "!gl-mb-2".to_string(),
            "js-settings-toggle".to_string(),
            "js-settings-toggle-trigger-only".to_string(),
        ]
        .join(" ")
    }

    /// Get button text and aria label
    fn get_button_text_and_aria_label(&self) -> (String, String) {
        let expanded = self.expanded.unwrap_or(false);
        let button_text = if expanded { "Collapse" } else { "Expand" };
        let aria_label = if expanded {
            format!("Collapse {}", self.heading)
        } else {
            format!("Expand {}", self.heading)
        };

        (button_text.to_string(), aria_label)
    }
}
