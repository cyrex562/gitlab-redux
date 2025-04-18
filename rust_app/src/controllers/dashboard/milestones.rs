use crate::controllers::dashboard::DashboardApplicationController;
use crate::models::milestone::Milestone;
use crate::services::{groups::GroupsFinder, milestones::MilestonesFinder};
use actix_web::{web, HttpResponse};

pub struct MilestonesController {
    base: DashboardApplicationController,
}

impl MilestonesController {
    pub fn new(base: DashboardApplicationController) -> Self {
        Self { base }
    }

    pub async fn index(&self, query: web::Query<MilestoneQuery>) -> HttpResponse {
        let groups = self.groups().await;
        let milestones = self.milestones(&query).await;

        if query.format == Some("json".to_string()) {
            return HttpResponse::Ok().json(milestones);
        }

        // HTML format
        let milestone_states = Milestone::states_count(&self.base.projects(), &groups);
        let paginated_milestones = milestones.paginate(query.page.unwrap_or(1));

        // TODO: Implement HTML template rendering
        HttpResponse::Ok().json(paginated_milestones)
    }

    async fn milestones(&self, query: &MilestoneQuery) -> Vec<Milestone> {
        let search_params = MilestoneSearchParams {
            state: query.state.clone(),
            search_title: query.search_title.clone(),
            group_ids: self.groups().await,
            project_ids: self.base.projects(),
        };

        MilestonesFinder::new(search_params).execute().await
    }

    async fn groups(&self) -> Vec<i32> {
        GroupsFinder::new(self.base.user.clone(), false)
            .execute()
            .await
            .iter()
            .map(|g| g.id)
            .collect()
    }
}

#[derive(serde::Deserialize)]
pub struct MilestoneQuery {
    state: Option<String>,
    search_title: Option<String>,
    page: Option<i32>,
    format: Option<String>,
}

struct MilestoneSearchParams {
    state: Option<String>,
    search_title: Option<String>,
    group_ids: Vec<i32>,
    project_ids: Vec<i32>,
}
