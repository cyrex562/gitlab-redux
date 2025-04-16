use super::button::ButtonComponent;
use std::collections::HashMap;

/// State options for the accordion item
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum AccordionItemState {
    Opened,
    Closed,
}

impl AccordionItemState {
    fn as_str(&self) -> &'static str {
        match self {
            AccordionItemState::Opened => "opened",
            AccordionItemState::Closed => "closed",
        }
    }
}

/// Renders a Pajamas compliant accordion item element
#[derive(Debug, Clone)]
pub struct AccordionItemComponent {
    /// The title of the accordion item
    title: Option<String>,
    /// The state of the accordion item
    state: AccordionItemState,
    /// Options for the button
    button_options: HashMap<String, String>,
    /// The content of the accordion item
    content: Option<String>,
}

impl AccordionItemComponent {
    /// Create a new accordion item component
    pub fn new(
        title: Option<String>,
        state: AccordionItemState,
        button_options: Option<HashMap<String, String>>,
    ) -> Self {
        Self {
            title,
            state,
            button_options: button_options.unwrap_or_default(),
            content: None,
        }
    }

    /// Set the content of the accordion item
    pub fn with_content(mut self, content: String) -> Self {
        self.content = Some(content);
        self
    }

    /// Get the icon for the accordion item
    fn icon(&self) -> &'static str {
        match self.state {
            AccordionItemState::Opened => "chevron-down",
            AccordionItemState::Closed => "chevron-right",
        }
    }

    /// Get the body class for the accordion item
    fn body_class(&self) -> HashMap<String, String> {
        let mut classes = HashMap::new();
        if self.state == AccordionItemState::Opened {
            classes.insert("class".to_string(), "show".to_string());
        }
        classes
    }

    /// Check if the accordion item is expanded
    fn expanded(&self) -> bool {
        self.state == AccordionItemState::Opened
    }

    /// Render the accordion item component
    pub fn render(&self) -> String {
        let mut html = String::new();

        // Start accordion item
        html.push_str(r#"<div class="gl-accordion-item">"#);

        // Render header
        html.push_str(r#"<h3 class="gl-accordion-item-header">"#);

        // Render button
        if let Some(title) = &self.title {
            let mut button = ButtonComponent::new();
            button.variant = super::button::ButtonVariant::Link;
            button.icon = Some(self.icon().to_string());
            button.icon_classes = Some("js-chevron-icon".to_string());
            button.content = Some(title.clone());

            // Add button options
            button.button_options = self.button_options.clone();
            button
                .button_options
                .insert("aria-controls".to_string(), "accordion-item".to_string());
            button
                .button_options
                .insert("aria-expanded".to_string(), self.expanded().to_string());

            html.push_str(&button.render());
        }

        html.push_str("</h3>");

        // Render body
        let body_classes = self.body_class();
        let body_class_str = body_classes
            .get("class")
            .map(|c| format!(" {}", c))
            .unwrap_or_default();

        html.push_str(&format!(
            r#"<div class="accordion-item gl-mt-3 gl-text-base collapse{}">"#,
            body_class_str
        ));

        // Render content
        if let Some(content) = &self.content {
            html.push_str(content);
        }

        html.push_str("</div></div>");
        html
    }
}
