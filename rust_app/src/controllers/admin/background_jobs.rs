use actix_web::{web, HttpResponse, Responder};
use serde::{Deserialize, Serialize};

use crate::controllers::admin::ApplicationController;

/// Controller for background jobs
pub struct BackgroundJobsController {
    /// The admin application controller
    app_controller: ApplicationController,
}

impl BackgroundJobsController {
    /// Create a new background jobs controller
    pub fn new() -> Self {
        Self {
            app_controller: ApplicationController::new(),
        }
    }
}
