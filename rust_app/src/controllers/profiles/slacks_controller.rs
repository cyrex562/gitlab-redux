// Ported from: orig_app/app/controllers/profiles/slacks_controller.rb
// Ported on: 2025-05-07
// This file implements the Profiles::SlacksController from the Ruby codebase.

use crate::controllers::concerns::integrations::slack_controller_settings::{
    SlackController, SlackControllerSettings,
};
use crate::controllers::profiles::application_controller::ProfilesApplicationController;
use actix_web::{get, post, web, HttpResponse, Responder};
use serde::Deserialize;

// Stub for current user and project lookup
type UserId = i64;
type ProjectId = i64;

#[derive(Deserialize)]
pub struct SlackLinkParams {
    pub project_id: ProjectId,
}

pub struct SlacksController {
    pub current_user_id: Option<UserId>,
}

impl SlacksController {
    pub fn new(current_user_id: Option<UserId>) -> Self {
        Self { current_user_id }
    }

    /// GET /profile/slack/edit
    #[get("/profile/slack/edit")]
    pub async fn edit(&self) -> impl Responder {
        // TODO: Fetch disabled projects for current user
        if self.current_user_id.is_some() {
            // let projects = ...
            HttpResponse::Ok().body("Edit Slack integration (stub)")
        } else {
            HttpResponse::Unauthorized().finish()
        }
    }

    /// POST /profile/slack/link
    #[post("/profile/slack/link")]
    pub async fn slack_link(&self, params: web::Json<SlackLinkParams>) -> impl Responder {
        // TODO: Find project and generate add_to_slack_link
        let project_id = params.project_id;
        // let link = ...
        let link = format!(
            "https://slack.com/oauth/authorize?project_id={}",
            project_id
        );
        HttpResponse::Ok().json(serde_json::json!({ "add_to_slack_link": link }))
    }
}
