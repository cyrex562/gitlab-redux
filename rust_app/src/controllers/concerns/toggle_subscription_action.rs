use crate::models::project::Project;
use crate::models::user::User;
use actix_web::{web, HttpResponse};

/// Module for handling subscription toggling functionality
pub trait ToggleSubscriptionAction {
    /// Toggle subscription for the current user
    fn toggle_subscription(&self, user: &User, project: &Project) -> HttpResponse {
        if let Some(user) = user {
            if let Some(project) = self.subscribable_project() {
                self.subscribable_resource()
                    .toggle_subscription(user, project);
                HttpResponse::Ok().finish()
            } else {
                HttpResponse::NotImplemented().finish()
            }
        } else {
            HttpResponse::Unauthorized().finish()
        }
    }

    /// Get the subscribable project
    fn subscribable_project(&self) -> Option<&Project>;

    /// Get the subscribable resource
    fn subscribable_resource(&self) -> &dyn SubscribableResource;
}

/// Trait for resources that can be subscribed to
pub trait SubscribableResource {
    /// Toggle subscription for a user on a project
    fn toggle_subscription(&self, user: &User, project: &Project);
}
