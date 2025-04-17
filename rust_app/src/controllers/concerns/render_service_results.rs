use actix_web::{HttpResponse, Responder};
use serde::Serialize;

#[derive(Debug, Serialize)]
pub struct ServiceResult<T> {
    pub http_status: u16,
    pub body: T,
}

#[derive(Debug, Serialize)]
pub struct ErrorResult {
    pub status: String,
    pub message: String,
}

pub trait RenderServiceResults {
    fn success_response<T: Serialize>(&self, result: ServiceResult<T>) -> impl Responder {
        HttpResponse::Ok()
            .status(result.http_status)
            .json(result.body)
    }

    fn continue_polling_response(&self) -> impl Responder {
        HttpResponse::NoContent().json(serde_json::json!({
            "status": "processing",
            "message": "Not ready yet. Try again later."
        }))
    }

    fn error_response(&self, result: ErrorResult) -> impl Responder {
        HttpResponse::BadRequest().json(result)
    }
}
