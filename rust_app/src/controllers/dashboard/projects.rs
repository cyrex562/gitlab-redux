// Ported from orig_app/app/controllers/dashboard/projects_controller.rb
use crate::controllers::concerns::{
    FiltersEvents, ParamsBackwardCompatibility, RendersMemberAccess, RendersProjectsList,
    SortingHelper, SortingPreference,
};
use crate::controllers::dashboard::DashboardApplicationController;
use crate::models::project::Project;
use crate::services::{events::EventCollection, projects::ProjectsFinder};
use actix_web::{web, HttpResponse};

pub struct ProjectsController {
    base: DashboardApplicationController,
}

impl ProjectsController {
    pub fn new(base: DashboardApplicationController) -> Self {
        Self { base }
    }

    pub async fn index(&self, query: web::Query<ProjectQuery>) -> HttpResponse {
        if query.personal == Some(true) {
            return HttpResponse::Found()
                .header("Location", "/dashboard/projects/personal")
                .finish();
        }

        if query.archived == Some("only".to_string()) {
            return HttpResponse::Found()
                .header("Location", "/dashboard/projects/inactive")
                .finish();
        }

        if query.format == Some("atom".to_string()) {
            let events = self.load_events(&query).await;
            // TODO: Implement atom feed rendering
            return HttpResponse::Ok().json(events);
        }

        // HTML format
        // TODO: Implement HTML template rendering
        HttpResponse::Ok().finish()
    }

    pub async fn starred(&self) -> HttpResponse {
        // TODO: Implement starred projects view
        HttpResponse::Ok().finish()
    }

    async fn load_events(&self, query: &ProjectQuery) -> Vec<Event> {
        let projects = ProjectsFinder::new(
            query.clone(),
            true, // non_public
            true, // not_aimed_for_deletion
            self.base.user.clone(),
        )
        .execute()
        .await;

        let events =
            EventCollection::new(projects, query.offset.unwrap_or(0), self.event_filter()).to_vec();

        // TODO: Implement event rendering service
        events
    }

    fn default_sort_order(&self) -> String {
        self.sort_value_name()
    }

    fn sorting_field(&self) -> String {
        Project::SORTING_PREFERENCE_FIELD.to_string()
    }

    fn set_sorting(&mut self, query: &mut ProjectQuery) {
        query.sort = Some(self.set_sort_order());
    }
}

#[derive(serde::Deserialize, Clone)]
pub struct ProjectQuery {
    personal: Option<bool>,
    archived: Option<String>,
    format: Option<String>,
    offset: Option<i32>,
    sort: Option<String>,
}

impl ParamsBackwardCompatibility for ProjectQuery {}
impl RendersMemberAccess for ProjectsController {}
impl RendersProjectsList for ProjectsController {}
impl SortingHelper for ProjectsController {}
impl SortingPreference for ProjectsController {}
impl FiltersEvents for ProjectsController {}
