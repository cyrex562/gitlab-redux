use actix_web::{web, HttpResponse, Responder};
use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::path::Path;
use std::sync::Arc;
use uuid::Uuid;

mod send_file_upload;
mod sends_blob;

pub use send_file_upload::SendFileUpload;
pub use sends_blob::{SendsBlob, Blob, LfsObject, Repository, Project, FileUploader};

// Starting with version 2, Markdown upload URLs use project / group IDs instead of paths
const ID_BASED_UPLOAD_PATH_VERSION: i32 = 2;

const UPLOAD_MOUNTS: [&str; 7] = ["avatar", "attachment", "file", "logo", "pwa_icon", "header_logo", "favicon", "screenshot"];

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Upload {
    pub id: i64,
    pub model_id: i64,
    pub model_type: String,
    pub uploader: String,
    pub path: String,
    pub size: i64,
    pub version: i32,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct UploadResult {
    pub link: HashMap<String, String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct UploadParams {
    pub file: web::Multipart,
    pub mounted_as: Option<String>,
    pub secret: Option<String>,
    pub filename: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct UploadService {
    pub model: Arc<dyn Uploadable>,
    pub file: web::Multipart,
    pub uploader_class: Box<dyn Uploader>,
    pub uploaded_by_user_id: Option<i64>,
}

pub trait Uploadable {
    fn id(&self) -> i64;
    fn type_name(&self) -> String;
}

pub trait Uploader {
    fn new(model: Arc<dyn Uploadable>, secret: Option<String>) -> Self where Self: Sized;
    fn model_valid(&self) -> bool;
    fn exists(&self) -> bool;
    fn path(&self) -> String;
    fn url(&self, params: HashMap<String, String>) -> String;
    fn file_storage(&self) -> bool;
    fn proxy_download_enabled(&self) -> bool;
    fn embeddable(&self) -> bool;
    fn pdf(&self) -> bool;
    fn image_safe_for_scaling(&self) -> bool;
    fn mounted_as(&self) -> Option<String>;
    fn versions(&self) -> HashMap<String, Box<dyn Uploader>>;
    fn upload_paths(&self, filename: &str) -> Vec<String>;
    fn workhorse_authorize(&self, has_length: bool, maximum_size: i64) -> Result<HashMap<String, String>, Box<dyn std::error::Error>>;
}

pub struct UploadsActionsHandler {
    db: Arc<dyn Database>,
    current_user: Option<User>,
    model: Arc<dyn Uploadable>,
    uploader_class: Box<dyn Uploader>,
}

impl UploadsActionsHandler {
    pub fn new(
        db: Arc<dyn Database>,
        current_user: Option<User>,
        model: Arc<dyn Uploadable>,
        uploader_class: Box<dyn Uploader>,
    ) -> Self {
        Self {
            db,
            current_user,
            model,
            uploader_class,
        }
    }

    pub async fn create(&self, params: web::Json<UploadParams>) -> impl Responder {
        let uploader = UploadService {
            model: self.model.clone(),
            file: params.file.clone(),
            uploader_class: self.uploader_class.clone(),
            uploaded_by_user_id: self.current_user.as_ref().map(|user| user.id),
        };

        match self.execute_upload(uploader).await {
            Ok(Some(result)) => {
                HttpResponse::Ok().json(result)
            }
            Ok(None) => {
                HttpResponse::UnprocessableEntity().json("Invalid file.")
            }
            Err(e) => {
                HttpResponse::InternalServerError().json(format!("Error: {}", e))
            }
        }
    }

    pub async fn show(&self, filename: &str) -> impl Responder {
        // Check for path traversal
        if !self.is_safe_path(filename) {
            return HttpResponse::BadRequest().finish();
        }

        let uploader = self.get_uploader().await;
        
        if !uploader.exists() {
            return HttpResponse::NotFound().finish();
        }

        let (ttl, directives) = self.cache_settings();
        
        // Set cache headers
        let mut response = HttpResponse::Ok();
        
        if let Some(ttl) = ttl {
            response.append_header(("Cache-Control", format!("max-age={}", ttl)));
        }
        
        if let Some(directives) = directives {
            if directives.private {
                response.append_header(("Cache-Control", "private"));
            }
            if directives.must_revalidate {
                response.append_header(("Cache-Control", "must-revalidate"));
            }
        }

        let file_uploader = self.find_file_uploader(&uploader, filename).await;
        
        if file_uploader.is_none() {
            return HttpResponse::NotFound().finish();
        }

        let file_uploader = file_uploader.unwrap();
        let content_disposition = self.content_disposition(&file_uploader);
        
        self.send_upload(
            file_uploader,
            HashMap::new(),
            HashMap::new(),
            Some(filename.to_string()),
            false,
            content_disposition,
        ).await
    }

    pub async fn authorize(&self) -> impl Responder {
        // Set workhorse internal API content type
        let mut response = HttpResponse::Ok();
        response.append_header(("Content-Type", "application/json"));
        
        let max_size = 10 * 1024 * 1024; // 10MB default
        
        match self.uploader_class.workhorse_authorize(false, max_size) {
            Ok(result) => response.json(result),
            Err(_) => HttpResponse::InternalServerError().json("Error uploading file"),
        }
    }

    async fn execute_upload(&self, service: UploadService) -> Result<Option<UploadResult>, Box<dyn std::error::Error>> {
        // TODO: Implement actual upload service
        Ok(Some(UploadResult {
            link: HashMap::from([
                ("url".to_string(), "https://example.com/uploads/file.jpg".to_string()),
                ("alt".to_string(), "file.jpg".to_string()),
                ("title".to_string(), "file.jpg".to_string()),
            ]),
        }))
    }

    async fn get_uploader(&self) -> Box<dyn Uploader> {
        if self.uploader_mounted() {
            // In a real implementation, this would get the uploader from the model
            self.uploader_class.clone()
        } else {
            self.build_uploader_from_upload().await.unwrap_or_else(|| self.uploader_class.clone())
        }
    }

    fn uploader_mounted(&self) -> bool {
        // In a real implementation, this would check if the uploader is mounted
        false
    }

    async fn build_uploader_from_upload(&self) -> Option<Box<dyn Uploader>> {
        let uploader = self.build_uploader()?;
        
        if !uploader.model_valid() {
            return None;
        }
        
        // In a real implementation, this would find the upload in the database
        Some(uploader)
    }

    fn build_uploader(&self) -> Option<Box<dyn Uploader>> {
        // In a real implementation, this would build an uploader from params
        Some(self.uploader_class.clone())
    }

    fn content_disposition(&self, uploader: &dyn Uploader) -> String {
        if uploader.embeddable() || uploader.pdf() {
            "inline".to_string()
        } else {
            "attachment".to_string()
        }
    }

    async fn find_file_uploader(&self, uploader: &dyn Uploader, filename: &str) -> Option<Box<dyn Uploader>> {
        let mut versions = uploader.versions();
        versions.insert("original".to_string(), uploader.clone());
        
        versions.into_values().find(|version| {
            // In a real implementation, this would check if the version's filename matches
            version.path().ends_with(filename)
        })
    }

    fn is_safe_path(&self, path: &str) -> bool {
        // TODO: Implement path traversal check
        !path.contains("..")
    }

    fn cache_settings(&self) -> (Option<i32>, Option<CacheDirectives>) {
        // In a real implementation, this would return cache settings
        (Some(3600), Some(CacheDirectives {
            private: true,
            must_revalidate: true,
        }))
    }

    async fn send_upload(
        &self,
        file_uploader: Box<dyn Uploader>,
        send_params: HashMap<String, String>,
        redirect_params: HashMap<String, String>,
        attachment: Option<String>,
        proxy: bool,
        disposition: String,
    ) -> impl Responder {
        let content_type = self.content_type_for(attachment.as_deref());
        
        let mut response = HttpResponse::Ok();
        
        if let Some(attachment) = attachment {
            let content_disposition = format!("{}; filename=\"{}\"", disposition, attachment);
            response.append_header(("Content-Disposition", content_disposition));
            
            // Handle JS files specially
            if attachment.ends_with(".js") {
                response.append_header(("Content-Type", "text/plain"));
            } else {
                response.append_header(("Content-Type", content_type));
            }
        }
        
        if self.image_scaling_request(&file_uploader) {
            // Handle image scaling
            let location = if file_uploader.file_storage() {
                file_uploader.path()
            } else {
                file_uploader.url(HashMap::new())
            };
            
            // In a real implementation, this would set headers for scaled image
            response.append_header(("X-Scaled-Image", location));
            response
        } else if file_uploader.file_storage() {
            // Send file directly
            response.append_header(("Content-Type", content_type));
            response.append_header(("Content-Disposition", format!("{}; filename=\"{}\"", disposition, attachment.unwrap_or_default())));
            response.body("File content would be here")
        } else if file_uploader.proxy_download_enabled() || proxy {
            // Proxy download
            let url = file_uploader.url(redirect_params);
            response.append_header(("X-Send-File", url));
            response
        } else {
            // Redirect to file URL
            let file_url = self.build_file_url(&file_uploader, redirect_params);
            HttpResponse::Found()
                .append_header(("Location", file_url))
                .finish()
        }
    }

    fn content_type_for(&self, attachment: Option<&str>) -> String {
        match attachment {
            Some(filename) => self.guess_content_type(filename),
            None => "application/octet-stream".to_string(),
        }
    }

    fn guess_content_type(&self, filename: &str) -> String {
        // In a real implementation, this would use a MIME type library
        match Path::new(filename).extension().and_then(|ext| ext.to_str()) {
            Some("jpg") | Some("jpeg") => "image/jpeg".to_string(),
            Some("png") => "image/png".to_string(),
            Some("gif") => "image/gif".to_string(),
            Some("pdf") => "application/pdf".to_string(),
            Some("js") => "text/javascript".to_string(),
            Some("css") => "text/css".to_string(),
            Some("html") => "text/html".to_string(),
            Some("txt") => "text/plain".to_string(),
            _ => "application/octet-stream".to_string(),
        }
    }

    fn image_scaling_request(&self, file_uploader: &dyn Uploader) -> bool {
        self.avatar_safe_for_scaling(file_uploader) || self.pwa_icon_safe_for_scaling(file_uploader)
    }

    fn pwa_icon_safe_for_scaling(&self, file_uploader: &dyn Uploader) -> bool {
        file_uploader.image_safe_for_scaling() &&
            self.mounted_as_pwa_icon(file_uploader) &&
            self.valid_image_scaling_width(&[16, 32, 64, 128, 192, 512])
    }

    fn avatar_safe_for_scaling(&self, file_uploader: &dyn Uploader) -> bool {
        file_uploader.image_safe_for_scaling() &&
            self.mounted_as_avatar(file_uploader) &&
            self.valid_image_scaling_width(&[32, 64, 128])
    }

    fn mounted_as_avatar(&self, file_uploader: &dyn Uploader) -> bool {
        file_uploader.mounted_as().map_or(false, |mounted_as| mounted_as == "avatar")
    }

    fn mounted_as_pwa_icon(&self, file_uploader: &dyn Uploader) -> bool {
        file_uploader.mounted_as().map_or(false, |mounted_as| mounted_as == "pwa_icon")
    }

    fn valid_image_scaling_width(&self, allowed_scalar_widths: &[i32]) -> bool {
        // In a real implementation, this would check if the width is in the allowed list
        allowed_scalar_widths.contains(&32)
    }

    fn build_file_url(&self, file_uploader: &dyn Uploader, redirect_params: HashMap<String, String>) -> String {
        // In a real implementation, this would build a CDN URL
        file_uploader.url(redirect_params)
    }
}

#[derive(Debug, Clone)]
pub struct CacheDirectives {
    pub private: bool,
    pub must_revalidate: bool,
}

// These would be defined elsewhere in the codebase
pub trait Database {}
pub struct User {
    pub id: i64,
    // Add other user fields as needed
} 