use actix_web::{web, HttpResponse, Responder};
use serde::{Deserialize, Serialize};
use std::sync::Arc;

use crate::models::awardable::Awardable;
use crate::services::award_emojis::toggle_service::ToggleService;

/// Module for handling toggle award emoji
pub trait ToggleAwardEmoji {
    /// Get the current user ID
    fn current_user_id(&self) -> i32;

    /// Get the awardable
    fn awardable(&self) -> Arc<dyn Awardable>;

    /// Get the emoji name from request parameters
    fn emoji_name(&self) -> String;

    /// Toggle award emoji
    async fn toggle_award_emoji(&self) -> impl Responder {
        // In a real implementation, this would check if the user is authenticated
        // authenticate_user!();

        let name = self.emoji_name();
        let awardable = self.awardable();
        let current_user_id = self.current_user_id();

        let service = ToggleService::new(awardable, name, current_user_id);
        match service.execute().await {
            Ok(result) => {
                if result.status == "success" {
                    HttpResponse::Ok().json(serde_json::json!({ "ok": true }))
                } else {
                    HttpResponse::Ok().json(serde_json::json!({ "ok": false }))
                }
            }
            Err(_) => HttpResponse::InternalServerError().json("Failed to toggle award emoji"),
        }
    }
}
