use crate::config::GitlabConfig;
use crate::controllers::ApplicationController;
use crate::utils::url_sanitizer::UrlSanitizer;
use actix_web::{web, HttpResponse};
use url::Url;

pub struct ExternalRedirectController {
    base: ApplicationController,
}

impl ExternalRedirectController {
    pub fn new(base: ApplicationController) -> Self {
        Self { base }
    }

    pub async fn index(&self, query: web::Query<RedirectQuery>) -> HttpResponse {
        if !self.should_handle_url(&query.url) {
            return HttpResponse::NotFound().finish();
        }

        if self.known_url(&query.url) {
            return HttpResponse::Found()
                .header("Location", &query.url)
                .finish();
        }

        // TODO: Implement fullscreen layout rendering
        HttpResponse::Ok().json(RedirectData {
            url: query.url.clone(),
            rel: query.rel.clone(),
        })
    }

    fn relme_keywords(&self, query: &RedirectQuery) -> Option<String> {
        query.rel.clone().map(|r| r.trim().to_string())
    }

    fn url_param(&self, query: &RedirectQuery) -> Option<String> {
        query.url.clone().map(|u| u.trim().to_string())
    }

    fn known_url(&self, url: &Option<String>) -> bool {
        if let Some(url_str) = url {
            if let Ok(uri) = Url::parse(url_str) {
                if let Some(host) = uri.host_str() {
                    return host == GitlabConfig::gitlab_url().host().unwrap_or("");
                }
            }
        }
        false
    }

    fn should_handle_url(&self, url: &Option<String>) -> bool {
        if let Some(url_str) = url {
            // To avoid lots of redirects, don't allow url to point to self
            UrlSanitizer::valid_web(url_str) && !url_str.starts_with(&self.base.request.base_url())
        } else {
            false
        }
    }
}

#[derive(serde::Deserialize)]
pub struct RedirectQuery {
    url: Option<String>,
    rel: Option<String>,
}

struct RedirectData {
    url: Option<String>,
    rel: Option<String>,
}
