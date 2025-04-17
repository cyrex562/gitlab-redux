use actix_web::{web, HttpRequest, HttpResponse};
use std::sync::Arc;

use crate::config::settings::Settings;
use crate::utils::version::Version;

/// Module for Content Security Policy patch
pub trait ContentSecurityPolicyPatch {
    /// Configure content security policy with context
    fn content_security_policy_with_context<F>(
        &self,
        enabled: bool,
        options: &CspOptions,
        f: F,
    ) -> HttpResponse
    where
        F: FnOnce(&mut ContentSecurityPolicy),
    {
        // Check if we're using a newer version of the framework
        if Version::current() >= Version::new(7, 2, 0) {
            // Log a warning that this method is deprecated
            log::warn!(
                "content_security_policy_with_context should only be used with framework < 7.2.
                Use content_security_policy instead."
            );
        }

        // Create a new CSP policy
        let mut policy = ContentSecurityPolicy::new();

        // Apply the configuration function
        f(&mut policy);

        // Set the CSP policy on the request
        self.set_content_security_policy(&policy);

        // Disable CSP if not enabled
        if !enabled {
            self.disable_content_security_policy();
        }

        HttpResponse::Ok().finish()
    }

    /// Set the content security policy
    fn set_content_security_policy(&self, policy: &ContentSecurityPolicy);

    /// Disable content security policy
    fn disable_content_security_policy(&self);

    /// Get the current content security policy
    fn current_content_security_policy(&self) -> ContentSecurityPolicy;
}

/// Content Security Policy options
pub struct CspOptions {
    /// Whether to apply the policy to all actions
    pub apply_to_all: bool,

    /// Actions to apply the policy to
    pub only: Vec<String>,

    /// Actions to exclude from the policy
    pub except: Vec<String>,
}

impl Default for CspOptions {
    fn default() -> Self {
        Self {
            apply_to_all: true,
            only: Vec::new(),
            except: Vec::new(),
        }
    }
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
