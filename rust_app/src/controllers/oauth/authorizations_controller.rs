// Ported from: orig_app/app/controllers/oauth/authorizations_controller.rb
// Ported on: 2025-05-04
//
// Handles OAuth authorization logic, extending Doorkeeper's authorization controller
// with GitLab-specific functionality.

use actix_web::{web, HttpRequest, HttpResponse};
use serde_json::json;
use std::sync::Arc;
use url::Url;

use crate::auth::doorkeeper::{PreAuth, Authorization, Application};
use crate::models::user::User;
use crate::services::audit::Auditor;
use crate::models::current::CurrentSettings;

pub struct AuthorizationsController {
    current_user: Option<Arc<User>>,
    pre_auth: Option<PreAuth>,
    doorkeeper_application: Option<Arc<Application>>,
}

pub struct Config {
    pub controller_path: &'static str,
}

impl Default for Config {
    fn default() -> Self {
        Self {
            controller_path: "/oauth/authorize",
        }
    }
}

pub fn configure(cfg: &mut web::ServiceConfig) {
    let config = Config::default();
    
    cfg.service(
        web::resource(config.controller_path)
            .route(web::get().to(AuthorizationsController::new_action))
    );
}

impl AuthorizationsController {
    pub async fn new_action(
        req: HttpRequest,
        session: web::Data<actix_session::Session>,
        data: web::Data<Arc<AuthorizationsController>>,
    ) -> HttpResponse {
        let mut controller = data.as_ref().clone();
        
        if let Some(pre_auth) = controller.pre_auth.as_ref() {
            if pre_auth.authorizable() {
                if controller.skip_authorization() || (controller.matching_token() && pre_auth.client().application().confidential()) {
                    let auth = controller.authorization().authorize();
                    if let Ok(redirect_uri) = Url::parse(&auth.redirect_uri) {
                        if let Ok(_) = session.remove("user_return_to") {
                            return HttpResponse::Ok().json(json!({
                                "view": "doorkeeper/authorizations/redirect",
                                "redirect_uri": redirect_uri,
                                "layout": false
                            }));
                        }
                    }
                } else {
                    if let Ok(redirect_uri) = Url::parse(&controller.authorization().authorize().redirect_uri) {
                        controller.allow_redirect_uri_form_action(redirect_uri.scheme());
                        return HttpResponse::Ok().json(json!({
                            "view": "doorkeeper/authorizations/new"
                        }));
                    }
                }
            }
            return HttpResponse::Ok().json(json!({
                "view": "doorkeeper/authorizations/error"
            }));
        }
        HttpResponse::BadRequest().finish()
    }

    // Private methods
    fn new(current_user: Option<Arc<User>>) -> Self {
        Self {
            current_user,
            pre_auth: None,
            doorkeeper_application: None,
        }
    }

    fn auth_user(&self) -> Option<Arc<User>> {
        self.current_user.clone()
    }

    fn audit_oauth_authorization(&self, response: &HttpResponse, req: &HttpRequest) -> Result<(), Box<dyn std::error::Error>> {
        if !response.status().is_success() && !response.status().is_redirection() {
            return Ok(());
        }

        if let Some(pre_auth) = self.pre_auth.as_ref() {
            if let Some(application) = pre_auth.client().application() {
                Auditor::audit(
                    "user_authorized_oauth_application",
                    self.current_user.as_ref().unwrap(),
                    self.current_user.as_ref().unwrap(),
                    &application,
                    "User authorized an OAuth application.",
                    json!({
                        "application_name": application.name,
                        "application_id": application.id,
                        "scopes": application.scopes
                    }),
                    Some(&req.connection_info().realip_remote_addr().unwrap_or("").to_string())
                )?;
            }
        }
        Ok(())
    }

    fn allow_redirect_uri_form_action(&self, redirect_uri_scheme: &str) {
        // TODO: Implement CSP form-action modification once we have CSP handling
    }

    fn skip_authorization(&self) -> bool {
        false // TODO: Implement skip authorization logic
    }

    fn matching_token(&self) -> bool {
        false // TODO: Implement matching token logic
    }

    fn authorization(&self) -> Authorization {
        Authorization // TODO: Implement proper authorization
    }

    fn action_name(&self) -> &str {
        "new"
    }

    fn downgrade_scopes(&mut self, params: &mut web::Query<std::collections::HashMap<String, String>>) {
        let auth_type = params.remove("gl_auth_type");
        if auth_type.as_deref() != Some("login") {
            return;
        }

        self.ensure_read_user_scope();

        if self.application_has_read_user_scope() {
            params.insert("scope".to_string(), crate::auth::Auth::READ_USER_SCOPE.to_string());
        }
    }

    fn ensure_read_user_scope(&mut self) {
        if self.application_has_read_user_scope() || !self.application_has_api_scope() {
            return;
        }
        self.add_read_user_scope();
    }

    fn add_read_user_scope(&mut self) {
        if let Some(app) = self.doorkeeper_application.clone() {
            let mut scopes = app.scopes.clone();
            scopes.push(crate::auth::Auth::READ_USER_SCOPE.to_string());
            
            // Create a new application with updated scopes
            let mut new_app = (*app).clone();
            let _ = new_app.update_scopes(scopes);
            self.doorkeeper_application = Some(Arc::new(new_app));
        }
    }

    fn doorkeeper_application(&mut self) -> Option<Arc<Application>> {
        if self.doorkeeper_application.is_none() {
            // TODO: Implement proper client lookup using client_id from params
            self.doorkeeper_application = None;
        }
        self.doorkeeper_application.clone()
    }

    fn application_has_read_user_scope(&self) -> bool {
        self.doorkeeper_application()
            .map(|app| app.includes_scope(crate::auth::Auth::READ_USER_SCOPE))
            .unwrap_or(false)
    }

    fn application_has_api_scope(&self) -> bool {
        self.doorkeeper_application()
            .map(|app| crate::auth::Auth::API_SCOPES.iter().any(|scope| app.includes_scope(scope)))
            .unwrap_or(false)
    }

    fn verify_confirmed_email(&self) -> Result<(), HttpResponse> {
        if self.current_user.as_ref().map_or(false, |u| u.confirmed()) {
            Ok(())
        } else {
            if let Some(pre_auth) = &mut self.pre_auth {
                pre_auth.set_error("unconfirmed_email");
            }
            Err(HttpResponse::Ok().json(json!({
                "view": "doorkeeper/authorizations/error"
            })))
        }
    }

    fn verify_admin_allowed(&self) -> Result<(), HttpResponse> {
        if self.disallow_connect() {
            Err(HttpResponse::Ok().json(json!({
                "view": "doorkeeper/authorizations/forbidden"
            })))
        } else {
            Ok(())
        }
    }

    fn disallow_connect(&self) -> bool {
        self.current_user.as_ref().map_or(false, |u| u.admin()) 
            && CurrentSettings::disable_admin_oauth_scopes()
            && self.dangerous_scopes()
    }

    fn dangerous_scopes(&self) -> bool {
        if let Some(app) = self.doorkeeper_application() {
            if app.trusted() {
                return false;
            }
            
            let dangerous = [
                crate::auth::Auth::API_SCOPE,
                crate::auth::Auth::READ_API_SCOPE,
                &crate::auth::Auth::ADMIN_SCOPES[..],
                &crate::auth::Auth::REPOSITORY_SCOPES[..],
                &crate::auth::Auth::REGISTRY_SCOPES[..],
            ].concat();
            
            dangerous.iter().any(|scope| app.includes_scope(scope))
        } else {
            false
        }
    }
}
