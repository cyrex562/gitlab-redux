// Ported from: orig_app/app/controllers/groups/harbor/application_controller.rb
// Ported on: 2025-05-01
// This file implements the Groups::Harbor::ApplicationController from the Ruby codebase.

use crate::controllers::concerns::harbor::HarborAccess;
use crate::controllers::groups::GroupsApplicationController;
use actix_web::{HttpResponse, Responder};

pub struct HarborApplicationController {
    pub base: GroupsApplicationController,
}

impl HarborApplicationController {
    pub fn new(base: GroupsApplicationController) -> Self {
        Self { base }
    }

    /// Authorize read access to the Harbor registry for the group.
    pub fn authorize_read_harbor_registry(
        &self,
        current_user: &crate::models::user::User,
        group: &crate::models::group::Group,
    ) -> Result<(), HttpResponse> {
        if !self.can_read_harbor_registry(current_user, group) {
            return Err(self.render_404());
        }
        Ok(())
    }

    fn can_read_harbor_registry(
        &self,
        _current_user: &crate::models::user::User,
        _group: &crate::models::group::Group,
    ) -> bool {
        // TODO: Implement actual permission check logic
        false
    }

    fn render_404(&self) -> HttpResponse {
        HttpResponse::NotFound().finish()
    }
}

impl HarborAccess for HarborApplicationController {}
