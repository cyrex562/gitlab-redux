use actix_web::{web, HttpResponse, Responder};
use serde::{Deserialize, Serialize};
use serde_json::json;

use crate::controllers::admin::ApplicationController;

/// Controller for managing broadcast messages in the admin area
pub struct BroadcastMessagesController {
    /// The admin application controller
    app_controller: ApplicationController,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct BroadcastMessage {
    id: i32,
    message: String,
    starts_at: Option<String>,
    ends_at: Option<String>,
    color: Option<String>,
    font: Option<String>,
    target_path: Option<String>,
    broadcast_type: String,
    dismissable: bool,
    show_in_cli: bool,
    target_access_levels: Vec<String>,
    created_at: String,
    updated_at: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct BroadcastMessageParams {
    message: String,
    starts_at: Option<String>,
    ends_at: Option<String>,
    color: Option<String>,
    font: Option<String>,
    target_path: Option<String>,
    broadcast_type: String,
    dismissable: Option<bool>,
    show_in_cli: Option<bool>,
    target_access_levels: Option<Vec<String>>,
}

impl BroadcastMessagesController {
    /// Create a new broadcast messages controller
    pub fn new() -> Self {
        Self {
            app_controller: ApplicationController::new(),
        }
    }

    /// Configure the routes for this controller
    pub fn configure(cfg: &mut web::ServiceConfig) {
        cfg.service(
            web::scope("/admin/broadcast_messages")
                .route("", web::get().to(Self::index))
                .route("", web::post().to(Self::create))
                .route("/{id}/edit", web::get().to(Self::edit))
                .route("/{id}", web::put().to(Self::update))
                .route("/{id}", web::delete().to(Self::destroy))
                .route("/preview", web::post().to(Self::preview)),
        );
    }

    /// List all broadcast messages
    pub async fn index() -> impl Responder {
        // TODO: Implement proper broadcast messages fetching
        // This is a placeholder implementation
        let broadcast_messages = vec![BroadcastMessage {
            id: 1,
            message: "Welcome to GitLab!".to_string(),
            starts_at: Some("2023-01-01T00:00:00Z".to_string()),
            ends_at: Some("2023-12-31T23:59:59Z".to_string()),
            color: Some("#000000".to_string()),
            font: Some("#FFFFFF".to_string()),
            target_path: Some("/".to_string()),
            broadcast_type: "banner".to_string(),
            dismissable: true,
            show_in_cli: false,
            target_access_levels: vec!["Guest".to_string()],
            created_at: "2023-01-01T00:00:00Z".to_string(),
            updated_at: "2023-01-01T00:00:00Z".to_string(),
        }];

        HttpResponse::Ok().json(broadcast_messages)
    }

    /// Create a new broadcast message
    pub async fn create(params: web::Json<BroadcastMessageParams>) -> impl Responder {
        // TODO: Implement proper broadcast message creation
        // This is a placeholder implementation
        let success = true;

        if success {
            HttpResponse::Ok().json(json!({
                "notice": "Broadcast Message was successfully created."
            }))
        } else {
            HttpResponse::BadRequest().json(json!({
                "errors": ["Failed to create broadcast message"]
            }))
        }
    }

    /// Show the edit broadcast message form
    pub async fn edit(path: web::Path<i32>) -> impl Responder {
        let id = path.into_inner();

        // TODO: Implement proper broadcast message fetching for edit
        // This is a placeholder implementation
        let broadcast_message = BroadcastMessage {
            id,
            message: "Welcome to GitLab!".to_string(),
            starts_at: Some("2023-01-01T00:00:00Z".to_string()),
            ends_at: Some("2023-12-31T23:59:59Z".to_string()),
            color: Some("#000000".to_string()),
            font: Some("#FFFFFF".to_string()),
            target_path: Some("/".to_string()),
            broadcast_type: "banner".to_string(),
            dismissable: true,
            show_in_cli: false,
            target_access_levels: vec!["Guest".to_string()],
            created_at: "2023-01-01T00:00:00Z".to_string(),
            updated_at: "2023-01-01T00:00:00Z".to_string(),
        };

        HttpResponse::Ok().json(broadcast_message)
    }

    /// Update a broadcast message
    pub async fn update(
        path: web::Path<i32>,
        params: web::Json<BroadcastMessageParams>,
    ) -> impl Responder {
        let id = path.into_inner();

        // TODO: Implement proper broadcast message updating
        // This is a placeholder implementation
        let success = true;

        if success {
            HttpResponse::Ok().json(json!({
                "notice": "Broadcast Message was successfully updated."
            }))
        } else {
            HttpResponse::BadRequest().json(json!({
                "errors": ["Failed to update broadcast message"]
            }))
        }
    }

    /// Delete a broadcast message
    pub async fn destroy(path: web::Path<i32>) -> impl Responder {
        let id = path.into_inner();

        // TODO: Implement proper broadcast message deletion
        // This is a placeholder implementation
        let success = true;

        if success {
            HttpResponse::Ok().finish()
        } else {
            HttpResponse::InternalServerError().finish()
        }
    }

    /// Preview a broadcast message
    pub async fn preview(params: web::Json<BroadcastMessageParams>) -> impl Responder {
        // TODO: Implement proper broadcast message preview
        // This is a placeholder implementation
        let message = params.message.clone();

        HttpResponse::Ok().body(format!("<div class='broadcast-message'>{}</div>", message))
    }
}
