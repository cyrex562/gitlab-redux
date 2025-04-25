// Ported from orig_app/app/controllers/concerns/snippets/send_blob.rb
// Provides send_snippet_blob for sending snippet blobs with correct headers and permissions.

use crate::models::snippets::sends_blob::SendsBlob;
use crate::models::snippets::{Blob, Snippet};
use actix_web::{HttpRequest, HttpResponse};

pub trait SendBlob: SendsBlob {
    fn request(&self) -> &HttpRequest;

    fn send_snippet_blob(&self, snippet: &Snippet, blob: &Blob) -> HttpResponse {
        self.workhorse_set_content_type();
        let inline = self.content_disposition() == "inline";
        let allow_caching = self.can_cache_blob(snippet);
        self.send_blob(&snippet.repository, blob, inline, allow_caching)
    }

    fn workhorse_set_content_type(&self) {
        // TODO: Set content type header for workhorse if needed
        // This is a stub for the workhorse_set_content_type! Ruby call
    }

    fn content_disposition(&self) -> String {
        let params = self.request().query_string();
        // Parse query string for "inline" param
        let inline_param = self.request().match_info().get("inline");
        match inline_param {
            Some("false") => "attachment".to_string(),
            _ => "inline".to_string(),
        }
    }

    fn can_cache_blob(&self, snippet: &Snippet) -> bool {
        // TODO: Implement permission check for caching
        // Placeholder: always false
        false
    }
}
