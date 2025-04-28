// Ported from: orig_app/app/controllers/concerns/render_service_results.rb
// This file implements the RenderServiceResults concern in Rust.

use actix_web::HttpResponse;
use serde::Serialize;

#[derive(Debug, Serialize)]
pub struct ServiceResult<T: Serialize> {
    pub http_status: u16,
    pub body: T,
}

pub trait RenderServiceResults {
    fn success_response<T: Serialize>(result: &ServiceResult<T>) -> HttpResponse {
        HttpResponse::build(
            actix_web::http::StatusCode::from_u16(result.http_status)
                .unwrap_or(actix_web::http::StatusCode::OK),
        )
        .json(&result.body)
    }

    fn continue_polling_response() -> HttpResponse {
        HttpResponse::NoContent().json(serde_json::json!({
            "status": "processing",
            "message": "Not ready yet. Try again later."
        }))
    }

    fn error_response(status: Option<u16>, message: &str, status_str: &str) -> HttpResponse {
        let http_status = status.unwrap_or(400);
        HttpResponse::build(
            actix_web::http::StatusCode::from_u16(http_status)
                .unwrap_or(actix_web::http::StatusCode::BAD_REQUEST),
        )
        .json(serde_json::json!({
            "status": status_str,
            "message": message
        }))
    }
}

// Example usage:
// impl RenderServiceResults for MyController {}
