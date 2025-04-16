use crate::auth::authorize;
use actix_web::{get, post, web, HttpResponse, Responder};

#[get("/admin/integrations")]
pub async fn index() -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("read_integration")?;

    // TODO: Implement Integration.all equivalent
    let integrations = Vec::new(); // Placeholder until we implement the model

    HttpResponse::Ok().json(integrations)
}

#[get("/admin/integrations/{id}")]
pub async fn show(path: web::Path<i32>) -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("read_integration")?;

    let integration_id = path.into_inner();

    // TODO: Implement Integration.find equivalent
    let integration = None; // Placeholder until we implement the model

    match integration {
        Some(i) => HttpResponse::Ok().json(i),
        None => HttpResponse::NotFound().finish(),
    }
}

#[post("/admin/integrations")]
pub async fn create(integration_data: web::Json<serde_json::Value>) -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("create_integration")?;

    // TODO: Implement Integration.create equivalent
    let integration = None; // Placeholder until we implement the model

    match integration {
        Some(i) => HttpResponse::Created().json(i),
        None => HttpResponse::BadRequest().finish(),
    }
}

#[post("/admin/integrations/{id}")]
pub async fn update(
    path: web::Path<i32>,
    integration_data: web::Json<serde_json::Value>,
) -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("update_integration")?;

    let integration_id = path.into_inner();

    // TODO: Implement Integration.find and update equivalent
    let integration = None; // Placeholder until we implement the model

    match integration {
        Some(i) => HttpResponse::Ok().json(i),
        None => HttpResponse::NotFound().finish(),
    }
}
