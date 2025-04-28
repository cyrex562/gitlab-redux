// Ported from: orig_app/app/controllers/concerns/project_stats_refresh_conflicts_guard.rb
// Ported on: 2025-04-28
//
// This file implements the ProjectStatsRefreshConflictsGuard concern in Rust.
// It provides a guard to reject requests if the project's build artifacts size is being refreshed.
//
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
