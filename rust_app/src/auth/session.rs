use super::AuthUser;
use actix_web::{web, HttpRequest, HttpResponse, Responder};
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct Session {
    pub user_id: Option<i32>,
    pub last_activity: chrono::DateTime<chrono::Utc>,
    pub remember_me: bool,
}

impl Session {
    pub fn new() -> Self {
        Self {
            user_id: None,
            last_activity: chrono::Utc::now(),
            remember_me: false,
        }
    }

    pub fn with_user(mut self, user: &AuthUser) -> Self {
        self.user_id = Some(user.id);
        self
    }

    pub fn with_remember_me(mut self, remember: bool) -> Self {
        self.remember_me = remember;
        self
    }

    pub fn update_activity(&mut self) {
        self.last_activity = chrono::Utc::now();
    }

    pub fn is_expired(&self) -> bool {
        let now = chrono::Utc::now();
        let expiry = if self.remember_me {
            chrono::Duration::days(30)
        } else {
            chrono::Duration::hours(24)
        };
        now - self.last_activity > expiry
    }
}

pub struct SessionManager;

impl SessionManager {
    pub fn new() -> Self {
        Self
    }

    pub fn create_session(&self, user: &AuthUser, remember_me: bool) -> Session {
        Session::new().with_user(user).with_remember_me(remember_me)
    }

    pub fn get_session(&self, req: &HttpRequest) -> Option<Session> {
        // TODO: Implement session retrieval from request
        None
    }

    pub fn destroy_session(&self, req: &HttpRequest) -> impl Responder {
        // TODO: Implement session destruction
        HttpResponse::Ok().finish()
    }

    pub fn validate_session(&self, session: &Session) -> bool {
        !session.is_expired()
    }
}
