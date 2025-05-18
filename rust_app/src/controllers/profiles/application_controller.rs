// Ported from: orig_app/app/controllers/profiles/application_controller.rb
// Ported on: 2025-05-05
// This file implements the Profiles::ApplicationController from the Ruby codebase.
// Intended layout: "profile"

/// Controller for profile-related base actions.
pub struct ProfilesApplicationController;

impl ProfilesApplicationController {
    /// Returns the layout name for profile pages.
    pub fn layout() -> &'static str {
        "profile"
    }
}
