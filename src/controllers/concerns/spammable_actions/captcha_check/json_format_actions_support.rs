use actix_web::HttpResponse;
use serde::Serialize;
use std::future::Future;

use crate::controllers::concerns::spammable_actions::captcha_check::Common;
use crate::models::spam::HasSpamActionResponseFields;
use crate::models::spammable::Spammable;

pub trait JsonFormatActionsSupport: Common + HasSpamActionResponseFields {
    fn with_captcha_check_json_format<F, Fut, R>(
        &self,
        spammable: &dyn Spammable,
        action: F,
    ) -> HttpResponse
    where
        F: FnOnce() -> Fut,
        Fut: Future<Output = R>,
        R: Into<HttpResponse>,
    {
        // NOTE: "409 - Conflict" seems to be the most appropriate HTTP status code for a response
        // which requires a CAPTCHA to be solved in order for the request to be resubmitted.
        // https://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html#sec10.4.10
        let captcha_render_lambda = || {
            let response_fields = self.spam_action_response_fields(spammable);
            HttpResponse::Conflict().json(response_fields)
        };

        self.with_captcha_check_common(spammable, captcha_render_lambda)
    }
}
