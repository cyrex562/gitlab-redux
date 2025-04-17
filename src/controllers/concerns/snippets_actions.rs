use actix_web::{web, HttpResponse, Responder};
use serde::Deserialize;

use crate::controllers::concerns::{
    NoteableMetadata, PaginatedCollection, ProductAnalyticsTracking, RendersBlob, RendersNotes,
    SnippetsSort,
};
use crate::models::blob::Blob;
use crate::models::discussion::Discussion;
use crate::models::note::Note;
use crate::models::snippet::Snippet;

pub trait SnippetsActions:
    RendersNotes
    + RendersBlob
    + PaginatedCollection
    + NoteableMetadata
    + SnippetsSort
    + ProductAnalyticsTracking
{
    fn snippet(&self) -> &Snippet;
    fn set_snippet(&mut self, snippet: Snippet);
    fn is_js_request(&self) -> bool;
    fn workhorse_set_content_type(&self);
    fn convert_line_endings(&self, content: &str, line_ending: &str) -> String;
    fn send_snippet_blob(&self, snippet: &Snippet, blob: &Blob) -> HttpResponse;
    fn sanitized_file_name(&self, name: &str) -> String;
    fn content_disposition(&self) -> String;

    fn edit(&self) -> impl Responder {
        // In a real implementation, we'd need to handle the view rendering
        HttpResponse::Ok().finish()
    }

    fn raw(&self) -> impl Responder {
        self.workhorse_set_content_type();

        let blob = self.blob();

        if let Some(snippet) = blob.snippet() {
            let data = self.convert_line_endings(&blob.data, "raw");
            let filename = self.sanitized_file_name(&blob.name);

            HttpResponse::Ok()
                .content_type("text/plain; charset=utf-8")
                .header("Content-Disposition", self.content_disposition())
                .body(data)
        } else {
            self.send_snippet_blob(self.snippet(), blob)
        }
    }

    fn show(&self) -> impl Responder {
        if self.is_js_request() {
            if self.snippet().embeddable() {
                self.conditionally_expand_blobs(self.blobs());
                // In a real implementation, we'd need to handle the view rendering
                HttpResponse::Ok().finish()
            } else {
                HttpResponse::NotFound().finish()
            }
        } else {
            let note = Note::new(self.snippet(), self.snippet().project());
            let discussions = self.snippet().discussions();
            let notes =
                self.prepare_notes_for_rendering(discussions.iter().flat_map(|d| d.notes()));

            // In a real implementation, we'd need to handle the view rendering
            HttpResponse::Ok().finish()
        }
    }

    fn blob(&self) -> &Blob {
        self.blobs().first().expect("No blobs found")
    }

    fn blobs(&self) -> Vec<&Blob> {
        if self.snippet().empty_repo() {
            vec![self.snippet().blob()]
        } else {
            self.snippet().blobs()
        }
    }
}
