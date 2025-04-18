use crate::models::Project;
use crate::utils::logging::ProjectStatsRefreshConflictsLogger;
use actix_web::{error::Error, HttpResponse};

pub trait ProjectStatsRefreshConflictsGuard {
    fn reject_if_build_artifacts_size_refreshing(&self, project: &Project) -> Result<(), Error>;
}

pub struct ProjectStatsRefreshConflictsGuardImpl;

impl ProjectStatsRefreshConflictsGuardImpl {
    pub fn new() -> Self {
        Self
    }
}

impl ProjectStatsRefreshConflictsGuard for ProjectStatsRefreshConflictsGuardImpl {
    fn reject_if_build_artifacts_size_refreshing(&self, project: &Project) -> Result<(), Error> {
        if project.refreshing_build_artifacts_size() {
            ProjectStatsRefreshConflictsLogger::warn_request_rejected_during_stats_refresh(
                project.id(),
            );

            Err(actix_web::error::ErrorConflict(
                "Action temporarily disabled. The project this pipeline belongs to is undergoing stats refresh.".to_string(),
            ))
        } else {
            Ok(())
        }
    }
}
