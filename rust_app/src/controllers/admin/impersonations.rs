use crate::auth::authorize;
use crate::models::user::User;
use actix_web::{get, post, web, HttpResponse, Responder};

#[get("/admin/impersonations")]
pub async fn index() -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("read_impersonation")?;

    // TODO: Implement Impersonation.all equivalent
    let impersonations = Vec::new(); // Placeholder until we implement the model

    HttpResponse::Ok().json(impersonations)
}

#[post("/admin/impersonations")]
pub async fn create(user_data: web::Json<User>) -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("create_impersonation")?;

    // TODO: Implement Impersonation.create equivalent
    let impersonation = None; // Placeholder until we implement the model

    match impersonation {
        Some(i) => HttpResponse::Created().json(i),
        None => HttpResponse::BadRequest().finish(),
    }
}

#[post("/admin/impersonations/{id}/stop")]
pub async fn stop(path: web::Path<i32>) -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("stop_impersonation")?;

    let impersonation_id = path.into_inner();

    // TODO: Implement Impersonation.find and stop equivalent
    let impersonation = None; // Placeholder until we implement the model

    match impersonation {
        Some(i) => {
            // TODO: Implement stop logic
            HttpResponse::Ok().json(i)
        }
        None => HttpResponse::NotFound().finish(),
    }
}
