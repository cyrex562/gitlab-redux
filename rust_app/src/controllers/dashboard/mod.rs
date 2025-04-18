use crate::controllers::concerns::{ControllerWithCrossProjectAccessCheck, RecordUserLastActivity};
use crate::models::user::User;
use actix_web::{web, HttpResponse};

pub mod groups;
pub mod labels;
pub mod milestones;
pub mod projects;
pub mod snippets;
pub mod todos;

pub struct DashboardApplicationController {
    user: User,
}

impl DashboardApplicationController {
    pub fn new(user: User) -> Self {
        Self { user }
    }

    pub fn projects(&self) -> Vec<i32> {
        // TODO: Implement authorized_projects.sorted_by_activity.non_archived
        // This is a placeholder that needs to be implemented with proper project filtering
        vec![]
    }
}

impl ControllerWithCrossProjectAccessCheck for DashboardApplicationController {
    fn check_cross_project_access(&self) -> bool {
        // TODO: Implement cross project access check
        true
    }
}

impl RecordUserLastActivity for DashboardApplicationController {
    fn record_user_activity(&self) {
        // TODO: Implement user activity recording
    }
}
