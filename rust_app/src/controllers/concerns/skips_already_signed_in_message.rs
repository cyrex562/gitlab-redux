use crate::config::settings::Settings;
use crate::models::user::User;
use actix_web::{web, HttpResponse};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

/// Module for handling already signed in message skipping
pub trait SkipsAlreadySignedInMessage {
    /// Get the current user ID
    fn user_id(&self) -> Option<i32>;

    /// Get the message ID
    fn message_id(&self) -> Option<i32>;

    /// Get the message type
    fn message_type(&self) -> Option<String>;

    /// Check if the user is already signed in
    fn is_user_signed_in(&self) -> bool {
        self.user_id().is_some()
    }

    /// Check if the message should be skipped
    fn should_skip_message(&self) -> bool {
        if !self.is_user_signed_in() {
            return false;
        }

        match self.message_type() {
            Some(msg_type) => matches!(msg_type.as_str(), "sign_in" | "registration"),
            None => false,
        }
    }

    /// Get message skip status
    fn get_message_skip_status(&self) -> HashMap<String, bool> {
        let mut status = HashMap::new();

        status.insert("user_signed_in".to_string(), self.is_user_signed_in());
        status.insert(
            "should_skip_message".to_string(),
            self.should_skip_message(),
        );

        status
    }

    /// Handle message skipping
    fn handle_message_skip(&self) -> Result<(), HttpResponse> {
        if self.should_skip_message() {
            return Err(HttpResponse::Ok().json(serde_json::json!({
                "message": "Message skipped - user already signed in"
            })));
        }
        Ok(())
    }
}
