use actix_web::{HttpRequest, HttpResponse};
use serde::Serialize;

pub struct Blob {
    pub rich_viewer: Option<Viewer>,
    pub auxiliary_viewer: Option<Viewer>,
    pub simple_viewer: Option<Viewer>,
}

pub struct Viewer {
    // Add viewer fields as needed
}

impl Blob {
    pub fn expand(&mut self) {
        // Implementation for expanding blob content
    }
}

#[derive(Serialize)]
pub struct BlobViewerResponse {
    html: String,
}

pub trait RendersBlob {
    fn blob_viewer_json(&self, req: &HttpRequest, blob: &Blob) -> Option<BlobViewerResponse> {
        let viewer = match req
            .query_string()
            .split('&')
            .find(|param| param.starts_with("viewer="))
            .and_then(|param| param.split('=').nth(1))
        {
            Some("rich") => blob.rich_viewer.as_ref(),
            Some("auxiliary") => blob.auxiliary_viewer.as_ref(),
            Some("none") => None,
            _ => blob.simple_viewer.as_ref(),
        };

        viewer.map(|v| BlobViewerResponse {
            html: self.view_to_html_string("projects/blob/_viewer", v, false),
        })
    }

    fn render_blob_json(&self, req: &HttpRequest, blob: &Blob) -> HttpResponse {
        match self.blob_viewer_json(req, blob) {
            Some(json) => HttpResponse::Ok().json(json),
            None => HttpResponse::NotFound().finish(),
        }
    }

    fn conditionally_expand_blob(&self, req: &HttpRequest, blob: &mut Blob) {
        self.conditionally_expand_blobs(req, &mut [blob]);
    }

    fn conditionally_expand_blobs(&self, req: &HttpRequest, blobs: &mut [&mut Blob]) {
        if req.query_string().contains("expanded=true") {
            for blob in blobs {
                blob.expand();
            }
        }
    }

    // Required method to be implemented by concrete types
    fn view_to_html_string(&self, template: &str, viewer: &Viewer, load_async: bool) -> String;
}
