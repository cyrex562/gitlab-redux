use std::sync::Arc;

use crate::auth::current_user_mode::CurrentUserMode;
use crate::models::user::User;

/// Module for initializing current user mode
pub trait InitializesCurrentUserMode {
    /// Get the current user mode
    fn current_user_mode(&self) -> Arc<CurrentUserMode> {
        if let Some(current_user) = self.current_user() {
            return Arc::new(CurrentUserMode::new(current_user));
        }

        // Return a default mode if no user is available
        Arc::new(CurrentUserMode::default())
    }

    // Required trait methods that need to be implemented by the controller
    fn current_user(&self) -> Option<Arc<User>>;
}
