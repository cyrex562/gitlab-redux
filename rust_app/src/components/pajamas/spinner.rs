use serde::{Deserialize, Serialize};
use std::collections::HashMap;

/// Color options for the spinner component
#[derive(Debug, Clone, Copy, PartialEq, Serialize, Deserialize)]
pub enum SpinnerColor {
    Light,
    Dark,
}

impl SpinnerColor {
    fn as_str(&self) -> &'static str {
        match self {
            SpinnerColor::Light => "light",
            SpinnerColor::Dark => "dark",
        }
    }
}

/// Size options for the spinner component
#[derive(Debug, Clone, Copy, PartialEq, Serialize, Deserialize)]
pub enum SpinnerSize {
    Sm,
    Md,
    Lg,
    Xl,
}

impl SpinnerSize {
    fn as_str(&self) -> &'static str {
        match self {
            SpinnerSize::Sm => "sm",
            SpinnerSize::Md => "md",
            SpinnerSize::Lg => "lg",
            SpinnerSize::Xl => "xl",
        }
    }
}

/// Renders a Pajamas compliant spinner component
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SpinnerComponent {
    /// The color of the spinner
    color: SpinnerColor,
    /// Whether the spinner is inline
    inline: bool,
    /// The label for the spinner
    label: Option<String>,
    /// The size of the spinner
    size: SpinnerSize,
    /// HTML options
    html_options: HashMap<String, String>,
}

impl SpinnerComponent {
    /// Create a new spinner component
    pub fn new(
        color: Option<SpinnerColor>,
        inline: bool,
        label: Option<String>,
        size: Option<SpinnerSize>,
        html_options: Option<HashMap<String, String>>,
    ) -> Self {
        Self {
            color: color.unwrap_or(SpinnerColor::Dark),
            inline,
            label: label.or_else(|| Some("Loading".to_string())),
            size: size.unwrap_or(SpinnerSize::Sm),
            html_options: html_options.unwrap_or_default(),
        }
    }

    /// Get the spinner class
    fn spinner_class(&self) -> String {
        format!(
            "gl-spinner gl-spinner-{} gl-spinner-{} !gl-align-text-bottom",
            self.size.as_str(),
            self.color.as_str()
        )
    }

    /// Get the HTML options
    fn html_options(&self) -> HashMap<String, String> {
        let mut options = self.html_options.clone();

        // Add default classes
        let container_class = "gl-spinner-container";
        if let Some(existing_class) = options.get("class") {
            options.insert(
                "class".to_string(),
                format!("{} {}", existing_class, container_class),
            );
        } else {
            options.insert("class".to_string(), container_class.to_string());
        }

        // Add role
        options.insert("role".to_string(), "status".to_string());

        options
    }

    /// Render the spinner component
    pub fn render(&self) -> String {
        let tag = if self.inline { "span" } else { "div" };
        let html_options = self.html_options();
        let html_options_str = html_options
            .iter()
            .map(|(key, value)| format!(" {}={}", key, value))
            .collect::<String>();

        let mut html = format!("<{}{}>", tag, html_options_str);

        // Add spinner
        html.push_str(&format!(
            r#"<span class="{}" aria-hidden="true"></span>"#,
            self.spinner_class()
        ));

        // Add label
        if let Some(label) = &self.label {
            html.push_str(&format!(
                r#"<span class="gl-sr-only !gl-absolute">{}</span>"#,
                label
            ));
        }

        html.push_str(&format!("</{}>", tag));
        html
    }
}
