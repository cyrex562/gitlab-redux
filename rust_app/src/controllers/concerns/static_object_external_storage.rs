// Ported from: orig_app/app/controllers/concerns/static_object_external_storage.rb
//
// This trait provides logic for redirecting to external storage and authenticating requests
// using an external storage token. Intended for use in Actix-web controllers.

use crate::config::settings::Settings;
use actix_web::{http::header, web, HttpRequest, HttpResponse};
use subtle::ConstantTimeEq;

pub trait StaticObjectExternalStorage {
    /// Checks if the request is from external storage by validating the token header.
    fn external_storage_request(&self, req: &HttpRequest) -> bool {
        let header_token = req.headers().get("X-Gitlab-External-Storage-Token");
        let header_token = match header_token.and_then(|h| h.to_str().ok()) {
            Some(token) if !token.is_empty() => token,
            _ => return false,
        };
        let settings = Settings::current();
        let external_storage_token = settings.static_objects_external_storage_auth_token.as_str();
        // Use constant-time comparison for security
        if header_token
            .as_bytes()
            .ct_eq(external_storage_token.as_bytes())
            .into()
        {
            true
        } else {
            // In Ruby: raise(Gitlab::Access::AccessDeniedError)
            // In Rust: could return false or handle error as needed
            false
        }
    }

    /// Redirects to external storage if the request is not already from external storage.
    fn redirect_to_external_storage(
        &self,
        req: &HttpRequest,
        project: &str,
    ) -> Option<HttpResponse> {
        if self.external_storage_request(req) {
            return None;
        }
        // Use the existing concern for URL generation
        let fullpath = req
            .uri()
            .path_and_query()
            .map(|pq| pq.as_str())
            .unwrap_or("");
        let url = crate::controllers::concerns::static_object_external_storage_url::StaticObjectExternalStorageUrl::get_external_storage_url(self, fullpath);
        if let Some(external_url) = url {
            Some(
                HttpResponse::Found()
                    .append_header((header::LOCATION, external_url))
                    .finish(),
            )
        } else {
            None
        }
    }
}
