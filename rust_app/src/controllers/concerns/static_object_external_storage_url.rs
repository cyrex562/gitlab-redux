use crate::config::settings::Settings;
use actix_web::{web, HttpResponse};
use url::Url;

/// Module for handling external storage URLs
pub trait StaticObjectExternalStorageUrl {
    /// Get the external storage URL for a given path
    fn get_external_storage_url(&self, path: &str) -> Option<String> {
        let settings = Settings::current();
        if !settings.static_objects_external_storage_enabled {
            return None;
        }

        let base_url = settings.static_objects_external_storage_url.as_str();
        let url = Url::parse(base_url).ok()?;
        let mut url = url.join(path).ok()?;

        // Add query parameters if they exist
        if let Some(query) = settings
            .static_objects_external_storage_query_params
            .as_ref()
        {
            url.set_query(Some(query));
        }

        Some(url.to_string())
    }

    /// Check if a URL is for external storage
    fn is_external_storage_url(&self, url: &str) -> bool {
        let settings = Settings::current();
        if !settings.static_objects_external_storage_enabled {
            return false;
        }

        settings
            .static_objects_hosts
            .iter()
            .any(|host| url.starts_with(host))
    }
}
