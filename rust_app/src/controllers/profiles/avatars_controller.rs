// Ported from: orig_app/app/controllers/profiles/avatars_controller.rb
// Ported on: 2025-05-05
// This file implements the Profiles::AvatarsController from the Ruby codebase.

use actix_web::{post, web, HttpResponse, Responder};

/// Controller for profile avatar actions, ported from Rails controller logic.
pub struct ProfilesAvatarsController;

impl ProfilesAvatarsController {
    /// Destroys the user's avatar and redirects to the user settings profile page.
    #[post("/profile/avatar/remove")]
    pub async fn destroy() -> impl Responder {
        // TODO: Fetch current user from session/context
        // TODO: Call Users::UpdateService equivalent to remove avatar
        // Example placeholder logic:
        // Remove avatar for current user
        // Save user changes
        // Redirect to user settings profile page
        HttpResponse::Found()
            .append_header(("Location", "/-/profile"))
            .finish()
    }
}
