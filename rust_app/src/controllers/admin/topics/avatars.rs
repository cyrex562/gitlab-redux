use actix_web::{get, post, web, HttpResponse, Responder};
use crate::auth::authorize;

#[get("/admin/topics/{id}/avatar")]
pub async fn show(path: web::Path<i32>) -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("read_topic_avatar")?;

    let topic_id = path.into_inner();
    
    // TODO: Implement Topic::Avatar.find equivalent
    let avatar = None; // Placeholder until we implement the model

    match avatar {
        Some(a) => HttpResponse::Ok().json(a),
        None => HttpResponse::NotFound().finish()
    }
}

#[post("/admin/topics/{id}/avatar")]
pub async fn create(path: web::Path<i32>, form: web::Form<serde_json::Value>) -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("create_topic_avatar")?;

    let topic_id = path.into_inner();
    
    // TODO: Implement Topic::Avatar.create equivalent
    let avatar = None; // Placeholder until we implement the model

    match avatar {
        Some(a) => HttpResponse::Created().json(a),
        None => HttpResponse::BadRequest().finish()
    }
}

#[post("/admin/topics/{id}/avatar/remove")]
pub async fn remove(path: web::Path<i32>) -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("remove_topic_avatar")?;

    let topic_id = path.into_inner();
    
    // TODO: Implement Topic::Avatar.find and remove equivalent
    let avatar = None; // Placeholder until we implement the model

    match avatar {
        Some(a) => {
            // TODO: Implement remove logic
            HttpResponse::Ok().json(a)
        },
        None => HttpResponse::NotFound().finish()
    }
} 