use actix_web::{HttpResponse, Responder};
use std::collections::HashMap;
use std::path::Path;

pub trait SendFileUpload {
    fn send_upload(
        &self,
        file_uploader: Box<dyn crate::uploads::Uploader>,
        send_params: HashMap<String, String>,
        redirect_params: HashMap<String, String>,
        attachment: Option<String>,
        proxy: bool,
        disposition: String,
    ) -> impl Responder;

    fn content_type_for(&self, attachment: Option<&str>) -> String;
    fn guess_content_type(&self, filename: &str) -> String;
    fn image_scaling_request(&self, file_uploader: &dyn crate::uploads::Uploader) -> bool;
    fn pwa_icon_safe_for_scaling(&self, file_uploader: &dyn crate::uploads::Uploader) -> bool;
    fn avatar_safe_for_scaling(&self, file_uploader: &dyn crate::uploads::Uploader) -> bool;
    fn mounted_as_avatar(&self, file_uploader: &dyn crate::uploads::Uploader) -> bool;
    fn mounted_as_pwa_icon(&self, file_uploader: &dyn crate::uploads::Uploader) -> bool;
    fn valid_image_scaling_width(&self, allowed_scalar_widths: &[i32]) -> bool;
    fn build_file_url(
        &self,
        file_uploader: &dyn crate::uploads::Uploader,
        redirect_params: HashMap<String, String>,
    ) -> String;
}

impl<T> SendFileUpload for T
where
    T: crate::uploads::UploadsActionsHandler,
{
    fn send_upload(
        &self,
        file_uploader: Box<dyn crate::uploads::Uploader>,
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
            response.append_header((
                "Content-Disposition",
                format!(
                    "{}; filename=\"{}\"",
                    disposition,
                    attachment.unwrap_or_default()
                ),
            ));
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

    fn image_scaling_request(&self, file_uploader: &dyn crate::uploads::Uploader) -> bool {
        self.avatar_safe_for_scaling(file_uploader) || self.pwa_icon_safe_for_scaling(file_uploader)
    }

    fn pwa_icon_safe_for_scaling(&self, file_uploader: &dyn crate::uploads::Uploader) -> bool {
        file_uploader.image_safe_for_scaling()
            && self.mounted_as_pwa_icon(file_uploader)
            && self.valid_image_scaling_width(&[16, 32, 64, 128, 192, 512])
    }

    fn avatar_safe_for_scaling(&self, file_uploader: &dyn crate::uploads::Uploader) -> bool {
        file_uploader.image_safe_for_scaling()
            && self.mounted_as_avatar(file_uploader)
            && self.valid_image_scaling_width(&[32, 64, 128])
    }

    fn mounted_as_avatar(&self, file_uploader: &dyn crate::uploads::Uploader) -> bool {
        file_uploader
            .mounted_as()
            .map_or(false, |mounted_as| mounted_as == "avatar")
    }

    fn mounted_as_pwa_icon(&self, file_uploader: &dyn crate::uploads::Uploader) -> bool {
        file_uploader
            .mounted_as()
            .map_or(false, |mounted_as| mounted_as == "pwa_icon")
    }

    fn valid_image_scaling_width(&self, allowed_scalar_widths: &[i32]) -> bool {
        // In a real implementation, this would check if the width is in the allowed list
        allowed_scalar_widths.contains(&32)
    }

    fn build_file_url(
        &self,
        file_uploader: &dyn crate::uploads::Uploader,
        redirect_params: HashMap<String, String>,
    ) -> String {
        // In a real implementation, this would build a CDN URL
        file_uploader.url(redirect_params)
    }
}
