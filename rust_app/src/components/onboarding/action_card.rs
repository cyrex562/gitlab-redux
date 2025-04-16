use serde::{Deserialize, Serialize};
use std::collections::HashMap;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ActionCardComponent {
    title: String,
    description: String,
    icon: Option<String>,
    href: Option<String>,
    variant: Option<String>,
    link_options: HashMap<String, String>,
    html_options: HashMap<String, String>,
    content: Option<String>,
}

impl ActionCardComponent {
    pub fn new(
        title: String,
        description: String,
        icon: Option<String>,
        href: Option<String>,
        variant: Option<String>,
        link_options: Option<HashMap<String, String>>,
        html_options: Option<HashMap<String, String>>,
        content: Option<String>,
    ) -> Self {
        Self {
            title,
            description,
            icon,
            href,
            variant,
            link_options: link_options.unwrap_or_default(),
            html_options: Self::format_options(html_options),
            content,
        }
    }

    pub fn render(&self) -> String {
        let mut html = format!("<div{}>", Self::format_attributes(&self.html_options));

        if self.has_link() {
            html.push_str(&format!(
                r#"<a class="gl-link action-card-title" href="{}"{}>"#,
                self.href.as_ref().unwrap(),
                Self::format_attributes(&self.link_options)
            ));
            html.push_str(&self.render_icon());
            html.push_str(&self.title);
            html.push_str(r#"<svg class="sprite-icon action-card-arrow"><use xlink:href="/assets/icons.svg#arrow-right"></use></svg>"#);
            html.push_str("</a>");
        } else {
            html.push_str(r#"<div class="action-card-title">"#);
            html.push_str(&self.render_icon());
            html.push_str(&self.title);
            html.push_str("</div>");
        }

        html.push_str(&format!(
            r#"<p class="action-card-text">{}</p>"#,
            self.description
        ));

        if let Some(content) = &self.content {
            html.push_str(&format!(
                r#"<div class="action-card-controls">{}</div>"#,
                content
            ));
        }

        html.push_str("</div>");
        html
    }

    fn has_link(&self) -> bool {
        self.href.is_some()
    }

    fn render_icon(&self) -> String {
        let icon = self.card_icon();
        format!(
            r#"<svg class="sprite-icon"><use xlink:href="/assets/icons.svg#{}</use></svg>"#,
            icon
        )
    }

    fn card_icon(&self) -> String {
        self.icon.clone().unwrap_or_else(|| {
            match self.variant.as_deref() {
                Some("primary") => "rocket",
                Some("secondary") => "information",
                _ => "information",
            }
            .to_string()
        })
    }

    fn format_options(options: Option<HashMap<String, String>>) -> HashMap<String, String> {
        let mut default_options = HashMap::new();
        default_options.insert("class".to_string(), "action-card".to_string());

        if let Some(mut opts) = options {
            if let Some(class) = opts.remove("class") {
                default_options.insert(
                    "class".to_string(),
                    format!("{} {}", default_options["class"], class),
                );
            }
            default_options.extend(opts);
        }

        default_options
    }

    fn format_attributes(options: &HashMap<String, String>) -> String {
        options
            .iter()
            .map(|(key, value)| format!(" {}={}", key, value))
            .collect::<Vec<String>>()
            .join("")
    }
}
