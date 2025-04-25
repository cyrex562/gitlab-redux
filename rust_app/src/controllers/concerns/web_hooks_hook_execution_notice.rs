// Ported from orig_app/app/controllers/concerns/web_hooks/hook_execution_notice.rb
// Provides a trait for setting hook execution flash notices
use crate::controllers::concerns::skips_already_signed_in_message::Flash;

/// Represents the result of a hook execution
pub struct HookExecutionResult {
    pub payload: HookExecutionPayload,
    pub message: String,
}

pub struct HookExecutionPayload {
    pub http_status: Option<u16>,
}

pub trait HookExecutionNotice {
    fn set_hook_execution_notice(&self, flash: &mut Flash, result: &HookExecutionResult);
}

pub struct HookExecutionNoticeHandler;

impl HookExecutionNotice for HookExecutionNoticeHandler {
    fn set_hook_execution_notice(&self, flash: &mut Flash, result: &HookExecutionResult) {
        let http_status = result.payload.http_status;
        let message = &result.message;
        if let Some(status) = http_status {
            if (200..400).contains(&status) {
                flash.messages.insert(
                    "notice".to_string(),
                    format!("Hook executed successfully: HTTP {}", status),
                );
            } else {
                flash.messages.insert(
                    "alert".to_string(),
                    format!(
                        "Hook executed successfully but returned HTTP {} {}",
                        status, message
                    ),
                );
            }
        } else {
            flash.messages.insert(
                "alert".to_string(),
                format!("Hook execution failed: {}", message),
            );
        }
    }
}
