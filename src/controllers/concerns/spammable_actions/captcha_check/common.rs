use actix_web::HttpResponse;
use std::future::Future;

use crate::models::spammable::Spammable;
use crate::utils::recaptcha::Recaptcha;

pub trait Common {
    fn with_captcha_check_common<F, Fut, R>(
        &self,
        spammable: &dyn Spammable,
        captcha_render_lambda: F,
    ) -> HttpResponse
    where
        F: FnOnce() -> HttpResponse,
        Fut: Future<Output = R>,
        R: Into<HttpResponse>;

    fn recaptcha(&self) -> &Recaptcha;
}

impl<T: Common> Common for T {
    fn with_captcha_check_common<F, Fut, R>(
        &self,
        spammable: &dyn Spammable,
        captcha_render_lambda: F,
    ) -> HttpResponse
    where
        F: FnOnce() -> HttpResponse,
        Fut: Future<Output = R>,
        R: Into<HttpResponse>,
    {
        // If the Spammable indicates that CAPTCHA is not necessary (either due to it not being flagged
        // as spam, or if spam/captcha is disabled for some reason), then we will go ahead and
        // return a success response.
        if !spammable.render_recaptcha() {
            return HttpResponse::Ok().finish();
        }

        // If we got here, we need to render the CAPTCHA instead of proceeding with the action.
        // We will present a CAPTCHA to be solved by executing the lambda which was passed
        // as the `captcha_render_lambda` argument. This lambda contains either the HTML-specific or
        // JSON-specific behavior to cause the CAPTCHA modal to be rendered.
        self.recaptcha().load_configurations();
        captcha_render_lambda()
    }
}
