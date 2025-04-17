use actix_web::{web, HttpRequest, HttpResponse};
use std::collections::HashMap;
use std::sync::Arc;

use crate::config::settings::Settings;
use crate::utils::helpers::Helpers;

/// Module for OneTrust Content Security Policy
pub trait OneTrustCSP {
    /// Configure content security policy for OneTrust
    fn configure_one_trust_csp(&self, policy: &mut ContentSecurityPolicy) {
        // Skip if OneTrust is not enabled and no directives are present
        if !self.helpers().one_trust_enabled() && policy.directives.is_empty() {
            return;
        }

        // Configure script-src directive
        let default_script_src = policy
            .directives
            .get("script-src")
            .or_else(|| policy.directives.get("default-src"))
            .cloned()
            .unwrap_or_default();

        let mut script_src_values = default_script_src;
        script_src_values.push("'unsafe-eval'".to_string());
        script_src_values.push("https://cdn.cookielaw.org".to_string());
        script_src_values.push("https://*.onetrust.com".to_string());

        policy
            .directives
            .insert("script-src".to_string(), script_src_values);

        // Configure connect-src directive
        let default_connect_src = policy
            .directives
            .get("connect-src")
            .or_else(|| policy.directives.get("default-src"))
            .cloned()
            .unwrap_or_default();

        let mut connect_src_values = default_connect_src;
        connect_src_values.push("https://cdn.cookielaw.org".to_string());
        connect_src_values.push("https://*.onetrust.com".to_string());

        policy
            .directives
            .insert("connect-src".to_string(), connect_src_values);
    }

    // Required trait methods that need to be implemented by the controller
    fn helpers(&self) -> &dyn Helpers;
}

/// Content Security Policy configuration
pub struct ContentSecurityPolicy {
    /// CSP directives
    pub directives: HashMap<String, Vec<String>>,
}

impl ContentSecurityPolicy {
    /// Create a new CSP configuration
    pub fn new() -> Self {
        Self {
            directives: HashMap::new(),
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
