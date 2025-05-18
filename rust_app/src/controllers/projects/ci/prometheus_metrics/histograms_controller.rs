// Ported from: orig_app/app/controllers/projects/ci/prometheus_metrics/histograms_controller.rb
// Ported on: 2025-05-07

use actix_web::{post, web, HttpResponse, Responder};
use serde::Deserialize;
use serde_json::json;

#[derive(Debug, Deserialize)]
pub struct Histogram {
    pub name: String,
    pub value: f64,
}

#[derive(Debug, Deserialize)]
pub struct HistogramsParams {
    pub histograms: Vec<Histogram>,
}

#[post("/projects/ci/prometheus_metrics/histograms")]
pub async fn create(params: web::Json<HistogramsParams>) -> impl Responder {
    // TODO: Replace with real service call
    // let result = ObserveHistogramsService::new(project, params.into_inner()).execute();
    let result = json!({
        "success": true,
        "payload": params.histograms
    });
    HttpResponse::Ok().json(result)
}
