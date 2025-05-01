// Ported from: orig_app/app/controllers/dashboard/groups_controller.rb
// Ported on: 2025-04-30
// This file implements the Dashboard::GroupsController from the Ruby codebase.
//
// See porting_log.txt for details.

use crate::controllers::dashboard::DashboardApplicationController;
use crate::services::groups::GroupsFinder;
use actix_web::{web, HttpResponse};

pub struct GroupsController {
    base: DashboardApplicationController,
}

impl GroupsController {
    pub fn new(base: DashboardApplicationController) -> Self {
        Self { base }
    }

    pub async fn index(&self) -> HttpResponse {
        // Skip cross project access check for index
        let groups = GroupsFinder::new(self.base.user.clone(), false)
            .execute()
            .await;

        // TODO: Implement render_group_tree
        // This is a placeholder that needs to be implemented with proper group tree rendering
        HttpResponse::Ok().json(groups)
    }
}
