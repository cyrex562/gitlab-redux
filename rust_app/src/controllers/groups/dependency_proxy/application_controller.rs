// Ported from: orig_app/app/controllers/groups/dependency_proxy/application_controller.rb

use crate::auth::{sign_in, AuthResult};
use crate::models::{DeployToken, PersonalAccessToken, User};
use crate::services::dependency_proxy::auth_token_service;
use actix_web::{HttpRequest, HttpResponse};

pub struct DependencyProxyApplicationController {
    authentication_result: Option<AuthResult>,
    personal_access_token: Option<PersonalAccessToken>,
}

impl DependencyProxyApplicationController {
    pub fn new() -> Self {
        Self {
            authentication_result: None,
            personal_access_token: None,
        }
    }

    pub fn authenticate_user_from_jwt_token(&mut self, req: &HttpRequest) -> HttpResponse {
        // Pseudocode for extracting token from header
        if let Some(token) = auth_token_service::extract_jwt_from_request(req) {
            self.authentication_result = Some(AuthResult::empty());
            let user_or_token = auth_token_service::user_or_token_from_jwt(&token);
            match user_or_token {
                Some(auth_token_service::UserOrToken::User(user)) => {
                    self.set_auth_result(user.clone(), "user");
                    if self.can_sign_in(&user) {
                        sign_in(&user);
                    }
                }
                Some(auth_token_service::UserOrToken::PersonalAccessToken(pat)) => {
                    self.set_auth_result(pat.user.clone(), "personal_access_token");
                    self.personal_access_token = Some(pat);
                }
                Some(auth_token_service::UserOrToken::DeployToken(dt)) => {
                    self.set_auth_result(dt, "deploy_token");
                }
                _ => {}
            }
        }
        if self.authenticated_user().is_none() {
            return self.request_bearer_token();
        }
        HttpResponse::Ok().finish()
    }

    fn authenticated_user(&self) -> Option<&User> {
        self.authentication_result
            .as_ref()
            .and_then(|ar| ar.actor())
    }

    fn request_bearer_token(&self) -> HttpResponse {
        // Set WWW-Authenticate header and return 401
        HttpResponse::Unauthorized()
            .append_header((
                "WWW-Authenticate",
                auth_token_service::registry_authenticate_header(),
            ))
            .body("")
    }

    fn can_sign_in(&self, user: &User) -> bool {
        !user.is_project_bot() && !user.is_service_account()
    }

    fn set_auth_result(&mut self, actor: User, type_: &str) {
        self.authentication_result = Some(AuthResult::new(actor, None, type_.to_string(), vec![]));
    }

    fn skip_session(&self, req: &HttpRequest) {
        // Pseudocode: mark session to be skipped
        // req.session().set_skip(true);
    }
}
