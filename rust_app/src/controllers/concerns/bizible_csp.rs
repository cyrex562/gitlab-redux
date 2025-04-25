// Ported from: orig_app/app/controllers/concerns/bizible_csp.rb
// Adds Bizible script to the Content Security Policy if enabled

use actix_web::{HttpRequest, HttpResponse, Responder};
use std::collections::HashMap;

pub trait BizibleCSP {
    fn set_bizible_csp_headers(
        &self,
        req: &HttpRequest,
        bizible_enabled: bool,
        csp_directives: &mut HashMap<String, Vec<String>>,
    ) -> HttpResponse;
}

pub struct BizibleCSPImpl;

impl BizibleCSPImpl {
    pub fn new() -> Self {
        BizibleCSPImpl
    }
}

impl BizibleCSP for BizibleCSPImpl {
    fn set_bizible_csp_headers(
        &self,
        _req: &HttpRequest,
        bizible_enabled: bool,
        csp_directives: &mut HashMap<String, Vec<String>>,
    ) -> HttpResponse {
        // Only update if Bizible is enabled or there are existing script-src/default-src directives
        let has_directives = csp_directives
            .get("script-src")
            .map(|v| !v.is_empty())
            .unwrap_or(false)
            || csp_directives
                .get("default-src")
                .map(|v| !v.is_empty())
                .unwrap_or(false);
        if !bizible_enabled && !has_directives {
            return HttpResponse::Ok().finish();
        }

        // Get script-src or fallback to default-src
        let default_script_src = csp_directives
            .get("script-src")
            .cloned()
            .or_else(|| csp_directives.get("default-src").cloned())
            .unwrap_or_else(Vec::new);
        let mut script_src_values = default_script_src;
        // Add Bizible values if not present
        if !script_src_values.contains(&"'unsafe-eval'".to_string()) {
            script_src_values.push("'unsafe-eval'".to_string());
        }
        if !script_src_values.contains(&"https://cdn.bizible.com/scripts/bizible.js".to_string()) {
            script_src_values.push("https://cdn.bizible.com/scripts/bizible.js".to_string());
        }
        csp_directives.insert("script-src".to_string(), script_src_values);

        // Format header
        let csp_header = csp_directives
            .iter()
            .map(|(k, v)| format!("{} {}", k, v.join(" ")))
            .collect::<Vec<_>>()
            .join("; ");
        HttpResponse::Ok()
            .header("Content-Security-Policy", csp_header)
            .finish()
    }
}
