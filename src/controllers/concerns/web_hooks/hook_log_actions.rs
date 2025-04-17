use actix_web::{web, HttpResponse, Responder};
use serde::Deserialize;

use crate::models::web_hook::WebHook;
use crate::models::web_hook_log::WebHookLog;
use crate::services::web_hooks::events::ResendService;

#[derive(Debug, Deserialize)]
pub struct HookLogParams {
    id: i32,
}

pub trait HookLogActions {
    fn hook(&self) -> &WebHook;
    fn current_user(&self) -> &User;

    fn show(&self, params: web::Path<HookLogParams>) -> impl Responder {
        // Hide search settings is handled by the view layer
        HttpResponse::Ok().finish()
    }

    fn retry(&self, params: web::Path<HookLogParams>) -> impl Responder {
        let hook_log = self.hook_log(&params);
        let result = self.execute_hook(hook_log);

        if result.is_success() {
            HttpResponse::Found()
                .header("Location", self.after_retry_redirect_path())
                .finish()
        } else {
            // In a real implementation, we'd need to handle flash messages
            HttpResponse::Found()
                .header("Location", self.after_retry_redirect_path())
                .finish()
        }
    }

    fn hook_log(&self, params: &HookLogParams) -> &WebHookLog {
        self.hook()
            .web_hook_logs()
            .find(params.id)
            .expect("Hook log not found")
    }

    fn execute_hook(&self, hook_log: &WebHookLog) -> HookExecutionResult {
        let service = ResendService::new(hook_log, self.current_user());
        let result = service.execute();
        self.set_hook_execution_notice(&result);
        result
    }

    fn after_retry_redirect_path(&self) -> String;
    fn set_hook_execution_notice(&self, result: &HookExecutionResult);
}

#[derive(Debug)]
pub struct HookExecutionResult {
    success: bool,
    message: String,
}

impl HookExecutionResult {
    pub fn is_success(&self) -> bool {
        self.success
    }

    pub fn message(&self) -> &str {
        &self.message
    }
}
