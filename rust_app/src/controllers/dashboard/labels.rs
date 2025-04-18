use crate::controllers::dashboard::DashboardApplicationController;
use crate::models::label::Label;
use crate::services::labels::LabelsFinder;
use actix_web::{web, HttpResponse};

pub struct LabelsController {
    base: DashboardApplicationController,
}

impl LabelsController {
    pub fn new(base: DashboardApplicationController) -> Self {
        Self { base }
    }

    pub async fn index(&self) -> HttpResponse {
        let labels = self.labels().await;
        HttpResponse::Ok().json(labels)
    }

    async fn labels(&self) -> Vec<Label> {
        let finder_params = LabelsFinderParams {
            project_ids: self.base.projects(),
        };

        LabelsFinder::new(self.base.user.clone(), finder_params)
            .execute()
            .await
        // TODO: Implement DISTINCT ON (labels.title) equivalent
        // This is a placeholder that needs to be implemented with proper SQL query
    }
}

struct LabelsFinderParams {
    project_ids: Vec<i32>,
}
