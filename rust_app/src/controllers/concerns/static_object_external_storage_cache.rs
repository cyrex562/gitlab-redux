use crate::config::settings::Settings;
use actix_web::{web, HttpResponse};
use std::time::{Duration, SystemTime};

/// Module for handling external storage caching
pub trait StaticObjectExternalStorageCache {
    /// Set cache headers for external storage responses
    fn set_external_storage_cache_headers(&self, response: &mut HttpResponse) {
        let settings = Settings::current();
        if !settings.static_objects_external_storage_enabled {
            return;
        }

        // Set cache control headers
        if let Some(max_age) = settings.static_objects_external_storage_cache_max_age {
            response.headers_mut().insert(
                "Cache-Control",
                format!("public, max-age={}", max_age).parse().unwrap(),
            );
        }

        // Set expires header
        if let Some(max_age) = settings.static_objects_external_storage_cache_max_age {
            if let Ok(expires) = SystemTime::now() + Duration::from_secs(max_age as u64) {
                response.headers_mut().insert(
                    "Expires",
                    expires
                        .duration_since(SystemTime::UNIX_EPOCH)
                        .unwrap()
                        .as_secs()
                        .to_string()
                        .parse()
                        .unwrap(),
                );
            }
        }

        // Set ETag header if enabled
        if settings.static_objects_external_storage_etag_enabled {
            if let Some(etag) = self.generate_etag() {
                response.headers_mut().insert("ETag", etag.parse().unwrap());
            }
        }
    }

    /// Generate an ETag for the current resource
    fn generate_etag(&self) -> Option<String> {
        None // Default implementation returns None, can be overridden
    }
}
