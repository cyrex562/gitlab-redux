use crate::models::user::User;
use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use uuid::Uuid;

#[derive(Debug, Serialize, Deserialize)]
pub struct Application {
    pub id: i32,
    pub name: String,
    pub uid: String,
    pub secret: String,
    pub redirect_uri: String,
    pub scopes: String,
    pub confidential: bool,
    pub owner_id: Option<i32>,
    pub owner_type: Option<String>,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
    pub trusted: bool,
}

impl Application {
    pub fn new(
        name: String,
        redirect_uri: String,
        scopes: String,
        confidential: bool,
        owner: Option<&User>,
    ) -> Self {
        Self {
            id: 0, // Will be set by database
            name,
            uid: Uuid::new_v4().to_string(),
            secret: Uuid::new_v4().to_string(),
            redirect_uri,
            scopes,
            confidential,
            owner_id: owner.map(|u| u.id),
            owner_type: owner.map(|_| "User".to_string()),
            created_at: Utc::now(),
            updated_at: Utc::now(),
            trusted: false,
        }
    }

    pub fn is_persisted(&self) -> bool {
        self.id > 0
    }

    pub fn renew_secret(&mut self) {
        self.secret = Uuid::new_v4().to_string();
        self.updated_at = Utc::now();
    }
}
