use crate::auth::authorize;
use actix_web::{get, post, web, HttpResponse, Responder};

#[get("/admin/runners")]
pub async fn index() -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("read_runner")?;

    // TODO: Implement Runner.all equivalent
    let runners = Vec::new(); // Placeholder until we implement the model

    HttpResponse::Ok().json(runners)
}

#[get("/admin/runners/{id}")]
pub async fn show(path: web::Path<i32>) -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("read_runner")?;

    let runner_id = path.into_inner();

    // TODO: Implement Runner.find equivalent
    let runner = None; // Placeholder until we implement the model

    match runner {
        Some(r) => HttpResponse::Ok().json(r),
        None => HttpResponse::NotFound().finish(),
    }
}

#[post("/admin/runners")]
pub async fn create(runner_data: web::Json<serde_json::Value>) -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("create_runner")?;

    // TODO: Implement Runner.create equivalent
    let runner = None; // Placeholder until we implement the model

    match runner {
        Some(r) => HttpResponse::Created().json(r),
        None => HttpResponse::BadRequest().finish(),
    }
}

#[post("/admin/runners/{id}")]
pub async fn update(
    path: web::Path<i32>,
    runner_data: web::Json<serde_json::Value>,
) -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("update_runner")?;

    let runner_id = path.into_inner();

    // TODO: Implement Runner.find and update equivalent
    let runner = None; // Placeholder until we implement the model

    match runner {
        Some(r) => HttpResponse::Ok().json(r),
        None => HttpResponse::NotFound().finish(),
    }
}

#[post("/admin/runners/{id}/pause")]
pub async fn pause(path: web::Path<i32>) -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("pause_runner")?;

    let runner_id = path.into_inner();

    // TODO: Implement Runner.find and pause equivalent
    let runner = None; // Placeholder until we implement the model

    match runner {
        Some(r) => {
            // TODO: Implement pause logic
            HttpResponse::Ok().json(r)
        }
        None => HttpResponse::NotFound().finish(),
    }
}

#[post("/admin/runners/{id}/resume")]
pub async fn resume(path: web::Path<i32>) -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("resume_runner")?;

    let runner_id = path.into_inner();

    // TODO: Implement Runner.find and resume equivalent
    let runner = None; // Placeholder until we implement the model

    match runner {
        Some(r) => {
            // TODO: Implement resume logic
            HttpResponse::Ok().json(r)
        }
        None => HttpResponse::NotFound().finish(),
    }
}
