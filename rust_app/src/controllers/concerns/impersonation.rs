use actix_web::{web, HttpRequest};
use once_cell::sync::Lazy;
use std::collections::HashSet;
use std::sync::Arc;

use crate::config::settings::Settings;
use crate::models::user::User;
use crate::utils::logger::AppLogger;
use crate::utils::session::Session;

/// Session keys that should be deleted during impersonation
static SESSION_KEYS_TO_DELETE: Lazy<HashSet<&'static str>> = Lazy::new(|| {
    let mut set = HashSet::new();
    set.insert("github_access_token");
    set.insert("gitea_access_token");
    set.insert("gitlab_access_token");
    set.insert("bitbucket_token");
    set.insert("bitbucket_refresh_token");
    set.insert("bitbucket_server_personal_access_token");
    set.insert("bulk_import_gitlab_access_token");
    set.insert("fogbugz_token");
    set.insert("cloud_platform_access_token");
    set
});

/// Module for handling user impersonation
pub trait Impersonation {
    /// Get the current user with impersonation information
    fn current_user_with_impersonation(&self) -> Option<Arc<User>> {
        let mut user = self.current_user()?;

        if let Some(impersonator) = self.impersonator() {
            user.set_impersonator(impersonator);
        }

        Some(user)
    }

    /// Check if impersonation is available
    fn check_impersonation_availability(&self) -> Result<(), String> {
        if self.impersonation_in_progress() {
            if !self.settings().impersonation_enabled() {
                self.stop_impersonation();
                return Err("Impersonation has been disabled".to_string());
            }
        }

        Ok(())
    }

    /// Stop impersonation
    fn stop_impersonation(&self) -> Option<Arc<User>> {
        self.log_impersonation_event();

        // Set the impersonator as the current user
        if let Some(impersonator) = self.impersonator() {
            self.set_current_user(impersonator);
        }

        // Clear impersonator from session
        self.session().remove("impersonator_id");

        // Clear access token session keys
        self.clear_access_token_session_keys();

        self.current_user()
    }

    /// Check if impersonation is in progress
    fn impersonation_in_progress(&self) -> bool {
        self.session().get::<i64>("impersonator_id").is_some()
    }

    /// Log impersonation event
    fn log_impersonation_event(&self) {
        if let (Some(impersonator), Some(current_user)) = (self.impersonator(), self.current_user())
        {
            AppLogger::info(&format!(
                "User {} has stopped impersonating {}",
                impersonator.username(),
                current_user.username()
            ));
        }
    }

    /// Clear access token session keys
    fn clear_access_token_session_keys(&self) {
        let session = self.session();
        let keys_to_delete: Vec<String> = session
            .keys()
            .filter(|key| SESSION_KEYS_TO_DELETE.contains(key.as_str()))
            .cloned()
            .collect();

        for key in keys_to_delete {
            session.remove(&key);
        }
    }

    /// Get the impersonator
    fn impersonator(&self) -> Option<Arc<User>> {
        if let Some(impersonator_id) = self.session().get::<i64>("impersonator_id") {
            return User::find(impersonator_id);
        }

        None
    }

    // Required trait methods that need to be implemented by the controller
    fn current_user(&self) -> Option<Arc<User>>;
    fn set_current_user(&self, user: Arc<User>);
    fn settings(&self) -> Arc<Settings>;
    fn session(&self) -> &Session;
}
