use actix_web::{web, HttpResponse, Responder};
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct Build {
    pub id: i32,
    pub debug_mode: bool,
}

pub struct AuthBuildTrace;

impl AuthBuildTrace {
    pub async fn authorize_read_build_trace(build: &Build, user_id: i32) -> impl Responder {
        if !Self::can_read_build_trace(user_id, build) {
            if build.debug_mode {
                return HttpResponse::Forbidden().json(serde_json::json!({
                    "message": "You must have developer or higher permissions in the associated project to view job logs when debug trace is enabled. To disable debug trace, set the 'CI_DEBUG_TRACE' and 'CI_DEBUG_SERVICES' variables to 'false' in your pipeline configuration or CI/CD settings. If you must view this job log, a project maintainer or owner must add you to the project with developer permissions or higher."
                }));
            } else {
                return HttpResponse::Forbidden().json(serde_json::json!({
                    "message": "The current user is not authorized to access the job log."
                }));
            }
        }

        HttpResponse::Ok().finish()
    }

    fn can_read_build_trace(user_id: i32, build: &Build) -> bool {
        // TODO: Implement permission check logic
        false
    }
} 