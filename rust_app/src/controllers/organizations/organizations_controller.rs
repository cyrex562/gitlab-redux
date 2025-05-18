// Ported from: orig_app/app/controllers/organizations/organizations_controller.rb
// Ported on: 2025-05-05

use actix_web::{get, web, HttpResponse, Responder};
use serde::Serialize;
use serde_json::json;

use crate::controllers::organizations::application_controller::OrganizationsApplicationController;
use crate::models::event::EventCollection;
use crate::models::organization::Organization;
use crate::models::project::Project;
use crate::models::group::Group;
use crate::services::projects::finder::ProjectsFinder;
use crate::services::organizations::groups_finder::GroupsFinder;
use crate::services::events::render_service::EventsRenderService;

const DEFAULT_RESOURCE_LIMIT: i32 = 1000;
const DEFAULT_ACTIVITY_EVENT_LIMIT: i32 = 20;

/// Controller for managing organizations
pub struct OrganizationsController {
    base: OrganizationsApplicationController,
}

impl OrganizationsController {
    pub fn new(base: OrganizationsApplicationController) -> Self {
        Self { base }
    }

    /// GET /organizations
    pub async fn index(&self) -> impl Responder {
        // Feature flag is checked in base controller
        if let Err(e) = self.base.check_feature_flag() {
            return e;
        }
        HttpResponse::Ok().finish()
    }

    /// GET /organizations/new
    pub async fn new(&self) -> impl Responder {
        if let Err(e) = self.base.check_feature_flag() {
            return e;
        }
        if let Err(e) = self.base.authorize_create_organization() {
            return e;
        }
        HttpResponse::Ok().finish()
    }

    /// GET /organizations/:id
    pub async fn show(&self) -> impl Responder {
        if let Err(e) = self.base.authorize_read_organization(None) {
            return e;
        }
        HttpResponse::Ok().finish()
    }

    /// GET /organizations/:id/activity
    pub async fn activity(&self, query: web::Query<ActivityQuery>) -> impl Responder {
        if let Err(e) = self.base.authorize_read_organization(None) {
            return e;
        }

        let limit = self.activity_query_limit(&query);
        let projects = self.projects();
        let groups = self.groups();

        // Load events with limit + 1 to determine if there's a next page
        let mut events = EventCollection::new(
            projects,
            Some(query.offset.unwrap_or(0)),
            None, // event_filter would go here
            limit + 1,
            Some(groups),
            true,
        ).to_a().await;

        let has_next_page = events.len() > limit as usize;
        if has_next_page {
            events.pop();
        }

        // Filter events visible to current user
        events.retain(|event| event.visible_to_user(self.base.current_user_id));

        // Render events
        let render_service = EventsRenderService::new(self.base.current_user_id);
        render_service.execute(&mut events).await;

        HttpResponse::Ok().json(json!({
            "events": events,
            "has_next_page": has_next_page
        }))
    }

    /// GET /organizations/:id/groups_and_projects
    pub async fn groups_and_projects(&self) -> impl Responder {
        if let Err(e) = self.base.authorize_read_organization(None) {
            return e;
        }
        HttpResponse::Ok().finish()
    }

    /// GET /organizations/:id/users
    pub async fn users(&self) -> impl Responder {
        if let Err(e) = self.base.authorize_read_organization_user(None) {
            return e;
        }
        HttpResponse::Ok().finish()
    }

    // Private helper methods

    fn activity_query_limit(&self, query: &ActivityQuery) -> i32 {
        if let Some(limit) = query.limit {
            if limit <= DEFAULT_ACTIVITY_EVENT_LIMIT {
                return limit;
            }
        }
        DEFAULT_ACTIVITY_EVENT_LIMIT
    }

    fn projects(&self) -> Vec<Project> {
        // TODO: Implement ProjectsFinder
        ProjectsFinder::new(
            params: json!({"organization": self.base.organization(None)}),
            current_user: self.base.current_user_id
        )
        .execute()
        .limit(DEFAULT_RESOURCE_LIMIT)
        .sorted_by_activity()
        .collect()
    }

    fn groups(&self) -> Vec<Group> {
        // TODO: Implement GroupsFinder
        GroupsFinder::new(
            current_user: self.base.current_user_id,
            params: json!({"organization": self.base.organization(None)})
        )
        .execute()
        .limit(DEFAULT_RESOURCE_LIMIT)
        .collect()
    }
}

#[derive(Debug, Serialize)]
struct ActivityQuery {
    limit: Option<i32>,
    offset: Option<i32>,
}