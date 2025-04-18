use crate::config::Settings;
use actix_web::{
    dev::ServiceRequest, error::Error, http::StatusCode, web::Data, HttpResponse, Responder,
};
use std::net::IpAddr;

pub trait RequiresAllowlistedMonitoringClient {
    fn validate_ip_allowlisted_or_valid_token(&self, req: &ServiceRequest) -> Result<(), Error>;
}

pub struct RequiresAllowlistedMonitoringClientImpl {
    settings: Data<Settings>,
}

impl RequiresAllowlistedMonitoringClientImpl {
    pub fn new(settings: Data<Settings>) -> Self {
        Self { settings }
    }

    fn client_ip_allowlisted(&self, req: &ServiceRequest) -> bool {
        // Always allow localhost in development
        if cfg!(debug_assertions)
            && req
                .connection_info()
                .peer_addr()
                .map_or(false, |addr| addr == "127.0.0.1")
        {
            return true;
        }

        let client_ip = req
            .connection_info()
            .peer_addr()
            .and_then(|addr| addr.parse::<IpAddr>().ok());

        if let Some(ip) = client_ip {
            self.settings
                .monitoring
                .ip_whitelist
                .iter()
                .any(|allowed_ip| {
                    allowed_ip
                        .parse::<IpAddr>()
                        .map(|allowed| allowed == ip)
                        .unwrap_or(false)
                })
        } else {
            false
        }
    }

    fn valid_token(&self, req: &ServiceRequest) -> bool {
        let token = req
            .query_string()
            .split('&')
            .find(|param| param.starts_with("token="))
            .and_then(|param| param.split('=').nth(1))
            .or_else(|| req.headers().get("TOKEN").and_then(|h| h.to_str().ok()))
            .unwrap_or("");

        if token.is_empty() {
            return false;
        }

        // Use constant-time comparison for security
        let expected_token = self.settings.health_check_access_token.as_str();
        if token.len() != expected_token.len() {
            return false;
        }

        token
            .bytes()
            .zip(expected_token.bytes())
            .all(|(a, b)| a == b)
    }

    fn render_404(&self) -> HttpResponse {
        HttpResponse::NotFound()
            .content_type("text/html")
            .body(include_str!("../../templates/errors/not_found.html"))
    }
}

impl RequiresAllowlistedMonitoringClient for RequiresAllowlistedMonitoringClientImpl {
    fn validate_ip_allowlisted_or_valid_token(&self, req: &ServiceRequest) -> Result<(), Error> {
        if self.client_ip_allowlisted(req) || self.valid_token(req) {
            Ok(())
        } else {
            Err(actix_web::error::ErrorForbidden("Access denied"))
        }
    }
}
