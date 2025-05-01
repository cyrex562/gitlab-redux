// Ported from: orig_app/app/controllers/jira_connect/installations_controller.rb
// Ported on: 2025-05-01
// This file implements the JiraConnect::InstallationsController from the Ruby codebase.
//
// See porting_log.txt for details.

use super::application_controller::JiraConnectApplicationController;
use actix_web::{get, put, web, HttpRequest, HttpResponse, Responder};
use serde::{Deserialize, Serialize};
use serde_json::json;

#[derive(Debug, Deserialize)]
pub struct InstallationParams {
    pub instance_url: Option<String>,
}

#[derive(Debug, Serialize)]
pub struct InstallationJson {
    pub gitlab_com: bool,
    pub instance_url: Option<String>,
}

pub struct InstallationsController {
    pub base: JiraConnectApplicationController,
}

impl InstallationsController {
    pub fn new(base: JiraConnectApplicationController) -> Self {
        Self { base }
    }

    /// GET /jira_connect/installations
    #[get("/jira_connect/installations")]
    pub async fn index(req: HttpRequest) -> impl Responder {
        // TODO: Get current_jira_installation from base
        let installation = get_current_jira_installation(&req);
        HttpResponse::Ok().json(installation_json(&installation))
    }

    /// PUT /jira_connect/installations
    #[put("/jira_connect/installations")]
    pub async fn update(req: HttpRequest, params: web::Json<InstallationParams>) -> impl Responder {
        // TODO: Update installation using a service, handle result
        let result = update_installation(&req, &params);
        if result.success {
            let installation = get_current_jira_installation(&req);
            HttpResponse::Ok().json(installation_json(&installation))
        } else {
            HttpResponse::UnprocessableEntity().json(json!({"errors": result.message}))
        }
    }
}

// --- Helper logic and stubs ---

#[derive(Debug)]
pub struct Installation {
    pub instance_url: Option<String>,
}

fn get_current_jira_installation(_req: &HttpRequest) -> Installation {
    // TODO: Real implementation should extract from request/session
    Installation {
        instance_url: Some("https://example.gitlab.com".to_string()),
    }
}

fn installation_json(installation: &Installation) -> InstallationJson {
    InstallationJson {
        gitlab_com: installation.instance_url.is_none()
            || installation.instance_url.as_ref().unwrap().is_empty(),
        instance_url: installation.instance_url.clone(),
    }
}

struct UpdateResult {
    success: bool,
    message: String,
}

fn update_installation(_req: &HttpRequest, _params: &InstallationParams) -> UpdateResult {
    // TODO: Call JiraConnectInstallations::UpdateService equivalent
    UpdateResult {
        success: true,
        message: String::new(),
    }
}
