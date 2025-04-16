use actix_web::{web, HttpRequest, HttpResponse, Responder};
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct ContentSecurityPolicy {
    pub default_src: Vec<String>,
    pub script_src: Vec<String>,
    pub style_src: Vec<String>,
    pub img_src: Vec<String>,
    pub connect_src: Vec<String>,
    pub font_src: Vec<String>,
    pub object_src: Vec<String>,
    pub media_src: Vec<String>,
    pub frame_src: Vec<String>,
    pub sandbox: Vec<String>,
    pub report_uri: Option<String>,
}

impl ContentSecurityPolicy {
    pub fn new() -> Self {
        Self {
            default_src: vec!["'self'".to_string()],
            script_src: vec!["'self'".to_string()],
            style_src: vec!["'self'".to_string()],
            img_src: vec!["'self'".to_string()],
            connect_src: vec!["'self'".to_string()],
            font_src: vec!["'self'".to_string()],
            object_src: vec!["'none'".to_string()],
            media_src: vec!["'self'".to_string()],
            frame_src: vec!["'self'".to_string()],
            sandbox: vec![],
            report_uri: None,
        }
    }

    pub fn to_header_value(&self) -> String {
        let mut directives = Vec::new();

        if !self.default_src.is_empty() {
            directives.push(format!("default-src {}", self.default_src.join("")));
        }
        if !self.script_src.is_empty() {
            directives.push(format!("script-src {}", self.script_src.join("")));
        }
        if !self.style_src.is_empty() {
            directives.push(format!("style-src {}", self.style_src.join("")));
        }
        if !self.img_src.is_empty() {
            directives.push(format!("img-src {}", self.img_src.join("")));
        }
        if !self.connect_src.is_empty() {
            directives.push(format!("connect-src {}", self.connect_src.join("")));
        }
        if !self.font_src.is_empty() {
            directives.push(format!("font-src {}", self.font_src.join("")));
        }
        if !self.object_src.is_empty() {
            directives.push(format!("object-src {}", self.object_src.join("")));
        }
        if !self.media_src.is_empty() {
            directives.push(format!("media-src {}", self.media_src.join("")));
        }
        if !self.frame_src.is_empty() {
            directives.push(format!("frame-src {}", self.frame_src.join("")));
        }
        if !self.sandbox.is_empty() {
            directives.push(format!("sandbox {}", self.sandbox.join("")));
        }
        if let Some(report_uri) = &self.report_uri {
            directives.push(format!("report-uri {}", report_uri));
        }

        directives.join("; ")
    }
}

pub struct SecurityPolicyManager;

impl SecurityPolicyManager {
    pub fn new() -> Self {
        Self
    }

    pub fn get_csp(&self) -> ContentSecurityPolicy {
        ContentSecurityPolicy::new()
    }

    pub fn apply_security_headers(&self, resp: &mut HttpResponse) {
        // TODO: Implement security headers
        resp.headers_mut().insert(
            "Content-Security-Policy",
            self.get_csp().to_header_value().parse().unwrap(),
        );
    }

    pub fn validate_request(&self, req: &HttpRequest) -> bool {
        // TODO: Implement request validation
        true
    }
}
