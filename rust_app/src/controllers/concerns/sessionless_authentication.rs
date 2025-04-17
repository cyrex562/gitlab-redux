use crate::config::settings::Settings;
use crate::models::user::User;
use actix_web::{web, HttpRequest, HttpResponse};
use jsonwebtoken::{encode, EncodingKey, Header};
use serde::{Deserialize, Serialize};
use std::cell::RefCell;
use std::collections::HashMap;
use std::sync::Arc;
use std::time::{Duration, SystemTime, UNIX_EPOCH};

use crate::auth::current_user_mode::CurrentUserMode;
use crate::auth::request_authenticator::RequestAuthenticator;

/// Module for handling sessionless authentication
pub trait SessionlessAuthentication {
    /// Authenticate a sessionless user
    fn authenticate_sessionless_user(&self, request_format: &str) -> Option<Arc<User>> {
        let user = self
            .request_authenticator()
            .find_sessionless_user(request_format);

        if let Some(user) = &user {
            self.sessionless_sign_in(user.clone());
        }

        user
    }

    /// Get the request authenticator
    fn request_authenticator(&self) -> RequestAuthenticator {
        RequestAuthenticator::new(self.request())
    }

    /// Check if the current user is a sessionless user
    fn is_sessionless_user(&self) -> bool {
        self.current_user().is_some() && self.sessionless_sign_in_flag()
    }

    /// Set the sessionless sign in flag
    fn set_sessionless_sign_in_flag(&self, value: bool) {
        self.sessionless_sign_in_flag_mut().set(value);
    }

    /// Get the sessionless sign in flag
    fn sessionless_sign_in_flag(&self) -> bool {
        *self.sessionless_sign_in_flag_mut().borrow()
    }

    /// Get a mutable reference to the sessionless sign in flag
    fn sessionless_sign_in_flag_mut(&self) -> &RefCell<bool>;

    /// Sign in a user without creating a session
    fn sessionless_sign_in(&self, user: Arc<User>) {
        // Set the sessionless sign in flag
        self.set_sessionless_sign_in_flag(true);

        if user.can_log_in_with_non_expired_password() {
            // Sign in the user without storing in session
            self.sign_in(user, false, "sessionless_sign_in");
        } else if self.request_authenticator().can_sign_in_bot(&user) {
            // Sign in the bot without callbacks
            self.sign_in_without_callbacks(user, "sessionless_sign_in");
        }
    }

    /// Bypass admin mode for sessionless authentication
    fn sessionless_bypass_admin_mode<F, R>(&self, f: F) -> R
    where
        F: FnOnce() -> R,
    {
        if !self.settings().admin_mode_enabled() {
            return f();
        }

        if let Some(current_user) = self.current_user() {
            return CurrentUserMode::bypass_session(current_user.id(), f);
        }

        f()
    }

    /// Get the current user
    fn current_user(&self) -> Option<Arc<User>>;

    /// Get the authentication token
    fn auth_token(&self) -> String;

    /// Get the token expiration time
    fn token_expiration(&self) -> Duration {
        Duration::from_secs(3600) // Default 1 hour
    }

    /// Get the token issuer
    fn token_issuer(&self) -> String {
        "gitlab".to_string()
    }

    /// Get the token audience
    fn token_audience(&self) -> String {
        "gitlab-api".to_string()
    }

    /// Authenticate without a session
    fn authenticate_without_session(&self) -> Result<Arc<User>, HttpResponse> {
        let token = match self.auth_token() {
            Some(token) => token,
            None => {
                return Err(HttpResponse::Unauthorized().json(json!({
                    "error": "Unauthorized",
                    "message": "Authentication token is required"
                })))
            }
        };

        // Verify token and get user
        match self.verify_token(&token) {
            Ok(user) => Ok(user),
            Err(_) => Err(HttpResponse::Unauthorized().json(json!({
                "error": "Unauthorized",
                "message": "Invalid authentication token"
            }))),
        }
    }

    /// Generate a sessionless authentication token
    fn generate_sessionless_token(&self, user: &Arc<User>) -> Result<String, HttpResponse> {
        let settings = Settings::current();
        let expiration = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_secs()
            + settings.sessionless_auth_expiration;

        let claims = SessionlessClaims {
            sub: user.id.to_string(),
            exp: expiration as usize,
            iat: SystemTime::now()
                .duration_since(UNIX_EPOCH)
                .unwrap()
                .as_secs() as usize,
        };

        encode(
            &Header::default(),
            &claims,
            &EncodingKey::from_secret(settings.sessionless_auth_secret.as_bytes()),
        )
        .map_err(|_| {
            HttpResponse::InternalServerError().json(json!({
                "error": "Internal Server Error",
                "message": "Failed to generate authentication token"
            }))
        })
    }

    /// Verify a sessionless authentication token
    fn verify_token(&self, token: &str) -> Result<Arc<User>, HttpResponse> {
        let settings = Settings::current();

        // TODO: Implement actual token verification
        // This would typically involve:
        // 1. Decoding and verifying the JWT token
        // 2. Extracting the user ID from the claims
        // 3. Loading the user from the database

        Err(HttpResponse::Unauthorized().json(json!({
            "error": "Unauthorized",
            "message": "Invalid authentication token"
        })))
    }

    /// Generate a new authentication token
    fn generate_auth_token(&self) -> Result<String, HttpResponse> {
        // TODO: Implement actual token generation
        // This would typically involve:
        // 1. Creating a JWT token
        // 2. Adding necessary claims
        // 3. Signing the token
        Ok("dummy_token".to_string())
    }

    /// Validate the authentication token
    fn validate_auth_token(&self) -> Result<bool, HttpResponse> {
        // TODO: Implement actual token validation
        // This would typically involve:
        // 1. Verifying the token signature
        // 2. Checking token expiration
        // 3. Validating claims
        Ok(true)
    }

    /// Get token claims
    fn get_token_claims(&self) -> Result<HashMap<String, String>, HttpResponse> {
        // TODO: Implement actual claims extraction
        // This would typically involve:
        // 1. Decoding the JWT token
        // 2. Extracting claims
        // 3. Validating claim format
        let mut claims = HashMap::new();

        claims.insert("iss".to_string(), self.token_issuer());
        claims.insert("aud".to_string(), self.token_audience());
        claims.insert(
            "exp".to_string(),
            (SystemTime::now() + self.token_expiration())
                .duration_since(SystemTime::UNIX_EPOCH)
                .unwrap_or(Duration::from_secs(0))
                .as_secs()
                .to_string(),
        );

        Ok(claims)
    }

    /// Enforce sessionless authentication
    fn enforce_sessionless_auth(&self) -> Result<(), HttpResponse> {
        if !self.validate_auth_token()? {
            return Err(HttpResponse::Unauthorized().json(serde_json::json!({
                "error": "Invalid authentication token"
            })));
        }
        Ok(())
    }

    /// Get authentication status
    fn get_auth_status(&self) -> Result<HashMap<String, bool>, HttpResponse> {
        let mut status = HashMap::new();

        status.insert("authenticated".to_string(), self.validate_auth_token()?);
        status.insert("token_valid".to_string(), self.validate_auth_token()?);

        Ok(status)
    }

    /// Sign in a user without creating a session
    fn sign_in(&self, user: Arc<User>, store: bool, message: &str);
    fn sign_in_without_callbacks(&self, user: Arc<User>, message: &str);
    fn settings(&self) -> Arc<Settings>;
    fn request(&self) -> &HttpRequest;
}

#[derive(Debug, Serialize, Deserialize)]
struct SessionlessClaims {
    sub: String,
    exp: usize,
    iat: usize,
}
