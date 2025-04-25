// Ported from: orig_app/app/controllers/concerns/one_trust_csp.rb
// This file implements the OneTrustCSP concern in Rust.
//
// Adds OneTrust-specific script-src and connect-src values to the Content Security Policy.

use crate::auth::security_policy::ContentSecurityPolicy;

/// Trait for OneTrustCSP integration
pub trait OneTrustCSP {
    /// Updates the given ContentSecurityPolicy with OneTrust values if enabled or if policy has directives
    fn apply_one_trust_csp(&self, csp: &mut ContentSecurityPolicy, one_trust_enabled: bool);
}

/// Implementation struct for OneTrustCSP
pub struct OneTrustCSPImpl;

impl OneTrustCSP for OneTrustCSPImpl {
    fn apply_one_trust_csp(&self, csp: &mut ContentSecurityPolicy, one_trust_enabled: bool) {
        // Only apply if OneTrust is enabled or there are existing directives
        let has_directives = !csp.script_src.is_empty() || !csp.default_src.is_empty();
        if !(one_trust_enabled || has_directives) {
            return;
        }

        // script-src: add 'unsafe-eval', https://cdn.cookielaw.org, https://*.onetrust.com
        let mut script_src = if !csp.script_src.is_empty() {
            csp.script_src.clone()
        } else {
            csp.default_src.clone()
        };
        for val in [
            "'unsafe-eval'",
            "https://cdn.cookielaw.org",
            "https://*.onetrust.com",
        ] {
            if !script_src.contains(&val.to_string()) {
                script_src.push(val.to_string());
            }
        }
        csp.script_src = script_src;

        // connect-src: add https://cdn.cookielaw.org, https://*.onetrust.com
        let mut connect_src = if !csp.connect_src.is_empty() {
            csp.connect_src.clone()
        } else {
            csp.default_src.clone()
        };
        for val in ["https://cdn.cookielaw.org", "https://*.onetrust.com"] {
            if !connect_src.contains(&val.to_string()) {
                connect_src.push(val.to_string());
            }
        }
        csp.connect_src = connect_src;
    }
}
