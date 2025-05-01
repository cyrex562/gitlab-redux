// Ported from: orig_app/app/controllers/groups/settings/repository_controller.rb
// This module handles group repository settings endpoints.
use axum::{Json, response::IntoResponse, extract::State, http::StatusCode};
use serde::{Deserialize, Serialize};
use crate::models::deploy_token::{DeployToken, DeployTokenParams};
use crate::services::deploy_tokens::CreateDeployTokenService;
use crate::AppState;

#[derive(Deserialize)]
pub struct CreateDeployTokenRequest {
    pub name: String,
    pub expires_at: Option<String>,
    pub read_repository: Option<bool>,
    pub read_registry: Option<bool>,
    pub write_registry: Option<bool>,
    pub read_package_registry: Option<bool>,
    pub write_package_registry: Option<bool>,
    pub username: Option<String>,
}

#[derive(Serialize)]
pub struct DeployTokenResponse {
    pub message: String,
    pub deploy_token: Option<DeployToken>,
}

pub async fn create_deploy_token(
    State(state): State<AppState>,
    Json(payload): Json<CreateDeployTokenRequest>,
) -> impl IntoResponse {
    // TODO: Implement authorization and group context extraction
    let service = CreateDeployTokenService::new(state.db.clone());
    let params = DeployTokenParams {
        name: payload.name,
        expires_at: payload.expires_at,
        read_repository: payload.read_repository.unwrap_or(false),
        read_registry: payload.read_registry.unwrap_or(false),
        write_registry: payload.write_registry.unwrap_or(false),
        read_package_registry: payload.read_package_registry.unwrap_or(false),
        write_package_registry: payload.write_package_registry.unwrap_or(false),
        username: payload.username,
    };
    match service.execute(params).await {
        Ok(token) => (
            StatusCode::CREATED,
            Json(DeployTokenResponse {
                message: "Your new group deploy token has been created.".to_string(),
                deploy_token: Some(token),
            }),
        ).into_response(),
        Err(e) => (
            StatusCode::BAD_REQUEST,
            Json(DeployTokenResponse {
                message: format!("Failed to create deploy token: {}", e),
                deploy_token: None,
            }),
        ).into_response(),
    }
}
