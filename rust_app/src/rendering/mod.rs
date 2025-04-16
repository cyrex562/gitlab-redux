pub mod blob;
pub mod commits;
pub mod notes;
pub mod projects;

use actix_web::{web, HttpResponse, Responder};
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct RenderOptions {
    pub viewer: Option<String>,
    pub expanded: bool,
    pub format: Option<String>,
}

impl Default for RenderOptions {
    fn default() -> Self {
        Self {
            viewer: None,
            expanded: false,
            format: None,
        }
    }
}

pub trait Renderable {
    fn render(&self, options: &RenderOptions) -> HttpResponse;
    fn render_json(&self, options: &RenderOptions) -> HttpResponse;
}

pub struct Renderer;

impl Renderer {
    pub fn new() -> Self {
        Self
    }

    pub fn render<T: Renderable>(&self, item: &T, options: &RenderOptions) -> HttpResponse {
        item.render(options)
    }

    pub fn render_json<T: Renderable>(&self, item: &T, options: &RenderOptions) -> HttpResponse {
        item.render_json(options)
    }

    pub fn render_collection<T: Renderable>(
        &self,
        items: &[T],
        options: &RenderOptions,
    ) -> HttpResponse {
        // TODO: Implement collection rendering
        HttpResponse::Ok().finish()
    }

    pub fn render_error(&self, error: &str) -> HttpResponse {
        HttpResponse::InternalServerError().json(serde_json::json!({
            "error": error
        }))
    }
}
