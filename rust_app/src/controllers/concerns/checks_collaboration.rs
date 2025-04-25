// Ported from orig_app/app/controllers/concerns/checks_collaboration.rb
// April 24, 2025

use crate::models::{project::Project, user::User};
use std::collections::HashMap;
use std::sync::{Arc, Mutex};

// Stub for UserAccess, to be expanded with real permission logic
pub struct UserAccess {
    pub user: Arc<User>,
    pub project: Arc<Project>,
}

impl UserAccess {
    pub fn new(user: Arc<User>, project: Arc<Project>) -> Self {
        UserAccess { user, project }
    }

    pub fn can_push_to_branch(&self, _ref_name: Option<&str>) -> bool {
        // TODO: Implement real branch push logic
        false
    }
}

pub trait ChecksCollaboration {
    fn can_collaborate_with_project(&self, project: Arc<Project>, ref_name: Option<&str>) -> bool;
    fn user_access(&self, project: Arc<Project>) -> Arc<UserAccess>;
}

// Example implementation for a controller
pub struct ChecksCollaborationImpl {
    pub current_user: Arc<User>,
    user_access_cache: Mutex<HashMap<uuid::Uuid, Arc<UserAccess>>>,
}

impl ChecksCollaborationImpl {
    pub fn new(current_user: Arc<User>) -> Self {
        Self {
            current_user,
            user_access_cache: Mutex::new(HashMap::new()),
        }
    }
}

impl ChecksCollaboration for ChecksCollaborationImpl {
    fn can_collaborate_with_project(&self, project: Arc<Project>, ref_name: Option<&str>) -> bool {
        // TODO: Replace with real permission checks
        // Simulate: can push code, can create MR and already forked, or can push to branch
        let access = self.user_access(project.clone());
        // Placeholder: always false
        access.can_push_to_branch(ref_name)
    }

    fn user_access(&self, project: Arc<Project>) -> Arc<UserAccess> {
        let mut cache = self.user_access_cache.lock().unwrap();
        cache
            .entry(project.id)
            .or_insert_with(|| Arc::new(UserAccess::new(self.current_user.clone(), project)))
            .clone()
    }
}
