use crate::config::settings::Settings;
use crate::models::user::User;
use actix_web::{web, HttpResponse};
use jsonwebtoken::{encode, EncodingKey, Header};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::time::{Duration, SystemTime, UNIX_EPOCH};

/// Module for handling sessionless authentication
pub trait SessionlessAuthentication {
    /// Get the current user
    fn current_user(&self) -> Option<&User>;

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
    fn authenticate_without_session(&self) -> Result<User, HttpResponse> {
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
    fn generate_sessionless_token(&self, user: &User) -> Result<String, HttpResponse> {
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
    fn verify_token(&self, token: &str) -> Result<User, HttpResponse> {
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
}

#[derive(Debug, Serialize, Deserialize)]
struct SessionlessClaims {
    sub: String,
    exp: usize,
    iat: usize,
}
