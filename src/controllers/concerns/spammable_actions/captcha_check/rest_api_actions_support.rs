use actix_web::HttpResponse;
use serde::Serialize;
use std::future::Future;

use crate::controllers::concerns::spammable_actions::captcha_check::Common;
use crate::models::spam::HasSpamActionResponseFields;
use crate::models::spammable::Spammable;

pub trait RestApiActionsSupport: Common + HasSpamActionResponseFields {
    fn render_structured_api_error<T: Serialize>(&self, fields: T, status: u16) -> HttpResponse;
    fn to_sentence(&self, messages: &[String]) -> String;

    fn with_captcha_check_rest_api<F, Fut, R>(
        &self,
        spammable: &dyn Spammable,
        action: F,
    ) -> HttpResponse
    where
        F: FnOnce() -> Fut,
        Fut: Future<Output = R>,
        R: Into<HttpResponse>,
    {
        // In the case of the REST API, the request is handled by Grape, so if there is a spam-related
        // error, we don't render directly, instead we will pass the error message and other necessary
        // fields to the Grape api error helper for it to handle.
        let captcha_render_lambda = || {
            let mut fields = self.spam_action_response_fields(spammable);

            // Remove the spam field
            fields.remove("spam");

            // NOTE: "409 - Conflict" seems to be the most appropriate HTTP status code for a response
            // which requires a CAPTCHA to be solved in order for the request to be resubmitted.
            // https://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html#sec10.4.10
            let status = 409;

            // NOTE: This nested 'error' key may not be consistent with all other API error responses,
            // because they are not currently consistent across different API endpoints
            // and models. Some (snippets) will nest errors in an errors key like this,
            // while others (issues) will return the model's errors hash without an errors key,
            // while still others just return a plain string error.
            // See https://gitlab.com/groups/gitlab-org/-/epics/5527#revisit-inconsistent-shape-of-error-responses-in-rest-api
            let error_message = spammable.errors().full_messages();
            fields.insert(
                "message",
                serde_json::json!({
                    "error": self.to_sentence(&error_message)
                }),
            );

            self.render_structured_api_error(fields, status)
        };

        self.with_captcha_check_common(spammable, captcha_render_lambda)
    }
}
