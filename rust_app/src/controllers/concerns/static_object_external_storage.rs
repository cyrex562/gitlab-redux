use crate::config::settings::Settings;
use crate::models::project::Project;
use actix_web::{web, HttpResponse};

/// Module for handling external storage of static objects
pub trait StaticObjectExternalStorage {
    /// Redirect to external storage if not already an external storage request
    fn redirect_to_external_storage(&self, project: Option<&Project>) -> HttpResponse {
        if !self.external_storage_request() {
            let settings = Settings::current();
            if settings.static_objects_external_storage_enabled {
                let url = self.external_storage_url_or_path(project);
                HttpResponse::Found().header("Location", url).finish()
            } else {
                HttpResponse::Ok().finish()
            }
        } else {
            HttpResponse::Ok().finish()
        }
    }

    /// Check if this is an external storage request
    fn external_storage_request(&self) -> bool {
        let header_token = self
            .request_headers()
            .get("X-Gitlab-External-Storage-Token")
            .and_then(|v| v.to_str().ok());

        if let Some(token) = header_token {
            let settings = Settings::current();
            let storage_token = settings.static_objects_external_storage_auth_token.as_str();

            // Use constant-time comparison to prevent timing attacks
            if !constant_time_eq(token.as_bytes(), storage_token.as_bytes()) {
                return false;
            }
            true
        } else {
            false
        }
    }

    /// Get external storage URL or path
    fn external_storage_url_or_path(&self, project: Option<&Project>) -> String;

    /// Get request headers
    fn request_headers(&self) -> &actix_web::http::header::HeaderMap;
}

/// Constant-time string comparison to prevent timing attacks
fn constant_time_eq(a: &[u8], b: &[u8]) -> bool {
    if a.len() != b.len() {
        return false;
    }
    a.iter().zip(b.iter()).fold(0, |acc, (x, y)| acc | (x ^ y)) == 0
}
