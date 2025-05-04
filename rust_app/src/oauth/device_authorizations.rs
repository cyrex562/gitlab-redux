// Ported from gitlab-redux/orig_app/app/controllers/oauth/device_authorizations_controller.rb
use axum::{
    extract::{Query, State},
    response::{IntoResponse, Response},
    Json,
};
use serde::Deserialize;
use crate::{
    error::Error,
    AppState,
    templates::RenderContext,
};

#[derive(Deserialize)]
pub struct ConfirmParams {
    user_code: String,
}

pub async fn index() -> impl IntoResponse {
    // Return different responses based on content type
    // HTML: render the device authorization template
    // JSON: return empty response with 204 status
    Response::builder()
        .header("Vary", "Accept")
        .status(204)
        .body(axum::body::Body::empty())
        .unwrap()
}

pub async fn confirm(
    State(state): State<AppState>,
    Query(params): Query<ConfirmParams>,
) -> Result<Response, Error> {
    // Find device grant by user code
    let device_grant = state.oauth_provider
        .find_device_grant(&params.user_code)
        .await?;
    
    let scopes = device_grant.map(|g| g.scopes).unwrap_or_default();
    
    let mut ctx = RenderContext::new();
    ctx.insert("scopes", &scopes);
    
    // Return different responses based on content type
    // HTML: render the authorization template
    // JSON: return empty response with 204 status
    Ok(Response::builder()
        .header("Vary", "Accept")
        .status(204)
        .body(axum::body::Body::empty())
        .unwrap())
}