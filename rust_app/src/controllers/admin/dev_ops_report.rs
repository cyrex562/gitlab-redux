use crate::auth::authorize;
use actix_web::{get, web, HttpResponse, Responder};

#[get("/admin/dev_ops_report")]
pub async fn show() -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("read_dev_ops_report")?;

    // TODO: Implement DevOpsReport::Metric.order(:created_at).last equivalent
    let metric = None; // Placeholder until we implement the model

    match metric {
        Some(m) => HttpResponse::Ok().json(m),
        None => HttpResponse::NotFound().finish(),
    }
}

// Helper functions that would be used in the controller
pub fn show_adoption() -> bool {
    // TODO: Implement show_adoption? equivalent
    false
}

pub fn should_track_devops_score() -> bool {
    // TODO: Implement should_track_devops_score? equivalent
    true
}

pub fn tracking_namespace_source() -> Option<String> {
    // TODO: Implement tracking_namespace_source equivalent
    None
}

pub fn tracking_project_source() -> Option<String> {
    // TODO: Implement tracking_project_source equivalent
    None
}
