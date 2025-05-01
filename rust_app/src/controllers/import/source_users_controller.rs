// Ported from: orig_app/app/controllers/import/source_users_controller.rb
// Ported on: 2025-05-01
// This file implements the Import::SourceUsersController from the Ruby codebase.
// See porting_log.txt for details.

use actix_web::{get, post, web, HttpResponse, Responder};
use serde::{Deserialize, Serialize};

#[derive(Debug, Deserialize)]
pub struct ReassignmentParams {
    pub reassignment_token: String,
}

pub struct SourceUsersController;

impl SourceUsersController {
    /// POST /import/source_users/accept
    #[post("/import/source_users/accept")]
    pub async fn accept(params: web::Json<ReassignmentParams>) -> impl Responder {
        // TODO: Implement AcceptReassignmentService logic
        // Placeholder: always success
        let success = true;
        if success {
            // TODO: flash banner('accept_invite')
            HttpResponse::Found().header("Location", "/").finish()
        } else {
            HttpResponse::Found()
                .header("Location", "/")
                .body("The invitation could not be accepted.")
        }
    }

    /// POST /import/source_users/decline
    #[post("/import/source_users/decline")]
    pub async fn decline(params: web::Json<ReassignmentParams>) -> impl Responder {
        // TODO: Implement RejectReassignmentService logic
        // Placeholder: always success
        let success = true;
        if success {
            // TODO: flash banner('reject_invite')
            HttpResponse::Found().header("Location", "/").finish()
        } else {
            HttpResponse::Found()
                .header("Location", "/")
                .body("The invitation could not be declined.")
        }
    }

    /// GET /import/source_users/show
    #[get("/import/source_users/show")]
    pub async fn show() -> impl Responder {
        // TODO: Implement show logic
        HttpResponse::Ok().body("Show source user invite page")
    }
}
