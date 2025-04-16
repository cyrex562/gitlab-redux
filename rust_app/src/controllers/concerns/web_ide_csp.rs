use actix_web::{web, HttpResponse};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

/// Module for handling Content Security Policy for Web IDE
pub trait WebIdeCsp {
    /// Get the CSP directives
    fn csp_directives(&self) -> HashMap<String, String> {
        let mut directives = HashMap::new();

        // Default CSP directives for Web IDE
        directives.insert(
            "default-src".to_string(),
            "'self' 'unsafe-inline' 'unsafe-eval'".to_string(),
        );
        directives.insert(
            "script-src".to_string(),
            "'self' 'unsafe-inline' 'unsafe-eval'".to_string(),
        );
        directives.insert(
            "style-src".to_string(),
            "'self' 'unsafe-inline'".to_string(),
        );
        directives.insert(
            "img-src".to_string(),
            "'self' data: blob: https:".to_string(),
        );
        directives.insert("connect-src".to_string(), "'self' wss: https:".to_string());
        directives.insert("font-src".to_string(), "'self' data: https:".to_string());
        directives.insert("object-src".to_string(), "'none'".to_string());
        directives.insert("media-src".to_string(), "'self'".to_string());
        directives.insert("frame-src".to_string(), "'self'".to_string());

        directives
    }

    /// Get the CSP header value
    fn csp_header_value(&self) -> String {
        self.csp_directives()
            .iter()
            .map(|(key, value)| format!("{} {}", key, value))
            .collect::<Vec<String>>()
            .join("; ")
    }

    /// Add CSP headers to response
    fn add_csp_headers(&self, response: &mut HttpResponse) -> Result<(), HttpResponse> {
        response.headers_mut().insert(
            "Content-Security-Policy",
            self.csp_header_value().parse().map_err(|e| {
                HttpResponse::InternalServerError().json(serde_json::json!({
                    "error": format!("Failed to parse CSP header: {}", e)
                }))
            })?,
        );
        Ok(())
    }

    /// Get CSP report URI
    fn csp_report_uri(&self) -> Option<String> {
        None // Implement based on your needs
    }

    /// Get CSP report only mode
    fn csp_report_only(&self) -> bool {
        false // Implement based on your needs
    }

    /// Get CSP nonce
    fn csp_nonce(&self) -> String {
        use rand::Rng;
        let mut rng = rand::thread_rng();
        let nonce: u32 = rng.gen();
        format!("nonce-{}", nonce)
    }

    /// Add CSP nonce to directives
    fn add_csp_nonce(&self, directives: &mut HashMap<String, String>) {
        let nonce = self.csp_nonce();

        // Add nonce to script-src and style-src
        if let Some(script_src) = directives.get_mut("script-src") {
            *script_src = format!("{} '{}'", script_src, nonce);
        }
        if let Some(style_src) = directives.get_mut("style-src") {
            *style_src = format!("{} '{}'", style_src, nonce);
        }
    }

    /// Get CSP frame ancestors
    fn csp_frame_ancestors(&self) -> String {
        "'self'".to_string()
    }

    /// Get CSP base URI
    fn csp_base_uri(&self) -> String {
        "'self'".to_string()
    }

    /// Get CSP form action
    fn csp_form_action(&self) -> String {
        "'self'".to_string()
    }
}
