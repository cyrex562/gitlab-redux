// Ported from: orig_app/app/controllers/oauth/applications_controller.rb
// Ported: 2025-05-01
//
// This controller handles user OAuth applications (Doorkeeper::Application) management.
// Actions: index, show, create, renew

use actix_web::{get, post, web, HttpRequest, HttpResponse, Responder};
use serde::{Deserialize, Serialize};
use serde_json::json;
use std::sync::Arc;

use crate::models::oauth::Application;
use crate::models::user::User;
use crate::services::applications::{ApplicationParams, CreateService};

#[derive(Debug, Serialize, Deserialize)]
pub struct CreateApplicationRequest {
    pub name: String,
    pub redirect_uri: String,
    pub scopes: String,
    pub confidential: bool,
}

#[get("")]
pub async fn index(req: HttpRequest) -> impl Responder {
    // Authenticate user first
    if !authenticate_user(&req).await {
        return HttpResponse::Unauthorized().finish();
    }

    let current_user = get_current_user(&req);
    if let Some(user) = current_user {
        set_index_vars(user).await
    } else {
        HttpResponse::Unauthorized().finish()
    }
}

#[post("")]
pub async fn create(
    req: HttpRequest,
    params: web::Json<CreateApplicationRequest>,
) -> impl Responder {
    // Authenticate user and verify OAuth applications are enabled
    if !authenticate_user(&req).await {
        return HttpResponse::Unauthorized().finish();
    }

    let current_user = get_current_user(&req);
    if let Some(user) = current_user {
        if !verify_user_oauth_applications_enabled().await {
            return HttpResponse::Found()
                .append_header(("Location", "/user/settings/profile"))
                .finish();
        }

        let service_params = ApplicationParams {
            name: params.name.clone(),
            redirect_uri: params.redirect_uri.clone(),
            scopes: params.scopes.clone(),
            confidential: params.confidential,
        };

        match CreateService::new(user.clone(), service_params)
            .execute(&req)
            .await
        {
            Ok(application) => {
                if application.is_persisted() {
                    HttpResponse::Ok()
                        .append_header(("X-Flash-Notice", "Application created successfully."))
                        .json(application)
                } else {
                    set_index_vars(user).await
                }
            }
            Err(_) => HttpResponse::UnprocessableEntity().finish(),
        }
    } else {
        HttpResponse::Unauthorized().finish()
    }
}

#[post("/{id}/renew")]
pub async fn renew(req: HttpRequest, path: web::Path<i32>) -> impl Responder {
    if !authenticate_user(&req).await {
        return HttpResponse::Unauthorized().finish();
    }

    let current_user = get_current_user(&req);
    if let Some(user) = current_user {
        match set_application(user, path.into_inner()).await {
            Ok(mut app) => {
                app.renew_secret();
                // TODO: Save updated application
                HttpResponse::Ok().json(json!({ "secret": app.secret }))
            }
            Err(_) => HttpResponse::NotFound().finish(),
        }
    } else {
        HttpResponse::Unauthorized().finish()
    }
}

// Private functions
async fn authenticate_user(_req: &HttpRequest) -> bool {
    // TODO: Implement actual authentication
    true
}

fn get_current_user(_req: &HttpRequest) -> Option<Arc<User>> {
    // TODO: Get actual user from request
    None
}

async fn verify_user_oauth_applications_enabled() -> bool {
    // TODO: Check settings if OAuth applications are enabled
    true
}

async fn set_index_vars(_user: Arc<User>) -> HttpResponse {
    // TODO: Set variables for index view
    HttpResponse::Ok().finish()
}

async fn set_application(_user: Arc<User>, _id: i32) -> Result<Application, &'static str> {
    // TODO: Find application by ID
    Err("Not found")
}

pub fn configure(cfg: &mut web::ServiceConfig) {
    cfg.service(
        web::scope("/oauth/applications")
            .service(index)
            .service(create)
            .service(renew),
    );
}
