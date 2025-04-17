use actix_web::{web, HttpResponse, Responder};
use serde::Deserialize;

use crate::models::snippet::Snippet;
use crate::models::blob::Blob;
use crate::models::repository::Repository;
use crate::utils::extracts_ref::RefExtractor;
use crate::controllers::concerns::snippets::SendBlob;
use crate::utils::strong_memoize::StrongMemoize;
use crate::utils::logger::AppLogger;

#[derive(Debug, Deserialize)]
pub struct BlobParams {
    id: Option<String>,
    ref_name: Option<String>,
    path: Option<String>,
    ref_type: Option<String>,
}

pub trait BlobsActions: SendBlob + StrongMemoize {
    fn snippet(&self) -> &Snippet;
    fn set_snippet(&mut self, snippet: Snippet);
    fn authorize_read_snippet(&self) -> Result<(), HttpResponse>;
    fn ensure_repository(&self) -> Result<(), HttpResponse>;
    fn ensure_blob(&self) -> Result<(), HttpResponse>;
    fn snippet_id(&self) -> i32;
    fn params(&self) -> &BlobParams;
    fn logger(&self) -> &AppLogger;

    fn raw(&self) -> impl Responder {
        if let Err(response) = self.authorize_read_snippet() {
            return response;
        }

        if let Err(response) = self.ensure_repository() {
            return response;
        }

        if let Err(response) = self.ensure_blob() {
            return response;
        }

        self.send_snippet_blob(self.snippet(), self.blob())
    }

    fn blob(&self) -> &Blob {
        self.strong_memoize("blob", || {
            let ref_extractor = RefExtractor::new(
                self.snippet(),
                self.params(),
            );
            
            if let Ok(commit) = ref_extractor.extract() {
                if let Some(path) = &self.params().path {
                    self.snippet().repository().blob_at(commit.id(), path)
                } else {
                    None
                }
            } else {
                None
            }
        })
    }

    fn ensure_blob(&self) -> Result<(), HttpResponse> {
        if self.blob().is_some() {
            Ok(())
        } else {
            Err(HttpResponse::NotFound().finish())
        }
    }

    fn ensure_repository(&self) -> Result<(), HttpResponse> {
        if self.snippet().repo_exists() {
            Ok(())
        } else {
            self.logger().error(
                "Snippet raw blob attempt with no repo",
                &[("snippet_id", &self.snippet_id())],
            );
            
            Err(HttpResponse::UnprocessableEntity().finish())
        }
    }
} 