// Ported from: orig_app/app/controllers/concerns/known_sign_in.rb
// This file implements the KnownSignIn concern in Rust, providing methods for verifying known sign-ins, checking known devices and IPs, updating cookies, and notifying users of unknown sign-ins.
// Ported on: 2025-04-25 13:31:44

use crate::models::user::User;
use crate::services::notification::NotificationService;
use crate::utils::visitor_location::VisitorLocation;
use actix_web::{web, HttpRequest};
use chrono::{DateTime, Duration, Utc};
use cookie::Cookie;
use std::collections::HashSet;

pub const KNOWN_SIGN_IN_COOKIE: &str = "known_sign_in";
const KNOWN_SIGN_IN_COOKIE_EXPIRY_DAYS: i64 = 14;

pub trait KnownSignIn {
    fn verify_known_sign_in(&self, user: &User, request: &HttpRequest) -> Option<Cookie<'static>>;
    fn known_remote_ip(&self, user: &User, request: &HttpRequest) -> bool;
    fn known_device(&self, user: &User, request: &HttpRequest) -> bool;
    fn update_cookie(&self, user: &User) -> Cookie<'static>;
    fn notify_user(&self, user: &User, request: &HttpRequest);
}

pub struct KnownSignInImpl {
    notification_service: NotificationService,
    settings: web::Data<crate::config::Settings>,
}

impl KnownSignInImpl {
    pub fn new(
        notification_service: NotificationService,
        settings: web::Data<crate::config::Settings>,
    ) -> Self {
        Self {
            notification_service,
            settings,
        }
    }
}

impl KnownSignIn for KnownSignInImpl {
    fn verify_known_sign_in(&self, user: &User, request: &HttpRequest) -> Option<Cookie<'static>> {
        if !self.settings.notify_on_unknown_sign_in || user.is_none() {
            return None;
        }

        if !self.known_device(user, request) && !self.known_remote_ip(user, request) {
            self.notify_user(user, request);
        }

        Some(self.update_cookie(user))
    }

    fn known_remote_ip(&self, user: &User, request: &HttpRequest) -> bool {
        let known_ips: HashSet<String> = user
            .sessions()
            .iter()
            .filter(|s| !s.is_impersonated())
            .map(|s| s.ip_address.clone())
            .chain(std::iter::once(user.last_sign_in_ip.clone()))
            .collect();

        known_ips.contains(
            &request
                .connection_info()
                .peer_addr()
                .unwrap_or_default()
                .to_string(),
        )
    }

    fn known_device(&self, user: &User, request: &HttpRequest) -> bool {
        request
            .cookies()
            .get(KNOWN_SIGN_IN_COOKIE)
            .and_then(|c| c.value().parse::<i64>().ok())
            .map_or(false, |id| id == user.id)
    }

    fn update_cookie(&self, user: &User) -> Cookie<'static> {
        let expiry = Utc::now() + Duration::days(KNOWN_SIGN_IN_COOKIE_EXPIRY_DAYS);

        Cookie::build(KNOWN_SIGN_IN_COOKIE, user.id.to_string())
            .secure(true)
            .http_only(true)
            .expires(expiry)
            .path("/")
            .finish()
    }

    fn notify_user(&self, user: &User, request: &HttpRequest) {
        let visitor_location = VisitorLocation::new(request);
        let remote_ip = request
            .connection_info()
            .peer_addr()
            .unwrap_or_default()
            .to_string();
        let current_sign_in_at = Utc::now();

        self.notification_service.unknown_sign_in(
            user,
            remote_ip,
            current_sign_in_at,
            visitor_location,
        );
    }
}
