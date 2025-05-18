// Ported from: orig_app/app/controllers/organizations/settings_controller.rb
// Ported on: 2025-05-05

use actix_web::{get, HttpResponse, Responder};

use crate::controllers::organizations::application_controller::OrganizationsApplicationController;

/// Controller for managing organization settings
pub struct SettingsController {
    base: OrganizationsApplicationController,
}

impl SettingsController {
    pub fn new(base: OrganizationsApplicationController) -> Self {
        Self { base }
    }

    /// GET /organizations/settings/general
    pub async fn general(&self) -> impl Responder {
        // Feature flag and authorization are checked in base controller
        if let Err(e) = self.base.authorize_admin_organization() {
            return e;
        }
        HttpResponse::Ok().finish()
    }
}
