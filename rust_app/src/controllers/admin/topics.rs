use crate::auth::authorize;
use actix_web::{get, post, web, HttpResponse, Responder};

#[get("/admin/topics")]
pub async fn index() -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("read_topic")?;

    // TODO: Implement Topic.all equivalent
    let topics = Vec::new(); // Placeholder until we implement the model

    HttpResponse::Ok().json(topics)
}

#[get("/admin/topics/{id}")]
pub async fn show(path: web::Path<i32>) -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("read_topic")?;

    let topic_id = path.into_inner();

    // TODO: Implement Topic.find equivalent
    let topic = None; // Placeholder until we implement the model

    match topic {
        Some(t) => HttpResponse::Ok().json(t),
        None => HttpResponse::NotFound().finish(),
    }
}

#[post("/admin/topics")]
pub async fn create(topic_data: web::Json<serde_json::Value>) -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("create_topic")?;

    // TODO: Implement Topic.create equivalent
    let topic = None; // Placeholder until we implement the model

    match topic {
        Some(t) => HttpResponse::Created().json(t),
        None => HttpResponse::BadRequest().finish(),
    }
}

#[post("/admin/topics/{id}")]
pub async fn update(
    path: web::Path<i32>,
    topic_data: web::Json<serde_json::Value>,
) -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("update_topic")?;

    let topic_id = path.into_inner();

    // TODO: Implement Topic.find and update equivalent
    let topic = None; // Placeholder until we implement the model

    match topic {
        Some(t) => HttpResponse::Ok().json(t),
        None => HttpResponse::NotFound().finish(),
    }
}
