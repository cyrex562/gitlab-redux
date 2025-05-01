// Ported from: orig_app/app/controllers/concerns/synchronize_broadcast_message_dismissals.rb
//
// Synchronize broadcast message dismissals for a user by setting cookies for each dismissal if not already present.
//
use crate::models::broadcast_message::BroadcastMessageDismissal;
use crate::models::user::User;
use actix_web::{
    cookie::{Cookie, CookieBuilder},
    http::header,
    HttpRequest, HttpResponse, Result,
};
use chrono::{DateTime, Utc};

pub trait SynchronizeBroadcastMessageDismissals {
    fn synchronize_broadcast_message_dismissals(
        &self,
        user: &User,
        req: &HttpRequest,
        resp: &mut HttpResponse,
    );
}

pub struct BroadcastMessageDismissalFinder;

impl BroadcastMessageDismissalFinder {
    pub fn new() -> Self {
        BroadcastMessageDismissalFinder
    }

    pub fn execute(&self, user: &User) -> Vec<BroadcastMessageDismissal> {
        // TODO: Replace with actual DB lookup
        BroadcastMessageDismissal::find_for_user(user)
    }
}

impl SynchronizeBroadcastMessageDismissals for BroadcastMessageDismissalFinder {
    fn synchronize_broadcast_message_dismissals(
        &self,
        user: &User,
        req: &HttpRequest,
        resp: &mut HttpResponse,
    ) {
        let dismissals = self.execute(user);
        for dismissal in dismissals {
            let cookie_key = dismissal.cookie_key();
            let has_cookie = req.cookie(&cookie_key).is_some();
            if !has_cookie {
                create_dismissal_cookie(resp, &dismissal);
            }
        }
    }
}

fn create_dismissal_cookie(resp: &mut HttpResponse, dismissal: &BroadcastMessageDismissal) {
    let expires = dismissal.expires_at();
    let cookie = CookieBuilder::new(dismissal.cookie_key(), "true")
        .expires(expires)
        .http_only(true)
        .finish();
    let _ = resp.add_cookie(&cookie);
}
