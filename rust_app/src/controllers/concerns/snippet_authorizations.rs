use crate::models::snippet::Snippet;
use crate::models::user::User;
use actix_web::{web, HttpRequest, HttpResponse};
use std::sync::Arc;

/// This trait provides authorization methods for snippet-related actions
pub trait SnippetAuthorizations {
    /// Authorize reading a snippet
    fn authorize_read_snippet(
        &self,
        snippet: &Snippet,
        current_user: Option<&User>,
    ) -> Result<(), HttpResponse>;

    /// Authorize updating a snippet
    fn authorize_update_snippet(
        &self,
        snippet: &Snippet,
        current_user: Option<&User>,
    ) -> Result<(), HttpResponse>;

    /// Authorize creating a snippet
    fn authorize_create_snippet(&self, current_user: Option<&User>) -> Result<(), HttpResponse>;
}

pub struct SnippetAuthorizationsHandler;

impl SnippetAuthorizationsHandler {
    pub fn new() -> Self {
        SnippetAuthorizationsHandler
    }
}

impl SnippetAuthorizations for SnippetAuthorizationsHandler {
    fn authorize_read_snippet(
        &self,
        snippet: &Snippet,
        current_user: Option<&User>,
    ) -> Result<(), HttpResponse> {
        if snippet.public || current_user.map_or(false, |user| user.can_read_snippet(snippet)) {
            Ok(())
        } else {
            Err(HttpResponse::NotFound().finish())
        }
    }

    fn authorize_update_snippet(
        &self,
        snippet: &Snippet,
        current_user: Option<&User>,
    ) -> Result<(), HttpResponse> {
        if current_user.map_or(false, |user| user.can_update_snippet(snippet)) {
            Ok(())
        } else {
            Err(HttpResponse::NotFound().finish())
        }
    }

    fn authorize_create_snippet(&self, current_user: Option<&User>) -> Result<(), HttpResponse> {
        if current_user.is_some() {
            Ok(())
        } else {
            Err(HttpResponse::NotFound().finish())
        }
    }
}

// These would be implemented in separate modules
pub mod models {
    pub mod snippet {
        pub struct Snippet {
            pub public: bool,
            // Add other fields as needed
        }
    }

    pub mod user {
        use super::snippet::Snippet;

        pub struct User {
            // Add fields as needed
        }

        impl User {
            pub fn can_read_snippet(&self, snippet: &Snippet) -> bool {
                // Implement permission check
                true
            }

            pub fn can_update_snippet(&self, snippet: &Snippet) -> bool {
                // Implement permission check
                true
            }
        }
    }
}
