// Ported from: orig_app/app/controllers/projects/google_cloud/base_controller.rb
// Ported on: 2025-05-19
// This file implements the Projects::GoogleCloud::BaseController from the Ruby codebase.

use actix_web::HttpResponse;

use crate::controllers::projects::application_controller::ProjectsApplicationController;

pub struct ProjectsGoogleCloudBaseController {
    base: ProjectsApplicationController,
}

impl ProjectsGoogleCloudBaseController {
    pub fn new(base: ProjectsApplicationController) -> Self {
        Self { base }
    }

    fn admin_project_google_cloud(
        &self,
        current_user: &crate::models::user::User,
        project: &crate::models::project::Project,
    ) -> Result<(), HttpResponse> {
        if !self
            .base
            .can(current_user, "admin_project_google_cloud", project)
        {
            self.track_event("error_invalid_user", None);
            return Err(self.base.access_denied());
        }
        Ok(())
    }

    fn google_oauth2_enabled(&self) -> Result<(), HttpResponse> {
        let config = crate::auth::oauth::provider::config_for("google_oauth2");
        if config.app_id.is_empty() || config.app_secret.is_empty() {
            self.track_event("error_google_oauth2_not_enabled", None);
            return Err(self.base.access_denied_with_message(
                "This GitLab instance not configured for Google Oauth2.",
            ));
        }
        Ok(())
    }

    // Other methods (validate_gcp_token, generate_session_key_redirect, etc.) would be ported similarly as needed.

    fn track_event(&self, action: &str, label: Option<&str>) {
        crate::tracking::event(
            std::any::type_name::<Self>(),
            action,
            label,
            &self.base.project(),
            &self.base.current_user(),
        );
    }
}
