use actix_web::{web, HttpResponse};
use std::sync::Arc;

use crate::{
    models::{oauth::OAuthApplication, user::User},
    services::{
        auth::AuthUtils,
        oauth::{OAuthConfig, OAuthScopes},
    },
    utils::{error::AppError, session::Session},
};

/// Module for handling OAuth applications
pub trait OauthApplications {
    /// Session key for created OAuth applications
    const CREATED_SESSION_KEY: &'static str = "oauth_applications_created";

    /// Prepare scopes for OAuth application
    fn prepare_scopes(&self) -> Result<(), AppError> {
        let mut params = self.params().clone();

        if let Some(doorkeeper_app) = params.get_mut("doorkeeper_application") {
            if let Some(scopes) = doorkeeper_app.get_mut("scopes") {
                if let Some(scopes_vec) = scopes.as_array() {
                    let scopes_str = scopes_vec
                        .iter()
                        .filter_map(|s| s.as_str())
                        .collect::<Vec<&str>>()
                        .join(" ");

                    *scopes = serde_json::Value::String(scopes_str);
                }
            }
        }

        Ok(())
    }

    /// Set created session
    fn set_created_session(&self) {
        self.session().set(Self::CREATED_SESSION_KEY, true);
    }

    /// Get created session
    fn get_created_session(&self) -> bool {
        self.session()
            .get::<bool>(Self::CREATED_SESSION_KEY)
            .unwrap_or(false)
    }

    /// Load scopes
    fn load_scopes(&self) -> OAuthScopes {
        let config = OAuthConfig::default();
        let mut scopes = config.scopes();

        // Remove restricted scopes
        scopes.retain(|scope| {
            !matches!(
                scope.as_str(),
                "ai_workflow" | "dynamic_user" | "self_rotate"
            )
        });

        OAuthScopes::from_vec(scopes)
    }

    /// Get permitted parameters
    fn permitted_params(&self) -> Vec<&'static str> {
        vec!["name", "redirect_uri", "scopes", "confidential"]
    }

    /// Get application parameters
    fn application_params(&self) -> Result<serde_json::Value, AppError> {
        let params = self.params();
        let doorkeeper_app = params.get("doorkeeper_application").ok_or_else(|| {
            AppError::BadRequest("Missing doorkeeper_application parameter".to_string())
        })?;

        let mut permitted = serde_json::Map::new();

        for param in self.permitted_params() {
            if let Some(value) = doorkeeper_app.get(param) {
                permitted.insert(param.to_string(), value.clone());
            }
        }

        Ok(serde_json::Value::Object(permitted))
    }

    // Required trait methods that need to be implemented by the controller
    fn current_user(&self) -> Option<&User>;
    fn params(&self) -> &serde_json::Value;
    fn session(&self) -> &dyn Session;
}
