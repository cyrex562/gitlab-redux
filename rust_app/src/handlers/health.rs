use axum::{http::StatusCode, response::IntoResponse, Json};
use serde_json::json;

pub async fn check() -> impl IntoResponse {
    Json(json!({
        "status": "ok",
        "timestamp": chrono::Utc::now().to_rfc3339()
    }))
}
