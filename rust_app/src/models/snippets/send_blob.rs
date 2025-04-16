use actix_web::{web, HttpResponse};
use crate::snippets::{Snippet, Blob};
use crate::models::user::User;

pub trait SendBlob {
    fn send_snippet_blob(&self, snippet: &Snippet, blob: &Blob, user: &User) -> HttpResponse {
        // TODO: Implement workhorse_set_content_type
        let content_disposition = self.content_disposition();
        
        // TODO: Implement send_blob with proper caching and inline handling
        HttpResponse::Ok()
            .content_type("application/octet-stream")
            .header("Content-Disposition", format!("{}; filename=\"{}\"", 
                content_disposition, 
                blob.path))
            .body(blob.content.clone())
    }

    fn content_disposition(&self) -> String {
        // TODO: Get from request params
        "inline".to_string()
    }
} 