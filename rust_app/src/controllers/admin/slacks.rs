use crate::auth::authorize;
use actix_web::{get, post, web, HttpResponse, Responder};

#[get("/admin/slacks")]
pub async fn index() -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("read_slack")?;

    // TODO: Implement Slack.all equivalent
    let slacks = Vec::new(); // Placeholder until we implement the model

    HttpResponse::Ok().json(slacks)
}

#[get("/admin/slacks/{id}")]
pub async fn show(path: web::Path<i32>) -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("read_slack")?;

    let slack_id = path.into_inner();

    // TODO: Implement Slack.find equivalent
    let slack = None; // Placeholder until we implement the model

    match slack {
        Some(s) => HttpResponse::Ok().json(s),
        None => HttpResponse::NotFound().finish(),
    }
}

#[post("/admin/slacks")]
pub async fn create(slack_data: web::Json<serde_json::Value>) -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("create_slack")?;

    // TODO: Implement Slack.create equivalent
    let slack = None; // Placeholder until we implement the model

    match slack {
        Some(s) => HttpResponse::Created().json(s),
        None => HttpResponse::BadRequest().finish(),
    }
}

#[post("/admin/slacks/{id}")]
pub async fn update(
    path: web::Path<i32>,
    slack_data: web::Json<serde_json::Value>,
) -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("update_slack")?;

    let slack_id = path.into_inner();

    // TODO: Implement Slack.find and update equivalent
    let slack = None; // Placeholder until we implement the model

    match slack {
        Some(s) => HttpResponse::Ok().json(s),
        None => HttpResponse::NotFound().finish(),
    }
}

#[post("/admin/slacks/{id}/test")]
pub async fn test(path: web::Path<i32>) -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("test_slack")?;

    let slack_id = path.into_inner();

    // TODO: Implement Slack.find and test equivalent
    let slack = None; // Placeholder until we implement the model

    match slack {
        Some(s) => {
            // TODO: Implement test logic
            HttpResponse::Ok().json(s)
        }
        None => HttpResponse::NotFound().finish(),
    }
}
