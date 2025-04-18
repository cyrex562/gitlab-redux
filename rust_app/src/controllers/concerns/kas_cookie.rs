use crate::config::Settings;
use crate::gitlab::kas::{Kas, UserAccess};
use crate::session::Session;
use actix_web::{web, HttpRequest, HttpResponse};
use cookie::Cookie;

pub const KAS_COOKIE_KEY: &str = "gitlab_kas";

pub trait KasCookie {
    fn set_kas_cookie(&self, session: &Session) -> Option<Cookie<'static>>;
    fn kas_url(&self) -> String;
    fn kas_ws_url(&self) -> String;
}

pub struct KasCookieImpl {
    settings: Settings,
    kas: Kas,
}

impl KasCookieImpl {
    pub fn new(settings: Settings, kas: Kas) -> Self {
        Self { settings, kas }
    }
}

impl KasCookie for KasCookieImpl {
    fn set_kas_cookie(&self, session: &Session) -> Option<Cookie<'static>> {
        if !UserAccess::is_enabled() {
            return None;
        }

        let public_session_id = session.public_id()?;
        let cookie_data = UserAccess::cookie_data(&public_session_id)?;

        Some(
            Cookie::build(KAS_COOKIE_KEY, cookie_data)
                .secure(true)
                .http_only(true)
                .path("/")
                .finish(),
        )
    }

    fn kas_url(&self) -> String {
        self.kas.tunnel_url()
    }

    fn kas_ws_url(&self) -> String {
        self.kas.tunnel_ws_url()
    }
}

// Content Security Policy implementation
impl KasCookieImpl {
    pub fn content_security_policy(&self) -> Option<String> {
        if !UserAccess::is_enabled() || !self.settings.content_security_policy_enabled {
            return None;
        }

        let kas_url = self.kas_url();
        let kas_ws_url = self.kas_ws_url();

        // Only add CSP if KAS is on a different host
        if kas_url
            .parse::<url::Url>()
            .ok()
            .and_then(|url| url.host_str().map(|h| h.to_string()))
            .map_or(false, |host| host == self.settings.gitlab_host)
        {
            return None;
        }

        Some(format!(
            "connect-src 'self' {} {};",
            kas_ws_url.trim_end_matches('/'),
            kas_url.trim_end_matches('/')
        ))
    }
}
