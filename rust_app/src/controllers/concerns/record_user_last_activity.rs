// Ported from: orig_app/app/controllers/concerns/record_user_last_activity.rb
// This trait provides controller logic to update the user's last activity on GET requests
// and to publish a member activity event after the action.

use crate::event_store::{ActivityEvent, EventStore};
use crate::models::{Group, Project, User};
use crate::services::users::ActivityService;
use actix_web::{HttpRequest, HttpResponse};
use chrono::Utc;

pub trait RecordUserLastActivity {
    // Must be implemented by the controller
    fn current_user(&self) -> Option<&User>;
    fn request(&self) -> &HttpRequest;
    fn group(&self) -> Option<&Group> {
        None
    }
    fn project(&self) -> Option<&Project> {
        None
    }
    fn is_db_read_only(&self) -> bool;

    // Call this at the start of GET requests
    fn set_user_last_activity(&self) {
        let req = self.request();
        if req.method() != actix_web::http::Method::GET {
            return;
        }
        if self.is_db_read_only() {
            return;
        }
        let user = match self.current_user() {
            Some(u) => u,
            None => return,
        };
        // TODO: add namespace & project context if needed
        let _ = ActivityService::new(user.clone()).execute();
    }

    // Call this after the action
    fn set_member_last_activity(&self) {
        let user = match self.current_user() {
            Some(u) => u,
            None => return,
        };
        let context = self
            .group()
            .map(|g| g as &dyn Persisted)
            .or_else(|| self.project().map(|p| p as &dyn Persisted));
        let context = match context {
            Some(c) if c.persisted() => c,
            _ => return,
        };
        let namespace_id = context.root_ancestor_id();
        let event = ActivityEvent {
            user_id: user.id,
            namespace_id,
            timestamp: Utc::now(),
        };
        let _ = EventStore::publish(event);
    }
}

// Helper trait for context objects
pub trait Persisted {
    fn persisted(&self) -> bool;
    fn root_ancestor_id(&self) -> i64;
}

impl Persisted for Group {
    fn persisted(&self) -> bool {
        self.persisted
    }
    fn root_ancestor_id(&self) -> i64 {
        self.root_ancestor_id
    }
}
impl Persisted for Project {
    fn persisted(&self) -> bool {
        self.persisted
    }
    fn root_ancestor_id(&self) -> i64 {
        self.root_ancestor_id
    }
}
