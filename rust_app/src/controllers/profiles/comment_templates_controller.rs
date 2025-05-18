// Ported from: orig_app/app/controllers/profiles/comment_templates_controller.rb
// Ported on: 2025-05-05
// This file implements the Profiles::CommentTemplatesController from the Ruby codebase.

use actix_web::{get, HttpResponse, Responder};

/// Controller for profile comment template actions, ported from Rails controller logic.
pub struct ProfilesCommentTemplatesController;

impl ProfilesCommentTemplatesController {
    /// Handler for showing comment templates settings.
    #[get("/profile/comment_templates")]
    pub async fn show() -> impl Responder {
        // In the Ruby code, @hide_search_settings is set to true before action.
        // Here, you would set this in the context/session if needed.
        // Feature category: user_profile (could be used for logging/metrics)
        HttpResponse::Ok().body("Comment Templates page (stub)")
    }
}
