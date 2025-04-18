use actix_web::{HttpResponse, Responder};
use crate::models::Project;

pub trait RepositorySettingsRedirect {
    fn redirect_to_repository_settings(&self, project: &Project, anchor: Option<&str>) -> impl Responder;
}

impl<T> RepositorySettingsRedirect for T
where
    T: actix_web::web::ServiceFactory<
        ServiceRequest = actix_web::dev::ServiceRequest,
        Response = actix_web::dev::ServiceResponse,
        Error = actix_web::Error,
        InitError = (),
    >,
{
    fn redirect_to_repository_settings(&self, project: &Project, anchor: Option<&str>) -> impl Responder {
        let mut path = format!("/projects/{}/settings/repository", project.id);
        if let Some(anchor) = anchor {
            path.push_str(&format!("#{}", anchor));
        }
        HttpResponse::Found().header("Location", path).finish()
    }
} 