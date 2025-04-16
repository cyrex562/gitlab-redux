use super::{RenderOptions, Renderable};
use actix_web::{web, HttpResponse};
use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct Project {
    pub id: i32,
    pub name: String,
    pub path: String,
    pub description: Option<String>,
    pub visibility: String,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
    pub last_activity_at: Option<DateTime<Utc>>,
    pub namespace: ProjectNamespace,
    pub statistics: Option<ProjectStatistics>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ProjectNamespace {
    pub id: i32,
    pub name: String,
    pub path: String,
    pub kind: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ProjectStatistics {
    pub commit_count: i32,
    pub storage_size: i64,
    pub repository_size: i64,
    pub wiki_size: i64,
    pub lfs_objects_size: i64,
    pub build_artifacts_size: i64,
    pub packages_size: i64,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ProjectCollection {
    pub projects: Vec<Project>,
    pub total_count: usize,
    pub current_page: i32,
    pub total_pages: i32,
}

impl Project {
    pub fn new(
        id: i32,
        name: String,
        path: String,
        namespace: ProjectNamespace,
        visibility: String,
        created_at: DateTime<Utc>,
        updated_at: DateTime<Utc>,
    ) -> Self {
        Self {
            id,
            name,
            path,
            description: None,
            visibility,
            created_at,
            updated_at,
            last_activity_at: None,
            namespace,
            statistics: None,
        }
    }

    pub fn with_description(mut self, description: String) -> Self {
        self.description = Some(description);
        self
    }

    pub fn with_statistics(mut self, statistics: ProjectStatistics) -> Self {
        self.statistics = Some(statistics);
        self
    }
}

impl Renderable for Project {
    fn render(&self, options: &RenderOptions) -> HttpResponse {
        // TODO: Implement project rendering
        HttpResponse::Ok().json(self)
    }

    fn render_json(&self, options: &RenderOptions) -> HttpResponse {
        self.render(options)
    }
}

impl Renderable for ProjectCollection {
    fn render(&self, options: &RenderOptions) -> HttpResponse {
        // TODO: Implement project collection rendering
        HttpResponse::Ok().json(self)
    }

    fn render_json(&self, options: &RenderOptions) -> HttpResponse {
        self.render(options)
    }
}

pub struct ProjectRenderer;

impl ProjectRenderer {
    pub fn new() -> Self {
        Self
    }

    pub fn render_project(&self, project: &Project, options: &RenderOptions) -> HttpResponse {
        project.render(options)
    }

    pub fn render_projects(
        &self,
        projects: &[Project],
        page: i32,
        per_page: i32,
        total_count: usize,
        options: &RenderOptions,
    ) -> HttpResponse {
        let total_pages = (total_count as f64 / per_page as f64).ceil() as i32;
        let collection = ProjectCollection {
            projects: projects.to_vec(),
            total_count,
            current_page: page,
            total_pages,
        };
        collection.render(options)
    }

    pub fn prepare_projects_for_rendering(&self, projects: &[Project]) -> Vec<Project> {
        // TODO: Implement project preparation for rendering
        projects.to_vec()
    }
}
