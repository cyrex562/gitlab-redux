// Ported from: orig_app/app/controllers/concerns/send_file_upload.rb
// Date ported: 2025-04-29
// This file implements the SendFileUpload concern from Ruby in Rust.
// See porting log for details.

use actix_web::{web, HttpRequest, HttpResponse, Responder};
use std::path::Path;
use std::collections::HashMap;
use mime_guess::MimeGuess;
use crate::workhorse::Workhorse;
use crate::storage::object_storage::cdn::FileUrl;

pub trait SendFileUpload {
    fn send_upload(
        &self,
        req: &HttpRequest,
        file_upload: &FileUpload,
        send_params: Option<HashMap<String, String>>,
        redirect_params: Option<HashMap<String, String>>,
        attachment: Option<String>,
        proxy: bool,
        disposition: &str,
    ) -> impl Responder;
    
    fn content_type_for(&self, attachment: Option<&String>) -> String;
    fn guess_content_type(&self, filename: &str) -> String;
}

pub struct SendFileUploadHandler;

impl SendFileUploadHandler {
    pub fn new() -> Self {
        SendFileUploadHandler
    }
    
    fn image_scaling_request(&self, file_upload: &FileUpload, req: &HttpRequest) -> bool {
        self.avatar_safe_for_scaling(file_upload, req) || 
        self.pwa_icon_safe_for_scaling(file_upload, req)
    }
    
    fn pwa_icon_safe_for_scaling(&self, file_upload: &FileUpload, req: &HttpRequest) -> bool {
        file_upload.image_safe_for_scaling() &&
        self.mounted_as_pwa_icon(file_upload) &&
        self.valid_image_scaling_width(req, &Appearance::ALLOWED_PWA_ICON_SCALER_WIDTHS)
    }
    
    fn avatar_safe_for_scaling(&self, file_upload: &FileUpload, req: &HttpRequest) -> bool {
        file_upload.image_safe_for_scaling() &&
        self.mounted_as_avatar(file_upload) &&
        self.valid_image_scaling_width(req, &Avatarable::ALLOWED_IMAGE_SCALER_WIDTHS)
    }
    
    fn mounted_as_avatar(&self, file_upload: &FileUpload) -> bool {
        file_upload.mounted_as().map_or(false, |m| m == "avatar")
    }
    
    fn mounted_as_pwa_icon(&self, file_upload: &FileUpload) -> bool {
        file_upload.mounted_as().map_or(false, |m| m == "pwa_icon")
    }
    
    fn valid_image_scaling_width(&self, req: &HttpRequest, allowed_scalar_widths: &[i32]) -> bool {
        if let Some(width_str) = req.query_string().get("width") {
            if let Ok(width) = width_str.parse::<i32>() {
                return allowed_scalar_widths.contains(&width);
            }
        }
        false
    }
}

impl SendFileUpload for SendFileUploadHandler {
    fn send_upload(
        &self,
        req: &HttpRequest,
        file_upload: &FileUpload,
        send_params: Option<HashMap<String, String>>,
        redirect_params: Option<HashMap<String, String>>,
        attachment: Option<String>,
        proxy: bool,
        disposition: &str,
    ) -> impl Responder {
        let mut send_params = send_params.unwrap_or_default();
        let mut redirect_params = redirect_params.unwrap_or_default();
        
        let content_type = self.content_type_for(attachment.as_ref());
        
        if let Some(attachment_name) = &attachment {
            let response_disposition = format!("{}; filename=\"{}\"", disposition, attachment_name);
            
            // Add response headers for cloud storage
            redirect_params.insert(
                "response-content-disposition".to_string(),
                response_disposition,
            );
            redirect_params.insert("response-content-type".to_string(), content_type.clone());
            
            // Handle .js files specially
            if Path::new(attachment_name).extension().map_or(false, |ext| ext == "js") {
                send_params.insert("content_type".to_string(), "text/plain".to_string());
            }
            
            send_params.insert("filename".to_string(), attachment_name.clone());
            send_params.insert("disposition".to_string(), disposition.to_string());
        }
        
        if self.image_scaling_request(file_upload, req) {
            let location = if file_upload.file_storage() {
                file_upload.path()
            } else {
                file_upload.url()
            };
            
            let width = req.query_string()
                .get("width")
                .and_then(|w| w.parse::<i32>().ok())
                .unwrap_or(0);
                
            let mut response = HttpResponse::Ok();
            let (header_name, header_value) = Workhorse::send_scaled_image(location, width, &content_type);
            response.append_header((header_name, header_value));
            return response.finish();
        } else if file_upload.file_storage() {
            // Send file directly
            let mut response = HttpResponse::Ok();
            for (key, value) in send_params {
                response.append_header((key, value));
            }
            response.append_header(("Content-Type", content_type));
            response.append_header(("Content-Disposition", format!("{}; filename=\"{}\"", 
                disposition, 
                attachment.unwrap_or_else(|| "file".to_string())
            )));
            return response.streaming(file_upload.stream());
        } else if file_upload.proxy_download_enabled() || proxy {
            // Use workhorse to proxy the download
            let mut response = HttpResponse::Ok();
            let (header_name, header_value) = Workhorse::send_url(file_upload.url_with_params(&redirect_params));
            response.append_header((header_name, header_value));
            return response.finish();
        } else {
            // Redirect to the file URL
            let file_url = FileUrl::new(
                file: file_upload.clone(),
                ip_address: req.connection_info().peer_addr().unwrap_or("unknown").to_string(),
                redirect_params: redirect_params,
            );
            return HttpResponse::Found()
                .append_header(("Location", file_url.url()))
                .finish();
        }
    }
    
    fn content_type_for(&self, attachment: Option<&String>) -> String {
        match attachment {
            Some(filename) => self.guess_content_type(filename),
            None => String::new(),
        }
    }
    
    fn guess_content_type(&self, filename: &str) -> String {
        let mime = MimeGuess::from_path(filename).first();
        mime.map_or_else(
            || "application/octet-stream".to_string(),
            |m| m.to_string(),
        )
    }
}

// These would be implemented in separate modules
pub mod workhorse {
    pub struct Workhorse;
    
    impl Workhorse {
        pub fn send_scaled_image(location: &str, width: i32, content_type: &str) -> (String, String) {
            // In a real implementation, this would generate the appropriate headers
            // for workhorse to handle image scaling
            ("X-Send-Scaled-Image".to_string(), format!("{}:{}:{}", location, width, content_type))
        }
        
        pub fn send_url(url: String) -> (String, String) {
            // In a real implementation, this would generate the appropriate headers
            // for workhorse to handle URL sending
            ("X-Send-URL".to_string(), url)
        }
    }
}

pub mod storage {
    pub mod object_storage {
        pub mod cdn {
            use std::collections::HashMap;
            
            pub struct FileUrl {
                file: FileUpload,
                ip_address: String,
                redirect_params: HashMap<String, String>,
            }
            
            impl FileUrl {
                pub fn new(
                    file: FileUpload,
                    ip_address: String,
                    redirect_params: HashMap<String, String>,
                ) -> Self {
                    FileUrl {
                        file,
                        ip_address,
                        redirect_params,
                    }
                }
                
                pub fn url(&self) -> String {
                    // In a real implementation, this would generate a signed URL
                    // for the file with the appropriate parameters
                    format!("https://cdn.example.com/{}", self.file.path())
                }
            }
        }
    }
}

// These would be implemented in separate modules
pub struct FileUpload {
    path: String,
    url: String,
    file_storage: bool,
    proxy_download_enabled: bool,
    mounted_as: Option<String>,
    image_safe_for_scaling: bool,
}

impl FileUpload {
    pub fn path(&self) -> &str {
        &self.path
    }
    
    pub fn url(&self) -> &str {
        &self.url
    }
    
    pub fn url_with_params(&self, params: &HashMap<String, String>) -> String {
        // In a real implementation, this would append the parameters to the URL
        let mut url = self.url.clone();
        if !params.is_empty() {
            url.push('?');
            let param_strings: Vec<String> = params
                .iter()
                .map(|(k, v)| format!("{}={}", k, v))
                .collect();
            url.push_str(&param_strings.join("&"));
        }
        url
    }
    
    pub fn file_storage(&self) -> bool {
        self.file_storage
    }
    
    pub fn proxy_download_enabled(&self) -> bool {
        self.proxy_download_enabled
    }
    
    pub fn mounted_as(&self) -> Option<&str> {
        self.mounted_as.as_deref()
    }
    
    pub fn image_safe_for_scaling(&self) -> bool {
        self.image_safe_for_scaling
    }
    
    pub fn stream(&self) -> impl std::io::Read {
        // In a real implementation, this would return a stream to the file
        // For now, we'll just return an empty reader
        std::io::empty()
    }
}

pub mod appearance {
    pub struct Appearance;
    
    impl Appearance {
        pub const ALLOWED_PWA_ICON_SCALER_WIDTHS: [i32; 3] = [192, 512, 1024];
    }
}

pub mod avatarable {
    pub struct Avatarable;
    
    impl Avatarable {
        pub const ALLOWED_IMAGE_SCALER_WIDTHS: [i32; 5] = [32, 64, 128, 256, 512];
    }
}