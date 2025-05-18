// Ported from: orig_app/app/controllers/projects/analytics/cycle_analytics/stages_controller.rb
use crate::utils::strong_memoize::StrongMemoizeImpl;
use axum::{http::StatusCode, response::IntoResponse, Json};
use serde::Serialize;

#[derive(Serialize)]
pub struct Stage {
    pub name: String,
}

#[derive(Serialize)]
pub struct Permissions {
    pub allowed_stages: Vec<String>,
}

pub struct StagesController {
    memoizer: StrongMemoizeImpl,
}

impl StagesController {
    pub fn new() -> Self {
        Self {
            memoizer: StrongMemoizeImpl::new(),
        }
    }

    pub async fn get_namespace(&self, project_namespace: String) -> String {
        project_namespace
    }

    pub async fn cycle_analytics_configuration(
        &self,
        stages: Vec<Stage>,
        permissions: &Permissions,
    ) -> Vec<Stage> {
        stages
            .into_iter()
            .filter(|stage| self.permitted_stage(stage, permissions))
            .collect()
    }

    pub async fn only_default_value_stream_is_allowed(
        &self,
        value_stream_id: &str,
        default_name: &str,
    ) -> Result<(), StatusCode> {
        if value_stream_id == default_name {
            Ok(())
        } else {
            Err(StatusCode::FORBIDDEN)
        }
    }

    pub fn permitted_stage(&self, stage: &Stage, permissions: &Permissions) -> bool {
        permissions.allowed_stages.contains(&stage.name)
    }

    pub fn permissions(&self, user: &str, project: &str) -> Permissions {
        // Placeholder: Replace with real permission logic
        self.memoizer
            .memoized(format!("{}:{}", user, project), || Permissions {
                allowed_stages: vec!["plan".to_string(), "code".to_string(), "test".to_string()],
            })
    }

    pub fn authorize_stage(
        &self,
        stage: &Stage,
        permissions: &Permissions,
    ) -> Result<(), StatusCode> {
        if self.permitted_stage(stage, permissions) {
            Ok(())
        } else {
            Err(StatusCode::FORBIDDEN)
        }
    }

    pub fn requests_default_value_stream(&self, value_stream_id: &str, default_name: &str) -> bool {
        value_stream_id == default_name
    }
}
