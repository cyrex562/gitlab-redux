// Ported from: orig_app/app/controllers/projects/error_tracking/projects_controller.rb
// Ported on: 2025-05-19
// This file implements the Projects::ErrorTracking::ProjectsController from the Ruby codebase.

use actix_web::{web, HttpResponse, Responder};
use serde_json::json;

use crate::controllers::projects::application_controller::ProjectsApplicationController;
use crate::serializers::error_tracking::project_serializer::ProjectSerializer;
use crate::services::error_tracking::list_projects_service::ListProjectsService;

#[derive(Debug, serde::Deserialize)]
pub struct ListProjectsParams {
    pub api_host: Option<String>,
    pub token: Option<String>,
}

pub struct ProjectsErrorTrackingProjectsController {
    base: ProjectsApplicationController,
}

impl ProjectsErrorTrackingProjectsController {
    pub fn new(base: ProjectsApplicationController) -> Self {
        Self { base }
    }

    /// GET /projects/:project_id/error_tracking/projects
    pub async fn index(&self, params: web::Query<ListProjectsParams>) -> impl Responder {
        if let Err(e) = self.base.authorize_admin_sentry() {
            return e;
        }
        let project = self.base.project();
        let current_user = self.base.current_user();
        let list_params = self.list_projects_params(&params);
        let service = ListProjectsService::new(project, current_user, &list_params);
        let result = service.execute();
        if result.status == "success" {
            let projects_json =
                ProjectSerializer::new(project, current_user).represent(&result.projects);
            HttpResponse::Ok().json(json!({ "projects": projects_json }))
        } else {
            HttpResponse::BadRequest().json(json!({ "message": result.message }))
        }
    }

    fn list_projects_params(&self, params: &ListProjectsParams) -> ListProjectsParams {
        ListProjectsParams {
            api_host: params.api_host.clone(),
            token: params.token.clone(),
        }
    }
}
