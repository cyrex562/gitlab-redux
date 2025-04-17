use actix_web::HttpResponse;
use serde::Deserialize;

use crate::models::flash_message::FlashMessage;

pub trait HookExecutionNotice {
    fn set_hook_execution_notice(&self, result: &HookExecutionResult);
    fn flash(&self) -> &FlashMessage;
    fn flash_mut(&mut self) -> &mut FlashMessage;
}

#[derive(Debug, Deserialize)]
pub struct HookExecutionResult {
    pub success: bool,
    pub message: String,
    pub payload: HookExecutionPayload,
}

#[derive(Debug, Deserialize)]
pub struct HookExecutionPayload {
    pub http_status: Option<i32>,
}

impl HookExecutionNotice for dyn HookExecutionNotice {
    fn set_hook_execution_notice(&self, result: &HookExecutionResult) {
        let http_status = result.payload.http_status;
        let message = &result.message;

        if let Some(status) = http_status {
            if status >= 200 && status < 400 {
                self.flash_mut()
                    .add_notice(format!("Hook executed successfully: HTTP {}", status));
            } else {
                self.flash_mut().add_alert(format!(
                    "Hook executed successfully but returned HTTP {} {}",
                    status, message
                ));
            }
        } else {
            self.flash_mut()
                .add_alert(format!("Hook execution failed: {}", message));
        }
    }
}
