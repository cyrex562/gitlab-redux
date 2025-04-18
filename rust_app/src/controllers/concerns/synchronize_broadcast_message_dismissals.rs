use crate::models::broadcast_message::BroadcastMessageDismissal;
use crate::models::user::User;
use actix_web::{web, HttpRequest, HttpResponse, Result};
use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct DismissalCookie {
    value: bool,
    expires_at: DateTime<Utc>,
}

pub trait SynchronizeBroadcastMessageDismissals {
    fn synchronize_broadcast_message_dismissals(
        &self,
        user: &User,
        req: &HttpRequest,
    ) -> Result<(), Box<dyn std::error::Error>>;
}

pub struct BroadcastMessageDismissalFinder {
    user: User,
}

impl BroadcastMessageDismissalFinder {
    pub fn new(user: User) -> Self {
        Self { user }
    }

    pub fn execute(&self) -> Vec<BroadcastMessageDismissal> {
        // This would typically query your database for dismissals
        // For now, returning an empty vector as placeholder
        Vec::new()
    }
}

impl<T> SynchronizeBroadcastMessageDismissals for T
where
    T: Fn(&User, &HttpRequest) -> Result<(), Box<dyn std::error::Error>>,
{
    fn synchronize_broadcast_message_dismissals(
        &self,
        user: &User,
        req: &HttpRequest,
    ) -> Result<(), Box<dyn std::error::Error>> {
        let finder = BroadcastMessageDismissalFinder::new(user.clone());
        let dismissals = finder.execute();

        for dismissal in dismissals {
            if !req.cookie(dismissal.cookie_key()).is_some() {
                create_dismissal_cookie(req, &dismissal)?;
            }
        }

        Ok(())
    }
}

fn create_dismissal_cookie(
    req: &HttpRequest,
    dismissal: &BroadcastMessageDismissal,
) -> Result<(), Box<dyn std::error::Error>> {
    let cookie = DismissalCookie {
        value: true,
        expires_at: dismissal.expires_at(),
    };

    // In a real implementation, this would set the cookie in the response
    // For now, we'll just log that we would set it
    println!(
        "Setting cookie {} with value {:?}",
        dismissal.cookie_key(),
        cookie
    );

    Ok(())
}
