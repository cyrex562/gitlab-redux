use actix_web::{HttpRequest, HttpResponse};
use std::net::IpAddr;
use std::sync::OnceLock;

pub struct Settings {
    pub monitoring: MonitoringSettings,
}

pub struct MonitoringSettings {
    pub ip_whitelist: Vec<String>,
}

impl Settings {
    pub fn get() -> &'static Settings {
        static INSTANCE: OnceLock<Settings> = OnceLock::new();
        INSTANCE.get_or_init(|| Settings {
            monitoring: MonitoringSettings {
                ip_whitelist: Vec::new(),
            },
        })
    }
}

pub struct CurrentSettings {
    pub health_check_access_token: String,
}

impl CurrentSettings {
    pub fn get() -> &'static CurrentSettings {
        static INSTANCE: OnceLock<CurrentSettings> = OnceLock::new();
        INSTANCE.get_or_init(|| CurrentSettings {
            health_check_access_token: String::new(),
        })
    }
}

pub trait RequiresAllowlistedMonitoringClient {
    fn validate_ip_allowlisted_or_valid_token(&self, req: &HttpRequest) -> HttpResponse {
        if self.client_ip_allowlisted(req) || self.valid_token(req) {
            HttpResponse::Ok().finish()
        } else {
            self.render_404()
        }
    }

    fn client_ip_allowlisted(&self, req: &HttpRequest) -> bool {
        // Always allow developers to access http://localhost:3000/-/metrics for
        // debugging purposes
        if cfg!(debug_assertions)
            && req
                .connection_info()
                .peer_addr()
                .map_or(false, |addr| addr.contains("127.0.0.1"))
        {
            return true;
        }

        let client_ip = req
            .connection_info()
            .peer_addr()
            .and_then(|addr| addr.parse::<IpAddr>().ok())
            .unwrap_or_else(|| "0.0.0.0".parse().unwrap());

        self.ip_allowlist().iter().any(|ip| ip.contains(&client_ip))
    }

    fn ip_allowlist(&self) -> Vec<IpAddr> {
        Settings::get()
            .monitoring
            .ip_whitelist
            .iter()
            .filter_map(|ip| ip.parse().ok())
            .collect()
    }

    fn valid_token(&self, req: &HttpRequest) -> bool {
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

        // Use a constant-time comparison to prevent timing attacks
        self.constant_time_compare(token, &CurrentSettings::get().health_check_access_token)
    }

    fn constant_time_compare(&self, a: &str, b: &str) -> bool {
        if a.len() != b.len() {
            return false;
        }

        let mut result = 0u8;
        for (x, y) in a.bytes().zip(b.bytes()) {
            result |= x ^ y;
        }
        result == 0
    }

    fn render_404(&self) -> HttpResponse {
        HttpResponse::NotFound()
            .content_type("text/html")
            .body(include_str!("../../../templates/errors/not_found.html"))
    }
}
