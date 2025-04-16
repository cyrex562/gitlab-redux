use crate::auth::authorize;
use actix_web::{get, post, web, HttpResponse, Responder};

#[get("/admin/jobs")]
pub async fn index() -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("read_job")?;

    // TODO: Implement Job.all equivalent
    let jobs = Vec::new(); // Placeholder until we implement the model

    HttpResponse::Ok().json(jobs)
}

#[get("/admin/jobs/{id}")]
pub async fn show(path: web::Path<i32>) -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("read_job")?;

    let job_id = path.into_inner();

    // TODO: Implement Job.find equivalent
    let job = None; // Placeholder until we implement the model

    match job {
        Some(j) => HttpResponse::Ok().json(j),
        None => HttpResponse::NotFound().finish(),
    }
}

#[post("/admin/jobs/{id}/cancel")]
pub async fn cancel(path: web::Path<i32>) -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("cancel_job")?;

    let job_id = path.into_inner();

    // TODO: Implement Job.find and cancel equivalent
    let job = None; // Placeholder until we implement the model

    match job {
        Some(j) => {
            // TODO: Implement cancel logic
            HttpResponse::Ok().json(j)
        }
        None => HttpResponse::NotFound().finish(),
    }
}

#[post("/admin/jobs/{id}/retry")]
pub async fn retry(path: web::Path<i32>) -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("retry_job")?;

    let job_id = path.into_inner();

    // TODO: Implement Job.find and retry equivalent
    let job = None; // Placeholder until we implement the model

    match job {
        Some(j) => {
            // TODO: Implement retry logic
            HttpResponse::Ok().json(j)
        }
        None => HttpResponse::NotFound().finish(),
    }
}
