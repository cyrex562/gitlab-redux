// Ported from: orig_app/app/controllers/projects/google_cloud/configuration_controller.rb
// Ported on: 2025-05-19
// This file implements the Projects::GoogleCloud::ConfigurationController from the Ruby codebase.

use actix_web::{web, HttpResponse, Responder};
use serde_json::json;

use crate::controllers::projects::google_cloud::base_controller::ProjectsGoogleCloudBaseController;

pub struct ProjectsGoogleCloudConfigurationController {
    base: ProjectsGoogleCloudBaseController,
}

impl ProjectsGoogleCloudConfigurationController {
    pub fn new(base: ProjectsGoogleCloudBaseController) -> Self {
        Self { base }
    }

    /// GET /projects/:project_id/google_cloud/configuration
    pub async fn index(&self) -> impl Responder {
        // The following is a simplified port. Actual URL helpers and service calls would be implemented as needed.
        let js_data = json!({
            "configurationUrl": "project_google_cloud_configuration_path", // placeholder
            "deploymentsUrl": "project_google_cloud_deployments_path", // placeholder
            "databasesUrl": "project_google_cloud_databases_path", // placeholder
            "serviceAccounts": [], // ServiceAccountsService logic omitted
            "createServiceAccountUrl": "project_google_cloud_service_accounts_path", // placeholder
            "emptyIllustrationUrl": "/assets/illustrations/empty-state/empty-pipeline-md.svg", // placeholder
            "configureGcpRegionsUrl": "project_google_cloud_gcp_regions_path", // placeholder
            "gcpRegions": [], // gcp_regions logic omitted
            "revokeOauthUrl": null // revoke_oauth_url logic omitted
        });
        // In a real implementation, js_data would be constructed with actual data.
        // Track event
        self.base.track_event("render_page", None);
        HttpResponse::Ok().json(js_data)
    }
}
