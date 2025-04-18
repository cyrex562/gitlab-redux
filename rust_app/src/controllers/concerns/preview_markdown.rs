use crate::{
    models::{Group, Project, Wiki},
    services::preview_markdown::PreviewMarkdownService,
    utils::browser::Browser,
};
use actix_web::{error::Error, web::Json, HttpResponse};
use serde_json::Value;

pub trait PreviewMarkdown {
    fn preview_markdown(&self, req: &ServiceRequest) -> Result<HttpResponse, Error>;
}

pub struct PreviewMarkdownImpl {
    browser: Arc<Browser>,
}

impl PreviewMarkdownImpl {
    pub fn new(browser: Arc<Browser>) -> Self {
        Self { browser }
    }

    fn resource_parent(&self, req: &ServiceRequest) -> Option<&dyn ResourceParent> {
        req.extensions()
            .get::<Project>()
            .map(|p| p as &dyn ResourceParent)
    }

    fn projects_filter_params(&self, req: &ServiceRequest) -> Value {
        serde_json::json!({
            "issuable_reference_expansion_enabled": true,
            "suggestions_filter_enabled": req.query_string().contains("preview_suggestions")
        })
    }

    fn timeline_events_filter_params(&self) -> Value {
        serde_json::json!({
            "issuable_reference_expansion_enabled": true,
            "pipeline": "incident_management/timeline_event"
        })
    }

    fn wikis_filter_params(&self, req: &ServiceRequest) -> Value {
        if let Some(wiki) = req.extensions().get::<Wiki>() {
            serde_json::json!({
                "pipeline": "wiki",
                "wiki": wiki,
                "page_slug": req.query_string().get("id"),
                "repository": wiki.repository(),
                "issuable_reference_expansion_enabled": true
            })
        } else {
            serde_json::json!({})
        }
    }

    fn markdown_service_params(&self, req: &ServiceRequest) -> Value {
        serde_json::to_value(req.query_string()).unwrap_or_default()
    }

    fn markdown_context_params(&self, req: &ServiceRequest) -> Value {
        let mut params = match req.path() {
            path if path.contains("/wikis") => {
                if let Some(wiki) = req.extensions().get::<Wiki>() {
                    if let Some(page) = wiki.find_page(req.query_string().get("id")) {
                        self.wikis_filter_params(req)
                    } else {
                        serde_json::json!({})
                    }
                } else {
                    serde_json::json!({})
                }
            }
            path if path.contains("/snippets") => serde_json::json!({
                "skip_project_check": true
            }),
            path if path.contains("/groups") => {
                if let Some(group) = req.extensions().get::<Group>() {
                    serde_json::json!({
                        "group": group,
                        "issuable_reference_expansion_enabled": true
                    })
                } else {
                    serde_json::json!({})
                }
            }
            path if path.contains("/projects") => self.projects_filter_params(req),
            path if path.contains("/timeline_events") => self.timeline_events_filter_params(),
            path if path.contains("/organizations") => serde_json::json!({
                "pipeline": "description"
            }),
            _ => serde_json::json!({}),
        };

        if let Some(path) = req.query_string().get("path") {
            params["requested_path"] = Value::String(path.to_string());
        }
        if let Some(ref_) = req.query_string().get("ref") {
            params["ref"] = Value::String(ref_.to_string());
        }
        params["allow_comments"] = Value::Bool(!self.browser.is_ie());

        params
    }
}

impl PreviewMarkdown for PreviewMarkdownImpl {
    fn preview_markdown(&self, req: &ServiceRequest) -> Result<HttpResponse, Error> {
        let result = PreviewMarkdownService::new(
            self.resource_parent(req),
            req.extensions().get::<User>().cloned(),
            self.markdown_service_params(req),
        )
        .execute()?;

        Ok(HttpResponse::Ok().json(serde_json::json!({
            "body": result.text,
            "references": {
                "users": result.users,
                "suggestions": result.suggestions,
                "commands": result.commands
            }
        })))
    }
}
