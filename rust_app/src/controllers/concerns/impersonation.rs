// Ported from orig_app/app/controllers/concerns/impersonation.rb on 2025-04-25
// This file implements the Impersonation concern from the Ruby codebase.

use actix_web::{web, HttpRequest, HttpResponse};
use once_cell::sync::Lazy;
use serde::{Deserialize, Serialize};
use std::collections::HashSet;
use std::sync::Arc;

use crate::models::user::User;
use crate::utils::logger::AppLogger;

pub trait Impersonation {
    fn current_user(&self) -> User;
    fn check_impersonation_availability(&self) -> Result<(), HttpResponse>;
    fn stop_impersonation(&self) -> User;
    fn impersonation_in_progress(&self) -> bool;
    fn log_impersonation_event(&self);
    fn clear_access_token_session_keys(&self);
    fn impersonator(&self) -> Option<User>;
}

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

pub struct ImpersonationImpl {
    session: web::Data<Session>,
    config: web::Data<Config>,
    logger: Arc<AppLogger>,
}

impl ImpersonationImpl {
    pub fn new(
        session: web::Data<Session>,
        config: web::Data<Config>,
        logger: Arc<AppLogger>,
    ) -> Self {
        Self {
            session,
            config,
            logger,
        }
    }
}

impl Impersonation for ImpersonationImpl {
    fn current_user(&self) -> User {
        let mut user = self.session.get_current_user();

        if let Some(impersonator) = self.impersonator() {
            user.set_impersonator(impersonator);
        }

        user
    }

    fn check_impersonation_availability(&self) -> Result<(), HttpResponse> {
        if !self.impersonation_in_progress() {
            return Ok(());
        }

        if !self.config.gitlab.impersonation_enabled {
            self.stop_impersonation();
            return Err(HttpResponse::Forbidden().body("Impersonation has been disabled"));
        }

        Ok(())
    }

    fn stop_impersonation(&self) -> User {
        self.log_impersonation_event();

        if let Some(impersonator) = self.impersonator() {
            self.session.set_user(impersonator);
        }

        self.session.remove("impersonator_id");
        self.clear_access_token_session_keys();

        self.current_user()
    }

    fn impersonation_in_progress(&self) -> bool {
        self.session.get::<i32>("impersonator_id").is_some()
    }

    fn log_impersonation_event(&self) {
        if let (Some(impersonator), Some(current_user)) =
            (self.impersonator(), self.session.get_current_user())
        {
            self.logger.info(&format!(
                "User {} has stopped impersonating {}",
                impersonator.username, current_user.username
            ));
        }
    }

    fn clear_access_token_session_keys(&self) {
        let session_keys: Vec<String> = self
            .session
            .keys()
            .filter(|key| SESSION_KEYS_TO_DELETE.contains(key.as_str()))
            .cloned()
            .collect();

        for key in session_keys {
            self.session.remove(&key);
        }
    }

    fn impersonator(&self) -> Option<User> {
        if let Some(impersonator_id) = self.session.get::<i32>("impersonator_id") {
            User::find(impersonator_id)
        } else {
            None
        }
    }
}
