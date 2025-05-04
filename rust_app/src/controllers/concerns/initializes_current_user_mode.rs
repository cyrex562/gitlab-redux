// Ported from: orig_app/app/controllers/concerns/initializes_current_user_mode.rb
// Provides a trait and implementation for initializing the current user mode.
use std::sync::Arc;

use crate::auth::current_user_mode::CurrentUserMode;
use crate::models::user::User;

pub trait InitializesCurrentUserMode {
    fn initialize_current_user_mode(&self, user: Option<Arc<User>>);
}

pub struct UserMode {
    pub admin_mode: bool,
}

impl UserMode {
    pub fn new(admin_mode: bool) -> Self {
        Self { admin_mode }
    }
}

pub struct InitializesCurrentUserModeImpl {
    current_user: User,
    current_user_mode: Option<Arc<CurrentUserMode>>,
}

impl InitializesCurrentUserModeImpl {
    pub fn new(current_user: User) -> Self {
        Self {
            current_user,
            current_user_mode: None,
        }
    }
}

impl InitializesCurrentUserMode for InitializesCurrentUserModeImpl {
    fn current_user_mode(&self) -> Arc<CurrentUserMode> {
        // In a real app, you may want to use interior mutability for caching
        Arc::new(CurrentUserMode::new(self.current_user.clone()))
    }
}
