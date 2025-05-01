// Ported from: orig_app/app/controllers/concerns/sourcegraph_decorator.rb
// This concern injects Sourcegraph settings and updates CSP if Sourcegraph is enabled.

use actix_web::{HttpRequest, HttpResponse};
use std::collections::HashMap;

pub trait SourcegraphDecorator {
    fn push_sourcegraph_gon(&self, req: &HttpRequest, gon: &mut HashMap<String, serde_json::Value>);
    fn update_csp_for_sourcegraph(
        &self,
        csp_directives: &mut HashMap<String, Vec<String>>,
        sourcegraph_enabled: bool,
        sourcegraph_url: &str,
        sourcegraph_public_only: bool,
        project_public: bool,
    );
}

pub struct SourcegraphDecoratorHandler;

impl SourcegraphDecoratorHandler {
    pub fn new() -> Self {
        SourcegraphDecoratorHandler
    }
}

impl SourcegraphDecorator for SourcegraphDecoratorHandler {
    fn push_sourcegraph_gon(
        &self,
        _req: &HttpRequest,
        gon: &mut HashMap<String, serde_json::Value>,
    ) {
        // Only push if Sourcegraph is enabled (should be checked by caller)
        gon.insert(
            "sourcegraph".to_string(),
            serde_json::json!({ "url": std::env::var("SOURCEGRAPH_URL").unwrap_or_default() }),
        );
    }

    fn update_csp_for_sourcegraph(
        &self,
        csp_directives: &mut HashMap<String, Vec<String>>,
        sourcegraph_enabled: bool,
        sourcegraph_url: &str,
        sourcegraph_public_only: bool,
        project_public: bool,
    ) {
        if !sourcegraph_enabled {
            return;
        }
        if sourcegraph_public_only && !project_public {
            return;
        }
        // Add connect-src for Sourcegraph
        let connect_src = csp_directives
            .entry("connect-src".to_string())
            .or_insert_with(|| {
                csp_directives
                    .get("default-src")
                    .cloned()
                    .unwrap_or_default()
            });
        let api_url = format!("{}/.api/", sourcegraph_url.trim_end_matches('/'));
        if !connect_src.contains(&api_url) {
            connect_src.push(api_url);
        }
    }
}
