// Ported from: orig_app/app/controllers/jira_connect/events_controller.rb
// Ported on: 2025-05-01
// This file implements the JiraConnect::EventsController from the Ruby codebase.
//
// See porting_log.txt for details.

use super::application_controller::JiraConnectApplicationController;
use actix_web::{post, web, HttpRequest, HttpResponse, Responder};
use serde_json::json;

pub struct EventsController {
    pub base: JiraConnectApplicationController,
}

impl EventsController {
    pub fn new(base: JiraConnectApplicationController) -> Self {
        Self { base }
    }

    /// POST /jira_connect/events/installed
    #[post("/jira_connect/events/installed")]
    pub async fn installed(req: HttpRequest) -> impl Responder {
        // TODO: Implement JWT verification and installation logic
        // Placeholder: always return 200 OK
        HttpResponse::Ok().finish()
    }

    /// POST /jira_connect/events/uninstalled
    #[post("/jira_connect/events/uninstalled")]
    pub async fn uninstalled(req: HttpRequest) -> impl Responder {
        // TODO: Implement uninstall logic
        // Placeholder: always return 200 OK
        HttpResponse::Ok().finish()
    }
}
