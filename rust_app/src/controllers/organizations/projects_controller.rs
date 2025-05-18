// Ported from: orig_app/app/controllers/organizations/projects_controller.rb
use actix_web::{get, web, HttpResponse, Responder};
use serde::Deserialize;
use std::sync::Arc;

use crate::auth::authorization::can;
use crate::models::{Project, User};

pub struct ProjectsController;

impl ProjectsController {
    pub async fn edit() -> impl Responder {
        // Edit action just renders the view
        HttpResponse::Ok().finish()
    }

    pub async fn find_project(namespace_id: String, id: String) -> Option<Arc<Project>> {
        let full_path = format!("{}/{}", namespace_id, id);
        // TODO: Implement actual project lookup by full path
        Project::find_by_full_path(&full_path).await
    }

    pub async fn authorize_project_view_edit_page(current_user: &User, project: &Project) -> bool {
        can(current_user, "view_edit_page", project).await
    }

    pub fn safe_params(params: web::Query<ProjectParams>) -> ProjectParams {
        params.into_inner()
    }
}

#[derive(serde::Deserialize)]
pub struct ProjectParams {
    pub id: String,
    pub namespace_id: String,
}

#[derive(Deserialize)]
pub struct PathParams {
    namespace_id: String,
    id: String,
}

pub fn configure(cfg: &mut web::ServiceConfig) {
    cfg.service(
        web::scope("/organizations/{organization_path}/projects")
            .service(web::resource("/{namespace_id}/{id}/edit").route(web::get().to(edit))),
    );
}

async fn edit(
    path: web::Path<PathParams>,
    current_user: web::ReqData<Arc<User>>,
) -> impl Responder {
    // Find project
    let project =
        match ProjectsController::find_project(path.namespace_id.clone(), path.id.clone()).await {
            Some(project) => project,
            None => return HttpResponse::NotFound().finish(),
        };

    // Check authorization
    if !ProjectsController::authorize_project_view_edit_page(&current_user, &project).await {
        return HttpResponse::Forbidden().finish();
    }

    // Just render the view
    HttpResponse::Ok().finish()
}
