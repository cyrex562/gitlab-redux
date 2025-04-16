use crate::auth::authorize;
use actix_web::{get, post, web, HttpResponse, Responder};

#[get("/admin/plan_limits")]
pub async fn index() -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("read_plan_limit")?;

    // TODO: Implement PlanLimit.all equivalent
    let plan_limits = Vec::new(); // Placeholder until we implement the model

    HttpResponse::Ok().json(plan_limits)
}

#[get("/admin/plan_limits/{id}")]
pub async fn show(path: web::Path<i32>) -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("read_plan_limit")?;

    let plan_limit_id = path.into_inner();

    // TODO: Implement PlanLimit.find equivalent
    let plan_limit = None; // Placeholder until we implement the model

    match plan_limit {
        Some(p) => HttpResponse::Ok().json(p),
        None => HttpResponse::NotFound().finish(),
    }
}

#[post("/admin/plan_limits")]
pub async fn create(plan_limit_data: web::Json<serde_json::Value>) -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("create_plan_limit")?;

    // TODO: Implement PlanLimit.create equivalent
    let plan_limit = None; // Placeholder until we implement the model

    match plan_limit {
        Some(p) => HttpResponse::Created().json(p),
        None => HttpResponse::BadRequest().finish(),
    }
}

#[post("/admin/plan_limits/{id}")]
pub async fn update(
    path: web::Path<i32>,
    plan_limit_data: web::Json<serde_json::Value>,
) -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("update_plan_limit")?;

    let plan_limit_id = path.into_inner();

    // TODO: Implement PlanLimit.find and update equivalent
    let plan_limit = None; // Placeholder until we implement the model

    match plan_limit {
        Some(p) => HttpResponse::Ok().json(p),
        None => HttpResponse::NotFound().finish(),
    }
}
