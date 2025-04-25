// Ported from: orig_app/app/controllers/concerns/issuable_collections_action.rb
// This module provides actions for issues, merge requests, and issues calendar collections.
// Ported on 2025-04-25.
//
// See porting log for details.

use crate::{
    models::{Issue, MergeRequest, Project, User},
    services::{issues::IssuesService, merge_requests::MergeRequestsService},
    utils::{pagination::PaginationHelper, sorting::SortingHelper},
};
use actix_web::{web, Error, HttpResponse};
use std::sync::Arc;

pub trait IssuableCollectionsAction {
    fn issues(
        &self,
        project: &Project,
        user: &User,
        params: web::Query<IssuableParams>,
    ) -> Result<HttpResponse, Error>;

    fn merge_requests(
        &self,
        project: &Project,
        user: &User,
        params: web::Query<IssuableParams>,
    ) -> Result<HttpResponse, Error>;

    fn issues_calendar(
        &self,
        project: &Project,
        user: &User,
        params: web::Query<IssuableParams>,
    ) -> Result<HttpResponse, Error>;
}

#[derive(serde::Deserialize)]
pub struct IssuableParams {
    pub state: Option<String>,
    pub sort: Option<String>,
    pub order_by: Option<String>,
    pub search: Option<String>,
    pub milestone_title: Option<String>,
    pub label_name: Option<Vec<String>>,
    pub assignee_username: Option<String>,
    pub author_username: Option<String>,
    pub scope: Option<String>,
    pub page: Option<i32>,
    pub per_page: Option<i32>,
}

pub struct IssuableCollectionsActionImpl {
    issues_service: Arc<IssuesService>,
    merge_requests_service: Arc<MergeRequestsService>,
    sorting_helper: Arc<SortingHelper>,
    pagination_helper: Arc<PaginationHelper>,
}

impl IssuableCollectionsActionImpl {
    pub fn new(
        issues_service: Arc<IssuesService>,
        merge_requests_service: Arc<MergeRequestsService>,
        sorting_helper: Arc<SortingHelper>,
        pagination_helper: Arc<PaginationHelper>,
    ) -> Self {
        Self {
            issues_service,
            merge_requests_service,
            sorting_helper,
            pagination_helper,
        }
    }

    fn find_issuables(
        &self,
        project: &Project,
        user: &User,
        params: &IssuableParams,
        issuable_type: &str,
    ) -> Result<Vec<Box<dyn Issuable>>, Error> {
        let state = params.state.as_deref().unwrap_or("opened");
        let sort = self
            .sorting_helper
            .parse_sort(params.sort.as_deref(), params.order_by.as_deref())?;

        match issuable_type {
            "issue" => {
                let issues = self.issues_service.list(
                    project,
                    user,
                    state,
                    &sort,
                    params.search.as_deref(),
                    params.milestone_title.as_deref(),
                    params.label_name.as_deref(),
                    params.assignee_username.as_deref(),
                    params.author_username.as_deref(),
                )?;
                Ok(issues
                    .into_iter()
                    .map(|i| Box::new(i) as Box<dyn Issuable>)
                    .collect())
            }
            "merge_request" => {
                let merge_requests = self.merge_requests_service.list(
                    project,
                    user,
                    state,
                    &sort,
                    params.search.as_deref(),
                    params.milestone_title.as_deref(),
                    params.label_name.as_deref(),
                    params.assignee_username.as_deref(),
                    params.author_username.as_deref(),
                )?;
                Ok(merge_requests
                    .into_iter()
                    .map(|mr| Box::new(mr) as Box<dyn Issuable>)
                    .collect())
            }
            _ => Err(actix_web::error::ErrorBadRequest("Invalid issuable type")),
        }
    }
}

impl IssuableCollectionsAction for IssuableCollectionsActionImpl {
    fn issues(
        &self,
        project: &Project,
        user: &User,
        params: web::Query<IssuableParams>,
    ) -> Result<HttpResponse, Error> {
        let issuables = self.find_issuables(project, user, &params, "issue")?;
        let paginated = self.pagination_helper.paginate(
            issuables,
            params.page.unwrap_or(1),
            params.per_page.unwrap_or(20),
        )?;

        Ok(HttpResponse::Ok().json(paginated))
    }

    fn merge_requests(
        &self,
        project: &Project,
        user: &User,
        params: web::Query<IssuableParams>,
    ) -> Result<HttpResponse, Error> {
        let issuables = self.find_issuables(project, user, &params, "merge_request")?;
        let paginated = self.pagination_helper.paginate(
            issuables,
            params.page.unwrap_or(1),
            params.per_page.unwrap_or(20),
        )?;

        Ok(HttpResponse::Ok().json(paginated))
    }

    fn issues_calendar(
        &self,
        project: &Project,
        user: &User,
        params: web::Query<IssuableParams>,
    ) -> Result<HttpResponse, Error> {
        let issues = self.issues_service.list_for_calendar(
            project,
            user,
            params.milestone_title.as_deref(),
            params.label_name.as_deref(),
        )?;

        Ok(HttpResponse::Ok().json(issues))
    }
}

pub trait Issuable: Send + Sync {
    fn id(&self) -> i64;
    fn title(&self) -> &str;
    fn state(&self) -> &str;
    fn created_at(&self) -> chrono::DateTime<chrono::Utc>;
    fn updated_at(&self) -> chrono::DateTime<chrono::Utc>;
}

impl Issuable for Issue {
    fn id(&self) -> i64 {
        self.id
    }
    fn title(&self) -> &str {
        &self.title
    }
    fn state(&self) -> &str {
        &self.state
    }
    fn created_at(&self) -> chrono::DateTime<chrono::Utc> {
        self.created_at
    }
    fn updated_at(&self) -> chrono::DateTime<chrono::Utc> {
        self.updated_at
    }
}

impl Issuable for MergeRequest {
    fn id(&self) -> i64 {
        self.id
    }
    fn title(&self) -> &str {
        &self.title
    }
    fn state(&self) -> &str {
        &self.state
    }
    fn created_at(&self) -> chrono::DateTime<chrono::Utc> {
        self.created_at
    }
    fn updated_at(&self) -> chrono::DateTime<chrono::Utc> {
        self.updated_at
    }
}
