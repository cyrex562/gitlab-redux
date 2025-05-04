// Ported from: orig_app/app/controllers/groups/settings/applications_controller.rb
// Ported on: 2025-05-01
// This file implements the Groups::Settings::ApplicationsController from the Ruby codebase.

use actix_web::{get, post, put, delete, web, HttpResponse, Responder};
use serde::{Serialize, Deserialize};
use crate::models::oauth::Application;

#[derive(Debug, Serialize, Deserialize)]
pub struct ApplicationResponse {
    applications: Vec<Application>,
}

#[get("")]
pub async fn index() -> impl Responder {
    let applications: Vec<Application> = vec![];
    HttpResponse::Ok().json(ApplicationResponse { applications })
}

#[get("/new")]
pub async fn new() -> impl Responder {
    let application: Option<Application> = None;
    HttpResponse::Ok().json(application)
}

#[get("/{id}/edit")]
pub async fn edit(_id: web::Path<i32>) -> impl Responder {
    HttpResponse::Ok().finish()
}

#[post("")]
pub async fn create(_params: web::Json<ApplicationParams>) -> impl Responder {
    HttpResponse::Ok().finish()
}

#[put("/{id}")]
pub async fn update(
    _id: web::Path<i32>,
    _params: web::Json<ApplicationParams>,
) -> impl Responder {
    HttpResponse::Ok().finish()
}

#[post("/{id}/renew")]
pub async fn renew(_id: web::Path<i32>) -> impl Responder {
    HttpResponse::Ok().finish()
}

#[delete("/{id}")]
pub async fn destroy(_id: web::Path<i32>) -> impl Responder {
    HttpResponse::Ok().finish()
}

pub fn configure(cfg: &mut web::ServiceConfig) {
    cfg.service(
        web::scope("/groups/{group_id}/settings/applications")
            .service(index)
            .service(new)
            .service(edit)
            .service(create)
            .service(update)
            .service(renew)
            .service(destroy)
    );
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ApplicationParams {
    pub name: String,
    pub redirect_uri: String,
    pub scopes: Option<String>,
    pub confidential: bool,
}
