use super::common::CaptchaCheckCommon;
use crate::spammable::{SpamParams, Spammable};
use actix_web::{web, HttpRequest, HttpResponse};
use std::future::Future;

pub trait HtmlFormatActionsSupport: CaptchaCheckCommon {
    fn with_captcha_check_html_format<F, Fut>(
        &self,
        spammable: &dyn Spammable,
        action: impl FnOnce() -> Fut,
    ) -> HttpResponse
    where
        F: Future<Output = HttpResponse>,
        Fut: Future<Output = HttpResponse>,
    {
        let captcha_render_lambda = || {
            // TODO: Implement proper captcha_check template rendering
            HttpResponse::Ok().json(serde_json::json!({
                "captcha_required": true
            }))
        };

        self.with_captcha_check_common(spammable, captcha_render_lambda, action)
    }

    fn convert_html_spam_params_to_headers(&self, req: &mut HttpRequest) {
        // TODO: Implement proper request parameter handling
        let captcha_response = req.query_string().get("g-recaptcha-response");
        let spam_log_id = req.query_string().get("spam_log_id");

        if captcha_response.is_some() || spam_log_id.is_some() {
            // TODO: Implement proper header setting
            // req.headers_mut().insert("X-GitLab-Captcha-Response", captcha_response);
            // req.headers_mut().insert("X-GitLab-Spam-Log-Id", spam_log_id);

            // TODO: Implement proper request context handling
            // Gitlab::RequestContext.instance.spam_params = SpamParams::new_from_request(req);
        }
    }
}
