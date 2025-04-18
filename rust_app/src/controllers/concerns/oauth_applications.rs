use crate::models::oauth_application::OAuthApplication;
use crate::models::user::User;
use crate::services::oauth_applications::{CreateService, DestroyService, UpdateService};
use crate::utils::boolean::to_boolean;
use actix_web::{web, HttpRequest, HttpResponse, Responder};
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct OAuthApplicationParams {
    pub name: String,
    pub redirect_uri: String,
    pub scopes: Vec<String>,
    pub confidential: Option<bool>,
}

pub trait OAuthApplications {
    fn index(&self, req: HttpRequest) -> impl Responder;
    fn new(&self, req: HttpRequest) -> impl Responder;
    fn create(&self, req: HttpRequest) -> impl Responder;
    fn edit(&self, req: HttpRequest) -> impl Responder;
    fn update(&self, req: HttpRequest) -> impl Responder;
    fn destroy(&self, req: HttpRequest) -> impl Responder;
}

pub struct OAuthApplicationsImpl {
    current_user: User,
    application: Option<OAuthApplication>,
}

impl OAuthApplicationsImpl {
    pub fn new(current_user: User, application: Option<OAuthApplication>) -> Self {
        Self {
            current_user,
            application,
        }
    }

    fn applications_path(&self) -> String {
        format!("/users/{}/oauth/applications", self.current_user.username)
    }

    fn application_path(&self, application: &OAuthApplication) -> String {
        format!("/oauth/applications/{}", application.id)
    }

    fn edit_application_path(&self, application: &OAuthApplication) -> String {
        format!("/oauth/applications/{}/edit", application.id)
    }

    fn render_form(&self, application: Option<&OAuthApplication>) -> String {
        // Implementation depends on your template engine
        // This is a placeholder that should be replaced with actual template rendering
        format!("Rendered form for application: {:?}", application)
    }
}

impl OAuthApplications for OAuthApplicationsImpl {
    fn index(&self, _req: HttpRequest) -> impl Responder {
        let applications = self.current_user.oauth_applications();
        HttpResponse::Ok().json(serde_json::json!({
            "applications": applications.iter().map(|app| {
                serde_json::json!({
                    "id": app.id,
                    "name": app.name,
                    "redirect_uri": app.redirect_uri,
                    "scopes": app.scopes,
                    "confidential": app.confidential,
                    "created_at": app.created_at,
                    "updated_at": app.updated_at
                })
            }).collect::<Vec<_>>()
        }))
    }

    fn new(&self, _req: HttpRequest) -> impl Responder {
        HttpResponse::Ok().json(serde_json::json!({
            "html": self.render_form(None)
        }))
    }

    fn create(&self, req: HttpRequest) -> impl Responder {
        let params = web::Json::<OAuthApplicationParams>::from_request(&req)
            .map_err(|_| HttpResponse::BadRequest())?;

        let application = CreateService::new(&self.current_user, &params.0).execute();

        if application.has_errors() {
            HttpResponse::UnprocessableEntity().json(serde_json::json!({
                "errors": application.errors().join(", ")
            }))
        } else {
            HttpResponse::Created()
                .header("Location", self.application_path(&application))
                .json(serde_json::json!({
                    "application": serde_json::json!({
                        "id": application.id,
                        "name": application.name,
                        "redirect_uri": application.redirect_uri,
                        "scopes": application.scopes,
                        "confidential": application.confidential,
                        "created_at": application.created_at,
                        "updated_at": application.updated_at
                    })
                }))
        }
    }

    fn edit(&self, req: HttpRequest) -> impl Responder {
        let application = self
            .application
            .as_ref()
            .ok_or_else(|| HttpResponse::NotFound())?;

        HttpResponse::Ok().json(serde_json::json!({
            "html": self.render_form(Some(application))
        }))
    }

    fn update(&self, req: HttpRequest) -> impl Responder {
        let application = self
            .application
            .as_ref()
            .ok_or_else(|| HttpResponse::NotFound())?;

        let params = web::Json::<OAuthApplicationParams>::from_request(&req)
            .map_err(|_| HttpResponse::BadRequest())?;

        let application = UpdateService::new(&self.current_user, &params.0).execute(application);

        if application.has_errors() {
            HttpResponse::UnprocessableEntity().json(serde_json::json!({
                "errors": application.errors().join(", ")
            }))
        } else {
            HttpResponse::Ok().json(serde_json::json!({
                "application": serde_json::json!({
                    "id": application.id,
                    "name": application.name,
                    "redirect_uri": application.redirect_uri,
                    "scopes": application.scopes,
                    "confidential": application.confidential,
                    "created_at": application.created_at,
                    "updated_at": application.updated_at
                })
            }))
        }
    }

    fn destroy(&self, req: HttpRequest) -> impl Responder {
        let application = self
            .application
            .as_ref()
            .ok_or_else(|| HttpResponse::NotFound())?;

        DestroyService::new(&self.current_user).execute(application);

        HttpResponse::SeeOther()
            .header("Location", self.applications_path())
            .json(serde_json::json!({
                "notice": "Application was successfully deleted."
            }))
    }
}
