// Ported from: orig_app/app/controllers/groups/custom_emoji_controller.rb
// Ported on: 2025-05-01
// This file implements the Groups::CustomEmojiController from the Ruby codebase.

use crate::controllers::groups::application_controller::GroupsApplicationController;
use actix_web::{web, HttpResponse, Responder};

/// Controller for custom emoji in groups (stub).
pub struct CustomEmojiController {
    base: GroupsApplicationController,
}

impl CustomEmojiController {
    pub fn new(base: GroupsApplicationController) -> Self {
        Self { base }
    }
    // No actions defined in the Ruby source.
}
