use crate::auth::authorize;
use actix_web::{get, post, web, HttpResponse, Responder};

#[get("/admin/instance_review")]
pub async fn show() -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("read_instance_review")?;

    // TODO: Implement InstanceReview.show equivalent
    let review_data = None; // Placeholder until we implement the model

    match review_data {
        Some(data) => HttpResponse::Ok().json(data),
        None => HttpResponse::NotFound().finish(),
    }
}

#[post("/admin/instance_review")]
pub async fn create(review_data: web::Json<serde_json::Value>) -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("create_instance_review")?;

    // TODO: Implement InstanceReview.create equivalent
    let review = None; // Placeholder until we implement the model

    match review {
        Some(r) => HttpResponse::Created().json(r),
        None => HttpResponse::BadRequest().finish(),
    }
}
