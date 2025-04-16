use actix_web::{web, HttpResponse, Responder};
use serde::{Deserialize, Serialize};
use serde_json::json;

use crate::controllers::admin::ApplicationController;

/// Controller for managing cohorts in the admin area
pub struct CohortsController {
    /// The admin application controller
    app_controller: ApplicationController,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct Cohort {
    id: i32,
    name: String,
    description: Option<String>,
    created_at: String,
    updated_at: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct CohortData {
    cohorts: Vec<Cohort>,
    total_users: i32,
    active_users: i32,
    inactive_users: i32,
}

impl CohortsController {
    /// Create a new cohorts controller
    pub fn new() -> Self {
        Self {
            app_controller: ApplicationController::new(),
        }
    }

    /// Configure the routes for this controller
    pub fn configure(cfg: &mut web::ServiceConfig) {
        cfg.service(
            web::scope("/admin/cohorts")
                .route("", web::get().to(Self::index)),
        );
    }

    /// List all cohorts
    pub async fn index() -> impl Responder {
        // TODO: Implement proper cohorts fetching
        // This is a placeholder implementation
        let cohorts = Self::load_cohorts().await;

        HttpResponse::Ok().json(cohorts)
    }

    /// Load cohorts data
    async fn load_cohorts() -> CohortData {
        // TODO: Implement proper cohorts loading with caching
        // This is a placeholder implementation
        let cohorts = vec![
            Cohort {
                id: 1,
                name: "January 2023".to_string(),
                description: Some("Users who signed up in January 2023".to_string()),
                created_at: "2023-01-01T00:00:00Z".to_string(),
                updated_at: "2023-01-01T00:00:00Z".to_string(),
            },
            Cohort {
                id: 2,
                name: "February 2023".to_string(),
                description: Some("Users who signed up in February 2023".to_string()),
                created_at: "2023-02-01T00:00:00Z".to_string(),
                updated_at: "2023-02-01T00:00:00Z".to_string(),
            },
        ];

        CohortData {
            cohorts,
            total_users: 100,
            active_users: 80,
            inactive_users: 20,
        }
    }
} 