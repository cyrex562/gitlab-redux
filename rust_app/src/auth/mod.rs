pub mod doorkeeper;
pub mod impersonation;
pub mod security_policy;
pub mod session;
pub mod two_factor;

use actix_web::{web, HttpRequest, HttpResponse, Responder};
use serde::{Deserialize, Serialize};
use std::sync::Arc;

pub use doorkeeper::{Application, Authorization, Client, PreAuth};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AuthUser {
    pub id: i32,
    pub username: String,
    pub email: String,
    pub two_factor_enabled: bool,
    pub admin: bool,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AuthSession {
    pub user: Option<AuthUser>,
    pub impersonator: Option<AuthUser>,
}

pub trait Authenticatable {
    fn authenticate(&self, req: &HttpRequest) -> Result<AuthSession, actix_web::Error>;
    fn authorize(&self, user: &AuthUser, resource: &str) -> bool;
}

pub struct Auth;

impl Auth {
    pub const READ_USER_SCOPE: &'static str = "read_user";
    pub const API_SCOPE: &'static str = "api";
    pub const READ_API_SCOPE: &'static str = "read_api";
    pub const ADMIN_SCOPES: [&'static str; 1] = ["admin"];
    pub const REPOSITORY_SCOPES: [&'static str; 2] = ["read_repository", "write_repository"];
    pub const REGISTRY_SCOPES: [&'static str; 2] = ["read_registry", "write_registry"];
    pub const API_SCOPES: [&'static str; 2] = ["api", "read_api"];

    pub fn new() -> Self {
        Self
    }

    pub fn authenticate_request(&self, req: &HttpRequest) -> Result<AuthSession, actix_web::Error> {
        // TODO: Implement request authentication
        Ok(AuthSession {
            user: None,
            impersonator: None,
        })
    }

    pub fn require_authentication(&self) -> impl Responder {
        HttpResponse::Unauthorized().finish()
    }

    pub fn require_admin(&self) -> impl Responder {
        HttpResponse::Forbidden().finish()
    }

    pub fn skip_authentication(&self) -> bool {
        false
    }
}
