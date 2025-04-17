use actix_web::{web, HttpRequest};
use std::sync::Arc;

use crate::config::database::Database;
use crate::events::event_store::EventStore;
use crate::events::users::activity_event::ActivityEvent;
use crate::models::group::Group;
use crate::models::project::Project;
use crate::models::user::User;
use crate::services::users::activity_service::ActivityService;

/// Module for recording user last activity
pub trait RecordUserLastActivity {
    /// Set the user's last activity
    fn set_user_last_activity(&self, req: &HttpRequest) {
        // Only record activity for GET requests
        if req.method() != "GET" {
            return;
        }

        // Don't record activity in read-only mode
        if Database::is_read_only() {
            return;
        }

        // Only record activity for authenticated users
        if let Some(current_user) = self.current_user() {
            // Record user activity
            ActivityService::new(current_user.clone()).execute();
        }
    }

    /// Set the member's last activity
    fn set_member_last_activity(&self, req: &HttpRequest) {
        // Only record activity for authenticated users
        if let Some(current_user) = self.current_user() {
            // Get the context (group or project)
            let context = self.get_context();

            if let Some(context) = context {
                // Only record activity for persisted contexts
                if context.is_persisted() {
                    // Get the root ancestor ID
                    let root_ancestor_id = context.root_ancestor_id();

                    // Publish activity event
                    let event = ActivityEvent::new(current_user.id(), root_ancestor_id);

                    EventStore::publish(event);
                }
            }
        }
    }

    /// Get the current context (group or project)
    fn get_context(&self) -> Option<Arc<dyn Context>> {
        self.group()
            .map(|g| g as Arc<dyn Context>)
            .or_else(|| self.project().map(|p| p as Arc<dyn Context>))
    }

    // Required trait methods that need to be implemented by the controller
    fn current_user(&self) -> Option<Arc<User>>;
    fn group(&self) -> Option<Arc<Group>>;
    fn project(&self) -> Option<Arc<Project>>;
}

/// Trait for context objects (groups and projects)
pub trait Context {
    /// Check if the context is persisted
    fn is_persisted(&self) -> bool;

    /// Get the root ancestor ID
    fn root_ancestor_id(&self) -> i64;
}
