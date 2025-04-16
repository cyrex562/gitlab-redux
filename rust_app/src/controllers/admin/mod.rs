use actix_web::{web, HttpResponse, Responder};
use serde::{Deserialize, Serialize};

mod abuse_reports;
mod applications;
mod background_jobs;
mod background_migrations;
mod batched_jobs;
mod broadcast_messages;
mod clusters;
mod cohorts;
mod dashboard;
mod deploy_keys;
mod dev_ops_report;
mod health_check;
mod impersonation_tokens;
mod system_info;
mod version_check;

pub use abuse_reports::AbuseReportsController;
pub use applications::ApplicationsController;
pub use background_jobs::BackgroundJobsController;
pub use background_migrations::BackgroundMigrationsController;
pub use batched_jobs::BatchedJobsController;
pub use broadcast_messages::BroadcastMessagesController;
pub use clusters::ClustersController;
pub use cohorts::CohortsController;
pub use dashboard::DashboardController;
pub use deploy_keys::*;
pub use dev_ops_report::*;
pub use health_check::HealthCheckController;
pub use impersonation_tokens::*;
pub use system_info::SystemInfoController;
pub use version_check::VersionCheckController;

/// Base controller for all admin controllers
///
/// Provides a base class for Admin controllers to subclass.
/// Automatically sets the layout and ensures an administrator is logged in.
pub struct ApplicationController;

impl ApplicationController {
    /// Create a new admin application controller
    pub fn new() -> Self {
        Self
    }

    /// Ensure the user is an administrator
    pub fn ensure_admin(&self) -> Result<(), impl Responder> {
        // TODO: Implement proper admin authentication
        // This is a placeholder implementation
        Ok(())
    }

    /// Get the admin layout
    pub fn get_layout(&self) -> &'static str {
        "admin"
    }
}
