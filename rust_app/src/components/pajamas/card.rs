use std::collections::HashMap;

/// Renders a Pajamas compliant card element
#[derive(Debug, Clone)]
pub struct CardComponent {
    /// Options for the card
    card_options: HashMap<String, String>,
    /// Options for the header
    header_options: HashMap<String, String>,
    /// Options for the body
    body_options: HashMap<String, String>,
    /// Options for the footer
    footer_options: HashMap<String, String>,
    /// The card structure
    card: HashMap<String, String>,
    /// The header content slot
    header_slot: Option<String>,
    /// The body content slot
    body_slot: Option<String>,
    /// The footer content slot
    footer_slot: Option<String>,
}

impl CardComponent {
    /// Create a new card component
    pub fn new(
        card: Option<HashMap<String, String>>,
        card_options: Option<HashMap<String, String>>,
        header_options: Option<HashMap<String, String>>,
        body_options: Option<HashMap<String, String>>,
        footer_options: Option<HashMap<String, String>>,
    ) -> Self {
        Self {
            card_options: card_options.unwrap_or_default(),
            header_options: header_options.unwrap_or_default(),
            body_options: body_options.unwrap_or_default(),
            footer_options: footer_options.unwrap_or_default(),
            card: card.unwrap_or_default(),
            header_slot: None,
            body_slot: None,
            footer_slot: None,
        }
    }

    /// Set the header content slot
    pub fn with_header_slot(mut self, header: String) -> Self {
        self.header_slot = Some(header);
        self
    }

    /// Set the body content slot
    pub fn with_body_slot(mut self, body: String) -> Self {
        self.body_slot = Some(body);
        self
    }

    /// Set the footer content slot
    pub fn with_footer_slot(mut self, footer: String) -> Self {
        self.footer_slot = Some(footer);
        self
    }

    /// Check if the header is present
    fn has_header(&self) -> bool {
        self.header_slot.is_some() || self.card.contains_key("header")
    }

    /// Get the header content
    fn get_header(&self) -> Option<String> {
        self.header_slot
            .clone()
            .or_else(|| self.card.get("header").cloned())
    }

    /// Get the body content
    fn get_body(&self) -> Option<String> {
        self.body_slot
            .clone()
            .or_else(|| self.card.get("body").cloned())
    }

    /// Check if the footer is present
    fn has_footer(&self) -> bool {
        self.footer_slot.is_some() || self.card.contains_key("footer")
    }

    /// Get the footer content
    fn get_footer(&self) -> Option<String> {
        self.footer_slot
            .clone()
            .or_else(|| self.card.get("footer").cloned())
    }

    /// Format the card options
    fn formatted_card_options(&self) -> String {
        self.card_options
            .iter()
            .map(|(key, value)| format!(" {}=\"{}\"", key, value))
            .collect::<Vec<String>>()
            .join("")
    }

    /// Format the header options
    fn formatted_header_options(&self) -> String {
        self.header_options
            .iter()
            .map(|(key, value)| format!(" {}=\"{}\"", key, value))
            .collect::<Vec<String>>()
            .join("")
    }

    /// Format the body options
    fn formatted_body_options(&self) -> String {
        self.body_options
            .iter()
            .map(|(key, value)| format!(" {}=\"{}\"", key, value))
            .collect::<Vec<String>>()
            .join("")
    }

    /// Format the footer options
    fn formatted_footer_options(&self) -> String {
        self.footer_options
            .iter()
            .map(|(key, value)| format!(" {}=\"{}\"", key, value))
            .collect::<Vec<String>>()
            .join("")
    }

    /// Render the card component
    pub fn render(&self) -> String {
        let mut html = String::new();
        let card_options = self.formatted_card_options();

        // Card container
        html.push_str(&format!(r#"<div class="gl-card"{}>"#, card_options));

        // Header
        if self.has_header() {
            let header_options = self.formatted_header_options();
            html.push_str(&format!(
                r#"<div class="gl-card-header"{}>"#,
                header_options
            ));

            if let Some(header) = self.get_header() {
                html.push_str(&header);
            }

            html.push_str("</div>");
        }

        // Body
        let body_options = self.formatted_body_options();
        html.push_str(&format!(r#"<div class="gl-card-body"{}>"#, body_options));

        if let Some(body) = self.get_body() {
            html.push_str(&body);
        }

        html.push_str("</div>");

        // Footer
        if self.has_footer() {
            let footer_options = self.formatted_footer_options();
            html.push_str(&format!(
                r#"<div class="gl-card-footer"{}>"#,
                footer_options
            ));

            if let Some(footer) = self.get_footer() {
                html.push_str(&footer);
            }

            html.push_str("</div>");
        }

        html.push_str("</div>");
        html
    }
}
