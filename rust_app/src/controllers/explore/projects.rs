use crate::controllers::concerns::{
    PageLimiter, ParamsBackwardCompatibility, RendersMemberAccess, RendersProjectsList,
    SortingHelper, SortingPreference,
};
use crate::controllers::explore::ExploreApplicationController;
use crate::models::project::Project;
use crate::services::projects::ProjectsFinder;
use actix_web::{web, HttpResponse};

const MIN_SEARCH_LENGTH: usize = 3;
const PAGE_LIMIT: i32 = 50;
const RSS_ENTRIES_LIMIT: i32 = 20;

pub struct ProjectsController {
    base: ExploreApplicationController,
}

impl ProjectsController {
    pub fn new(base: ExploreApplicationController) -> Self {
        Self { base }
    }

    pub async fn index(&self, query: web::Query<ProjectQuery>) -> HttpResponse {
        self.show_alert_if_search_is_disabled(&query);
        let projects = self.load_projects(&query).await;

        if query.format == Some("json".to_string()) {
            return HttpResponse::Ok().json(projects);
        }

        // HTML format
        // TODO: Implement HTML template rendering
        HttpResponse::Ok().json(projects)
    }

    pub async fn trending(&self, query: web::Query<ProjectQuery>) -> HttpResponse {
        let mut query = query.clone();
        query.trending = Some(true);
        self.index(web::Query::from(query)).await
    }

    pub async fn starred(&self, query: web::Query<ProjectQuery>) -> HttpResponse {
        let mut projects = self.load_projects(&query).await;
        projects.sort_by(|a, b| b.star_count.cmp(&a.star_count));

        if query.format == Some("json".to_string()) {
            return HttpResponse::Ok().json(projects);
        }

        // HTML format
        // TODO: Implement HTML template rendering
        HttpResponse::Ok().json(projects)
    }

    pub async fn topics(&self, query: web::Query<ProjectQuery>) -> HttpResponse {
        let project_counts = self.load_project_counts().await;
        let topics = self.load_topics(&query).await;

        // TODO: Implement HTML template rendering
        HttpResponse::Ok().json((project_counts, topics))
    }

    pub async fn topic(&self, query: web::Query<ProjectQuery>) -> HttpResponse {
        let topic = self.load_topic(&query).await;

        if topic.is_none() {
            return HttpResponse::NotFound().finish();
        }

        let mut query = query.clone();
        query.topic = topic.map(|t| t.name);
        let projects = self.load_projects(&query).await;

        if query.format == Some("atom".to_string()) {
            let limited_projects = projects
                .iter()
                .take(RSS_ENTRIES_LIMIT as usize)
                .cloned()
                .collect::<Vec<_>>();
            // TODO: Implement atom feed rendering
            return HttpResponse::Ok().json(limited_projects);
        }

        // HTML format
        // TODO: Implement HTML template rendering
        HttpResponse::Ok().json(projects)
    }

    async fn load_project_counts(&self) -> ProjectCounts {
        let all_user_projects = ProjectsFinder::new(
            self.base.user.clone(),
            Some(true), // non_public
            None,
            None,
        )
        .execute()
        .await;

        let all_starred_projects = ProjectsFinder::new(
            self.base.user.clone(),
            None,
            Some(true), // starred
            None,
        )
        .execute()
        .await;

        ProjectCounts {
            all_user_projects,
            all_starred_projects,
        }
    }

    async fn load_projects(&self, query: &ProjectQuery) -> Vec<Project> {
        let project_counts = self.load_project_counts().await;

        let finder_params = ProjectFinderParams {
            minimum_search_length: MIN_SEARCH_LENGTH,
            not_aimed_for_deletion: true,
            current_organization: self.current_organization().await,
        };

        let mut projects =
            ProjectsFinder::new(self.base.user.clone(), finder_params, query.clone())
                .execute()
                .await;

        projects = self.preload_associations(projects);
        projects = projects.paginate(query.page.unwrap_or(1)).without_count();

        self.prepare_projects_for_rendering(projects)
    }

    async fn load_topics(&self, query: &ProjectQuery) -> Vec<Topic> {
        Projects::TopicsFinder::new(
            query.search.clone(),
            self.current_organization().await.map(|o| o.id),
        )
        .execute()
        .await
        .page(query.page.unwrap_or(1))
        .without_count()
    }

    async fn load_topic(&self, query: &ProjectQuery) -> Option<Topic> {
        let topic_name = if Feature::enabled("explore_topics_cleaned_path") {
            urlencoding::decode(&query.topic_name)?.into_owned()
        } else {
            query.topic_name.clone()?
        };

        let organization = self.current_organization().await?;

        Projects::Topic::for_organization(organization.id)
            .find_by_name_case_insensitive(&topic_name)
            .await
    }

    fn preload_associations(&self, projects: Vec<Project>) -> Vec<Project> {
        // TODO: Implement proper preloading
        projects
    }

    fn set_sorting(&mut self, query: &mut ProjectQuery) {
        query.sort = Some(self.set_sort_order());
    }

    fn default_sort_order(&self) -> String {
        self.sort_value_latest_activity()
    }

    fn sorting_field(&self) -> String {
        Project::SORTING_PREFERENCE_FIELD.to_string()
    }

    async fn current_organization(&self) -> Option<Organization> {
        Current::organization().await
    }

    fn show_alert_if_search_is_disabled(&self, query: &ProjectQuery) {
        if self.base.user.is_some()
            || (query.name.is_none() && query.search.is_none())
            || !self.is_html_request()
            || !Feature::disabled("disable_anonymous_project_search", "ops")
        {
            return;
        }

        // TODO: Implement flash message
    }
}

#[derive(serde::Deserialize, Clone)]
pub struct ProjectQuery {
    name: Option<String>,
    search: Option<String>,
    format: Option<String>,
    page: Option<i32>,
    sort: Option<String>,
    trending: Option<bool>,
    topic: Option<String>,
    topic_name: Option<String>,
}

struct ProjectCounts {
    all_user_projects: Vec<Project>,
    all_starred_projects: Vec<Project>,
}

struct ProjectFinderParams {
    minimum_search_length: usize,
    not_aimed_for_deletion: bool,
    current_organization: Option<Organization>,
}

impl PageLimiter for ProjectsController {}
impl ParamsBackwardCompatibility for ProjectQuery {}
impl RendersMemberAccess for ProjectsController {}
impl RendersProjectsList for ProjectsController {}
impl SortingHelper for ProjectsController {}
impl SortingPreference for ProjectsController {}
