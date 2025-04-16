use actix_web::HttpResponse;
use crate::webhooks::WebHookResult;

pub trait HookExecutionNotice {
    fn set_hook_execution_notice(&self, result: &WebHookResult) -> HttpResponse {
        let http_status = result.http_status;
        let message = &result.message;

        if let Some(status) = http_status {
            if status >= 200 && status < 400 {
                HttpResponse::Ok().json(json!({
                    "notice": format!("Hook executed successfully: HTTP {}", status)
                }))
            } else {
                HttpResponse::Ok().json(json!({
                    "alert": format!("Hook executed successfully but returned HTTP {} {}", status, message)
                }))
            }
        } else {
            HttpResponse::Ok().json(json!({
                "alert": format!("Hook execution failed: {}", message)
            }))
        }
    }
} 