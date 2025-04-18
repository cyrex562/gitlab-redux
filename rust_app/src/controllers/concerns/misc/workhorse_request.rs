use actix_web::{web, HttpRequest, HttpResponse, error::ErrorUnauthorized};
use serde::{Deserialize, Serialize};
use std::sync::Arc;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct WorkhorseRequest {
    token: String,
}

impl WorkhorseRequest {
    pub fn new(token: String) -> Self {
        Self { token }
    }
}

pub struct WorkhorseRequestHandler {
    db: Arc<sqlx::PgPool>,
}

impl WorkhorseRequestHandler {
    pub fn new(db: Arc<sqlx::PgPool>) -> Self {
        Self { db }
    }

    pub async fn verify_workhorse_api(&self, req: &HttpRequest) -> Result<(), actix_web::Error> {
        // TODO: Implement workhorse API verification
        let token = req.headers()
            .get("X-Gitlab-Workhorse")
            .and_then(|h| h.to_str().ok())
            .ok_or_else(|| ErrorUnauthorized("Missing workhorse token"))?;

        if token != "true" {
            return Err(ErrorUnauthorized("Invalid workhorse token"));
        }

        Ok(())
    }
} 