// Ported from: orig_app/app/controllers/concerns/continue_params.rb
use actix_web::web::Query;
use actix_web::HttpRequest;
use serde::Deserialize;
use std::collections::HashMap;

use super::internal_redirect::InternalRedirect;

#[derive(Debug, Deserialize)]
pub struct ContinueParamsData {
    pub to: Option<String>,
    pub notice: Option<String>,
    pub notice_now: Option<String>,
}

pub trait ContinueParams: InternalRedirect {
    fn continue_params(&self, req: &HttpRequest) -> HashMap<String, String> {
        // Try to extract the 'continue' query params
        let query = Query::<HashMap<String, String>>::from_query(req.query_string());
        if let Ok(params) = query {
            if let Some(continue_param) = params.get("continue") {
                // Parse the inner continue param as a query string
                let inner_query = Query::<ContinueParamsData>::from_query(continue_param);
                if let Ok(inner) = inner_query {
                    let mut map = HashMap::new();
                    if let Some(to) = &inner.to {
                        if let Some(safe_to) = self.safe_redirect_path(Some(to)) {
                            map.insert("to".to_string(), safe_to);
                        }
                    }
                    if let Some(notice) = &inner.notice {
                        map.insert("notice".to_string(), notice.clone());
                    }
                    if let Some(notice_now) = &inner.notice_now {
                        map.insert("notice_now".to_string(), notice_now.clone());
                    }
                    return map;
                }
            }
        }
        HashMap::new()
    }
}
