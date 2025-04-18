use actix_web::{web, HttpRequest, HttpResponse, Result};
use std::sync::Arc;

// Define the Settings trait
pub trait Settings: Send + Sync {
    fn static_objects_external_storage_enabled(&self) -> bool;
    fn static_objects_external_storage_url(&self) -> Option<String>;
}

// Define the ContentSecurityPolicy trait
pub trait ContentSecurityPolicy: Send + Sync {
    fn directives(&self) -> Option<CSPDirectives>;
    fn set_connect_src(&mut self, values: Vec<String>);
}

// Define the CSPDirectives struct
pub struct CSPDirectives {
    pub connect_src: Option<Vec<String>>,
    pub default_src: Option<Vec<String>>,
}

// Define the StaticObjectExternalStorageCSP trait
pub trait StaticObjectExternalStorageCSP {
    fn configure_content_security_policy(&self, policy: &mut dyn ContentSecurityPolicy);
    fn get_settings(&self) -> Arc<dyn Settings>;
}

// Define the StaticObjectExternalStorageCSPHandler struct
pub struct StaticObjectExternalStorageCSPHandler {
    settings: Arc<dyn Settings>,
}

impl StaticObjectExternalStorageCSPHandler {
    pub fn new(settings: Arc<dyn Settings>) -> Self {
        StaticObjectExternalStorageCSPHandler { settings }
    }
}

// Implement the StaticObjectExternalStorageCSP trait for StaticObjectExternalStorageCSPHandler
impl StaticObjectExternalStorageCSP for StaticObjectExternalStorageCSPHandler {
    fn configure_content_security_policy(&self, policy: &mut dyn ContentSecurityPolicy) {
        if let Some(directives) = policy.directives() {
            if !self.settings.static_objects_external_storage_enabled() {
                return;
            }

            let default_connect_src = directives
                .connect_src
                .clone()
                .or_else(|| directives.default_src.clone());

            if let Some(default_connect_src) = default_connect_src {
                let mut connect_src_values = default_connect_src;

                if let Some(url) = self.settings.static_objects_external_storage_url() {
                    connect_src_values.push(url);
                }

                policy.set_connect_src(connect_src_values);
            }
        }
    }

    fn get_settings(&self) -> Arc<dyn Settings> {
        self.settings.clone()
    }
}
