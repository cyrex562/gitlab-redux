use crate::auth::authorize;
use actix_web::{get, post, web, HttpResponse, Responder};

#[get("/admin/application_settings/appearances")]
pub async fn show() -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("read_appearance")?;

    // TODO: Implement Appearance.show equivalent
    let appearance = None; // Placeholder until we implement the model

    match appearance {
        Some(a) => HttpResponse::Ok().json(a),
        None => HttpResponse::NotFound().finish(),
    }
}

#[post("/admin/application_settings/appearances")]
pub async fn update(appearance_data: web::Json<serde_json::Value>) -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("update_appearance")?;

    // TODO: Implement Appearance.update equivalent
    let appearance = None; // Placeholder until we implement the model

    match appearance {
        Some(a) => HttpResponse::Ok().json(a),
        None => HttpResponse::BadRequest().finish(),
    }
}

#[post("/admin/application_settings/appearances/logo")]
pub async fn update_logo(form: web::Form<serde_json::Value>) -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("update_appearance_logo")?;

    // TODO: Implement Appearance.update_logo equivalent
    let appearance = None; // Placeholder until we implement the model

    match appearance {
        Some(a) => HttpResponse::Ok().json(a),
        None => HttpResponse::BadRequest().finish(),
    }
}

#[post("/admin/application_settings/appearances/header_logo")]
pub async fn update_header_logo(form: web::Form<serde_json::Value>) -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("update_appearance_header_logo")?;

    // TODO: Implement Appearance.update_header_logo equivalent
    let appearance = None; // Placeholder until we implement the model

    match appearance {
        Some(a) => HttpResponse::Ok().json(a),
        None => HttpResponse::BadRequest().finish(),
    }
}

#[post("/admin/application_settings/appearances/favicon")]
pub async fn update_favicon(form: web::Form<serde_json::Value>) -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("update_appearance_favicon")?;

    // TODO: Implement Appearance.update_favicon equivalent
    let appearance = None; // Placeholder until we implement the model

    match appearance {
        Some(a) => HttpResponse::Ok().json(a),
        None => HttpResponse::BadRequest().finish(),
    }
}
