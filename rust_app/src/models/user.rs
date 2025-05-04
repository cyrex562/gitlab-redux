use std::sync::Arc;
use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct User {
    pub id: i32,
    pub username: String,
    pub email: String,
    pub confirmed_at: Option<chrono::DateTime<chrono::Utc>>,
    pub admin: bool,
}

impl User {
    pub fn confirmed(&self) -> bool {
        self.confirmed_at.is_some()
    }

    pub fn admin(&self) -> bool {
        self.admin
    }

    pub fn can_log_in(&self) -> bool {
        self.confirmed()
    }

    pub fn valid_password(&self, _password: &str) -> bool {
        // TODO: Implement password validation
        false
    }

    pub fn increment_failed_attempts(&self) {
        // TODO: Implement failed attempts tracking
    }
}