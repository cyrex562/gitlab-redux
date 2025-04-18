use crate::utils::{device_detector::SafeDeviceDetector, metrics::Metrics};
use actix_web::{dev::ServiceRequest, error::Error, HttpResponse};
use std::sync::Arc;

#[derive(Debug)]
pub enum PageLimiterError {
    PageLimitNotANumber,
    PageLimitNotSensible,
    PageOutOfBounds(i32),
}

impl std::fmt::Display for PageLimiterError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            Self::PageLimitNotANumber => write!(f, "Page limit must be a number"),
            Self::PageLimitNotSensible => write!(f, "Page limit must be greater than 0"),
            Self::PageOutOfBounds(max) => write!(f, "Page number exceeds maximum of {}", max),
        }
    }
}

impl std::error::Error for PageLimiterError {}

pub trait PageLimiter {
    fn limit_pages(&self, max_page_number: i32) -> Result<(), Error>;
    fn default_page_out_of_bounds_response(&self) -> HttpResponse;
}

pub struct PageLimiterImpl {
    metrics: Arc<Metrics>,
}

impl PageLimiterImpl {
    pub fn new(metrics: Arc<Metrics>) -> Self {
        Self { metrics }
    }

    fn check_page_number(
        &self,
        req: &ServiceRequest,
        max_page_number: i32,
    ) -> Result<(), PageLimiterError> {
        if max_page_number <= 0 {
            return Err(PageLimiterError::PageLimitNotSensible);
        }

        if let Some(page) = req.query_string().get("page") {
            if let Ok(page_num) = page.parse::<i32>() {
                if page_num > max_page_number {
                    self.record_page_limit_interception(req);
                    return Err(PageLimiterError::PageOutOfBounds(max_page_number));
                }
            }
        }

        Ok(())
    }

    fn record_page_limit_interception(&self, req: &ServiceRequest) {
        let device_detector = SafeDeviceDetector::new(req.headers().get("user-agent"));

        self.metrics
            .counter("gitlab_page_out_of_bounds")
            .with_label("controller", req.path())
            .with_label("action", req.method().as_str())
            .with_label("bot", device_detector.is_bot().to_string())
            .increment();
    }
}

impl PageLimiter for PageLimiterImpl {
    fn limit_pages(&self, max_page_number: i32) -> Result<(), Error> {
        self.check_page_number(req, max_page_number)
            .map_err(|e| actix_web::error::ErrorBadRequest(e.to_string()))
    }

    fn default_page_out_of_bounds_response(&self) -> HttpResponse {
        HttpResponse::BadRequest().finish()
    }
}
