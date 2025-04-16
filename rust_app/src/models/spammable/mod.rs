pub mod akismet_mark_as_spam_action;
pub mod captcha_check;

use crate::models::user::User;
use actix_web::{web, HttpResponse};
use serde::{Deserialize, Serialize};

// TODO: Implement Spammable trait
pub trait Spammable {
    fn render_recaptcha(&self) -> bool;
    fn spammable_entity_type(&self) -> String;
    fn errors(&self) -> &Vec<String>;
}

// TODO: Implement SpamLog model
pub struct SpamLog {
    pub id: i32,
    // Add other necessary fields
}

// TODO: Implement SpamParams model
pub struct SpamParams {
    pub captcha_response: Option<String>,
    pub spam_log_id: Option<i32>,
    // Add other necessary fields
}

// TODO: Implement SpamActionResponseFields trait
pub trait SpamActionResponseFields {
    fn spam_action_response_fields(&self, spammable: &dyn Spammable) -> serde_json::Value;
}
