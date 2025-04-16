use actix_web::{web, HttpResponse, Responder, get};
use std::sync::Arc;
use crate::controllers::ApplicationController;

/// This controller implements /.well-known paths that have no better home.
///
/// Other controllers also implement /.well-known/* paths.
pub struct WellKnownController {
    app_controller: Arc<ApplicationController>,
}

impl WellKnownController {
    pub fn new(app_controller: Arc<ApplicationController>) -> Self {
        Self { app_controller }
    }

    /// Returns the security.txt content if configured
    #[get("/.well-known/security.txt")]
    pub async fn security_txt(&self) -> impl Responder {
        // TODO: Implement getting security.txt content from settings
        // For now, we'll just return a 404
        self.app_controller.render_404().await
    }
} 