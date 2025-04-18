use crate::controllers::concerns::{GitlabNoteableMetadata, PaginatedCollection, SnippetsSort};
use crate::controllers::dashboard::DashboardApplicationController;
use crate::services::{snippets::CountService, snippets::SnippetsFinder};
use actix_web::{web, HttpResponse};

pub struct SnippetsController {
    base: DashboardApplicationController,
}

impl SnippetsController {
    pub fn new(base: DashboardApplicationController) -> Self {
        Self { base }
    }

    pub async fn index(&self, query: web::Query<SnippetQuery>) -> HttpResponse {
        let snippet_counts =
            CountService::new(self.base.user.clone(), Some(self.base.user.clone()))
                .execute()
                .await;

        let snippets = SnippetsFinder::new(
            self.base.user.clone(),
            Some(self.base.user.clone()),
            query.scope.clone(),
            self.sort_param(),
        )
        .execute()
        .await
        .page(query.page.unwrap_or(1))
        .inc_author()
        .inc_projects_namespace_route()
        .inc_statistics();

        if let Some(redirect) = self.redirect_out_of_range(&snippets) {
            return redirect;
        }

        let noteable_meta_data = self.noteable_meta_data(&snippets, "Snippet");

        // TODO: Implement HTML template rendering
        HttpResponse::Ok().json(snippets)
    }
}

#[derive(serde::Deserialize)]
pub struct SnippetQuery {
    scope: Option<String>,
    page: Option<i32>,
}

impl PaginatedCollection for SnippetsController {}
impl GitlabNoteableMetadata for SnippetsController {}
impl SnippetsSort for SnippetsController {}
