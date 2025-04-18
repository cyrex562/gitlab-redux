use actix_web::{web, HttpRequest, HttpResponse};
use std::collections::HashMap;
use url::Url;

pub trait WebIdeCSP {
    fn include_web_ide_csp(&self, req: &HttpRequest) -> Result<(), Box<dyn std::error::Error>>;
}

pub struct WebIdeCSPHandler;

impl WebIdeCSPHandler {
    pub fn new() -> Self {
        WebIdeCSPHandler
    }
}

impl WebIdeCSP for WebIdeCSPHandler {
    fn include_web_ide_csp(&self, req: &HttpRequest) -> Result<(), Box<dyn std::error::Error>> {
        // Get the current CSP directives from the request
        let mut csp_directives = req
            .headers()
            .get("content-security-policy")
            .and_then(|h| h.to_str().ok())
            .map(|h| parse_csp_directives(h))
            .unwrap_or_default();

        if csp_directives.is_empty() {
            return Ok(());
        }

        // Build the webpack URL
        let base_uri = req.uri().to_string();
        let mut url = Url::parse(&base_uri)?;
        let relative_url_root =
            std::env::var("GITLAB_RELATIVE_URL_ROOT").unwrap_or_else(|_| "/".to_string());
        url.set_path(&format!(
            "{}/assets/webpack/",
            relative_url_root.trim_end_matches('/')
        ));
        url.set_query(None);
        let webpack_url = url.to_string();

        // Update frame-src directive
        let default_src = csp_directives
            .get("default-src")
            .cloned()
            .unwrap_or_default();

        let frame_src = csp_directives
            .entry("frame-src".to_string())
            .or_insert_with(Vec::new);

        frame_src.extend(default_src.iter().cloned());
        frame_src.push(webpack_url.clone());
        frame_src.push("https://*.web-ide.gitlab-static.net/".to_string());

        // Update worker-src directive
        let worker_src = csp_directives
            .entry("worker-src".to_string())
            .or_insert_with(Vec::new);

        worker_src.extend(default_src);
        worker_src.push(webpack_url);

        // Set the updated CSP header
        let csp_header = format_csp_directives(&csp_directives);
        req.headers_mut()
            .insert("content-security-policy", csp_header.parse()?);

        Ok(())
    }
}

fn parse_csp_directives(header: &str) -> HashMap<String, Vec<String>> {
    let mut directives = HashMap::new();

    for directive in header.split(';') {
        if let Some((name, values)) = directive.split_once(' ') {
            let values: Vec<String> = values.split_whitespace().map(|s| s.to_string()).collect();
            directives.insert(name.to_string(), values);
        }
    }

    directives
}

fn format_csp_directives(directives: &HashMap<String, Vec<String>>) -> String {
    directives
        .iter()
        .map(|(name, values)| format!("{} {}", name, values.join(" ")))
        .collect::<Vec<_>>()
        .join("; ")
}
