use actix_web::{web, HttpRequest, HttpResponse, Result};
use std::sync::Arc;

// Define the Project trait
pub trait Project: Send + Sync {
    fn id(&self) -> i64;
    fn name(&self) -> &str;
}

// Define the Settings trait
pub trait Settings: Send + Sync {
    fn static_objects_external_storage_auth_token(&self) -> Option<String>;
    fn static_objects_external_storage_enabled(&self) -> bool;
    fn static_objects_external_storage_url(&self) -> Option<String>;
}

// Define the StaticObjectExternalStorage trait
pub trait StaticObjectExternalStorage {
    fn redirect_to_external_storage(&self, req: &HttpRequest) -> Result<HttpResponse>;
    fn external_storage_request(&self, req: &HttpRequest) -> bool;
    fn external_storage_url_or_path(&self, path: &str, project: Option<&dyn Project>) -> String;
    fn get_settings(&self) -> Arc<dyn Settings>;
    fn get_project(&self) -> Option<Arc<dyn Project>>;
    fn get_request_path(&self, req: &HttpRequest) -> String;
    fn get_request_header(&self, req: &HttpRequest, name: &str) -> Option<String>;
}

// Define the StaticObjectExternalStorageHandler struct
pub struct StaticObjectExternalStorageHandler {
    settings: Arc<dyn Settings>,
}

impl StaticObjectExternalStorageHandler {
    pub fn new(settings: Arc<dyn Settings>) -> Self {
        StaticObjectExternalStorageHandler { settings }
    }
}

// Implement the StaticObjectExternalStorage trait for StaticObjectExternalStorageHandler
impl StaticObjectExternalStorage for StaticObjectExternalStorageHandler {
    fn redirect_to_external_storage(&self, req: &HttpRequest) -> Result<HttpResponse> {
        if self.external_storage_request(req) {
            return Ok(HttpResponse::Ok().finish());
        }

        let path = self.get_request_path(req);
        let project = self.get_project();
        let url = self.external_storage_url_or_path(&path, project.as_deref());

        Ok(HttpResponse::Found().header("Location", url).finish())
    }

    fn external_storage_request(&self, req: &HttpRequest) -> bool {
        if let Some(header_token) = self.get_request_header(req, "X-Gitlab-External-Storage-Token")
        {
            if let Some(external_storage_token) =
                self.settings.static_objects_external_storage_auth_token()
            {
                // Use a secure comparison to prevent timing attacks
                return secure_compare(&header_token, &external_storage_token);
            }
        }

        false
    }

    fn external_storage_url_or_path(&self, path: &str, project: Option<&dyn Project>) -> String {
        if let Some(url) = self.settings.static_objects_external_storage_url() {
            format!("{}{}", url, path)
        } else {
            path.to_string()
        }
    }

    fn get_settings(&self) -> Arc<dyn Settings> {
        self.settings.clone()
    }

    fn get_project(&self) -> Option<Arc<dyn Project>> {
        // This would be implemented by the concrete class
        None
    }

    fn get_request_path(&self, req: &HttpRequest) -> String {
        // This would be implemented by the concrete class
        req.uri().path().to_string()
    }

    fn get_request_header(&self, req: &HttpRequest, name: &str) -> Option<String> {
        // This would be implemented by the concrete class
        req.headers()
            .get(name)
            .and_then(|v| v.to_str().ok())
            .map(|s| s.to_string())
    }
}

// Helper function to perform a secure comparison of two strings
fn secure_compare(a: &str, b: &str) -> bool {
    if a.len() != b.len() {
        return false;
    }

    let mut result = 0u8;
    for (a_byte, b_byte) in a.bytes().zip(b.bytes()) {
        result |= a_byte ^ b_byte;
    }

    result == 0
}
