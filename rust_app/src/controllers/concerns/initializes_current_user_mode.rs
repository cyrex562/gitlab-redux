use actix_web::web;
use std::sync::Arc;

use crate::auth::current_user_mode::CurrentUserMode;
use crate::models::user::User;

pub trait InitializesCurrentUserMode {
    fn current_user_mode(&self) -> Arc<CurrentUserMode>;
}

pub struct InitializesCurrentUserModeImpl {
    current_user: User,
}

impl InitializesCurrentUserModeImpl {
    pub fn new(current_user: User) -> Self {
        Self { current_user }
    }
}

impl InitializesCurrentUserMode for InitializesCurrentUserModeImpl {
    fn current_user_mode(&self) -> Arc<CurrentUserMode> {
        Arc::new(CurrentUserMode::new(self.current_user.clone()))
    }
}
