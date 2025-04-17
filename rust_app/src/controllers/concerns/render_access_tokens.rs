use actix_web::{HttpRequest, HttpResponse};
use chrono::NaiveDate;
use icalendar::{Calendar, Component, Event};
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct AccessToken {
    pub name: String,
    pub expires_at: String,
    // Add other fields as needed
}

pub trait RenderAccessTokens {
    fn active_access_tokens(&self, req: &HttpRequest) -> (Vec<AccessToken>, usize) {
        let page = self.page(req);
        let tokens = self.find_tokens("active", "expires_asc");
        let size = tokens.len();

        let paginated_tokens = self.paginate_tokens(tokens, page);
        self.add_pagination_headers(req, &paginated_tokens);

        (paginated_tokens, size)
    }

    fn inactive_access_tokens(&self) -> Vec<AccessToken> {
        self.find_tokens("inactive", "updated_at_desc")
    }

    fn add_pagination_headers(&self, req: &HttpRequest, tokens: &[AccessToken]) {
        // Implementation would depend on your pagination system
        // This is a placeholder for the pagination header logic
    }

    fn page(&self, req: &HttpRequest) -> i32 {
        // Extract page from query parameters, default to 1
        req.query_string()
            .split('&')
            .find(|param| param.starts_with("page="))
            .and_then(|param| param.split('=').nth(1))
            .and_then(|page| page.parse().ok())
            .unwrap_or(1)
    }

    fn expiry_ics(&self, tokens: &[AccessToken]) -> String {
        let mut cal = Calendar::new();

        for token in tokens {
            let mut event = Event::new();

            // Parse the date string and create an event
            if let Ok(date) = NaiveDate::parse_from_str(&token.expires_at, "%Y-%m-%d") {
                event.summary(&format!("Token {} expires today", token.name));
                event.dtstart(date);
                event.dtend(date);
                cal.push(event);
            }
        }

        cal.to_string()
    }

    // Required methods to be implemented by concrete types
    fn find_tokens(&self, state: &str, sort: &str) -> Vec<AccessToken>;
    fn paginate_tokens(&self, tokens: Vec<AccessToken>, page: i32) -> Vec<AccessToken>;
}
