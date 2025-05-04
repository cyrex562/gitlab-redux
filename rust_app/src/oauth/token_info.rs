// Ported from gitlab-redux/orig_app/app/controllers/oauth/token_info_controller.rb
use axum::{
    extract::{State, TypedHeader},
    response::Response,
    Json,
};
use headers::Authorization;
use oauth2::TokenResponse;
use serde_json::{json, Value};
use crate::{
    auth::enforce_two_factor::EnforceTwoFactor,
    error::Error,
    AppState,
};

pub async fn show(
    enforce_2fa: EnforceTwoFactor,
    State(state): State<AppState>,
    TypedHeader(auth): TypedHeader<Authorization>,
) -> Result<Response, Error> {
    if let Err(e) = enforce_2fa.check() {
        return Err(e.into());
    }

    let token = state.oauth_provider.validate_token(auth.token()).await?;
    
    if token.is_accessible() {
        let mut token_json: Value = serde_json::to_value(&token)?;
        
        // Maintain backwards compatibility by adding scopes and expires_in_seconds
        if let Value::Object(ref mut map) = token_json {
            if let Some(scope) = map.get("scope").cloned() {
                map.insert("scopes".to_string(), scope);
            }
            if let Some(expires_in) = map.get("expires_in").cloned() {
                map.insert("expires_in_seconds".to_string(), expires_in);
            }
        }
        
        Ok(Json(token_json).into_response())
    } else {
        let error = state.oauth_provider.invalid_token_response();
        Ok(Response::builder()
            .status(error.status())
            .headers(error.headers().clone())
            .body(Json(error.body()).into_response().into_body())
            .unwrap())
    }
}