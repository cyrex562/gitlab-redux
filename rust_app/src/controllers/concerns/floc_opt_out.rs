use actix_web::{web, HttpRequest, HttpResponse};
use std::sync::Arc;

use crate::config::settings::Settings;

/// Module for FLoC opt-out
pub trait FlocOptOut {
    /// Check if FLoC is enabled
    fn floc_enabled(&self) -> bool {
        self.settings().floc_enabled()
    }
    
    /// Set the FLoC opt-out header
    fn set_floc_opt_out_header(&self, response: &mut HttpResponse) {
        // Only set the header if FLoC is not enabled
        if !self.floc_enabled() {
            response.headers_mut().insert(
                "Permissions-Policy",
                "interest-cohort=()".parse().unwrap()
            );
        }
    }
    
    // Required trait methods that need to be implemented by the controller
    fn settings(&self) -> Arc<Settings>;
} 