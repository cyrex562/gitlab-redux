// Ported from orig_app/app/controllers/concerns/snippets/blobs_actions.rb
// Provides actions for handling snippet blobs (raw download, existence checks, etc.)

use crate::models::snippets::send_blob::SendBlob;
use crate::models::snippets::RefExtractor;
use crate::models::snippets::{Blob, Snippet};
use actix_web::{HttpRequest, HttpResponse};

pub trait SnippetsBlobsActions: SendBlob {
    // Should return a reference to the current snippet
    fn snippet(&self) -> &Snippet;
    // Should return a mutable reference to the current request (for param extraction)
    fn request(&self) -> &HttpRequest;

    // Main entrypoint for GET /raw
    fn raw(&self) -> HttpResponse {
        // Authorization and existence checks are assumed to be handled by middleware or wrappers
        match (self.ensure_repository(), self.ensure_blob()) {
            (Ok(_), Ok(blob)) => self.send_snippet_blob(self.snippet(), &blob),
            (Err(resp), _) | (_, Err(resp)) => resp,
        }
    }

    // Extract the blob from the request parameters
    fn blob(&self) -> Option<Blob> {
        let params = self.request().match_info();
        let id = params.get("id");
        let ref_name = params.get("ref");
        let path = params.get("path");
        let ref_type = params.get("ref_type");
        let ref_extractor = RefExtractor {
            id: id.map(|s| s.to_string()),
            ref_name: ref_name.map(|s| s.to_string()),
            path: path.map(|s| s.to_string()),
            ref_type: ref_type.map(|s| s.to_string()),
        };
        // In a real implementation, extract commit and get blob from repo
        if let Some(commit_id) = &ref_extractor.id {
            // TODO: Lookup blob in repository using commit_id and path
            Some(Blob {
                path: ref_extractor.path.clone().unwrap_or_default(),
                content: Vec::new(), // Placeholder
            })
        } else {
            None
        }
    }

    // Ensure the blob exists, or return a 404 response
    fn ensure_blob(&self) -> Result<Blob, HttpResponse> {
        self.blob().ok_or_else(|| HttpResponse::NotFound().finish())
    }

    // Ensure the repository exists, or return a 422 response
    fn ensure_repository(&self) -> Result<(), HttpResponse> {
        if self.snippet().repo_exists() {
            Ok(())
        } else {
            // Log error in real implementation
            HttpResponse::UnprocessableEntity().finish().into()
        }
    }

    // Optionally: extract snippet_id from params
    fn snippet_id(&self) -> Option<String> {
        self.request()
            .match_info()
            .get("snippet_id")
            .map(|s| s.to_string())
    }
}
