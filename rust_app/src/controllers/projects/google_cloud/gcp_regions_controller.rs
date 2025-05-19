// Ported from: orig_app/app/controllers/projects/google_cloud/gcp_regions_controller.rb
// Ported on: 2025-05-19
// This file implements the Projects::GoogleCloud::GcpRegionsController from the Ruby codebase.

use actix_web::{web, HttpResponse, Responder};
use serde_json::json;

use crate::controllers::projects::google_cloud::base_controller::ProjectsGoogleCloudBaseController;

const AVAILABLE_REGIONS: &[&str] = &[
    "asia-east1",
    "asia-northeast1",
    "asia-southeast1",
    "europe-north1",
    "europe-west1",
    "europe-west4",
    "us-central1",
    "us-east1",
    "us-east4",
    "us-west1",
];

pub struct ProjectsGoogleCloudGcpRegionsController {
    base: ProjectsGoogleCloudBaseController,
}

impl ProjectsGoogleCloudGcpRegionsController {
    pub fn new(base: ProjectsGoogleCloudBaseController) -> Self {
        Self { base }
    }

    /// GET /projects/:project_id/google_cloud/gcp_regions
    pub async fn index(&self) -> impl Responder {
        let js_data = json!({
            "availableRegions": AVAILABLE_REGIONS,
            "refs": [], // refs logic omitted
            "cancelPath": "project_google_cloud_configuration_path" // placeholder
        });
        self.base.track_event("render_form", None);
        HttpResponse::Ok().json(js_data)
    }

    /// POST /projects/:project_id/google_cloud/gcp_regions
    pub async fn create(&self) -> impl Responder {
        // Service call and redirect logic omitted for brevity
        self.base.track_event("configure_region", None);
        HttpResponse::Ok().json(json!({"message": "GCP region configured"}))
    }
}
