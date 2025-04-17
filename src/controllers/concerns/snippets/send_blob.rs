use actix_web::{web, HttpResponse};
use serde::Deserialize;

use crate::models::snippet::Snippet;
use crate::models::blob::Blob;
use crate::models::repository::Repository;
use crate::models::user::User;
use crate::controllers::concerns::SendsBlob;

#[derive(Debug, Deserialize)]
pub struct SendBlobParams {
    inline: Option<bool>,
}

pub trait SendBlob: SendsBlob {
    fn workhorse_set_content_type(&self);
    fn params(&self) -> &SendBlobParams;
    fn current_user(&self) -> Option<&User>;

    fn send_snippet_blob(&self, snippet: &Snippet, blob: &Blob) -> HttpResponse {
        self.workhorse_set_content_type();

        let inline = self.content_disposition() == "inline";
        let allow_caching = self.current_user()
            .map(|user| user.can_cache_blob(snippet))
            .unwrap_or(false);

        self.send_blob(
            snippet.repository(),
            blob,
            inline,
            allow_caching,
        )
    }

    fn content_disposition(&self) -> &'static str {
        if self.params().inline == Some(false) {
            "attachment"
        } else {
            "inline"
        }
    }
} 