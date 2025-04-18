use crate::controllers::concerns::GitlabInternalEventsTracking;
use crate::controllers::dashboard::DashboardApplicationController;
use actix_web::{web, HttpResponse};

pub struct TodosController {
    base: DashboardApplicationController,
}

impl TodosController {
    pub fn new(base: DashboardApplicationController) -> Self {
        Self { base }
    }

    pub async fn index(&self) -> HttpResponse {
        self.track_internal_event("view_todo_list", self.base.user.clone());

        self.push_frontend_feature_flag("todos_bulk_actions", self.base.user.clone());

        // TODO: Implement HTML template rendering
        HttpResponse::Ok().finish()
    }
}

impl GitlabInternalEventsTracking for TodosController {}
