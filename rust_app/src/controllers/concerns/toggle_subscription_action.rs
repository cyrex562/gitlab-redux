use actix_web::{web, HttpResponse, Responder};
use serde::{Deserialize, Serialize};
use std::sync::Arc;

use crate::models::project::Project;
use crate::models::subscribable::Subscribable;

/// Module for handling toggle subscription action
pub trait ToggleSubscriptionAction {
    /// Get the current user ID
    fn current_user_id(&self) -> Option<i32>;

    /// Get the subscribable project
    fn subscribable_project(&self) -> Arc<Project>;

    /// Get the subscribable resource
    fn subscribable_resource(&self) -> Arc<dyn Subscribable>;

    /// Toggle subscription
    async fn toggle_subscription(&self) -> impl Responder {
        if let Some(current_user_id) = self.current_user_id() {
            let subscribable_resource = self.subscribable_resource();
            let subscribable_project = self.subscribable_project();

            match subscribable_resource
                .toggle_subscription(current_user_id, subscribable_project)
                .await
            {
                Ok(_) => HttpResponse::Ok().finish(),
                Err(_) => HttpResponse::InternalServerError().finish(),
            }
        } else {
            HttpResponse::Unauthorized().finish()
        }
    }
}

/// Trait for resources that can be subscribed to
pub trait SubscribableResource {
    /// Toggle subscription for a user on a project
    fn toggle_subscription(&self, user: &User, project: &Project);
}
