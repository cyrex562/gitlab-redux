// Ported from: orig_app/app/controllers/oauth/device_authorizations_controller.rb
// Ported: 2025-05-04
//
// Handles device authorization grant flows (Doorkeeper::DeviceAuthorizationGrant::DeviceAuthorizationsController).
// Actions: index, confirm

use actix_web::{get, post, web, HttpResponse, ResponseError};
use serde::{Deserialize, Serialize};
use serde_json::json;
use thiserror::Error;

#[derive(Debug, Error)]
pub enum DeviceAuthError {
    #[error("Invalid user code")]
    InvalidUserCode,
    #[error("Database error: {0}")]
    DatabaseError(String),
    #[error("Internal server error")]
    InternalError,
}

impl ResponseError for DeviceAuthError {
    fn error_response(&self) -> HttpResponse {
        match self {
            DeviceAuthError::InvalidUserCode => HttpResponse::BadRequest().json(json!({
                "error": "invalid_user_code",
                "error_description": "The user code is invalid"
            })),
            DeviceAuthError::DatabaseError(_) => HttpResponse::InternalServerError().json(json!({
                "error": "server_error",
                "error_description": "A database error occurred"
            })),
            DeviceAuthError::InternalError => HttpResponse::InternalServerError().json(json!({
                "error": "server_error",
                "error_description": "An internal error occurred"
            })),
        }
    }
}

#[derive(Debug, Deserialize)]
pub struct ConfirmRequest {
    user_code: String,
}

#[derive(Debug, Serialize)]
pub struct DeviceGrant {
    scopes: String,
}

#[derive(Debug, Serialize)]
pub struct IndexResponse {
    view: String,
}

#[derive(Debug, Serialize)]
pub struct ConfirmResponse {
    view: String,
    scopes: String,
}

pub struct DeviceAuthorizationsController;

impl DeviceAuthorizationsController {
    pub fn new() -> Self {
        Self
    }

    /// GET /oauth/device_authorizations
    /// Returns HTML template for device authorization or no content for JSON requests
    #[get("/oauth/device_authorizations")]
    pub async fn index(
        content_type: Option<web::Header<String>>,
    ) -> Result<HttpResponse, DeviceAuthError> {
        match content_type.as_deref().map(|h| h.as_str()) {
            Some("application/json") => Ok(HttpResponse::NoContent().finish()),
            _ => Ok(HttpResponse::Ok()
                .content_type("text/html")
                .json(IndexResponse {
                    view: "doorkeeper/device_authorization_grant/index".to_string(),
                })),
        }
    }

    /// POST /oauth/device_authorizations/confirm
    /// Confirms a device authorization request using a user code
    #[post("/oauth/device_authorizations/confirm")]
    pub async fn confirm(
        req: web::Json<ConfirmRequest>,
        content_type: Option<web::Header<String>>,
    ) -> Result<HttpResponse, DeviceAuthError> {
        // Look up device grant using user code
        let device_grant = find_device_grant(&req.user_code).await?;

        match content_type.as_deref().map(|h| h.as_str()) {
            Some("application/json") => Ok(HttpResponse::NoContent().finish()),
            _ => Ok(HttpResponse::Ok()
                .content_type("text/html")
                .json(ConfirmResponse {
                    view: "doorkeeper/device_authorization_grant/authorize".to_string(),
                    scopes: device_grant.scopes,
                })),
        }
    }
}

// Database interaction placeholder
async fn find_device_grant(user_code: &str) -> Result<DeviceGrant, DeviceAuthError> {
    // TODO: Implement actual database lookup
    // For now returning empty scopes as placeholder
    if user_code.is_empty() {
        return Err(DeviceAuthError::InvalidUserCode);
    }

    Ok(DeviceGrant {
        scopes: String::new(),
    })
}

// Future improvements needed:
// 1. Implement proper device grant model and database interactions
// 2. Add proper user authentication checks
// 3. Add proper template rendering
// 4. Add user session handling
// 5. Add tests
