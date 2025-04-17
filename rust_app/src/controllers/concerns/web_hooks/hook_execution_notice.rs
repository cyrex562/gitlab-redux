use actix_web::{HttpRequest, HttpResponse};
use serde_json::json;

pub struct HookResult {
    pub http_status: Option<i32>,
    pub message: String,
}

pub trait HookExecutionNotice {
    fn set_hook_execution_notice(&self, result: &HookResult) -> HttpResponse {
        let (status, message) = match result.http_status {
            Some(status) if (200..400).contains(&status) => (
                HttpResponse::Ok(),
                format!("Hook executed successfully: HTTP {}", status),
            ),
            Some(status) => (
                HttpResponse::Ok(),
                format!(
                    "Hook executed successfully but returned HTTP {} {}",
                    status, result.message
                ),
            ),
            None => (
                HttpResponse::InternalServerError(),
                format!("Hook execution failed: {}", result.message),
            ),
        };

        status.json(json!({
            "notice": message,
            "status": status.status().as_u16()
        }))
    }
}
