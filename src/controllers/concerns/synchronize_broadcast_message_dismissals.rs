use actix_web::{web, HttpRequest, HttpResponse, Result};
use chrono::{DateTime, Utc};
use std::sync::Arc;

// Define the User trait
pub trait User: Send + Sync {
    fn id(&self) -> i64;
    fn username(&self) -> &str;
}

// Define the BroadcastMessageDismissal struct
pub struct BroadcastMessageDismissal {
    pub id: i64,
    pub user_id: i64,
    pub broadcast_message_id: i64,
    pub cookie_key: String,
    pub expires_at: DateTime<Utc>,
}

// Define the BroadcastMessageDismissalFinder trait
pub trait BroadcastMessageDismissalFinder: Send + Sync {
    fn execute(&self, user: &dyn User) -> Vec<BroadcastMessageDismissal>;
}

// Define the Cookie trait
pub trait Cookie: Send + Sync {
    fn get(&self, key: &str) -> Option<String>;
    fn set(&mut self, key: &str, value: String, expires_at: DateTime<Utc>);
}

// Define the SynchronizeBroadcastMessageDismissals trait
pub trait SynchronizeBroadcastMessageDismissals {
    fn synchronize_broadcast_message_dismissals(&self, user: &dyn User) -> Result<()>;
    fn create_dismissal_cookie(&self, dismissal: &BroadcastMessageDismissal) -> Result<()>;
    fn get_dismissal_finder(&self) -> Arc<dyn BroadcastMessageDismissalFinder>;
    fn get_cookies(&self) -> Arc<dyn Cookie>;
}

// Define the SynchronizeBroadcastMessageDismissalsHandler struct
pub struct SynchronizeBroadcastMessageDismissalsHandler {
    dismissal_finder: Arc<dyn BroadcastMessageDismissalFinder>,
    cookies: Arc<dyn Cookie>,
}

impl SynchronizeBroadcastMessageDismissalsHandler {
    pub fn new(
        dismissal_finder: Arc<dyn BroadcastMessageDismissalFinder>,
        cookies: Arc<dyn Cookie>,
    ) -> Self {
        SynchronizeBroadcastMessageDismissalsHandler {
            dismissal_finder,
            cookies,
        }
    }
}

// Implement the SynchronizeBroadcastMessageDismissals trait for SynchronizeBroadcastMessageDismissalsHandler
impl SynchronizeBroadcastMessageDismissals for SynchronizeBroadcastMessageDismissalsHandler {
    fn synchronize_broadcast_message_dismissals(&self, user: &dyn User) -> Result<()> {
        let dismissals = self.dismissal_finder.execute(user);

        for dismissal in dismissals {
            if self.cookies.get(&dismissal.cookie_key).is_none() {
                self.create_dismissal_cookie(&dismissal)?;
            }
        }

        Ok(())
    }

    fn create_dismissal_cookie(&self, dismissal: &BroadcastMessageDismissal) -> Result<()> {
        self.cookies.set(
            &dismissal.cookie_key,
            "true".to_string(),
            dismissal.expires_at,
        );

        Ok(())
    }

    fn get_dismissal_finder(&self) -> Arc<dyn BroadcastMessageDismissalFinder> {
        self.dismissal_finder.clone()
    }

    fn get_cookies(&self) -> Arc<dyn Cookie> {
        self.cookies.clone()
    }
}
