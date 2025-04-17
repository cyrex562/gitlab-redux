use actix_web::{web, HttpRequest, HttpResponse};
use serde::Deserialize;

use crate::models::spammable::Spammable;
use crate::models::spam::SpamParams;
use crate::utils::request_context::RequestContext;
use crate::controllers::concerns::spammable_actions::captcha_check::Common;

#[derive(Debug, Deserialize)]
pub struct CaptchaParams {
    #[serde(rename = "g-recaptcha-response")]
    recaptcha_response: Option<String>,
    spam_log_id: Option<i32>,
}

pub trait HtmlFormatActionsSupport: Common {
    fn request(&self) -> &HttpRequest;
    fn request_mut(&mut self) -> &mut HttpRequest;
    fn request_context(&self) -> &RequestContext;
    fn request_context_mut(&mut self) -> &mut RequestContext;
    fn render_captcha_check(&self) -> HttpResponse;

    fn with_captcha_check_html_format<F, Fut, R>(
        &self,
        spammable: &dyn Spammable,
        action: F,
    ) -> HttpResponse
    where
        F: FnOnce() -> Fut,
        Fut: std::future::Future<Output = R>,
        R: Into<HttpResponse>,
    {
        let captcha_render_lambda = || self.render_captcha_check();
        self.with_captcha_check_common(spammable, captcha_render_lambda)
    }

    fn convert_html_spam_params_to_headers(&mut self) {
        let params = web::Form::<CaptchaParams>::from_request(self.request())
            .unwrap_or_else(|_| web::Form(CaptchaParams {
                recaptcha_response: None,
                spam_log_id: None,
            }));

        let mut headers_changed = false;

        if let Some(recaptcha_response) = &params.recaptcha_response {
            self.request_mut()
                .headers_mut()
                .insert("X-GitLab-Captcha-Response", recaptcha_response.parse().unwrap());
            headers_changed = true;
        }

        if let Some(spam_log_id) = params.spam_log_id {
            self.request_mut()
                .headers_mut()
                .insert("X-GitLab-Spam-Log-Id", spam_log_id.to_string().parse().unwrap());
            headers_changed = true;
        }

        // Reset the spam_params on the request context, since they have changed mid-request
        if headers_changed {
            let spam_params = SpamParams::new_from_request(self.request());
            self.request_context_mut().set_spam_params(spam_params);
        }
    }
} 