// Ported from: orig_app/app/controllers/jira_connect/application_controller.rb
// This file implements the JiraConnect::ApplicationController base logic for JWT authentication and installation context.

use actix_web::{web, HttpRequest, HttpResponse, Responder};
use serde::Deserialize;
use std::sync::Arc;

// Placeholder for memoization (strong_memoize equivalent)
use once_cell::sync::OnceCell;

// Placeholder types for JWT and Installation
#[derive(Clone)]
pub struct JiraConnectInstallation {
    pub client_key: String,
    pub shared_secret: String,
    // ... other fields ...
}

#[derive(Clone)]
pub struct JwtSymmetric {
    pub token: String,
}

impl JwtSymmetric {
    pub fn new(token: String) -> Self {
        Self { token }
    }
    pub fn iss_claim(&self) -> Option<String> {
        // TODO: Parse JWT and extract iss claim
        None
    }
    pub fn sub_claim(&self) -> Option<String> {
        // TODO: Parse JWT and extract sub claim
        None
    }
    pub fn valid(&self, _shared_secret: &str) -> bool {
        // TODO: Validate JWT signature
        false
    }
    pub fn verify_qsh_claim(&self, _url: &str, _method: &str, _base_url: &str) -> bool {
        // TODO: Implement QSH claim verification
        false
    }
    pub fn verify_context_qsh_claim(&self) -> bool {
        // TODO: Implement context QSH claim verification
        false
    }
}

pub struct JiraConnectApplicationController {
    pub current_jira_installation: OnceCell<Option<JiraConnectInstallation>>,
    pub jwt: OnceCell<Option<JwtSymmetric>>,
}

impl JiraConnectApplicationController {
    pub fn new() -> Self {
        Self {
            current_jira_installation: OnceCell::new(),
            jwt: OnceCell::new(),
        }
    }

    pub fn verify_atlassian_jwt(&self, req: &HttpRequest) -> Result<(), HttpResponse> {
        if !self.atlassian_jwt_valid(req) {
            return Err(HttpResponse::Forbidden().finish());
        }
        Ok(())
    }

    pub fn verify_qsh_claim(
        &self,
        req: &HttpRequest,
        jira_connect_base_url: &str,
    ) -> Result<(), HttpResponse> {
        let jwt = self.jwt(req);
        if req
            .headers()
            .get("Content-Type")
            .map_or(false, |v| v == "application/json")
        {
            if let Some(jwt) = jwt {
                if jwt.verify_context_qsh_claim() {
                    return Ok(());
                }
            }
        }
        if let Some(jwt) = jwt {
            if jwt.verify_qsh_claim(
                req.uri().to_string().as_str(),
                req.method().as_str(),
                jira_connect_base_url,
            ) {
                return Ok(());
            }
        }
        Err(HttpResponse::Forbidden().finish())
    }

    pub fn atlassian_jwt_valid(&self, req: &HttpRequest) -> bool {
        let installation = self.installation_from_jwt(req);
        if installation.is_none() {
            return false;
        }
        let jwt = self.jwt(req);
        if let (Some(jwt), Some(installation)) = (jwt, installation) {
            jwt.valid(&installation.shared_secret)
        } else {
            false
        }
    }

    pub fn installation_from_jwt(&self, req: &HttpRequest) -> Option<JiraConnectInstallation> {
        self.current_jira_installation
            .get_or_init(|| {
                let jwt = self.jwt(req)?;
                let iss = jwt.iss_claim()?;
                // TODO: Lookup installation by client_key (iss)
                // JiraConnectInstallation::find_by_client_key(iss)
                None
            })
            .clone()
    }

    pub fn jira_user(&self, req: &HttpRequest) -> Option<String> {
        // Placeholder: return user info string
        let installation = self.installation_from_jwt(req)?;
        let jwt = self.jwt(req)?;
        let sub = jwt.sub_claim()?;
        // TODO: installation.client.user_info(sub)
        None
    }

    pub fn jwt(&self, req: &HttpRequest) -> Option<JwtSymmetric> {
        self.jwt
            .get_or_init(|| {
                let token = Self::auth_token(req)?;
                Some(JwtSymmetric::new(token))
            })
            .clone()
    }

    pub fn auth_token(req: &HttpRequest) -> Option<String> {
        // Try to get from query param 'jwt'
        if let Some(jwt) = req.query_string().split('&').find_map(|kv| {
            let mut parts = kv.splitn(2, '=');
            if let (Some(k), Some(v)) = (parts.next(), parts.next()) {
                if k == "jwt" {
                    return Some(v.to_string());
                }
            }
            None
        }) {
            return Some(jwt);
        }
        // Try to get from Authorization header
        if let Some(auth) = req.headers().get("Authorization") {
            if let Ok(auth_str) = auth.to_str() {
                if let Some((_, token)) = auth_str.split_once(' ') {
                    return Some(token.to_string());
                }
            }
        }
        None
    }
}
