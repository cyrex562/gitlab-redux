use actix_web::{web, HttpResponse};
use std::sync::Arc;

use crate::{
    models::user::User,
    services::{database::DatabaseService, logging::LoggingService},
    utils::{error::AppError, i18n::I18n},
};

/// Module for handling graceful timeouts
pub trait GracefulTimeoutHandling {
    /// Handle query timeout
    fn handle_query_timeout(&self, error: &AppError) -> Result<HttpResponse, AppError> {
        // Check if request is JSON
        if !self.is_json_request() {
            return Err(error.clone());
        }

        // Log exception
        self.log_exception(error);

        // Return error response
        Ok(HttpResponse::RequestTimeout().json(serde_json::json!({
            "error": I18n::t("There is too much data to calculate. Please change your selection.")
        })))
    }

    /// Check if request is JSON
    fn is_json_request(&self) -> bool {
        self.request()
            .headers()
            .get("accept")
            .and_then(|v| v.to_str().ok())
            .map(|v| v.contains("application/json"))
            .unwrap_or(false)
    }

    /// Log exception
    fn log_exception(&self, error: &AppError) {
        LoggingService::error(
            "Query timeout",
            &[
                ("error", error.to_string()),
                ("controller", self.controller_name()),
                ("action", self.action_name()),
            ],
        );
    }

    // Required trait methods that need to be implemented by the controller
    fn request(&self) -> &web::HttpRequest;
    fn controller_name(&self) -> String;
    fn action_name(&self) -> String;
}
