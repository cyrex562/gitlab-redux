use actix_web::{web, HttpResponse, Responder};
use serde::{Deserialize, Serialize};
use serde_json::json;

use crate::controllers::admin::ApplicationController;

/// Controller for batched jobs
pub struct BatchedJobsController {
    /// The admin application controller
    app_controller: ApplicationController,
}

/// Query parameters for the show action
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ShowQuery {
    /// The database to use
    pub database: Option<String>,
}

impl BatchedJobsController {
    /// Create a new batched jobs controller
    pub fn new() -> Self {
        Self {
            app_controller: ApplicationController::new(),
        }
    }

    /// Handle the show action
    pub async fn show(&self, id: web::Path<i64>, query: web::Query<ShowQuery>) -> impl Responder {
        // TODO: Implement proper database support
        let database = query.database.clone().unwrap_or_else(|| "main".to_string());

        // TODO: Implement proper job fetching
        let job = json!({
            "id": id.into_inner(),
            "status": "failed",
            "error_message": "Something went wrong"
        });

        let transition_logs = vec![];

        let response = json!({
            "job": job,
            "transition_logs": transition_logs
        });

        HttpResponse::Ok().json(response)
    }
}
