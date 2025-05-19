// Ported from: orig_app/app/controllers/projects/google_cloud/databases_controller.rb
// Ported on: 2025-05-19
// This file implements the Projects::GoogleCloud::DatabasesController from the Ruby codebase.

use actix_web::{web, HttpResponse, Responder};
use serde_json::json;

use crate::controllers::projects::google_cloud::base_controller::ProjectsGoogleCloudBaseController;

pub struct ProjectsGoogleCloudDatabasesController {
    base: ProjectsGoogleCloudBaseController,
}

impl ProjectsGoogleCloudDatabasesController {
    pub fn new(base: ProjectsGoogleCloudBaseController) -> Self {
        Self { base }
    }

    /// GET /projects/:project_id/google_cloud/databases
    pub async fn index(&self) -> impl Responder {
        // The following is a simplified port. Actual URL helpers and service calls would be implemented as needed.
        let js_data = json!({
            "configurationUrl": "project_google_cloud_configuration_path", // placeholder
            "deploymentsUrl": "project_google_cloud_deployments_path", // placeholder
            "databasesUrl": "project_google_cloud_databases_path", // placeholder
            "cloudsqlPostgresUrl": "new_project_google_cloud_database_path_postgres", // placeholder
            "cloudsqlMysqlUrl": "new_project_google_cloud_database_path_mysql", // placeholder
            "cloudsqlSqlserverUrl": "new_project_google_cloud_database_path_sqlserver", // placeholder
            "cloudsqlInstances": [], // GetCloudsqlInstancesService logic omitted
            "emptyIllustrationUrl": "/assets/illustrations/empty-state/empty-pipeline-md.svg" // placeholder
        });
        self.base.track_event("render_page", None);
        HttpResponse::Ok().json(js_data)
    }

    /// GET /projects/:project_id/google_cloud/databases/new
    pub async fn new(&self, product: String) -> impl Responder {
        // The following is a simplified port. Actual logic for gcp_projects, refs, etc. would be implemented as needed.
        let js_data = json!({
            "gcpProjects": [], // gcp_projects logic omitted
            "refs": [], // refs logic omitted
            "cancelPath": "project_google_cloud_databases_path", // placeholder
            "formTitle": self.form_title(&product),
            "formDescription": self.description(&product),
            "databaseVersions": [], // CloudsqlHelper::VERSIONS logic omitted
            "tiers": [] // CloudsqlHelper::TIERS logic omitted
        });
        self.base.track_event("render_form", None);
        HttpResponse::Ok().json(js_data)
    }

    // POST /projects/:project_id/google_cloud/databases
    pub async fn create(&self) -> impl Responder {
        // Service call and flash/redirect logic omitted for brevity
        // Track event examples:
        self.base.track_event("create_cloudsql_instance", None);
        HttpResponse::Ok().json(json!({"message": "Cloud SQL instance creation request successful. Expected resolution time is ~5 minutes."}))
    }

    fn form_title(&self, product: &str) -> &str {
        match product {
            "postgres" => "Cloud SQL for Postgres",
            "mysql" => "Cloud SQL for MySQL",
            _ => "Cloud SQL for SQL Server",
        }
    }

    fn description(&self, product: &str) -> &str {
        match product {
            "postgres" => "Cloud SQL instances are fully managed, relational PostgreSQL databases. Google handles replication, patch management, and database management to ensure availability and performance.",
            "mysql" => "Cloud SQL instances are fully managed, relational MySQL databases. Google handles replication, patch management, and database management to ensure availability and performance.",
            _ => "Cloud SQL instances are fully managed, relational SQL Server databases. Google handles replication, patch management, and database management to ensure availability and performance.",
        }
    }
}
