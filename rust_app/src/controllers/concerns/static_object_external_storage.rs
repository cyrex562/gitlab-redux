use actix_web::{web, HttpRequest, HttpResponse};
use std::sync::Arc;

use crate::config::settings::Settings;
use crate::models::project::Project;
use crate::utils::security::SecurityUtils;

/// Module for handling static object external storage
pub trait StaticObjectExternalStorage {
    /// Redirect to external storage
    fn redirect_to_external_storage(&self, req: &HttpRequest) -> HttpResponse {
        if self.external_storage_request(req) {
            return HttpResponse::Ok().finish();
        }

        let project = self.project();
        let full_path = req.uri().path().to_string();
        let redirect_url = self.external_storage_url_or_path(&full_path, project);

        HttpResponse::Found()
            .header("Location", redirect_url)
            .finish()
    }

    /// Check if request is from external storage
    fn external_storage_request(&self, req: &HttpRequest) -> bool {
        let header_token = req
            .headers()
            .get("X-Gitlab-External-Storage-Token")
            .and_then(|h| h.to_str().ok());

        if let Some(header_token) = header_token {
            let settings = self.settings();
            let external_storage_token = settings.static_objects_external_storage_auth_token();

            if SecurityUtils::secure_compare(header_token, external_storage_token) {
                return true;
            }

            // If tokens don't match, raise access denied error
            panic!("Access denied: Invalid external storage token");
        }

        false
    }

    // Required trait methods that need to be implemented by the controller
    fn project(&self) -> Option<Arc<dyn Project>>;
    fn settings(&self) -> Arc<Settings>;
    fn external_storage_url_or_path(&self, path: &str, project: Option<Arc<dyn Project>>)
        -> String;
}

/// Trait for projects
pub trait Project: Send + Sync {
    /// Get project ID
    fn id(&self) -> i32;

    /// Get project path
    fn path(&self) -> String;

    /// Get project name
    fn name(&self) -> String;
}
