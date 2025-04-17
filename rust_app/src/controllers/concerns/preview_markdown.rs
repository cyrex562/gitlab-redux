use actix_web::{web, HttpResponse};
use std::sync::Arc;

use crate::{
    models::{group::Group, project::Project, user::User, wiki::Wiki},
    services::{
        markdown::{MarkdownContext, MarkdownFilterParams, PreviewMarkdownService},
        suggestion::SuggestionSerializer,
    },
    utils::{browser::BrowserUtils, error::AppError},
};

/// Module for handling markdown preview
pub trait PreviewMarkdown {
    /// Preview markdown content
    fn preview_markdown(&self) -> Result<HttpResponse, AppError> {
        let current_user = self.current_user();
        let resource_parent = self.resource_parent();
        let markdown_service_params = self.markdown_service_params();

        // Create preview service
        let preview_service =
            PreviewMarkdownService::new(resource_parent, current_user, markdown_service_params);

        // Execute preview service
        let result = preview_service.execute()?;

        // Get markdown context params
        let markdown_context_params = self.markdown_context_params();

        // Render markdown
        let body = self
            .view_context()
            .markdown(&result.text, markdown_context_params);
        let commands = self.view_context().markdown(&result.commands);

        // Create response
        let response = serde_json::json!({
            "body": body,
            "references": {
                "users": result.users,
                "suggestions": SuggestionSerializer::new().represent_diff(&result.suggestions),
                "commands": commands,
            }
        });

        Ok(HttpResponse::Ok().json(response))
    }

    /// Get resource parent
    fn resource_parent(&self) -> &Project {
        self.project()
    }

    /// Get projects filter params
    fn projects_filter_params(&self) -> MarkdownFilterParams {
        MarkdownFilterParams {
            issuable_reference_expansion_enabled: true,
            suggestions_filter_enabled: self.params().get("preview_suggestions").is_some(),
        }
    }

    /// Get timeline events filter params
    fn timeline_events_filter_params(&self) -> MarkdownFilterParams {
        MarkdownFilterParams {
            issuable_reference_expansion_enabled: true,
            pipeline: Some("incident_management/timeline_event".to_string()),
        }
    }

    /// Get wikis filter params
    fn wikis_filter_params(&self) -> MarkdownFilterParams {
        MarkdownFilterParams {
            pipeline: Some("wiki".to_string()),
            wiki: Some(self.wiki()),
            page_slug: self.params().get("id").cloned(),
            repository: Some(self.wiki().repository()),
            issuable_reference_expansion_enabled: true,
        }
    }

    /// Get markdown service params
    fn markdown_service_params(&self) -> &serde_json::Value {
        self.params()
    }

    /// Get markdown context params
    fn markdown_context_params(&self) -> MarkdownContext {
        let mut context = match self.controller_name().as_str() {
            "wikis" => {
                let wiki_page = self.wiki().find_page(self.params().get("id"));
                self.wikis_filter_params()
            }
            "snippets" => MarkdownFilterParams {
                skip_project_check: Some(true),
                ..Default::default()
            },
            "groups" => MarkdownFilterParams {
                group: Some(self.group()),
                issuable_reference_expansion_enabled: true,
                ..Default::default()
            },
            "projects" => self.projects_filter_params(),
            "timeline_events" => self.timeline_events_filter_params(),
            "organizations" => MarkdownFilterParams {
                pipeline: Some("description".to_string()),
                ..Default::default()
            },
            _ => MarkdownFilterParams::default(),
        };

        // Merge common params
        context.requested_path = self.params().get("path").cloned().or_else(|| {
            self.wiki()
                .find_page(self.params().get("id"))
                .map(|p| p.path)
        });
        context.ref_param = self.params().get("ref").cloned();

        // Disable comments in markdown for IE browsers
        context.allow_comments = !BrowserUtils::is_ie();

        context
    }

    // Required trait methods that need to be implemented by the controller
    fn current_user(&self) -> Option<&User>;
    fn project(&self) -> &Project;
    fn wiki(&self) -> &Wiki;
    fn group(&self) -> &Group;
    fn params(&self) -> &serde_json::Value;
    fn controller_name(&self) -> String;
    fn view_context(&self) -> &dyn ViewContext;
}

/// Trait for view context
pub trait ViewContext {
    fn markdown(&self, text: &str, context: MarkdownContext) -> String;
}
