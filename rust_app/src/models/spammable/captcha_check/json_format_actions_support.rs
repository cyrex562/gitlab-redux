use super::common::CaptchaCheckCommon;
use crate::spammable::{SpamActionResponseFields, Spammable};
use actix_web::{web, HttpResponse};
use std::future::Future;

pub trait JsonFormatActionsSupport: CaptchaCheckCommon + SpamActionResponseFields {
    fn with_captcha_check_json_format<F, Fut>(
        &self,
        spammable: &dyn Spammable,
        action: impl FnOnce() -> Fut,
    ) -> HttpResponse
    where
        F: Future<Output = HttpResponse>,
        Fut: Future<Output = HttpResponse>,
    {
        let captcha_render_lambda = || {
            // Return 409 Conflict status for CAPTCHA requirement
            HttpResponse::Conflict().json(self.spam_action_response_fields(spammable))
        };

        self.with_captcha_check_common(spammable, captcha_render_lambda, action)
    }
}
