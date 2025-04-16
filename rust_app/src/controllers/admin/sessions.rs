use crate::auth::authorize;
use actix_web::{get, post, web, HttpResponse, Responder};

#[get("/admin/sessions")]
pub async fn index() -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("read_session")?;

    // TODO: Implement Session.all equivalent
    let sessions = Vec::new(); // Placeholder until we implement the model

    HttpResponse::Ok().json(sessions)
}

#[get("/admin/sessions/{id}")]
pub async fn show(path: web::Path<i32>) -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("read_session")?;

    let session_id = path.into_inner();

    // TODO: Implement Session.find equivalent
    let session = None; // Placeholder until we implement the model

    match session {
        Some(s) => HttpResponse::Ok().json(s),
        None => HttpResponse::NotFound().finish(),
    }
}

#[post("/admin/sessions/{id}/destroy")]
pub async fn destroy(path: web::Path<i32>) -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("destroy_session")?;

    let session_id = path.into_inner();

    // TODO: Implement Session.find and destroy equivalent
    let session = None; // Placeholder until we implement the model

    match session {
        Some(s) => {
            // TODO: Implement destroy logic
            HttpResponse::Ok().json(s)
        }
        None => HttpResponse::NotFound().finish(),
    }
}
