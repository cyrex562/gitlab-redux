// Ported from: orig_app/app/controllers/projects/error_tracking/stack_traces_controller.rb
// Ported on: 2025-05-19
// This file implements the Projects::ErrorTracking::StackTracesController from the Ruby codebase.

use actix_web::{web, HttpResponse, Responder};
use serde_json::json;

use crate::controllers::projects::error_tracking::base_controller::ProjectsErrorTrackingBaseController;
use crate::decorators::error_tracking::stack_trace_highlight_decorator::StackTraceHighlightDecorator;
use crate::serializers::error_tracking::error_event_serializer::ErrorEventSerializer;
use crate::services::error_tracking::issue_latest_event_service::IssueLatestEventService;

#[derive(Debug, serde::Deserialize)]
pub struct StackTracesParams {
    pub issue_id: String,
}

pub struct ProjectsErrorTrackingStackTracesController {
    base: ProjectsErrorTrackingBaseController,
}

impl ProjectsErrorTrackingStackTracesController {
    pub fn new(base: ProjectsErrorTrackingBaseController) -> Self {
        Self { base }
    }

    /// GET /projects/:project_id/error_tracking/stack_traces
    pub async fn index(&self, params: web::Query<StackTracesParams>) -> impl Responder {
        if let Err(e) = self.base.authorize_read_sentry_issue() {
            return e;
        }
        // Set polling interval header (handled in middleware or response builder)
        // self.base.set_polling_interval(&mut response);
        let project = self.base.base.project();
        let current_user = self.base.base.current_user();
        let result =
            IssueLatestEventService::new(project, current_user, &params.issue_id).execute();
        if result.status == "success" {
            let highlighted = StackTraceHighlightDecorator::decorate(&result.latest_event);
            let error_json =
                ErrorEventSerializer::new(project, current_user).represent(&highlighted);
            HttpResponse::Ok().json(json!({ "error": error_json }))
        } else {
            HttpResponse::BadRequest().json(json!({ "message": result.message }))
        }
    }
}
