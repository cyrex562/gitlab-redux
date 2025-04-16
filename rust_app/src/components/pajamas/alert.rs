use super::button::ButtonComponent;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

/// Alert variant options
#[derive(
    Debug, Clone, Copy, PartialEq, Eq, Debug, Clone, Copy, PartialEq, Serialize, Deserialize,
)]
pub enum AlertVariant {
    Info,
    Warning,
    Success,
    Danger,
    Tip,
}

impl AlertVariant {
    fn as_str(&self) -> &'static str {
        match self {
            AlertVariant::Info => "info",
            AlertVariant::Warning => "warning",
            AlertVariant::Success => "success",
            AlertVariant::Danger => "danger",
            AlertVariant::Tip => "tip",
        }
    }

    fn icon(&self) -> &'static str {
        match self {
            AlertVariant::Info => "information-o",
            AlertVariant::Warning => "warning",
            AlertVariant::Success => "check-circle",
            AlertVariant::Danger => "error",
            AlertVariant::Tip => "bulb",
        }
    }
}

/// Renders a GlAlert root element
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AlertComponent {
    /// Alert title
    title: Option<String>,
    /// Alert variant
    variant: AlertVariant,
    /// Whether the alert is dismissible
    dismissible: bool,
    /// Whether to show the icon
    show_icon: bool,
    /// Alert options
    alert_options: HashMap<String, String>,
    /// Close button options
    close_button_options: HashMap<String, String>,
    /// Alert body content
    body: Option<String>,
    /// Alert actions content
    actions: Option<String>,
}

impl AlertComponent {
    /// Create a new alert component
    pub fn new(
        title: Option<String>,
        variant: Option<AlertVariant>,
        dismissible: bool,
        show_icon: bool,
        alert_options: Option<HashMap<String, String>>,
        close_button_options: Option<HashMap<String, String>>,
        body: Option<String>,
        actions: Option<String>,
    ) -> Self {
        Self {
            title,
            variant: variant.unwrap_or(AlertVariant::Info),
            dismissible,
            show_icon,
            alert_options: alert_options.unwrap_or_default(),
            close_button_options: close_button_options.unwrap_or_default(),
            body,
            actions,
        }
    }

    /// Render the alert component
    pub fn render(&self) -> String {
        let mut html = format!(
            r#"<div role="alert" class="gl-alert {}"{}>"#,
            self.base_class(),
            self.format_attributes(&self.alert_options)
        );

        // Render icon
        if self.show_icon {
            html.push_str(&format!(
                r#"<div class="gl-alert-icon-container">{}</div>"#,
                self.render_icon()
            ));
        }

        // Render dismiss button
        if self.dismissible {
            html.push_str(&self.render_dismiss_button());
        }

        // Render content
        html.push_str(r#"<div class="gl-alert-content" role="alert">"#);

        // Render title
        if let Some(title) = &self.title {
            html.push_str(&format!(r#"<h2 class="gl-alert-title">{}</h2>"#, title));
        }

        // Render body
        if let Some(body) = &self.body {
            html.push_str(&format!(r#"<div class="gl-alert-body">{}</div>"#, body));
        }

        // Render actions
        if let Some(actions) = &self.actions {
            html.push_str(&format!(
                r#"<div class="gl-alert-actions">{}</div>"#,
                actions
            ));
        }

        html.push_str("</div></div>");
        html
    }

    /// Generate base class string
    fn base_class(&self) -> String {
        let mut classes = vec![format!("gl-alert-{}", self.variant.as_str())];

        if !self.dismissible {
            classes.push("gl-alert-not-dismissible".to_string());
        }
        if !self.show_icon {
            classes.push("gl-alert-no-icon".to_string());
        }
        if self.title.is_some() {
            classes.push("gl-alert-has-title".to_string());
        }

        classes.join(" ")
    }

    /// Render the icon
    fn render_icon(&self) -> String {
        let icon_class = if self.title.is_none() {
            "gl-alert-icon gl-alert-icon-no-title"
        } else {
            "gl-alert-icon"
        };

        format!(
            r#"<svg class="sprite-icon {}" aria-hidden="true"><use xlink:href="/assets/icons.svg#{}"></use></svg>"#,
            icon_class,
            self.variant.icon()
        )
    }

    /// Render the dismiss button
    fn render_dismiss_button(&self) -> String {
        let mut button_options = self.close_button_options.clone();
        let mut class = "js-close gl-dismiss-btn".to_string();
        if let Some(existing_class) = button_options.get("class") {
            class = format!("{} {}", class, existing_class);
        }
        button_options.insert("class".to_string(), class);
        button_options.insert("aria-label".to_string(), "Dismiss".to_string());

        format!(
            r#"<button type="button"{}><svg class="sprite-icon" aria-hidden="true"><use xlink:href="/assets/icons.svg#close"></use></svg></button>"#,
            self.format_attributes(&button_options)
        )
    }

    /// Format HTML attributes
    fn format_attributes(&self, options: &HashMap<String, String>) -> String {
        options
            .iter()
            .map(|(key, value)| format!(" {}={}", key, value))
            .collect::<Vec<String>>()
            .join("")
    }
}
