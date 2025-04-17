use crate::config::settings::Settings;
use actix_web::{web, HttpResponse};
use std::collections::HashMap;
use std::sync::Arc;

/// Module for handling CSP for static object external storage
pub trait StaticObjectExternalStorageCsp {
    /// Set CSP headers for external storage
    fn set_external_storage_csp(&self) -> HttpResponse {
        let settings = Settings::current();
        if !settings.static_objects_external_storage_enabled {
            return HttpResponse::Ok().finish();
        }

        // Build CSP directives
        let mut directives = Vec::new();

        // Default directives
        directives.push("default-src 'none'".to_string());
        directives.push("script-src 'none'".to_string());
        directives.push("style-src 'none'".to_string());

        // Image and font sources
        let image_src = format!(
            "img-src 'self' data: {}",
            settings.static_objects_hosts.join(" ")
        );
        directives.push(image_src);

        let font_src = format!(
            "font-src 'self' data: {}",
            settings.static_objects_hosts.join(" ")
        );
        directives.push(font_src);

        // Connect sources
        let connect_src = format!(
            "connect-src 'self' {}",
            settings.static_objects_hosts.join(" ")
        );
        directives.push(connect_src);

        // Add any additional CSP directives from settings
        if let Some(additional_directives) = &settings.static_objects_csp_directives {
            directives.extend(additional_directives.clone());
        }

        // Join directives with semicolons
        let csp = directives.join("; ");

        // Create response with CSP header
        HttpResponse::Ok()
            .header("Content-Security-Policy", csp)
            .finish()
    }

    /// Get CSP settings
    fn get_csp_settings(&self) -> HashMap<String, String> {
        let mut settings = HashMap::new();
        let settings = Settings::current();

        settings.insert(
            "static_objects_external_storage_enabled".to_string(),
            settings.static_objects_external_storage_enabled.to_string(),
        );

        if let Some(directives) = &settings.static_objects_csp_directives {
            settings.insert(
                "static_objects_csp_directives".to_string(),
                directives.join(", "),
            );
        }

        settings
    }
}

/// Module for handling static object external storage CSP
pub trait StaticObjectExternalStorageCSP {
    /// Configure content security policy for external storage
    fn configure_csp(&self, policy: &mut ContentSecurityPolicy) {
        if policy.directives.is_empty() {
            return;
        }

        let settings = self.settings();
        if !settings.static_objects_external_storage_enabled() {
            return;
        }

        let default_connect_src = policy
            .directives
            .get("connect-src")
            .or_else(|| policy.directives.get("default-src"))
            .cloned()
            .unwrap_or_default();

        let mut connect_src_values = default_connect_src;
        connect_src_values.push(settings.static_objects_external_storage_url());

        policy
            .directives
            .insert("connect-src".to_string(), connect_src_values);
    }

    // Required trait methods that need to be implemented by the controller
    fn settings(&self) -> Arc<Settings>;
}

/// Content Security Policy configuration
pub struct ContentSecurityPolicy {
    /// CSP directives
    pub directives: std::collections::HashMap<String, Vec<String>>,
}

impl ContentSecurityPolicy {
    /// Create a new CSP configuration
    pub fn new() -> Self {
        Self {
            directives: std::collections::HashMap::new(),
        }
    }

    /// Add a directive
    pub fn add_directive(&mut self, name: &str, values: Vec<String>) {
        self.directives.insert(name.to_string(), values);
    }

    /// Get a directive
    pub fn get_directive(&self, name: &str) -> Option<&Vec<String>> {
        self.directives.get(name)
    }

    /// Remove a directive
    pub fn remove_directive(&mut self, name: &str) {
        self.directives.remove(name);
    }

    /// Clear all directives
    pub fn clear(&mut self) {
        self.directives.clear();
    }
}
