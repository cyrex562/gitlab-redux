use actix_web::{web, HttpResponse, Responder};
use serde::{Deserialize, Serialize};
use std::sync::Arc;

use crate::models::issuable::Issuable;
use crate::models::todo::Todo;
use crate::services::todo_service::TodoService;
use crate::services::todos_finder::TodosFinder;

/// Module for handling todos actions
pub trait TodosActions {
    /// Get the current user ID
    fn current_user_id(&self) -> i32;

    /// Get the issuable
    fn issuable(&self) -> Arc<dyn Issuable>;

    /// Create a todo
    async fn create(&self) -> impl Responder {
        let todo_service = TodoService::new();
        let issuable = self.issuable();
        let current_user_id = self.current_user_id();

        match todo_service.mark_todo(issuable, current_user_id).await {
            Ok(todo) => {
                let todos_finder = TodosFinder::new(current_user_id, "pending");
                let count = todos_finder.execute().await.len();

                HttpResponse::Ok().json(serde_json::json!({
                    "count": count,
                    "delete_path": format!("/dashboard/todos/{}", todo.id())
                }))
            }
            Err(_) => HttpResponse::InternalServerError().json("Failed to create todo"),
        }
    }
}
