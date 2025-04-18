use actix_web::{web, HttpRequest, HttpResponse, Responder};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

pub trait OneTrustCSP {
    fn set_csp_headers(&self, req: &HttpRequest) -> impl Responder;
}

pub struct OneTrustCSPImpl {
    nonce: String,
    report_only: bool,
}

impl OneTrustCSPImpl {
    pub fn new(nonce: String, report_only: bool) -> Self {
        Self { nonce, report_only }
    }

    fn csp_directives(&self) -> HashMap<String, String> {
        let mut directives = HashMap::new();

        // Default directives
        directives.insert("default-src".to_string(), "'self'".to_string());
        directives.insert(
            "script-src".to_string(),
            format!(
                "'self' 'nonce-{}' 'unsafe-inline' 'unsafe-eval'",
                self.nonce
            ),
        );
        directives.insert(
            "style-src".to_string(),
            "'self' 'unsafe-inline'".to_string(),
        );
        directives.insert(
            "img-src".to_string(),
            "'self' data: blob: https:".to_string(),
        );
        directives.insert("font-src".to_string(), "'self' data: https:".to_string());
        directives.insert("connect-src".to_string(), "'self' https: wss:".to_string());
        directives.insert("frame-src".to_string(), "'self' https:".to_string());
        directives.insert("object-src".to_string(), "'none'".to_string());
        directives.insert("base-uri".to_string(), "'self'".to_string());
        directives.insert("form-action".to_string(), "'self'".to_string());
        directives.insert("frame-ancestors".to_string(), "'none'".to_string());
        directives.insert("upgrade-insecure-requests".to_string(), "".to_string());

        directives
    }

    fn format_csp_header(&self) -> String {
        let directives = self.csp_directives();
        directives
            .iter()
            .map(|(key, value)| {
                if value.is_empty() {
                    key.clone()
                } else {
                    format!("{} {}", key, value)
                }
            })
            .collect::<Vec<String>>()
            .join("; ")
    }
}

impl OneTrustCSP for OneTrustCSPImpl {
    fn set_csp_headers(&self, req: &HttpRequest) -> impl Responder {
        let header_name = if self.report_only {
            "Content-Security-Policy-Report-Only"
        } else {
            "Content-Security-Policy"
        };

        HttpResponse::Ok()
            .header(header_name, self.format_csp_header())
            .finish()
    }
}
