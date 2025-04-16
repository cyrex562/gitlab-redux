#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum ButtonCategory {
    Primary,
    Secondary,
    Tertiary,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum ButtonVariant {
    Default,
    Confirm,
    Danger,
    Dashed,
    Link,
    Reset,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum ButtonSize {
    Small,
    Medium,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum ButtonType {
    Button,
    Reset,
    Submit,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum ButtonTarget {
    Self_,
    Blank,
    Parent,
    Top,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum ButtonMethod {
    Get,
    Post,
    Put,
    Delete,
    Patch,
}

pub struct ButtonComponent {
    category: ButtonCategory,
    variant: ButtonVariant,
    size: ButtonSize,
    button_type: ButtonType,
    disabled: bool,
    loading: bool,
    block: bool,
    label: bool,
    selected: bool,
    icon: Option<String>,
    href: Option<String>,
    form: bool,
    target: Option<ButtonTarget>,
    method: Option<ButtonMethod>,
    button_options: std::collections::HashMap<String, String>,
    button_text_classes: Option<String>,
    icon_classes: Option<String>,
    content: Option<String>,
    icon_content: Option<String>,
}

impl ButtonComponent {
    pub fn new() -> Self {
        Self {
            category: ButtonCategory::Primary,
            variant: ButtonVariant::Default,
            size: ButtonSize::Medium,
            button_type: ButtonType::Button,
            disabled: false,
            loading: false,
            block: false,
            label: false,
            selected: false,
            icon: None,
            href: None,
            form: false,
            target: None,
            method: None,
            button_options: std::collections::HashMap::new(),
            button_text_classes: None,
            icon_classes: None,
            content: None,
            icon_content: None,
        }
    }

    fn button_class(&self) -> String {
        let mut classes = vec!["gl-button", "btn"];

        if self.disabled || self.loading {
            classes.push("disabled");
        }
        if self.selected {
            classes.push("selected");
        }
        if self.block {
            classes.push("btn-block");
        }
        if self.label {
            classes.push("btn-label");
        }
        if self.icon.is_some() && self.content.is_none() {
            classes.push("btn-icon");
        }

        classes.push(match self.size {
            ButtonSize::Small => "btn-sm",
            ButtonSize::Medium => "btn-md",
        });

        let variant_class = match self.variant {
            ButtonVariant::Default => "btn-default",
            ButtonVariant::Confirm => "btn-confirm",
            ButtonVariant::Danger => "btn-danger",
            ButtonVariant::Dashed => "btn-dashed",
            ButtonVariant::Link => "btn-link",
            ButtonVariant::Reset => "btn-gl-reset",
        };
        classes.push(variant_class);

        let non_category_variants = [
            ButtonVariant::Dashed,
            ButtonVariant::Link,
            ButtonVariant::Reset,
        ];
        if !non_category_variants.contains(&self.variant)
            && self.category != ButtonCategory::Primary
        {
            let category_class = match self.category {
                ButtonCategory::Primary => "",
                ButtonCategory::Secondary => "secondary",
                ButtonCategory::Tertiary => "tertiary",
            };
            if !category_class.is_empty() {
                classes.push(&format!("{}-{}", variant_class, category_class));
            }
        }

        if let Some(class) = self.button_options.get("class") {
            classes.push(class);
        }

        classes.join(" ")
    }

    fn is_link(&self) -> bool {
        self.href.is_some()
    }

    fn is_form(&self) -> bool {
        self.href.is_some() && self.form
    }

    fn base_attributes(&self) -> std::collections::HashMap<String, String> {
        let mut attributes = std::collections::HashMap::new();

        if self.disabled || self.loading {
            attributes.insert("disabled".to_string(), "disabled".to_string());
            attributes.insert("aria-disabled".to_string(), "true".to_string());
        }

        if !self.is_link() {
            attributes.insert(
                "type".to_string(),
                match self.button_type {
                    ButtonType::Button => "button".to_string(),
                    ButtonType::Reset => "reset".to_string(),
                    ButtonType::Submit => "submit".to_string(),
                },
            );
        }

        if self.is_link() && self.target == Some(ButtonTarget::Blank) {
            let rel = self
                .button_options
                .get("rel")
                .map(|r| format!("{} noopener noreferrer", r))
                .unwrap_or_else(|| "noopener noreferrer".to_string());
            attributes.insert("rel".to_string(), rel);
        }

        attributes
    }

    pub fn render(&self) -> String {
        let mut html = String::new();
        let show_icon = !self.loading || self.content.is_some();
        let tag = if self.label { "span" } else { "button" };

        let mut content = String::new();
        if self.loading {
            content.push_str(&format!(
                r#"<span class="gl-button-icon gl-button-loading-indicator">{}</span>"#,
                // TODO: Implement gl_loading_icon helper
                "Loading..."
            ));
        }

        if show_icon {
            if let Some(icon) = &self.icon {
                content.push_str(&format!(
                    r#"<span class="gl-icon gl-button-icon {}">{}</span>"#,
                    self.icon_classes.as_deref().unwrap_or(""),
                    // TODO: Implement sprite_icon helper
                    icon
                ));
            } else if let Some(icon_content) = &self.icon_content {
                content.push_str(icon_content);
            }
        }

        if let Some(text) = &self.content {
            content.push_str(&format!(
                r#"<span class="gl-button-text {}">{}</span>"#,
                self.button_text_classes.as_deref().unwrap_or(""),
                text
            ));
        }

        let attributes = self.base_attributes();
        let mut attr_str = String::new();
        for (key, value) in attributes {
            attr_str.push_str(&format!(r#" {}="{}""#, key, value));
        }

        if self.is_form() {
            // TODO: Implement button_to helper
            html.push_str(&format!(
                r#"<form action="{}" method="post">{}</form>"#,
                self.href.as_deref().unwrap_or(""),
                content
            ));
        } else if self.is_link() {
            html.push_str(&format!(
                r#"<a href="{}" class="{}"{}>{}</a>"#,
                self.href.as_deref().unwrap_or(""),
                self.button_class(),
                attr_str,
                content
            ));
        } else {
            html.push_str(&format!(
                r#"<{} class="{}"{}>{}</{}>"#,
                tag,
                self.button_class(),
                attr_str,
                content,
                tag
            ));
        }

        html
    }
}

impl Default for ButtonComponent {
    fn default() -> Self {
        Self::new()
    }
}
