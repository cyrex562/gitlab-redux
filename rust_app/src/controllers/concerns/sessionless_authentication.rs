// Ported from: orig_app/app/controllers/concerns/sessionless_authentication.rb
// This file implements sessionless authentication logic for PAT, RSS, and static object tokens.

use crate::auth::request_authenticator::RequestAuthenticator;
use crate::models::user::User;
use crate::settings::ApplicationSettings;
use actix_web::{web, HttpRequest, HttpResponse};
use std::sync::Arc;

pub trait SessionlessAuthentication {
    fn authenticate_sessionless_user(
        &self,
        req: &HttpRequest,
        request_format: &str,
    ) -> Result<Option<Arc<User>>, HttpResponse>;
    fn request_authenticator(&self, req: &HttpRequest) -> RequestAuthenticator;
    fn sessionless_user(&self, current_user: Option<Arc<User>>) -> bool;
    fn sessionless_sign_in(&self, user: Arc<User>) -> Result<(), HttpResponse>;
    fn sessionless_bypass_admin_mode<F, R>(
        &self,
        current_user: Option<Arc<User>>,
        settings: &ApplicationSettings,
        f: F,
    ) -> Result<R, HttpResponse>
    where
        F: FnOnce() -> Result<R, HttpResponse>;
}

pub struct SessionlessAuthenticationHandler {
    settings: Arc<ApplicationSettings>,
}

impl SessionlessAuthenticationHandler {
    pub fn new(settings: Arc<ApplicationSettings>) -> Self {
        SessionlessAuthenticationHandler { settings }
    }
}

impl SessionlessAuthentication for SessionlessAuthenticationHandler {
    fn authenticate_sessionless_user(
        &self,
        req: &HttpRequest,
        request_format: &str,
    ) -> Result<Option<Arc<User>>, HttpResponse> {
        let authenticator = self.request_authenticator(req);
        let user = authenticator.find_sessionless_user(request_format)?;

        if let Some(user) = &user {
            self.sessionless_sign_in(user.clone())?;
        }

        Ok(user)
    }

    fn request_authenticator(&self, req: &HttpRequest) -> RequestAuthenticator {
        RequestAuthenticator::new(req)
    }

    fn sessionless_user(&self, current_user: Option<Arc<User>>) -> bool {
        // In a real implementation, this would check if the current user was authenticated
        // via a sessionless method (PAT, RSS token, etc.)
        // For now, we'll just return false
        false
    }

    fn sessionless_sign_in(&self, user: Arc<User>) -> Result<(), HttpResponse> {
        // In a real implementation, this would sign in the user without creating a session
        // For now, we'll just return Ok
        Ok(())
    }

    fn sessionless_bypass_admin_mode<F, R>(
        &self,
        current_user: Option<Arc<User>>,
        settings: &ApplicationSettings,
        f: F,
    ) -> Result<R, HttpResponse>
    where
        F: FnOnce() -> Result<R, HttpResponse>,
    {
        if !settings.admin_mode {
            return f();
        }

        if let Some(user) = current_user {
            // In a real implementation, this would bypass admin mode for the user
            // For now, we'll just call the function
            f()
        } else {
            f()
        }
    }
}

// This would be implemented in a separate module
pub mod auth {
    pub mod request_authenticator {
        use crate::models::user::User;
        use actix_web::HttpRequest;
        use std::sync::Arc;

        pub struct RequestAuthenticator {
            req: HttpRequest,
        }

        impl RequestAuthenticator {
            pub fn new(req: HttpRequest) -> Self {
                RequestAuthenticator { req }
            }

            pub fn find_sessionless_user(
                &self,
                request_format: &str,
            ) -> Result<Option<Arc<User>>, actix_web::Error> {
                // In a real implementation, this would check for various authentication methods:
                // - Personal Access Tokens
                // - RSS Tokens
                // - Static Object Tokens
                // For now, we'll just return None
                Ok(None)
            }
        }
    }
}

// This would be implemented in a separate module
pub mod settings {
    use serde::{Deserialize, Serialize};
    use std::sync::Arc;

    #[derive(Debug, Clone, Serialize, Deserialize)]
    pub struct ApplicationSettings {
        pub admin_mode: bool,
    }

    impl ApplicationSettings {
        pub fn new() -> Self {
            ApplicationSettings { admin_mode: false }
        }
    }
}

// This would be implemented in a separate module
pub mod models {
    pub mod user {
        use std::sync::Arc;

        pub struct User {
            pub id: i64,
            pub username: String,
            pub password_expired: bool,
        }

        impl User {
            pub fn id(&self) -> i64 {
                self.id
            }

            pub fn username(&self) -> &str {
                &self.username
            }

            pub fn can_log_in_with_non_expired_password(&self) -> bool {
                !self.password_expired
            }
        }
    }
}
