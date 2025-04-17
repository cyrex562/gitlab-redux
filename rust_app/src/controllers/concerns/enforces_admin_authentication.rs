use actix_web::{web, HttpRequest, HttpResponse};
use std::sync::Arc;

use crate::config::settings::Settings;
use crate::models::user::User;
use crate::utils::session::Session;

/// Module for enforcing admin authentication
pub trait EnforcesAdminAuthentication {
    /// Authenticate that the current user is an admin
    fn authenticate_admin(&self, req: &HttpRequest) -> HttpResponse {
        // Check if user is logged in and is an admin
        if let Some(current_user) = self.current_user() {
            if !current_user.is_admin() {
                return HttpResponse::NotFound().finish();
            }

            // Check if admin mode is enabled in settings
            if self.settings().admin_mode_enabled() {
                // Check if user is in admin mode
                if !self.current_user_mode().is_admin_mode() {
                    // Request admin mode
                    self.current_user_mode().request_admin_mode();

                    // Store the current location for redirect after admin mode is enabled
                    if self.storable_location(req) {
                        self.store_location_for_redirect(req.uri().path());
                    }

                    // Redirect to admin session page
                    return HttpResponse::Found()
                        .header("Location", "/admin/session/new")
                        .header("X-Flash-Notice", "Re-authentication required")
                        .finish();
                }
            }

            return HttpResponse::Ok().finish();
        }

        HttpResponse::NotFound().finish()
    }

    /// Check if the current location can be stored for redirect
    fn storable_location(&self, req: &HttpRequest) -> bool {
        req.path() != "/admin/session/new"
    }

    /// Store the current location for redirect
    fn store_location_for_redirect(&self, path: &str) {
        self.session().insert("redirect", path);
    }

    /// Authorize a specific ability
    fn authorize_ability(&self, ability: &str) -> HttpResponse {
        if let Some(current_user) = self.current_user() {
            if current_user.is_admin() {
                return self.authenticate_admin(self.request());
            }

            if !current_user.can(ability) {
                return HttpResponse::NotFound().finish();
            }
        }

        HttpResponse::Ok().finish()
    }

    // Required trait methods that need to be implemented by the controller
    fn current_user(&self) -> Option<Arc<User>>;
    fn current_user_mode(&self) -> &dyn UserMode;
    fn settings(&self) -> Arc<Settings>;
    fn session(&self) -> &Session;
    fn request(&self) -> &HttpRequest;
}

/// Trait for user mode functionality
pub trait UserMode {
    /// Check if the user is in admin mode
    fn is_admin_mode(&self) -> bool;

    /// Request admin mode
    fn request_admin_mode(&self);
}
