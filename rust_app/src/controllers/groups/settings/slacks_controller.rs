// Ported from: orig_app/app/controllers/groups/settings/slacks_controller.rb
// Date ported: 2025-05-01

use crate::controllers::concerns::integrations::slack_controller_settings::SlackController;
use crate::controllers::concerns::integrations::SlackControllerSettings;
use actix_web::{web, HttpResponse, Responder};

// Controller for group-level Slack integration settings
pub struct SlacksController {
    group_id: i64,
    current_user_id: i64,
    oauth_state: String,
}

impl SlacksController {
    pub fn new(group_id: i64, current_user_id: i64, oauth_state: String) -> Self {
        Self {
            group_id,
            current_user_id,
            oauth_state,
        }
    }

    // GET /groups/{group_id}/settings/slack
    pub async fn show(&self) -> impl Responder {
        // TODO: Authorize admin group!
        // TODO: Render group_settings layout
        HttpResponse::Ok().body("Slack integration settings page (group)")
    }

    // POST /groups/{group_id}/settings/slack/auth
    pub async fn slack_auth(&self) -> impl Responder {
        let slack_controller = SlackController::new(self.group_id, self.oauth_state.clone());
        match slack_controller.slack_auth().await {
            Ok(result) => {
                if !result.success {
                    // TODO: Set flash alert with result.message
                }
                // TODO: Set session[:slack_install_success]
                self.redirect_to_integration_page().await
            }
            Err(e) => HttpResponse::InternalServerError().body(format!("OAuth error: {}", e)),
        }
    }

    // Helper: redirect to integration page
    pub async fn redirect_to_integration_page(&self) -> HttpResponse {
        // TODO: Find or create integration, then redirect
        HttpResponse::Found()
            .header(
                "Location",
                format!("/groups/{}/settings/integrations/slack", self.group_id),
            )
            .finish()
    }

    // Helper: installation service (stub)
    pub fn installation_service(&self) {
        // TODO: Implement SlackInstallation::GroupService equivalent
    }
}
