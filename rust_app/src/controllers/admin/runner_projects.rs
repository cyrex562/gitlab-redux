use crate::auth::authorize;
use actix_web::{get, post, web, HttpResponse, Responder};

#[get("/admin/runner_projects")]
pub async fn index() -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("read_runner_project")?;

    // TODO: Implement RunnerProject.all equivalent
    let runner_projects = Vec::new(); // Placeholder until we implement the model

    HttpResponse::Ok().json(runner_projects)
}

#[get("/admin/runner_projects/{id}")]
pub async fn show(path: web::Path<i32>) -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("read_runner_project")?;

    let runner_project_id = path.into_inner();

    // TODO: Implement RunnerProject.find equivalent
    let runner_project = None; // Placeholder until we implement the model

    match runner_project {
        Some(r) => HttpResponse::Ok().json(r),
        None => HttpResponse::NotFound().finish(),
    }
}

#[post("/admin/runner_projects")]
pub async fn create(runner_project_data: web::Json<serde_json::Value>) -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("create_runner_project")?;

    // TODO: Implement RunnerProject.create equivalent
    let runner_project = None; // Placeholder until we implement the model

    match runner_project {
        Some(r) => HttpResponse::Created().json(r),
        None => HttpResponse::BadRequest().finish(),
    }
}

#[post("/admin/runner_projects/{id}")]
pub async fn update(
    path: web::Path<i32>,
    runner_project_data: web::Json<serde_json::Value>,
) -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("update_runner_project")?;

    let runner_project_id = path.into_inner();

    // TODO: Implement RunnerProject.find and update equivalent
    let runner_project = None; // Placeholder until we implement the model

    match runner_project {
        Some(r) => HttpResponse::Ok().json(r),
        None => HttpResponse::NotFound().finish(),
    }
}
