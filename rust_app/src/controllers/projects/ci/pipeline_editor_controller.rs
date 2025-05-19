// Ported from: orig_app/app/controllers/projects/ci/pipeline_editor_controller.rb
// Ported on: 2025-05-19
// This file implements the Projects::Ci::PipelineEditorController from the Ruby codebase.

use actix_web::{web, HttpResponse, Responder};

use crate::controllers::projects::application_controller::ProjectsApplicationController;

pub struct ProjectsCiPipelineEditorController {
    base: ProjectsApplicationController,
}

impl ProjectsCiPipelineEditorController {
    pub fn new(base: ProjectsApplicationController) -> Self {
        Self { base }
    }

    /// GET /projects/:project_id/ci/pipeline_editor
    pub async fn show(&self) -> impl Responder {
        // No-op, just return 200 OK
        HttpResponse::Ok().finish()
    }

    fn check_can_collaborate(&self) -> Result<(), HttpResponse> {
        let project = self.base.project();
        if !self.base.can_collaborate_with_project(&project) {
            return Err(HttpResponse::NotFound().finish());
        }
        Ok(())
    }
}
