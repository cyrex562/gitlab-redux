// Ported from orig_app/app/controllers/concerns/hotlink_interceptor.rb
use actix_web::{web, HttpRequest, HttpResponse, Responder};
use serde::{Deserialize, Serialize};
use std::collections::HashSet;
use url::Url;

pub trait HotlinkInterceptor {
    fn intercept_hotlink(&self, req: &HttpRequest) -> impl Responder;
}

pub struct HotlinkInterceptorImpl {
    allowed_domains: HashSet<String>,
    asset_path: String,
}

impl HotlinkInterceptorImpl {
    pub fn new(allowed_domains: HashSet<String>, asset_path: String) -> Self {
        Self {
            allowed_domains,
            asset_path,
        }
    }

    fn is_allowed_referer(&self, referer: Option<&str>) -> bool {
        if let Some(referer) = referer {
            if let Ok(url) = Url::parse(referer) {
                if let Some(host) = url.host_str() {
                    return self.allowed_domains.contains(host);
                }
            }
        }
        false
    }

    fn is_asset_request(&self, path: &str) -> bool {
        path.starts_with(&self.asset_path)
    }
}

impl HotlinkInterceptor for HotlinkInterceptorImpl {
    fn intercept_hotlink(&self, req: &HttpRequest) -> impl Responder {
        let path = req.path();

        if !self.is_asset_request(path) {
            return HttpResponse::Ok().finish();
        }

        let referer = req.headers().get("referer").and_then(|h| h.to_str().ok());

        if self.is_allowed_referer(referer) {
            HttpResponse::Ok().finish()
        } else {
            HttpResponse::Forbidden().body("Hotlinking not allowed")
        }
    }
}
