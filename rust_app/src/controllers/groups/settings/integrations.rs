// Ported from: orig_app/app/controllers/groups/settings/integrations_controller.rb
// Ported on: 2025-05-01
// This file implements the Groups::Settings::IntegrationsController from the Ruby codebase.

use crate::services::integrations::{actions::Actions, finder::Finder};
use actix_web::{web, HttpResponse, Responder};
use serde::Deserialize;

pub struct GroupsSettingsIntegrationsController {
    actions: Actions,
}

impl GroupsSettingsIntegrationsController {
    pub fn new() -> Self {
        Self {
            actions: Actions::new(),
        }
    }

    /// GET /groups/:group_id/settings/integrations
    pub async fn index(&self, group_id: i32) -> impl Responder {
        // Find or initialize all non-project-specific integrations for the group
        let integrations = Finder::find_by_group_id(group_id);
        // Sort by title (case-insensitive)
        let mut integrations_sorted = integrations;
        integrations_sorted.sort_by(|a, b| a.title.to_lowercase().cmp(&b.title.to_lowercase()));
        HttpResponse::Ok().json(integrations_sorted)
    }

    /// GET /groups/:group_id/settings/integrations/:id/edit
    pub async fn edit(&self, integration_id: i32, group_id: i32) -> impl Responder {
        // Find the integration and its default
        let integration = Finder::find_by_id(integration_id);
        let default_integration = integration
            .as_ref()
            .map(|intg| Finder::find_default_integration(intg.integration_type(), group_id));
        // Call the shared edit logic
        match integration {
            Some(intg) => self.actions.edit(&intg).await,
            None => HttpResponse::NotFound().finish(),
        }
    }

    // Private helper (not an endpoint)
    fn find_or_initialize_non_project_specific_integration(
        &self,
        name: &str,
        group_id: i32,
    ) -> Option<crate::services::integrations::actions::Integration> {
        Finder::find_by_properties(
            [
                ("name".to_string(), name.to_string()),
                ("group_id".to_string(), group_id.to_string()),
            ]
            .iter()
            .cloned()
            .collect(),
        )
        .into_iter()
        .next()
    }
}
