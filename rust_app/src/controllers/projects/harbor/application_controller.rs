// Ported from: orig_app/app/controllers/projects/harbor/application_controller.rb
// Ported on: 2025-05-19
// This file implements the Projects::Harbor::ApplicationController from the Ruby codebase.

use crate::controllers::projects::application_controller::ProjectsApplicationController;

pub struct ProjectsHarborApplicationController {
    base: ProjectsApplicationController,
}

impl ProjectsHarborApplicationController {
    pub fn new(base: ProjectsApplicationController) -> Self {
        Self { base }
    }

    fn authorize_read_harbor_registry(
        &self,
        current_user: &crate::models::user::User,
        project: &crate::models::project::Project,
    ) -> Result<(), actix_web::HttpResponse> {
        if !self.base.can(current_user, "read_harbor_registry", project) {
            return Err(self.base.render_404());
        }
        Ok(())
    }
}
