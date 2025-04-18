use crate::models::user::User;
use crate::settings::ApplicationSettings;
use actix_web::{web, HttpRequest, HttpResponse};
use std::sync::Arc;

pub trait SearchRateLimitable {
    fn check_search_rate_limit(
        &self,
        req: &HttpRequest,
        current_user: Option<Arc<User>>,
    ) -> Result<(), HttpResponse>;
    fn safe_search_scope(&self, req: &HttpRequest) -> Option<String>;
}

pub struct SearchRateLimitableHandler {
    settings: Arc<ApplicationSettings>,
}

impl SearchRateLimitableHandler {
    pub fn new(settings: Arc<ApplicationSettings>) -> Self {
        SearchRateLimitableHandler { settings }
    }
}

impl SearchRateLimitable for SearchRateLimitableHandler {
    fn check_search_rate_limit(
        &self,
        req: &HttpRequest,
        current_user: Option<Arc<User>>,
    ) -> Result<(), HttpResponse> {
        if let Some(user) = current_user {
            // For authenticated users, apply rate limits based on search scope
            let scope = self.safe_search_scope(req);
            let scope_key = scope.as_deref().unwrap_or("global");

            // Check if user is in allowlist
            let is_allowed = self
                .settings
                .search_rate_limit_allowlist
                .iter()
                .any(|allowed_user| allowed_user == &user.username());

            if is_allowed {
                return Ok(());
            }

            // Apply rate limiting based on user and scope
            self.check_rate_limit(
                "search_rate_limit",
                &[user.id().to_string(), scope_key.to_string()],
                self.settings.search_rate_limit,
            )?;
        } else {
            // For unauthenticated users, apply rate limits based on IP
            let ip = req.connection_info().peer_addr().unwrap_or("unknown");

            self.check_rate_limit(
                "search_rate_limit_unauthenticated",
                &[ip.to_string()],
                self.settings.search_rate_limit_unauthenticated,
            )?;
        }

        Ok(())
    }

    fn safe_search_scope(&self, req: &HttpRequest) -> Option<String> {
        // Extract scope from query parameters
        let query = req.query_string();
        let params: std::collections::HashMap<String, String> =
            form_urlencoded::parse(query.as_bytes())
                .map(|(k, v)| (k.into_owned(), v.into_owned()))
                .collect();

        // Check if scope parameter exists and is not abusive
        if let Some(scope) = params.get("scope") {
            if !self.is_abusive_search(scope) {
                return Some(scope.clone());
            }
        }

        None
    }
}

impl SearchRateLimitableHandler {
    fn check_rate_limit(
        &self,
        key: &str,
        scope: &[String],
        limit: i32,
    ) -> Result<(), HttpResponse> {
        // In a real implementation, this would use Redis or another rate limiting mechanism
        // For now, we'll just simulate rate limiting

        // This is a placeholder implementation
        // In a real app, you would use a rate limiting library like governor or governor-ratelimit

        // For demonstration purposes, we'll just return Ok
        Ok(())
    }

    fn is_abusive_search(&self, scope: &str) -> bool {
        // Check if the search scope is abusive
        // This could include checks for:
        // - Excessive length
        // - Invalid characters
        // - SQL injection attempts
        // - etc.

        // For now, we'll just check for excessive length
        scope.len() > 1000
    }
}

// This would be implemented in a separate module
pub mod settings {
    use serde::{Deserialize, Serialize};
    use std::sync::Arc;

    #[derive(Debug, Clone, Serialize, Deserialize)]
    pub struct ApplicationSettings {
        pub search_rate_limit: i32,
        pub search_rate_limit_unauthenticated: i32,
        pub search_rate_limit_allowlist: Vec<String>,
    }

    impl ApplicationSettings {
        pub fn new() -> Self {
            ApplicationSettings {
                search_rate_limit: 60,                 // 60 requests per minute
                search_rate_limit_unauthenticated: 30, // 30 requests per minute
                search_rate_limit_allowlist: Vec::new(),
            }
        }
    }
}
