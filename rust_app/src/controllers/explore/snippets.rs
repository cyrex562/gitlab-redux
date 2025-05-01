// Ported from: orig_app/app/controllers/explore/snippets_controller.rb
// Ported on: 2025-04-30
// This file implements the Explore::SnippetsController from the Ruby codebase.
//
// See porting_log.txt for details.

use crate::controllers::concerns::GitlabNoteableMetadata;
use crate::controllers::explore::ExploreApplicationController;
use crate::services::snippets::SnippetsFinder;
use actix_web::{web, HttpResponse};

pub struct SnippetsController {
    base: ExploreApplicationController,
}

impl SnippetsController {
    pub fn new(base: ExploreApplicationController) -> Self {
        Self { base }
    }

    pub async fn index(&self, query: web::Query<SnippetQuery>) -> HttpResponse {
        let snippets = SnippetsFinder::new(
            self.base.user.clone(),
            true, // explore
            None,
            None,
        )
        .execute()
        .await
        .page(query.page.unwrap_or(1))
        .without_count()
        .inc_author()
        .inc_statistics();

        let noteable_meta_data = self.noteable_meta_data(&snippets, "Snippet");

        // TODO: Implement HTML template rendering
        HttpResponse::Ok().json(snippets)
    }
}

#[derive(serde::Deserialize)]
pub struct SnippetQuery {
    page: Option<i32>,
}

impl GitlabNoteableMetadata for SnippetsController {}
