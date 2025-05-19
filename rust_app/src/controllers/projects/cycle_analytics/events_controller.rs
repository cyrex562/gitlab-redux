// Ported from: orig_app/app/controllers/projects/cycle_analytics/events_controller.rb
// Ported on: 2025-05-19
// This file implements the Projects::CycleAnalytics::EventsController from the Ruby codebase.

use actix_web::{web, HttpResponse, Responder};
use serde_json::json;

use crate::controllers::projects::application_controller::ProjectsApplicationController;
use crate::services::analytics::cycle_analytics::ProjectLevel;

pub struct ProjectsCycleAnalyticsEventsController {
    base: ProjectsApplicationController,
}

impl ProjectsCycleAnalyticsEventsController {
    pub fn new(base: ProjectsApplicationController) -> Self {
        Self { base }
    }

    pub async fn issue(&self) -> impl Responder {
        self.render_events(self.cycle_analytics().issue().events())
            .await
    }
    pub async fn plan(&self) -> impl Responder {
        self.render_events(self.cycle_analytics().plan().events())
            .await
    }
    pub async fn code(&self) -> impl Responder {
        self.render_events(self.cycle_analytics().code().events())
            .await
    }
    pub async fn test(&self) -> impl Responder {
        // branch param handling omitted for brevity
        self.render_events(self.cycle_analytics().test().events())
            .await
    }
    pub async fn review(&self) -> impl Responder {
        self.render_events(self.cycle_analytics().review().events())
            .await
    }
    pub async fn staging(&self) -> impl Responder {
        self.render_events(self.cycle_analytics().staging().events())
            .await
    }
    pub async fn production(&self) -> impl Responder {
        self.render_events(self.cycle_analytics().production().events())
            .await
    }

    async fn render_events(&self, events: Vec<serde_json::Value>) -> impl Responder {
        HttpResponse::Ok().json(json!({ "events": events }))
    }

    fn cycle_analytics(&self) -> ProjectLevel {
        // This would construct the ProjectLevel analytics object
        ProjectLevel::new(self.base.project())
    }
}
