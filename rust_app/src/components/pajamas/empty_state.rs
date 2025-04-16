use super::button::ButtonComponent;
use std::collections::HashMap;

/// Renders a Pajamas compliant empty state
#[derive(Debug, Clone)]
pub struct EmptyStateComponent {
    /// Whether the empty state is compact
    compact: bool,
    /// The title of the empty state
    title: Option<String>,
    /// The SVG path for the illustration
    svg_path: Option<String>,
    /// The primary button text
    primary_button_text: Option<String>,
    /// The primary button link
    primary_button_link: Option<String>,
    /// The primary button options
    primary_button_options: HashMap<String, String>,
    /// The secondary button text
    secondary_button_text: Option<String>,
    /// The secondary button link
    secondary_button_link: Option<String>,
    /// The secondary button options
    secondary_button_options: HashMap<String, String>,
    /// The empty state options
    empty_state_options: HashMap<String, String>,
    /// The description content slot
    description_slot: Option<String>,
}

impl EmptyStateComponent {
    /// Create a new empty state component
    pub fn new(
        compact: bool,
        title: Option<String>,
        svg_path: Option<String>,
        primary_button_text: Option<String>,
        primary_button_link: Option<String>,
        primary_button_options: Option<HashMap<String, String>>,
        secondary_button_text: Option<String>,
        secondary_button_link: Option<String>,
        secondary_button_options: Option<HashMap<String, String>>,
        empty_state_options: Option<HashMap<String, String>>,
    ) -> Self {
        Self {
            compact,
            title,
            svg_path,
            primary_button_text,
            primary_button_link,
            primary_button_options: primary_button_options.unwrap_or_default(),
            secondary_button_text,
            secondary_button_link,
            secondary_button_options: secondary_button_options.unwrap_or_default(),
            empty_state_options: empty_state_options.unwrap_or_default(),
            description_slot: None,
        }
    }

    /// Set the description content slot
    pub fn with_description_slot(mut self, description: String) -> Self {
        self.description_slot = Some(description);
        self
    }

    /// Get the empty state class
    fn empty_state_class(&self) -> String {
        if self.compact {
            "gl-flex-row".to_string()
        } else {
            "gl-text-center gl-flex-col".to_string()
        }
    }

    /// Get the image class
    fn image_class(&self) -> String {
        if self.compact {
            "gl-hidden sm:gl-block gl-px-4".to_string()
        } else {
            "gl-max-w-full".to_string()
        }
    }

    /// Get the content wrapper class
    fn content_wrapper_class(&self) -> String {
        if self.compact {
            "gl-grow gl-basis-0 gl-px-4".to_string()
        } else {
            "gl-m-auto gl-p-5".to_string()
        }
    }

    /// Get the title class
    fn title_class(&self) -> String {
        if self.compact {
            "h5".to_string()
        } else {
            "h4".to_string()
        }
    }

    /// Get the button wrapper class
    fn button_wrapper_class(&self) -> String {
        if self.compact {
            "".to_string()
        } else {
            "gl-justify-center".to_string()
        }
    }

    /// Format the empty state options
    fn formatted_empty_state_options(&self) -> String {
        self.empty_state_options
            .iter()
            .map(|(key, value)| format!(" {}=\"{}\"", key, value))
            .collect::<Vec<String>>()
            .join("")
    }

    /// Render the primary button
    fn render_primary_button(&self) -> Option<String> {
        if let (Some(text), Some(link)) = (&self.primary_button_text, &self.primary_button_link) {
            let mut options = self.primary_button_options.clone();
            options.insert("class".to_string(), "!gl-ml-0".to_string());

            let options_str = options
                .iter()
                .map(|(key, value)| format!(" {}=\"{}\"", key, value))
                .collect::<Vec<String>>()
                .join("");

            Some(format!(
                r#"<a href="{}" class="btn btn-confirm"{}>{}</a>"#,
                link, options_str, text
            ))
        } else {
            None
        }
    }

    /// Render the secondary button
    fn render_secondary_button(&self) -> Option<String> {
        if let (Some(text), Some(link)) = (&self.secondary_button_text, &self.secondary_button_link)
        {
            let mut options = self.secondary_button_options.clone();

            // Add !gl-ml-0 class if primary button is not present
            if self.primary_button_text.is_none() {
                options.insert("class".to_string(), "!gl-ml-0".to_string());
            }

            let options_str = options
                .iter()
                .map(|(key, value)| format!(" {}=\"{}\"", key, value))
                .collect::<Vec<String>>()
                .join("");

            Some(format!(
                r#"<a href="{}" class="btn btn-default"{}>{}</a>"#,
                link, options_str, text
            ))
        } else {
            None
        }
    }

    /// Render the empty state component
    pub fn render(&self) -> String {
        let mut html = String::new();
        let empty_state_class = self.empty_state_class();
        let empty_state_options = self.formatted_empty_state_options();

        // Empty state container
        html.push_str(&format!(
            r#"<section class="gl-flex gl-empty-state {}"{}>"#,
            empty_state_class, empty_state_options
        ));

        // SVG illustration
        if let Some(svg_path) = &self.svg_path {
            let image_class = self.image_class();
            html.push_str(&format!(
                r#"<div class="{}"><img src="{}" alt="" class="gl-dark-invert-keep-hue"></div>"#,
                image_class, svg_path
            ));
        }

        // Content wrapper
        let content_wrapper_class = self.content_wrapper_class();
        html.push_str(&format!(
            r#"<div class="gl-empty-state-content gl-mx-auto gl-my-0 {}">"#,
            content_wrapper_class
        ));

        // Title
        if let Some(title) = &self.title {
            let title_class = self.title_class();
            html.push_str(&format!(
                r#"<h1 class="gl-text-size-h-display gl-leading-36 gl-mt-0 gl-mb-0 {}">{}</h1>"#,
                title_class, title
            ));
        }

        // Description
        if let Some(description) = &self.description_slot {
            html.push_str(&format!(
                r#"<p class="gl-mt-4 gl-mb-0" data-testid="empty-state-description">{}</p>"#,
                description
            ));
        }

        // Buttons
        if self.primary_button_text.is_some() || self.secondary_button_text.is_some() {
            let button_wrapper_class = self.button_wrapper_class();
            html.push_str(&format!(
                r#"<div class="gl-flex gl-flex-wrap gl-mt-5 gl-gap-3 {}">"#,
                button_wrapper_class
            ));

            if let Some(primary_button) = self.render_primary_button() {
                html.push_str(&primary_button);
            }

            if let Some(secondary_button) = self.render_secondary_button() {
                html.push_str(&secondary_button);
            }

            html.push_str("</div>");
        }

        html.push_str("</div></section>");
        html
    }
}
