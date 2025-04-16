use std::collections::HashMap;

/// Renders a Pajamas compliant broadcast banner
#[derive(Debug, Clone)]
pub struct BroadcastBannerComponent {
    /// The message to display
    message: String,
    /// The ID of the banner
    id: String,
    /// The theme of the banner
    theme: String,
    /// Whether the banner is dismissable
    dismissable: bool,
    /// The expiration date of the banner
    expire_date: String,
    /// The cookie key for the banner
    cookie_key: String,
    /// The dismissal path for the banner
    dismissal_path: Option<String>,
    /// The test ID for the button
    button_testid: Option<String>,
    /// The banner data
    banner: Option<String>,
}

impl BroadcastBannerComponent {
    /// Create a new broadcast banner component
    pub fn new(
        message: String,
        id: String,
        theme: String,
        dismissable: bool,
        expire_date: String,
        cookie_key: String,
        dismissal_path: Option<String>,
        button_testid: Option<String>,
        banner: Option<String>,
    ) -> Self {
        Self {
            message,
            id,
            theme,
            dismissable,
            expire_date,
            cookie_key,
            dismissal_path,
            button_testid,
            banner,
        }
    }

    /// Render the broadcast banner component
    pub fn render(&self) -> String {
        let mut html = String::new();

        // Banner container
        let banner_data = self
            .banner
            .as_ref()
            .map(|b| format!(" data-broadcast-banner=\"{}\"", b))
            .unwrap_or_default();

        html.push_str(&format!(
            r#"<div class="gl-broadcast-message banner js-broadcast-notification-{} {}" role="alert" data-testid="banner-broadcast-message"{}>"#,
            self.id, self.theme, banner_data
        ));

        // Banner content
        html.push_str(r#"<div class="gl-broadcast-message-content">"#);

        // Banner icon
        html.push_str(r#"<div class="gl-broadcast-message-icon">"#);
        html.push_str(r#"<svg class="s16" data-testid="bullhorn-icon"><use xlink:href="/assets/icons.svg#bullhorn"></use></svg>"#);
        html.push_str("</div>");

        // Banner text
        html.push_str(r#"<div class="gl-broadcast-message-text">"#);
        html.push_str(r#"<h2 class="gl-sr-only">Admin message</h2>"#);
        html.push_str(&format!("{}", self.message));
        html.push_str("</div>");

        html.push_str("</div>");

        // Dismiss button
        if self.dismissable {
            let button_testid = self
                .button_testid
                .as_ref()
                .map(|t| format!(" data-testid=\"{}\"", t))
                .unwrap_or_default();

            let dismissal_path = self
                .dismissal_path
                .as_ref()
                .map(|p| format!(" data-dismissal-path=\"{}\"", p))
                .unwrap_or_default();

            html.push_str(&format!(
                r#"<button class="btn btn-tertiary gl-broadcast-message-dismiss js-dismiss-current-broadcast-notification" aria-label="Close" data-id="{}" data-expire-date="{}" data-cookie-key="{}"{}>{}"#,
                self.id,
                self.expire_date,
                self.cookie_key,
                button_testid + &dismissal_path,
                r#"<svg class="s16" data-testid="close-icon"><use xlink:href="/assets/icons.svg#close"></use></svg>"#
            ));
            html.push_str("</button>");
        }

        html.push_str("</div>");
        html
    }
}
