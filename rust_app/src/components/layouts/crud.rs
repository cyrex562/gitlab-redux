use std::collections::HashMap;
use crate::components::Component;

/// Component for displaying CRUD (Create, Read, Update, Delete) interfaces
pub struct CrudComponent {
    title: String,
    description: Option<String>,
    count: Option<i32>,
    icon: Option<String>,
    icon_class: Option<String>,
    toggle_text: Option<String>,
    options: HashMap<String, String>,
    count_options: HashMap<String, String>,
    body_options: HashMap<String, String>,
    form_options: HashMap<String, String>,
    toggle_options: HashMap<String, String>,
    footer_options: HashMap<String, String>,
    is_collapsible: bool,
    container_tag: String,
    description_slot: Option<String>,
    actions_slot: Option<String>,
    body_slot: Option<String>,
    form_slot: Option<String>,
    footer_slot: Option<String>,
    pagination_slot: Option<String>,
}

impl Component for CrudComponent {
    fn render(&self) -> String {
        let id = self.generate_id();
        let options_attrs = self.format_options_attrs();
        let body_options_attrs = self.format_body_options_attrs();
        let form_options_attrs = self.format_form_options_attrs();
        let footer_options_attrs = self.format_footer_options_attrs();
        let toggle_button_options_attrs = self.format_toggle_button_options_attrs();
        
        // Build the HTML
        let mut html = format!(
            r#"<{} id="{}" {}>",
            self.container_tag, id, options_attrs
        );
        
        // Header
        html.push_str(&self.render_header());
        
        // Form
        if let Some(form) = &self.form_slot {
            html.push_str(&format!(
                r#"<div class="gl-p-5 gl-pt-4 gl-bg-subtle gl-border-b gl-border-section" {}>{}</div>"#,
                form_options_attrs, form
            ));
        }
        
        // Body
        html.push_str(&format!(
            r#"<div class="crud-body gl-mx-5 gl-my-4" {}>{}</div>"#,
            body_options_attrs, self.body_slot.as_ref().unwrap_or(&String::new())
        ));
        
        // Pagination
        if let Some(pagination) = &self.pagination_slot {
            html.push_str(&format!(
                r#"<div class="crud-pagination gl-flex gl-justify-center gl-p-5 gl-border-t gl-border-t-section" data-testid="crud-pagination">{}</div>"#,
                pagination
            ));
        }
        
        // Footer
        if let Some(footer) = &self.footer_slot {
            html.push_str(&format!(
                r#"<div class="gl-px-5 gl-py-4 gl-bg-section gl-border-t gl-border-section gl-rounded-b-base" {}>{}</div>"#,
                footer_options_attrs, footer
            ));
        }
        
        html.push_str(&format!("</{}>", self.container_tag));
        
        html
    }
}

impl CrudComponent {
    /// Create a new CRUD component
    /// 
    /// # Arguments
    /// 
    /// * `title` - The title of the CRUD component
    /// * `description` - Optional description
    /// * `count` - Optional count to display
    /// * `icon` - Optional icon name
    /// * `icon_class` - Optional icon class
    /// * `toggle_text` - Optional toggle button text
    /// * `options` - Optional HTML attributes for the container
    /// * `count_options` - Optional HTML attributes for the count
    /// * `body_options` - Optional HTML attributes for the body
    /// * `form_options` - Optional HTML attributes for the form
    /// * `toggle_options` - Optional HTML attributes for the toggle button
    /// * `footer_options` - Optional HTML attributes for the footer
    /// * `is_collapsible` - Whether the component is collapsible
    /// * `container_tag` - The HTML tag to use for the container
    pub fn new(
        title: String,
        description: Option<String>,
        count: Option<i32>,
        icon: Option<String>,
        icon_class: Option<String>,
        toggle_text: Option<String>,
        options: HashMap<String, String>,
        count_options: HashMap<String, String>,
        body_options: HashMap<String, String>,
        form_options: HashMap<String, String>,
        toggle_options: HashMap<String, String>,
        footer_options: HashMap<String, String>,
        is_collapsible: bool,
        container_tag: String,
    ) -> Self {
        Self {
            title,
            description,
            count,
            icon,
            icon_class,
            toggle_text,
            options,
            count_options,
            body_options,
            form_options,
            toggle_options,
            footer_options,
            is_collapsible,
            container_tag,
            description_slot: None,
            actions_slot: None,
            body_slot: None,
            form_slot: None,
            footer_slot: None,
            pagination_slot: None,
        }
    }
    
    /// Set the description slot content
    pub fn with_description(mut self, description: String) -> Self {
        self.description_slot = Some(description);
        self
    }
    
    /// Set the actions slot content
    pub fn with_actions(mut self, actions: String) -> Self {
        self.actions_slot = Some(actions);
        self
    }
    
    /// Set the body slot content
    pub fn with_body(mut self, body: String) -> Self {
        self.body_slot = Some(body);
        self
    }
    
    /// Set the form slot content
    pub fn with_form(mut self, form: String) -> Self {
        self.form_slot = Some(form);
        self
    }
    
    /// Set the footer slot content
    pub fn with_footer(mut self, footer: String) -> Self {
        self.footer_slot = Some(footer);
        self
    }
    
    /// Set the pagination slot content
    pub fn with_pagination(mut self, pagination: String) -> Self {
        self.pagination_slot = Some(pagination);
        self
    }
    
    /// Check if description slot is set
    pub fn has_description(&self) -> bool {
        self.description_slot.is_some() || self.description.is_some()
    }
    
    /// Check if form slot is set
    pub fn has_form(&self) -> bool {
        self.form_slot.is_some()
    }
    
    /// Check if footer slot is set
    pub fn has_footer(&self) -> bool {
        self.footer_slot.is_some()
    }
    
    /// Check if pagination slot is set
    pub fn has_pagination(&self) -> bool {
        self.pagination_slot.is_some()
    }
    
    /// Generate a unique ID for the component
    fn generate_id(&self) -> String {
        self.title.to_lowercase()
            .replace(" ", "-")
            .chars()
            .filter(|c| c.is_alphanumeric() || *c == '-')
            .collect()
    }
    
    /// Format options attributes
    fn format_options_attrs(&self) -> String {
        let mut classes = vec![
            "crud".to_string(),
            "gl-bg-subtle".to_string(),
            "gl-border".to_string(),
            "gl-border-section".to_string(),
            "gl-rounded-base".to_string(),
        ];
        
        if self.toggle_text.is_some() {
            classes.push("js-toggle-container".to_string());
        }
        
        if self.is_collapsible {
            classes.push("js-crud-collapsible-section".to_string());
        }
        
        let mut attrs = self.options.clone();
        
        // Add data-testid if not present
        if !attrs.contains_key("data-testid") {
            attrs.insert("data-testid".to_string(), "haml-crud".to_string());
        }
        
        // Add classes
        if let Some(class) = attrs.remove("class") {
            attrs.insert("class".to_string(), format!("{} {}", class, classes.join(" ")));
        } else {
            attrs.insert("class".to_string(), classes.join(" "));
        }
        
        // Format attributes
        attrs.iter()
            .map(|(k, v)| format!("{}={}", k, v))
            .collect::<Vec<String>>()
            .join(" ")
    }
    
    /// Format body options attributes
    fn format_body_options_attrs(&self) -> String {
        let mut classes = Vec::new();
        
        if !self.has_footer() {
            classes.push("gl-rounded-b-base".to_string());
        }
        
        if self.is_collapsible {
            classes.push("js-crud-collapsible-content".to_string());
        }
        
        let mut attrs = self.body_options.clone();
        
        // Add data-testid if not present
        if !attrs.contains_key("data-testid") {
            attrs.insert("data-testid".to_string(), "crud-body".to_string());
        }
        
        // Add classes
        if let Some(class) = attrs.remove("class") {
            attrs.insert("class".to_string(), format!("{} {}", class, classes.join(" ")));
        } else {
            attrs.insert("class".to_string(), classes.join(" "));
        }
        
        // Format attributes
        attrs.iter()
            .map(|(k, v)| format!("{}={}", k, v))
            .collect::<Vec<String>>()
            .join(" ")
    }
    
    /// Format form options attributes
    fn format_form_options_attrs(&self) -> String {
        let mut classes = Vec::new();
        
        if self.toggle_text.is_some() {
            classes.push("js-toggle-content".to_string());
        }
        
        if self.toggle_text.is_some() && !self.form_options.contains_key("form_errors") {
            classes.push("gl-hidden".to_string());
        }
        
        if self.is_collapsible {
            classes.push("js-crud-collapsible-content".to_string());
        }
        
        let mut attrs = self.form_options.clone();
        
        // Add data-testid if not present
        if !attrs.contains_key("data-testid") {
            attrs.insert("data-testid".to_string(), "crud-form".to_string());
        }
        
        // Add classes
        if let Some(class) = attrs.remove("class") {
            attrs.insert("class".to_string(), format!("{} {}", class, classes.join(" ")));
        } else {
            attrs.insert("class".to_string(), classes.join(" "));
        }
        
        // Format attributes
        attrs.iter()
            .map(|(k, v)| format!("{}={}", k, v))
            .collect::<Vec<String>>()
            .join(" ")
    }
    
    /// Format footer options attributes
    fn format_footer_options_attrs(&self) -> String {
        let mut classes = Vec::new();
        
        if self.is_collapsible {
            classes.push("js-crud-collapsible-content".to_string());
        }
        
        let mut attrs = self.footer_options.clone();
        
        // Add data-testid if not present
        if !attrs.contains_key("data-testid") {
            attrs.insert("data-testid".to_string(), "crud-footer".to_string());
        }
        
        // Add classes
        if let Some(class) = attrs.remove("class") {
            attrs.insert("class".to_string(), format!("{} {}", class, classes.join(" ")));
        } else {
            attrs.insert("class".to_string(), classes.join(" "));
        }
        
        // Format attributes
        attrs.iter()
            .map(|(k, v)| format!("{}={}", k, v))
            .collect::<Vec<String>>()
            .join(" ")
    }
    
    /// Format toggle button options attributes
    fn format_toggle_button_options_attrs(&self) -> String {
        let classes = vec!["js-toggle-button".to_string(), "js-toggle-content".to_string()];
        
        let mut attrs = self.toggle_options.clone();
        
        // Add data-testid if not present
        if !attrs.contains_key("data-testid") {
            attrs.insert("data-testid".to_string(), "crud-action-toggle".to_string());
        }
        
        // Add classes
        if let Some(class) = attrs.remove("class") {
            attrs.insert("class".to_string(), format!("{} {}", class, classes.join(" ")));
        } else {
            attrs.insert("class".to_string(), classes.join(" "));
        }
        
        // Format attributes
        attrs.iter()
            .map(|(k, v)| format!("{}={}", k, v))
            .collect::<Vec<String>>()
            .join(" ")
    }
    
    /// Get icon classes
    fn get_icon_classes(&self) -> String {
        let mut classes = vec!["gl-fill-icon-subtle".to_string()];
        
        if let Some(icon_class) = &self.icon_class {
            classes.push(icon_class.clone());
        }
        
        classes.join(" ")
    }
    
    /// Render the header section
    fn render_header(&self) -> String {
        let mut header = format!(
            r#"<header class="crud-header gl-flex gl-flex-wrap gl-justify-between gl-gap-x-5 gl-gap-y-2 gl-px-5 gl-py-4 gl-bg-section gl-border-b gl-border-section gl-rounded-t-base gl-relative{}">"#,
            if self.is_collapsible { " gl-pr-10" } else { "" }
        );
        
        // Title section
        header.push_str(r#"<div class="gl-flex gl-flex-col gl-self-center">"#);
        
        // Title with count
        header.push_str(&format!(
            r#"<h2 class="gl-text-base gl-font-bold gl-leading-normal gl-inline-flex gl-gap-3 gl-m-0" data-testid="crud-title">{}</h2>"#,
            self.title
        ));
        
        if let Some(count) = self.count {
            header.push_str(&format!(
                r#"<span class="gl-inline-flex gl-items-center gl-gap-2 gl-text-sm gl-text-subtle" data-testid="crud-count">"#,
            ));
            
            if let Some(icon) = &self.icon {
                header.push_str(&format!(
                    r#"<span class="{}"><img src="/assets/icons/{}.svg" alt="" /></span>"#,
                    self.get_icon_classes(), icon
                ));
            }
            
            header.push_str(&format!(
                r#"<span {}>{}</span></span>"#,
                self.count_options.iter()
                    .map(|(k, v)| format!("{}={}", k, v))
                    .collect::<Vec<String>>()
                    .join(" "),
                count
            ));
        }
        
        // Description
        if self.has_description() {
            let description = self.description_slot.as_ref().or(self.description.as_ref()).unwrap();
            header.push_str(&format!(
                r#"<div class="gl-text-sm gl-text-subtle gl-leading-normal gl-mt-2 gl-mb-0" data-testid="crud-description">{}</div>"#,
                description
            ));
        }
        
        header.push_str("</div>");
        
        // Actions
        header.push_str(r#"<div class="gl-flex gl-gap-3 gl-items-baseline" data-testid="crud-actions">"#);
        
        if let Some(actions) = &self.actions_slot {
            header.push_str(actions);
        }
        
        // Toggle button
        if let Some(toggle_text) = &self.toggle_text {
            header.push_str(&format!(
                r#"<button class="btn btn-sm btn-default {}" {}>{}</button>"#,
                self.format_toggle_button_options_attrs(),
                toggle_text
            ));
        }
        
        // Collapsible button
        if self.is_collapsible {
            header.push_str(r#"<div class="gl-border-l gl-absolute gl-right-5 gl-top-4 gl-h-6 gl-border-l-section gl-pl-3">"#);
            header.push_str(r#"<button class="btn btn-sm btn-tertiary -gl-mr-2 js-crud-collapsible-button has-tooltip" title="Collapse" aria-expanded="true" aria-controls="" data-collapse-title="Expand" data-expand-title="Collapse">"#);
            header.push_str(r#"<span class="!-gl-mx-2 js-crud-collapsible-collapse"><img src="/assets/icons/chevron-lg-up.svg" alt="" /></span>"#);
            header.push_str(r#"<span class="!-gl-mx-2 gl-hidden js-crud-collapsible-expand"><img src="/assets/icons/chevron-lg-down.svg" alt="" /></span>"#);
            header.push_str("</button></div>");
        }
        
        header.push_str("</div></header>");
        
        header
    }
} 