// Ported from: orig_app/app/controllers/concerns/skips_already_signed_in_message.rb
// Ported on: 2025-04-29
// This concern can be included in devise controllers to skip showing an "already signed in" warning on registrations and logins

use crate::auth::devise::DeviseController;
use crate::i18n::I18n;
use actix_web::{web, HttpRequest, HttpResponse};
use std::sync::Arc;

/// This trait can be implemented by devise controllers to skip showing an "already signed in"
/// warning on registrations and logins
pub trait SkipsAlreadySignedInMessage {
    /// Skip the authentication check for new and create actions
    fn skip_authentication_for_registration(&self, req: &HttpRequest) -> bool;

    /// Require no authentication without showing a flash message
    fn require_no_authentication_without_flash(
        &self,
        req: &HttpRequest,
    ) -> Result<(), HttpResponse>;
}

pub struct SkipsAlreadySignedInMessageHandler {
    i18n: Arc<I18n>,
}

impl SkipsAlreadySignedInMessageHandler {
    pub fn new(i18n: Arc<I18n>) -> Self {
        SkipsAlreadySignedInMessageHandler { i18n }
    }
}

impl SkipsAlreadySignedInMessage for SkipsAlreadySignedInMessageHandler {
    fn skip_authentication_for_registration(&self, req: &HttpRequest) -> bool {
        // Check if the current path is for new or create actions
        let path = req.path();
        path.ends_with("/new") || path.ends_with("/create")
    }

    fn require_no_authentication_without_flash(
        &self,
        req: &HttpRequest,
    ) -> Result<(), HttpResponse> {
        // First, require no authentication (this would normally set a flash message)
        let result = DeviseController::require_no_authentication(req);

        // If the flash message is "already authenticated", clear it
        if let Ok(flash) = req.extensions().get::<Flash>() {
            if let Some(alert) = flash.get("alert") {
                if alert == self.i18n.t("devise.failure.already_authenticated") {
                    flash.remove("alert");
                }
            }
        }

        result
    }
}

// This would be implemented in a separate module
pub mod auth {
    pub mod devise {
        use actix_web::HttpRequest;

        pub struct DeviseController;

        impl DeviseController {
            pub fn require_no_authentication(req: &HttpRequest) -> Result<(), actix_web::Error> {
                // In a real implementation, this would check if the user is already authenticated
                // and return an error if they are
                // For now, we'll just return Ok
                Ok(())
            }
        }
    }
}

// This would be implemented in a separate module
pub mod i18n {
    use std::sync::Arc;

    pub struct I18n;

    impl I18n {
        pub fn t(&self, key: &str) -> String {
            // In a real implementation, this would translate the key
            // For now, we'll just return the key
            key.to_string()
        }
    }
}

// This would be implemented in a separate module
pub struct Flash {
    messages: std::collections::HashMap<String, String>,
}

impl Flash {
    pub fn get(&self, key: &str) -> Option<&String> {
        self.messages.get(key)
    }

    pub fn remove(&mut self, key: &str) {
        self.messages.remove(key);
    }
}
