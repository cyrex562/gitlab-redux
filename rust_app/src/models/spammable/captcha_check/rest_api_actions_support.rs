use actix_web::{web, HttpResponse};
use crate::spammable::{Spammable, SpamActionResponseFields};
use super::common::CaptchaCheckCommon;
use std::future::Future;

pub trait RestApiActionsSupport: CaptchaCheckCommon + SpamActionResponseFields {
    fn with_captcha_check_rest_api<F, Fut>(
        &self,
        spammable: &dyn Spammable,
        action: impl FnOnce() -> Fut,
    ) -> HttpResponse 
    where
        F: Future<Output = HttpResponse>,
        Fut: Future<Output = HttpResponse>,
    {
        let captcha_render_lambda = || {
            let mut fields = self.spam_action_response_fields(spammable);
            
            // Remove spam field as it's not needed in API response
            if let serde_json::Value::Object(ref mut map) = fields {
                map.remove("spam");
            }
            
            // Add error message
            let error_message = spammable.errors().join(", ");
            if let serde_json::Value::Object(ref mut map) = fields {
                map.insert(
                    "message".to_string(), 
                    serde_json::json!({ "error": error_message })
                );
            }
            
            // Return 409 Conflict status for CAPTCHA requirement
            HttpResponse::Conflict().json(fields)
        };
        
        self.with_captcha_check_common(spammable, captcha_render_lambda, action)
    }
} 