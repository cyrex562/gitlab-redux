use std::collections::HashMap;

/// Renders a Pajamas compliant breadcrumb navigation
#[derive(Debug, Clone)]
pub struct BreadcrumbComponent {
    /// HTML options for the breadcrumb
    html_options: HashMap<String, String>,
    /// The breadcrumb items
    items: Vec<BreadcrumbItemComponent>,
}

impl BreadcrumbComponent {
    /// Create a new breadcrumb component
    pub fn new(html_options: Option<HashMap<String, String>>) -> Self {
        Self {
            html_options: html_options.unwrap_or_default(),
            items: Vec::new(),
        }
    }

    /// Add a breadcrumb item
    pub fn with_item(mut self, item: BreadcrumbItemComponent) -> Self {
        self.items.push(item);
        self
    }

    /// Add multiple breadcrumb items
    pub fn with_items(mut self, items: Vec<BreadcrumbItemComponent>) -> Self {
        self.items.extend(items);
        self
    }

    /// Format the HTML options
    fn formatted_html_options(&self) -> String {
        self.html_options
            .iter()
            .map(|(key, value)| format!(" {}=\"{}\"", key, value))
            .collect::<Vec<String>>()
            .join("")
    }

    /// Render the breadcrumb component
    pub fn render(&self) -> String {
        let mut html = String::new();
        let html_options = self.formatted_html_options();

        // Navigation element
        html.push_str(&format!(
            r#"<nav class="gl-breadcrumbs" aria-label="Breadcrumbs"{}>"#,
            html_options
        ));

        // Breadcrumb list
        html.push_str(r#"<ul class="gl-breadcrumb-list breadcrumb js-breadcrumbs-list">"#);

        // Breadcrumb items
        for item in &self.items {
            html.push_str(&item.render());
        }

        html.push_str("</ul></nav>");
        html
    }
}

/// Renders a Pajamas compliant breadcrumb item
#[derive(Debug, Clone)]
pub struct BreadcrumbItemComponent {
    /// The href for the breadcrumb item
    href: String,
    /// The text for the breadcrumb item
    text: String,
}

impl BreadcrumbItemComponent {
    /// Create a new breadcrumb item component
    pub fn new(href: String, text: String) -> Self {
        Self { href, text }
    }

    /// Render the breadcrumb item component
    pub fn render(&self) -> String {
        format!(
            r#"<li class="gl-breadcrumb-item"><a href="{}">{}</a></li>"#,
            self.href, self.text
        )
    }
}
