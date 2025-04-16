use crate::config::settings::Settings;
use actix_web::{web, HttpResponse};
use std::collections::HashMap;

/// Module for handling CSP for static object external storage
pub trait StaticObjectExternalStorageCsp {
    /// Set CSP headers for external storage
    fn set_external_storage_csp(&self) -> HttpResponse {
        let settings = Settings::current();
        if !settings.static_objects_external_storage_enabled {
            return HttpResponse::Ok().finish();
        }

        // Build CSP directives
        let mut directives = Vec::new();

        // Default directives
        directives.push("default-src 'none'".to_string());
        directives.push("script-src 'none'".to_string());
        directives.push("style-src 'none'".to_string());

        // Image and font sources
        let image_src = format!(
            "img-src 'self' data: {}",
            settings.static_objects_hosts.join(" ")
        );
        directives.push(image_src);

        let font_src = format!(
            "font-src 'self' data: {}",
            settings.static_objects_hosts.join(" ")
        );
        directives.push(font_src);

        // Connect sources
        let connect_src = format!(
            "connect-src 'self' {}",
            settings.static_objects_hosts.join(" ")
        );
        directives.push(connect_src);

        // Add any additional CSP directives from settings
        if let Some(additional_directives) = &settings.static_objects_csp_directives {
            directives.extend(additional_directives.clone());
        }

        // Join directives with semicolons
        let csp = directives.join("; ");

        // Create response with CSP header
        HttpResponse::Ok()
            .header("Content-Security-Policy", csp)
            .finish()
    }

    /// Get CSP settings
    fn get_csp_settings(&self) -> HashMap<String, String> {
        let mut settings = HashMap::new();
        let settings = Settings::current();

        settings.insert(
            "static_objects_external_storage_enabled".to_string(),
            settings.static_objects_external_storage_enabled.to_string(),
        );

        if let Some(directives) = &settings.static_objects_csp_directives {
            settings.insert(
                "static_objects_csp_directives".to_string(),
                directives.join(", "),
            );
        }

        settings
    }
}
