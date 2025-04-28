// Ported from: orig_app/app/controllers/concerns/render_access_tokens.rb
//
// This module provides functions for rendering access tokens, including active/inactive tokens,
// pagination headers, and ICS calendar export for expiries.

use actix_web::{web, HttpRequest, HttpResponse};
use chrono::NaiveDate;
use serde::{Deserialize, Serialize};
use std::sync::Arc;

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct AccessToken {
    pub id: i32,
    pub name: String,
    pub scopes: Vec<String>,
    pub expires_at: Option<String>,
    pub created_at: String,
    pub state: String, // "active" or "inactive"
}

pub struct RenderAccessTokensHandler {
    pub current_user: Option<Arc<User>>,
}

impl RenderAccessTokensHandler {
    pub fn new(current_user: Option<Arc<User>>) -> Self {
        RenderAccessTokensHandler { current_user }
    }

    pub fn active_access_tokens(&self, page: usize, per_page: usize) -> (Vec<AccessToken>, usize) {
        let tokens = self.fetch_access_tokens_by_state("active");
        let size = tokens.len();
        let paged = tokens
            .into_iter()
            .skip((page - 1) * per_page)
            .take(per_page)
            .collect();
        (paged, size)
    }

    pub fn inactive_access_tokens(&self) -> Vec<AccessToken> {
        self.fetch_access_tokens_by_state("inactive")
    }

    pub fn add_pagination_headers(
        &self,
        resp: &mut HttpResponse,
        page: usize,
        per_page: usize,
        total: usize,
    ) {
        let total_pages = (total + per_page - 1) / per_page;
        resp.headers_mut()
            .insert("X-Page", page.to_string().parse().unwrap());
        resp.headers_mut()
            .insert("X-Per-Page", per_page.to_string().parse().unwrap());
        resp.headers_mut()
            .insert("X-Total", total.to_string().parse().unwrap());
        resp.headers_mut()
            .insert("X-Total-Pages", total_pages.to_string().parse().unwrap());
    }

    pub fn page(&self, req: &HttpRequest) -> usize {
        req.query_string()
            .split('&')
            .find_map(|kv| {
                let mut split = kv.split('=');
                if split.next()? == "page" {
                    split.next()?.parse().ok()
                } else {
                    None
                }
            })
            .unwrap_or(1)
    }

    pub fn expiry_ics(&self, tokens: &[AccessToken]) -> String {
        let mut ical = String::from("BEGIN:VCALENDAR\nVERSION:2.0\n");
        for token in tokens {
            if let Some(ref expires_at) = token.expires_at {
                if let Ok(date) = NaiveDate::parse_from_str(expires_at, "%Y-%m-%d") {
                    let date_str = date.format("%Y%m%d").to_string();
                    ical.push_str(&format!(
                        "BEGIN:VEVENT\nDTSTART:{}\nDTEND:{}\nSUMMARY:Token {} expires today\nEND:VEVENT\n",
                        date_str, date_str, token.name
                    ));
                }
            }
        }
        ical.push_str("END:VCALENDAR\n");
        ical
    }

    fn fetch_access_tokens_by_state(&self, state: &str) -> Vec<AccessToken> {
        // TODO: Replace with real DB query
        Vec::new()
    }
}

pub struct User {
    pub id: i32,
    // ...
}
