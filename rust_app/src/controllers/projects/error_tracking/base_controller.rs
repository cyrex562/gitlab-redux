// Ported from: orig_app/app/controllers/projects/error_tracking/base_controller.rb
// Ported on: 2025-05-19
// This file implements the Projects::ErrorTracking::BaseController from the Ruby codebase.

use actix_web::HttpResponse;

use crate::controllers::projects::application_controller::ProjectsApplicationController;

pub struct ProjectsErrorTrackingBaseController {
    base: ProjectsApplicationController,
}

impl ProjectsErrorTrackingBaseController {
    pub fn new(base: ProjectsApplicationController) -> Self {
        Self { base }
    }

    pub fn set_polling_interval(&self, response: &mut actix_web::HttpResponseBuilder) {
        // Set polling interval header (1000ms)
        response.insert_header(("X-Poll-Interval", "1000"));
    }
}
