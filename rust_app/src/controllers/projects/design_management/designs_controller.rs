// Ported from: orig_app/app/controllers/projects/design_management/designs_controller.rb
// Ported on: 2025-05-19
// This file implements the Projects::DesignManagement::DesignsController from the Ruby codebase.

use actix_web::HttpResponse;

use crate::controllers::projects::application_controller::ProjectsApplicationController;

pub struct ProjectsDesignManagementDesignsController {
    base: ProjectsApplicationController,
}

impl ProjectsDesignManagementDesignsController {
    pub fn new(base: ProjectsApplicationController) -> Self {
        Self { base }
    }

    fn authorize_read_design(
        &self,
        current_user: &crate::models::user::User,
        design: &crate::models::design::Design,
    ) -> Result<(), HttpResponse> {
        if !self.base.can(current_user, "read_design", design) {
            return Err(self.base.access_denied());
        }
        Ok(())
    }

    fn design(&self, design_id: &str) -> crate::models::design::Design {
        self.base.project().designs().find(design_id)
    }

    fn sha(&self, sha_param: Option<String>) -> Option<String> {
        sha_param.filter(|s| !s.is_empty())
    }
}
