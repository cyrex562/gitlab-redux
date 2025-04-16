use crate::components::Component;
use std::collections::HashMap;

/// Component for displaying page headings
pub struct PageHeadingComponent {
    heading: String,
    description: Option<String>,
    options: HashMap<String, String>,
    heading_slot: Option<String>,
    actions_slot: Option<String>,
    description_slot: Option<String>,
}

impl Component for PageHeadingComponent {
    fn render(&self) -> String {
        let heading = self.heading_slot.as_ref().unwrap_or(&self.heading);
        let description = self.description_slot.as_ref().or(self.description.as_ref());
        let actions = self.actions_slot.as_ref();

        // Convert options to HTML attributes
        let options_str = self
            .options
            .iter()
            .map(|(k, v)| format!("{}={}", k, v))
            .collect::<Vec<String>>()
            .join(" ");

        // Build the HTML
        let mut html = format!(
            r#"<div class="gl-flex gl-flex-wrap gl-items-center gl-justify-between gl-gap-y-3 gl-my-5 {}">
                <div class="gl-flex gl-flex-wrap md:gl-flex-nowrap gl-justify-between gl-gap-x-5 gl-gap-y-3 gl-w-full">
                    <h1 class="!gl-m-0 gl-heading-1" data-testid="page-heading">
                        {}
                    </h1>"#,
            options_str, heading
        );

        // Add actions if present
        if let Some(actions) = actions {
            html.push_str(&format!(
                r#"<div class="page-heading-actions gl-self-start md:gl-mt-1 lg:gl-mt-2 gl-flex gl-flex-wrap gl-items-center gl-gap-3 gl-w-full sm:gl-w-auto gl-shrink-0" data-testid="page-heading-actions">
                    {}
                </div>"#,
                actions
            ));
        }

        html.push_str("</div>");

        // Add description if present
        if let Some(description) = description {
            html.push_str(&format!(
                r#"<div class="gl-w-full gl-text-subtle" data-testid="page-heading-description">
                    {}
                </div>"#,
                description
            ));
        }

        html.push_str("</div>");

        html
    }
}

impl PageHeadingComponent {
    /// Create a new page heading component
    ///
    /// # Arguments
    ///
    /// * `heading` - The heading text
    /// * `description` - Optional description text
    /// * `options` - Optional HTML attributes
    pub fn new(
        heading: String,
        description: Option<String>,
        options: HashMap<String, String>,
    ) -> Self {
        Self {
            heading,
            description,
            options,
            heading_slot: None,
            actions_slot: None,
            description_slot: None,
        }
    }

    /// Set the heading slot content
    pub fn with_heading(mut self, heading: String) -> Self {
        self.heading_slot = Some(heading);
        self
    }

    /// Set the actions slot content
    pub fn with_actions(mut self, actions: String) -> Self {
        self.actions_slot = Some(actions);
        self
    }

    /// Set the description slot content
    pub fn with_description(mut self, description: String) -> Self {
        self.description_slot = Some(description);
        self
    }

    /// Check if actions slot is set
    pub fn has_actions(&self) -> bool {
        self.actions_slot.is_some()
    }

    /// Check if description slot is set
    pub fn has_description(&self) -> bool {
        self.description_slot.is_some() || self.description.is_some()
    }
}
