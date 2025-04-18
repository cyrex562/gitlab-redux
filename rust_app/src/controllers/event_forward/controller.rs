use crate::controllers::base_action::BaseActionController;
use crate::features::Feature;
use crate::tracking::Tracking;
use actix_web::{web, HttpResponse};
use log::info;
use serde_json::Value;

pub struct EventForwardController {
    base: BaseActionController,
}

impl EventForwardController {
    pub fn new(base: BaseActionController) -> Self {
        Self { base }
    }

    pub async fn forward(&self, payload: web::Json<Value>) -> HttpResponse {
        if !Feature::enabled("collect_product_usage_events", "instance") {
            return HttpResponse::NotFound().finish();
        }

        self.process_events(payload).await;
        HttpResponse::Ok().finish()
    }

    async fn process_events(&self, payload: web::Json<Value>) {
        let tracker = Tracking::tracker();

        if let Some(data) = payload.get("data").and_then(|d| d.as_array()) {
            for event in data {
                tracker.emit_event_payload(event).await;
            }

            info!("Enqueued events for forwarding. Count: {}", data.len());
        }
    }
}
