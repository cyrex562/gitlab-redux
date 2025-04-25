// Ported from Ruby: app/controllers/concerns/spammable_actions/captcha_check/html_format_actions_support.rb
// Ported on 2025-04-24
// This module provides support for HTML form submissions that require spam/captcha checks.
// Only use for HTML (not JSON/AJAX) form actions.

use actix_web::{HttpRequest, HttpResponse};

// Placeholder for shared logic, to be implemented or imported as needed.
pub trait CaptchaCheckCommon {
    fn with_captcha_check_common<F, R>(
        &self,
        spammable: &str,
        captcha_render: F,
        block: impl FnOnce() -> R,
    ) -> R
    where
        F: Fn() -> HttpResponse;
}

pub trait HtmlFormatActionsSupport: CaptchaCheckCommon {
    // Should be called before create/update actions to convert params to headers
    fn convert_html_spam_params_to_headers(&self, req: &mut HttpRequest);

    // Wraps logic with captcha check, rendering captcha page if needed
    fn with_captcha_check_html_format<R>(
        &self,
        spammable: &str,
        req: &mut HttpRequest,
        block: impl FnOnce() -> R,
    ) -> R {
        let captcha_render = || HttpResponse::Ok().body("Render captcha_check page");
        self.with_captcha_check_common(spammable, captcha_render, block)
    }
}

// Example implementation for a controller struct
impl<T: CaptchaCheckCommon> HtmlFormatActionsSupport for T {
    fn convert_html_spam_params_to_headers(&self, req: &mut HttpRequest) {
        // Pseudocode: extract params and set headers
        // In Actix, params are usually in query or form, headers in req.headers_mut()
        // let params = ...;
        // if let Some(recaptcha) = params.get("g-recaptcha-response") {
        //     req.headers_mut().insert("X-GitLab-Captcha-Response", recaptcha.parse().unwrap());
        // }
        // if let Some(spam_log_id) = params.get("spam_log_id") {
        //     req.headers_mut().insert("X-GitLab-Spam-Log-Id", spam_log_id.parse().unwrap());
        // }
        // Reset spam params on request context if needed
    }
}
