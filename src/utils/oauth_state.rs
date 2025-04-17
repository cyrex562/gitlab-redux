use serde::{Deserialize, Serialize};
use std::time::{Duration, SystemTime};

#[derive(Debug, Serialize, Deserialize)]
pub struct OAuthState {
    pub state: String,
    pub created_at: SystemTime,
    pub expires_at: SystemTime,
}

impl OAuthState {
    pub fn new(state: String) -> Self {
        let now = SystemTime::now();
        let expires_at = now + Duration::from_secs(3600); // 1 hour expiration

        Self {
            state,
            created_at: now,
            expires_at,
        }
    }

    pub fn is_valid(&self) -> bool {
        if let Ok(now) = SystemTime::now().duration_since(self.created_at) {
            now < Duration::from_secs(3600)
        } else {
            false
        }
    }
}

impl Default for OAuthState {
    fn default() -> Self {
        Self::new(String::new())
    }
}
