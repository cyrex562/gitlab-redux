// Ported from: orig_app/app/controllers/projects/ci/lints_controller.rb
// Ported on: 2025-05-19
// This file implements the Projects::Ci::LintsController from the Ruby codebase.

use actix_web::{web, HttpResponse, Responder};
use serde::Deserialize;
use serde_json::json;

use crate::controllers::projects::application_controller::ProjectsApplicationController;
use crate::serializers::ci::lint::ResultSerializer;
use crate::services::ci::lint::LintService;

#[derive(Debug, Deserialize)]
pub struct LintParams {
    pub content: String,
    pub dry_run: Option<bool>,
}

pub struct ProjectsCiLintsController {
    base: ProjectsApplicationController,
}

impl ProjectsCiLintsController {
    pub fn new(base: ProjectsApplicationController) -> Self {
        Self { base }
    }

    /// GET /projects/:project_id/ci/lint
    pub async fn show(&self) -> impl Responder {
        // No-op, just return 200 OK
        HttpResponse::Ok().finish()
    }

    /// POST /projects/:project_id/ci/lint
    pub async fn create(&self, params: web::Json<LintParams>) -> impl Responder {
        if let Err(e) = self.base.authorize_create_pipeline() {
            return e;
        }
        let content = &params.content;
        let dry_run = params.dry_run.unwrap_or(false);
        let project = self.base.project();
        let current_user = self.base.current_user();

        let result = LintService::new(project, current_user).validate(content, dry_run);
        let json = ResultSerializer::new().represent(&result);
        HttpResponse::Ok().json(json)
    }
}
