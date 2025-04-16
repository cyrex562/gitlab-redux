use super::button::ButtonComponent;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

/// Banner variant options
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum BannerVariant {
    Introduction,
    Promotion,
}

impl BannerVariant {
    fn as_str(&self) -> &'static str {
        match self {
            BannerVariant::Introduction => "introduction",
            BannerVariant::Promotion => "promotion",
        }
    }
}

/// Renders a Pajamas compliant banner element
#[derive(Debug, Clone)]
pub struct BannerComponent {
    /// The text for the button
    button_text: String,
    /// The link for the button
    button_link: String,
    /// The variant of the banner
    variant: BannerVariant,
    /// The path to the SVG image
    svg_path: Option<String>,
    /// Options for the banner
    banner_options: HashMap<String, String>,
    /// Options for the button
    button_options: HashMap<String, String>,
    /// Options for the close button
    close_options: HashMap<String, String>,
    /// The title content
    title: Option<String>,
    /// The illustration content
    illustration: Option<String>,
    /// The primary action content
    primary_action: Option<String>,
    /// The actions content
    actions: Vec<String>,
    /// The main content
    content: Option<String>,
}

impl BannerComponent {
    /// Create a new banner component
    pub fn new(
        button_text: Option<String>,
        button_link: Option<String>,
        variant: BannerVariant,
        svg_path: Option<String>,
        banner_options: Option<HashMap<String, String>>,
        button_options: Option<HashMap<String, String>>,
        close_options: Option<HashMap<String, String>>,
    ) -> Self {
        Self {
            button_text: button_text.unwrap_or_else(|| "OK".to_string()),
            button_link: button_link.unwrap_or_else(|| "#".to_string()),
            variant,
            svg_path,
            banner_options: banner_options.unwrap_or_default(),
            button_options: button_options.unwrap_or_default(),
            close_options: Self::format_close_options(close_options),
            title: None,
            illustration: None,
            primary_action: None,
            actions: Vec::new(),
            content: None,
        }
    }

    /// Format the close button options
    fn format_close_options(options: Option<HashMap<String, String>>) -> HashMap<String, String> {
        let mut close_options = options.unwrap_or_default();
        close_options.insert("class".to_string(), "js-close gl-banner-close".to_string());
        close_options
    }

    /// Set the title content
    pub fn with_title(mut self, title: String) -> Self {
        self.title = Some(title);
        self
    }

    /// Set the illustration content
    pub fn with_illustration(mut self, illustration: String) -> Self {
        self.illustration = Some(illustration);
        self
    }

    /// Set the primary action content
    pub fn with_primary_action(mut self, primary_action: String) -> Self {
        self.primary_action = Some(primary_action);
        self
    }

    /// Add an action
    pub fn with_action(mut self, action: String) -> Self {
        self.actions.push(action);
        self
    }

    /// Set the main content
    pub fn with_content(mut self, content: String) -> Self {
        self.content = Some(content);
        self
    }

    /// Get the banner class
    fn banner_class(&self) -> String {
        if self.is_introduction() {
            "gl-banner-introduction".to_string()
        } else {
            String::new()
        }
    }

    /// Get the close button variant
    fn close_button_variant(&self) -> super::button::ButtonVariant {
        if self.is_introduction() {
            super::button::ButtonVariant::Confirm
        } else {
            super::button::ButtonVariant::Default
        }
    }

    /// Check if the banner is an introduction
    fn is_introduction(&self) -> bool {
        self.variant == BannerVariant::Introduction
    }

    /// Render the banner component
    pub fn render(&self) -> String {
        let mut html = String::new();

        // Start banner
        let banner_class = self.banner_class();
        let banner_class_str = if banner_class.is_empty() {
            String::new()
        } else {
            format!(" {}", banner_class)
        };

        let mut banner_attrs = self.banner_options.clone();
        banner_attrs.insert(
            "class".to_string(),
            format!(
                "gl-banner gl-card gl-pl-6 gl-pr-8 gl-py-6{}",
                banner_class_str
            ),
        );

        let banner_attrs_str = banner_attrs
            .iter()
            .map(|(key, value)| format!(" {}=\"{}\"", key, value))
            .collect::<Vec<String>>()
            .join("");

        html.push_str(&format!("<div{}>", banner_attrs_str));

        // Banner body
        html.push_str(r#"<div class="gl-card-body gl-flex !gl-p-0">"#);

        // Illustration
        if let Some(illustration) = &self.illustration {
            html.push_str(r#"<div class="gl-banner-illustration">"#);
            html.push_str(illustration);
            html.push_str("</div>");
        } else if let Some(svg_path) = &self.svg_path {
            html.push_str(&format!(
                r#"<div class="gl-banner-illustration"><img src="{}" alt=""></div>"#,
                svg_path
            ));
        }

        // Content
        html.push_str(r#"<div class="gl-banner-content">"#);

        // Title
        if let Some(title) = &self.title {
            html.push_str(&format!(r#"<h2 class="gl-banner-title">{}</h2>"#, title));
        }

        // Main content
        if let Some(content) = &self.content {
            html.push_str(content);
        }

        // Primary action
        if let Some(primary_action) = &self.primary_action {
            html.push_str(primary_action);
        } else {
            let mut button = ButtonComponent::new();
            button.variant = super::button::ButtonVariant::Confirm;
            button.href = Some(self.button_link.clone());
            button.content = Some(self.button_text.clone());
            button.button_options = self.button_options.clone();
            button
                .button_options
                .insert("class".to_string(), "js-close-callout".to_string());
            html.push_str(&button.render());
        }

        // Actions
        for action in &self.actions {
            html.push_str(action);
        }

        html.push_str("</div></div>");

        // Close button
        let mut close_button = ButtonComponent::new();
        close_button.category = super::button::ButtonCategory::Tertiary;
        close_button.variant = self.close_button_variant();
        close_button.size = super::button::ButtonSize::Small;
        close_button.icon = Some("close".to_string());
        close_button.button_options = self.close_options.clone();
        html.push_str(&close_button.render());

        html.push_str("</div>");
        html
    }
}
