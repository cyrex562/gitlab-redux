use actix_web::HttpResponse;
use log::warn;

pub struct Project {
    pub id: i32,
}

impl Project {
    pub fn refreshing_build_artifacts_size(&self) -> bool {
        // Implementation would depend on your project state management
        false
    }
}

pub struct ProjectStatsRefreshConflictsLogger;

impl ProjectStatsRefreshConflictsLogger {
    pub fn warn_request_rejected_during_stats_refresh(project_id: i32) {
        warn!(
            "Request rejected during stats refresh for project {}",
            project_id
        );
    }
}

pub trait ProjectStatsRefreshConflictsGuard {
    fn reject_if_build_artifacts_size_refreshing(&self, project: &Project) -> Option<HttpResponse> {
        if project.refreshing_build_artifacts_size() {
            ProjectStatsRefreshConflictsLogger::warn_request_rejected_during_stats_refresh(
                project.id,
            );

            Some(
                HttpResponse::Conflict()
                    .content_type("text/plain")
                    .body("Action temporarily disabled. The project this pipeline belongs to is undergoing stats refresh.")
            )
        } else {
            None
        }
    }
}
