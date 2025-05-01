// Ported from: orig_app/app/controllers/import/gitlab_groups_controller.rb
// Ported on: 2025-05-01
// This file implements the Import::GitlabGroupsController from the Ruby codebase.
//
// See porting_log.txt for details.

use crate::controllers::concerns::check_rate_limit::CheckRateLimit;
use crate::controllers::concerns::{
    UploadedFile, WorkhorseAuthorization, WorkhorseAuthorizationHandler,
};
use actix_web::{post, web, HttpRequest, HttpResponse, Responder};
use serde::{Deserialize, Serialize};

#[derive(Debug, Deserialize)]
pub struct ImportGroupParams {
    pub path: String,
    pub name: String,
    pub parent_id: Option<i64>,
    pub file: Option<UploadedFile>,
}

pub struct GitlabGroupsController;

impl GitlabGroupsController {
    /// POST /import/gitlab/groups
    #[post("/import/gitlab/groups")]
    pub async fn create(req: HttpRequest, params: web::Json<ImportGroupParams>) -> impl Responder {
        // Rate limit check (placeholder)
        // In real code, use CheckRateLimit trait
        // if rate_limited { return error response }

        // File validation
        let file_valid = match &params.file {
            Some(file) => WorkhorseAuthorizationHandler::default().file_is_valid(file),
            None => false,
        };
        if !file_valid {
            return HttpResponse::Found()
                .header("Location", "/groups/new#import-group-pane")
                .body("Unable to process group import file");
        }

        // Group creation logic (placeholder)
        // In real code, call service to create group and handle import
        let group_created = true; // Simulate success
        let import_scheduled = true; // Simulate async import
        let group_name = &params.name;

        if group_created {
            if import_scheduled {
                HttpResponse::Found()
                    .header("Location", format!("/groups/{}", group_name))
                    .body(format!("Group '{}' is being imported.", group_name))
            } else {
                HttpResponse::Found()
                    .header("Location", format!("/groups/{}", group_name))
                    .body("Group import could not be scheduled")
            }
        } else {
            HttpResponse::Found()
                .header("Location", "/groups/new#import-group-pane")
                .body("Group could not be imported: <error details>")
        }
    }
}
