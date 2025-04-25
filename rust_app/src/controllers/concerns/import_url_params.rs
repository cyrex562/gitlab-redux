// Ported from: orig_app/app/controllers/concerns/import_url_params.rb
//
// This file provides ImportUrlParams trait and implementation for handling import URL parameters.
//
// See also: orig_app/lib/gitlab/url_sanitizer.rb for the original Ruby URL sanitizer logic.
//
// ---
use actix_web::web;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

use crate::utils::url_sanitizer::UrlSanitizer;

#[derive(Debug, Serialize, Deserialize)]
pub struct ProjectParams {
    pub import_url: Option<String>,
    pub import_url_user: Option<String>,
    pub import_url_password: Option<String>,
}

pub trait ImportUrlParams {
    fn import_url_params(&self, params: &ProjectParams) -> HashMap<String, String>;
    fn import_params_to_full_url(&self, params: &ProjectParams) -> String;
}

pub struct ImportUrlParamsImpl;

impl ImportUrlParamsImpl {
    pub fn new() -> Self {
        Self
    }
}

impl ImportUrlParams for ImportUrlParamsImpl {
    fn import_url_params(&self, params: &ProjectParams) -> HashMap<String, String> {
        let mut result = HashMap::new();

        if let Some(import_url) = &params.import_url {
            result.insert(
                "import_url".to_string(),
                self.import_params_to_full_url(params),
            );
            // We need to set import_type because attempting to retry an import by URL
            // could leave a stale value around. This would erroneously cause an importer
            // (e.g. import/export) to run.
            result.insert("import_type".to_string(), "git".to_string());
        }

        result
    }

    fn import_params_to_full_url(&self, params: &ProjectParams) -> String {
        let credentials =
            if params.import_url_user.is_some() || params.import_url_password.is_some() {
                Some((
                    params.import_url_user.clone().unwrap_or_default(),
                    params.import_url_password.clone().unwrap_or_default(),
                ))
            } else {
                None
            };

        UrlSanitizer::new(params.import_url.as_deref().unwrap_or(""), credentials).full_url()
    }
}
