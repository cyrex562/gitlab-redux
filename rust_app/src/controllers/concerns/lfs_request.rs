// Ported from: orig_app/app/controllers/concerns/lfs_request.rb
// This file implements the LfsRequest concern in Rust.

use actix_web::{HttpRequest, HttpResponse, Responder};
use chrono::Utc;
use serde_json::json;

const CONTENT_TYPE: &str = "application/vnd.git-lfs+json";

pub trait LfsRequest {
    // These methods/fields are expected to be implemented by the handler struct:
    // fn container(&self) -> Option<&Container>;
    // fn project(&self) -> Option<&Project>;
    // fn user(&self) -> Option<&User>;
    // fn deploy_token(&self) -> Option<&DeployToken>;
    // fn authentication_result(&self) -> &AuthenticationResult;
    // fn can(&self, object: &dyn Any, action: &str, subject: &dyn Any) -> bool;
    // fn ci(&self) -> bool;
    // fn download_request(&self) -> bool;
    // fn upload_request(&self) -> bool;
    // fn has_authentication_ability(&self, ability: &str) -> bool;
    // fn help_url(&self) -> String;

    fn require_lfs_enabled(&self, req: &HttpRequest) -> Option<HttpResponse> {
        if !self.lfs_enabled() {
            Some(HttpResponse::NotImplemented()
                .content_type(CONTENT_TYPE)
                .json(json!({
                    "message": "Git LFS is not enabled on this GitLab server, contact your admin.",
                    "documentation_url": self.help_url()
                })))
        } else {
            None
        }
    }

    fn lfs_check_access(&self, req: &HttpRequest) -> Option<HttpResponse> {
        if !self.container_lfs_enabled() {
            return Some(self.render_lfs_not_found());
        }
        if self.download_request() && self.lfs_download_access() {
            return None;
        }
        if self.upload_request() && self.lfs_upload_access() {
            return None;
        }
        if self.lfs_download_access() {
            Some(self.render_lfs_forbidden())
        } else {
            Some(self.render_lfs_not_found())
        }
    }

    fn render_lfs_forbidden(&self) -> HttpResponse {
        HttpResponse::Forbidden()
            .content_type(CONTENT_TYPE)
            .json(json!({
                "message": "Access forbidden. Check your access level.",
                "documentation_url": self.help_url()
            }))
    }

    fn render_lfs_not_found(&self) -> HttpResponse {
        HttpResponse::NotFound()
            .content_type(CONTENT_TYPE)
            .json(json!({
                "message": "Not found.",
                "documentation_url": self.help_url()
            }))
    }

    // The following methods should be implemented or delegated to the handler struct:
    fn lfs_enabled(&self) -> bool;
    fn container_lfs_enabled(&self) -> bool;
    fn lfs_download_access(&self) -> bool;
    fn lfs_upload_access(&self) -> bool;
    fn download_request(&self) -> bool;
    fn upload_request(&self) -> bool;
    fn help_url(&self) -> String;
}

// Handler struct and integration with Actix Web would be implemented in the actual controller.
