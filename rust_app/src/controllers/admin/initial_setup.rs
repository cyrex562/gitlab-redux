use crate::auth::authorize;
use crate::models::user::User;
use actix_web::{get, post, web, HttpResponse, Responder};

#[get("/admin/initial_setup")]
pub async fn show() -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("read_initial_setup")?;

    // TODO: Implement InitialSetup.show equivalent
    let setup_data = None; // Placeholder until we implement the model

    match setup_data {
        Some(data) => HttpResponse::Ok().json(data),
        None => HttpResponse::NotFound().finish(),
    }
}

#[post("/admin/initial_setup")]
pub async fn create(user_data: web::Json<User>) -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("create_initial_setup")?;

    // TODO: Implement InitialSetup.create equivalent
    let setup = None; // Placeholder until we implement the model

    match setup {
        Some(s) => HttpResponse::Created().json(s),
        None => HttpResponse::BadRequest().finish(),
    }
}
