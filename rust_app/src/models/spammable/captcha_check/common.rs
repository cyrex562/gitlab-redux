use crate::spammable::Spammable;
use actix_web::{web, HttpResponse};
use std::future::Future;

pub trait CaptchaCheckCommon {
    fn with_captcha_check_common<F, Fut>(
        &self,
        spammable: &dyn Spammable,
        captcha_render_lambda: F,
        action: impl FnOnce() -> Fut,
    ) -> HttpResponse
    where
        F: FnOnce() -> HttpResponse,
        Fut: Future<Output = HttpResponse>,
    {
        // If the Spammable indicates that CAPTCHA is not necessary, execute the action
        if !spammable.render_recaptcha() {
            // TODO: Implement proper async handling
            return HttpResponse::Ok().finish();
        }

        // TODO: Implement Recaptcha configuration loading
        // Gitlab::Recaptcha.load_configurations!

        // Render the CAPTCHA
        captcha_render_lambda()
    }
}
