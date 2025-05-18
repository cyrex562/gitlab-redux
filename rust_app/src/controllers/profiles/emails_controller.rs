// Ported from: orig_app/app/controllers/profiles/emails_controller.rb
// Ported on: 2025-05-05
// Handles user profile email management (list, add, remove, resend confirmation)

use actix_web::{delete, get, post, web, HttpResponse, Responder};
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct EmailParams {
    pub email: String,
}

#[get("/profile/emails")]
pub async fn index() -> impl Responder {
    // TODO: Fetch current user, list primary and secondary emails
    HttpResponse::Ok().body("Emails index (stub)")
}

#[post("/profile/emails")]
pub async fn create(params: web::Json<EmailParams>) -> impl Responder {
    // TODO: Add new email for current user, handle errors and flash messages
    let _email = &params.email;
    HttpResponse::Found()
        .header("Location", "/profile/emails")
        .finish()
}

#[delete("/profile/emails/{id}")]
pub async fn destroy(path: web::Path<i64>) -> impl Responder {
    // TODO: Remove email with id for current user
    let _email_id = path.into_inner();
    HttpResponse::Found()
        .header("Location", "/profile/emails")
        .finish()
}

#[post("/profile/emails/{id}/resend_confirmation")]
pub async fn resend_confirmation_instructions(path: web::Path<i64>) -> impl Responder {
    // TODO: Resend confirmation instructions for email with id
    let _email_id = path.into_inner();
    HttpResponse::Found()
        .header("Location", "/profile/emails")
        .finish()
}

// TODO: Integrate with user session, flash messages, and email services as needed.
