// Ported from: orig_app/app/controllers/profiles/chat_names_controller.rb
// Ported on: 2025-05-05
// This file implements the Profiles::ChatNamesController from the Ruby codebase.
//
// See porting_log.txt for details.

use actix_web::{delete, get, post, web, HttpResponse, Responder};

pub struct ChatNamesController;

impl ChatNamesController {
    /// GET /profiles/chat_names
    #[get("/profiles/chat_names")]
    pub async fn index() -> impl Responder {
        // TODO: Fetch chat names for current user
        HttpResponse::Ok().body("List of chat names")
    }

    /// GET /profiles/chat_names/new
    #[get("/profiles/chat_names/new")]
    pub async fn new() -> impl Responder {
        // TODO: Show new chat name form
        HttpResponse::Ok().body("New chat name form")
    }

    /// POST /profiles/chat_names
    #[post("/profiles/chat_names")]
    pub async fn create() -> impl Responder {
        // TODO: Create new chat name for current user
        HttpResponse::Ok().body("Chat name created")
    }

    /// DELETE /profiles/chat_names/deny
    #[delete("/profiles/chat_names/deny")]
    pub async fn deny() -> impl Responder {
        // TODO: Deny chat name authorization
        HttpResponse::Ok().body("Chat name authorization denied")
    }

    /// DELETE /profiles/chat_names/{id}
    #[delete("/profiles/chat_names/{id}")]
    pub async fn destroy(path: web::Path<i32>) -> impl Responder {
        // TODO: Delete chat name by id
        HttpResponse::Ok().body(format!("Deleted chat name with id {}", path.into_inner()))
    }
}
