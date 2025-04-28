// Ported from: orig_app/app/controllers/concerns/renders_blob.rb (Ruby)
// This file implements blob rendering logic as described in the Ruby concern.
// Ported on 2025-04-28.

use super::{RenderOptions, Renderable};
use actix_web::{web, HttpResponse};
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct Blob {
    pub path: String,
    pub content: String,
    pub size: usize,
    pub mode: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct BlobViewer {
    pub html: String,
    pub type: String,
}

impl Blob {
    pub fn new(path: String, content: String, size: usize, mode: String) -> Self {
        Self {
            path,
            content,
            size,
            mode,
        }
    }

    pub fn rich_viewer(&self) -> Option<BlobViewer> {
        // TODO: Implement rich viewer
        None
    }

    pub fn auxiliary_viewer(&self) -> Option<BlobViewer> {
        // TODO: Implement auxiliary viewer
        None
    }

    pub fn simple_viewer(&self) -> Option<BlobViewer> {
        // TODO: Implement simple viewer
        None
    }

    pub fn expand(&mut self) {
        // TODO: Implement blob expansion
    }
}

impl Renderable for Blob {
    fn render(&self, options: &RenderOptions) -> HttpResponse {
        let viewer = match options.viewer.as_deref() {
            Some("rich") => self.rich_viewer(),
            Some("auxiliary") => self.auxiliary_viewer(),
            Some("none") => None,
            _ => self.simple_viewer(),
        };

        match viewer {
            Some(viewer) => HttpResponse::Ok().json(viewer),
            None => HttpResponse::NotFound().finish(),
        }
    }

    fn render_json(&self, options: &RenderOptions) -> HttpResponse {
        self.render(options)
    }
}

pub struct BlobRenderer;

impl BlobRenderer {
    pub fn new() -> Self {
        Self
    }

    pub fn render_blob(&self, blob: &Blob, options: &RenderOptions) -> HttpResponse {
        blob.render(options)
    }

    pub fn render_blob_json(&self, blob: &Blob, options: &RenderOptions) -> HttpResponse {
        blob.render_json(options)
    }

    pub fn conditionally_expand_blob(&self, blob: &mut Blob, options: &RenderOptions) {
        if options.expanded {
            blob.expand();
        }
    }

    pub fn conditionally_expand_blobs(&self, blobs: &mut [Blob], options: &RenderOptions) {
        if options.expanded {
            for blob in blobs {
                blob.expand();
            }
        }
    }
}