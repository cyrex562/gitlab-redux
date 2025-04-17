use actix_web::web;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

#[derive(Debug, Serialize, Deserialize)]
pub struct ContinueParams {
    pub to: Option<String>,
    pub notice: Option<String>,
    pub notice_now: Option<String>,
}

pub trait ContinueParams {
    fn continue_params(&self, req: &web::HttpRequest) -> HashMap<String, String> {
        let params: web::Query<HashMap<String, String>> =
            web::Query::from_query(req.query_string()).unwrap_or_default();
        let continue_params = params.get("continue").cloned().unwrap_or_default();

        let mut result = HashMap::new();
        if let Some(to) = continue_params.get("to") {
            if let Some(safe_path) = self.safe_redirect_path(Some(to)) {
                result.insert("to".to_string(), safe_path);
            }
        }
        if let Some(notice) = continue_params.get("notice") {
            result.insert("notice".to_string(), notice.clone());
        }
        if let Some(notice_now) = continue_params.get("notice_now") {
            result.insert("notice_now".to_string(), notice_now.clone());
        }

        result
    }

    // Required method to be implemented by concrete types
    fn safe_redirect_path(&self, path: Option<&str>) -> Option<String>;
}
