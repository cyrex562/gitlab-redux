use serde::{Deserialize, Serialize};
use std::collections::HashMap;

/// Variant options for the badge component
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum BadgeVariant {
    Muted,
    Neutral,
    Info,
    Success,
    Warning,
    Danger,
    Tier,
}

impl BadgeVariant {
    fn as_str(&self) -> &'static str {
        match self {
            BadgeVariant::Muted => "muted",
            BadgeVariant::Neutral => "neutral",
            BadgeVariant::Info => "info",
            BadgeVariant::Success => "success",
            BadgeVariant::Warning => "warning",
            BadgeVariant::Danger => "danger",
            BadgeVariant::Tier => "tier",
        }
    }
}

/// Renders a Pajamas compliant badge element
#[derive(Debug, Clone)]
pub struct BadgeComponent {
    /// The text content of the badge
    text: Option<String>,
    /// The icon name
    icon: Option<String>,
    /// Additional CSS classes for the icon
    icon_classes: Vec<String>,
    /// Whether to show only the icon
    icon_only: bool,
    /// The href for the badge link
    href: Option<String>,
    /// The variant of the badge
    variant: BadgeVariant,
    /// Additional HTML options
    html_options: HashMap<String, String>,
    /// The content slot
    content: Option<String>,
}

impl BadgeComponent {
    /// Create a new badge component
    pub fn new(
        text: Option<String>,
        icon: Option<String>,
        icon_classes: Option<Vec<String>>,
        icon_only: bool,
        href: Option<String>,
        variant: BadgeVariant,
        html_options: Option<HashMap<String, String>>,
    ) -> Self {
        Self {
            text,
            icon,
            icon_classes: icon_classes.unwrap_or_default(),
            icon_only,
            href,
            variant,
            html_options: html_options.unwrap_or_default(),
            content: None,
        }
    }

    /// Set the content slot
    pub fn with_content(mut self, content: String) -> Self {
        self.content = Some(content);
        self
    }

    /// Get the badge classes
    fn badge_classes(&self) -> String {
        let mut classes = vec![
            "gl-badge".to_string(),
            "badge".to_string(),
            "badge-pill".to_string(),
            format!("badge-{}", self.variant.as_str()),
        ];

        if self.icon_only {
            classes.push("!gl-px-2".to_string());
        }

        classes.join(" ")
    }

    /// Get the icon classes
    fn icon_classes(&self) -> String {
        let mut classes = vec!["gl-icon".to_string(), "gl-badge-icon".to_string()];
        classes.extend(self.icon_classes.clone());

        if self.circular_icon() {
            classes.push("-gl-ml-2".to_string());
        }

        classes.join(" ")
    }

    /// Check if the badge is icon-only
    fn icon_only(&self) -> bool {
        self.icon_only
    }

    /// Check if the badge is a link
    fn is_link(&self) -> bool {
        self.href.is_some()
    }

    /// Get the text content
    fn text(&self) -> Option<String> {
        self.content.clone().or_else(|| self.text.clone())
    }

    /// Check if the badge has an icon
    fn has_icon(&self) -> bool {
        self.icon_only || self.icon.is_some()
    }

    /// Check if the icon is circular
    fn circular_icon(&self) -> bool {
        if let Some(icon) = &self.icon {
            matches!(icon.as_str(), "issue-open-m" | "issue-close")
        } else {
            false
        }
    }

    /// Get the HTML options
    fn html_options(&self) -> HashMap<String, String> {
        let mut options = self.html_options.clone();
        options.insert("class".to_string(), self.badge_classes());

        if self.icon_only {
            if let Some(text) = self.text() {
                options.insert("aria-label".to_string(), text);
                options.insert("role".to_string(), "img".to_string());
            }
        }

        options
    }

    /// Render the badge component
    pub fn render(&self) -> String {
        let mut html = String::new();
        let options = self.html_options();
        let options_str = options
            .iter()
            .map(|(key, value)| format!(" {}=\"{}\"", key, value))
            .collect::<Vec<String>>()
            .join("");

        if self.is_link() {
            html.push_str(&format!(
                r#"<a href="{}"{}>"#,
                self.href.as_ref().unwrap(),
                options_str
            ));
        } else {
            html.push_str(&format!("<span{}>", options_str));
        }

        if self.has_icon() {
            html.push_str(&format!(
                r#"<svg class="sprite-icon {}" aria-hidden="true"><use xlink:href="/assets/icons.svg#{}"></use></svg>"#,
                self.icon_classes(),
                self.icon.as_ref().unwrap()
            ));
        }

        if let Some(text) = self.text() {
            html.push_str(&format!(
                r#"<span class="gl-badge-content">{}</span>"#,
                text
            ));
        }

        if self.is_link() {
            html.push_str("</a>");
        } else {
            html.push_str("</span>");
        }

        html
    }
}
