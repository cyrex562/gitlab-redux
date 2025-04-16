use actix_web::{get, post, web, HttpResponse, Responder};
use crate::auth::authorize;

#[get("/admin/labels")]
pub async fn index() -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("read_label")?;

    // TODO: Implement Label.all equivalent
    let labels = Vec::new(); // Placeholder until we implement the model

    HttpResponse::Ok().json(labels)
}

#[get("/admin/labels/{id}")]
pub async fn show(path: web::Path<i32>) -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("read_label")?;

    let label_id = path.into_inner();
    
    // TODO: Implement Label.find equivalent
    let label = None; // Placeholder until we implement the model

    match label {
        Some(l) => HttpResponse::Ok().json(l),
        None => HttpResponse::NotFound().finish()
    }
}

#[post("/admin/labels")]
pub async fn create(label_data: web::Json<serde_json::Value>) -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("create_label")?;

    // TODO: Implement Label.create equivalent
    let label = None; // Placeholder until we implement the model

    match label {
        Some(l) => HttpResponse::Created().json(l),
        None => HttpResponse::BadRequest().finish()
    }
}

#[post("/admin/labels/{id}")]
pub async fn update(path: web::Path<i32>, label_data: web::Json<serde_json::Value>) -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("update_label")?;

    let label_id = path.into_inner();
    
    // TODO: Implement Label.find and update equivalent
    let label = None; // Placeholder until we implement the model

    match label {
        Some(l) => HttpResponse::Ok().json(l),
        None => HttpResponse::NotFound().finish()
    }
} 