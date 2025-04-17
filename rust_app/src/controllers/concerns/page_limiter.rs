use actix_web::{web, HttpResponse};
use std::sync::Arc;

use crate::{
    models::user::User,
    services::{
        device_detection::DeviceDetector, logging::LoggingService, metrics::MetricsService,
    },
    utils::{error::AppError, i18n::I18n},
};

/// Custom errors for page limiting
#[derive(Debug)]
pub enum PageLimiterError {
    PageLimitNotANumber,
    PageLimitNotSensible,
    PageOutOfBounds(i32),
}

impl std::fmt::Display for PageLimiterError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            PageLimiterError::PageLimitNotANumber => write!(f, "Page limit must be a number"),
            PageLimiterError::PageLimitNotSensible => {
                write!(f, "Page limit must be greater than 0")
            }
            PageLimiterError::PageOutOfBounds(max) => {
                write!(f, "Page number exceeds maximum of {}", max)
            }
        }
    }
}

impl std::error::Error for PageLimiterError {}

/// Module for handling page limiting
pub trait PageLimiter {
    /// Limit pages to a maximum number
    fn limit_pages(&self, max_page_number: i32) -> Result<(), AppError> {
        self.check_page_number(max_page_number)?;
        Ok(())
    }

    /// Check if page number is valid
    fn check_page_number(&self, max_page_number: i32) -> Result<(), AppError> {
        // Validate max page number
        if max_page_number <= 0 {
            return Err(AppError::from(PageLimiterError::PageLimitNotSensible));
        }

        // Get page from params
        let page = self.page_param();
        if page.is_none() {
            return Ok(());
        }

        let page = page.unwrap();
        if page > max_page_number {
            self.record_page_limit_interception();
            return Err(AppError::from(PageLimiterError::PageOutOfBounds(
                max_page_number,
            )));
        }

        Ok(())
    }

    /// Default response for page out of bounds
    fn default_page_out_of_bounds_response(&self) -> HttpResponse {
        HttpResponse::BadRequest().finish()
    }

    /// Record page limit interception in metrics
    fn record_page_limit_interception(&self) {
        let user_agent = self
            .request()
            .headers()
            .get("user-agent")
            .and_then(|v| v.to_str().ok())
            .unwrap_or("");

        let detector = DeviceDetector::new(user_agent);
        let is_bot = detector.is_bot();

        MetricsService::increment_counter(
            "gitlab_page_out_of_bounds",
            &[
                ("controller", self.controller_name()),
                ("action", self.action_name()),
                ("bot", &is_bot.to_string()),
            ],
        );
    }

    /// Get page from params
    fn page_param(&self) -> Option<i32> {
        self.request()
            .query_string()
            .split('&')
            .find(|param| param.starts_with("page="))
            .and_then(|param| param.split('=').nth(1))
            .and_then(|page| page.parse::<i32>().ok())
    }

    // Required trait methods that need to be implemented by the controller
    fn request(&self) -> &web::HttpRequest;
    fn controller_name(&self) -> String;
    fn action_name(&self) -> String;
}
