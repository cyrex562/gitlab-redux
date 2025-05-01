// Ported from orig_app/app/controllers/concerns/sends_blob.rb on 2025-04-29
//
// This file implements the SendsBlob concern logic from the original Ruby code.
//

use crate::controllers::concerns::send_file_upload::SendFileUpload;
use crate::models::blob::Blob;
use crate::models::lfs_object::LfsObject;
use crate::models::project::Project;
use crate::models::repository::Repository;
use actix_web::{web, HttpRequest, HttpResponse, Responder};
use std::sync::Arc;
use std::time::{Duration, SystemTime};

pub trait SendsBlob {
    fn send_blob(
        &self,
        req: &HttpRequest,
        repository: &Repository,
        blob: Option<&Blob>,
        inline: bool,
        allow_caching: bool,
    ) -> impl Responder;
}

pub struct SendsBlobHandler {
    send_file_upload: Arc<dyn SendFileUpload>,
}

impl SendsBlobHandler {
    pub fn new(send_file_upload: Arc<dyn SendFileUpload>) -> Self {
        SendsBlobHandler { send_file_upload }
    }

    fn cached_blob(
        &self,
        req: &HttpRequest,
        blob: &Blob,
        allow_caching: bool,
        ref_name: Option<&str>,
        commit_id: Option<&str>,
    ) -> bool {
        // Check if the blob is already cached
        let etag = blob.id();

        // Check if the request has a matching ETag
        if let Some(if_none_match) = req.headers().get("if-none-match") {
            if if_none_match
                .to_str()
                .ok()
                .map_or(false, |etag_str| etag_str == etag)
            {
                return true;
            }
        }

        // Determine cache time based on whether this is a commit SHA or a branch/tag
        let cache_time =
            if ref_name.is_some() && commit_id.is_some() && ref_name.unwrap() == commit_id.unwrap()
            {
                // This is a link to a commit by its commit SHA, so it's immutable
                Blob::CACHE_TIME_IMMUTABLE
            } else {
                // A branch or tag points at this blob, so it may change over time
                Blob::CACHE_TIME
            };

        // Set cache headers
        let mut response = HttpResponse::Ok();
        response.append_header(("ETag", etag));
        response.append_header(("Cache-Control", format!(
            "max-age={}, {}must-revalidate, stale-if-error=300, stale-while-revalidate=60, s-maxage=60",
            cache_time.as_secs(),
            if allow_caching { "public, " } else { "private, " }
        )));

        // Calculate expiration time
        let expires_at = SystemTime::now() + cache_time;
        response.append_header((
            "Expires",
            format!(
                "{}",
                expires_at
                    .duration_since(SystemTime::UNIX_EPOCH)
                    .unwrap()
                    .as_secs()
            ),
        ));

        false
    }

    fn send_lfs_object(&self, req: &HttpRequest, blob: &Blob, project: &Project) -> impl Responder {
        if let Some(lfs_object) = self.find_lfs_object(blob) {
            if lfs_object.project_allowed_access(project) {
                return self.send_file_upload.send_upload(
                    req,
                    &lfs_object.file,
                    None,
                    None,
                    Some(blob.name().to_string()),
                    false,
                    "attachment",
                );
            }
        }

        HttpResponse::NotFound().finish()
    }

    fn find_lfs_object(&self, blob: &Blob) -> Option<LfsObject> {
        if let Some(oid) = blob.lfs_oid() {
            if let Some(lfs_object) = LfsObject::find_by_oid(oid) {
                if lfs_object.file_exists() {
                    return Some(lfs_object);
                }
            }
        }

        None
    }

    fn send_git_blob(
        &self,
        req: &HttpRequest,
        repository: &Repository,
        blob: &Blob,
        inline: bool,
    ) -> impl Responder {
        // In a real implementation, this would stream the git blob
        // For now, we'll just return a placeholder response
        let mut response = HttpResponse::Ok();
        response.append_header(("Content-Type", blob.mime_type()));
        response.append_header((
            "Content-Disposition",
            format!(
                "{}; filename=\"{}\"",
                if inline { "inline" } else { "attachment" },
                blob.name()
            ),
        ));
        response.body(blob.data())
    }
}

impl SendsBlob for SendsBlobHandler {
    fn send_blob(
        &self,
        req: &HttpRequest,
        repository: &Repository,
        blob: Option<&Blob>,
        inline: bool,
        allow_caching: bool,
    ) -> impl Responder {
        if let Some(blob) = blob {
            let mut response = HttpResponse::Ok();
            response.append_header(("X-Content-Type-Options", "nosniff"));

            // Extract ref and commit from query parameters
            let ref_name = req.query_string().get("ref").map(|s| s.to_string());
            let commit_id = req.query_string().get("commit").map(|s| s.to_string());

            // Check if the blob is already cached
            if self.cached_blob(
                req,
                blob,
                allow_caching,
                ref_name.as_deref(),
                commit_id.as_deref(),
            ) {
                return HttpResponse::NotModified().finish();
            }

            // Send the blob based on its storage type
            if blob.stored_externally() {
                return self.send_lfs_object(req, blob, repository.project());
            } else {
                return self.send_git_blob(req, repository, blob, inline);
            }
        } else {
            HttpResponse::NotFound().finish()
        }
    }
}

// These would be implemented in separate modules
pub mod models {
    pub mod repository {
        use crate::models::project::Project;
        use std::sync::Arc;

        pub struct Repository {
            project: Arc<Project>,
        }

        impl Repository {
            pub fn project(&self) -> &Project {
                &self.project
            }
        }
    }

    pub mod blob {
        use std::time::Duration;

        pub struct Blob {
            id: String,
            name: String,
            mime_type: String,
            data: Vec<u8>,
            lfs_oid: Option<String>,
            stored_externally: bool,
        }

        impl Blob {
            pub const CACHE_TIME: Duration = Duration::from_secs(60 * 60); // 1 hour
            pub const CACHE_TIME_IMMUTABLE: Duration = Duration::from_secs(60 * 60 * 24 * 7); // 1 week

            pub fn id(&self) -> &str {
                &self.id
            }

            pub fn name(&self) -> &str {
                &self.name
            }

            pub fn mime_type(&self) -> &str {
                &self.mime_type
            }

            pub fn data(&self) -> &[u8] {
                &self.data
            }

            pub fn lfs_oid(&self) -> Option<&str> {
                self.lfs_oid.as_deref()
            }

            pub fn stored_externally(&self) -> bool {
                self.stored_externally
            }
        }
    }

    pub mod project {
        use std::sync::Arc;

        pub struct Project {
            id: i64,
            name: String,
        }

        impl Project {
            pub fn id(&self) -> i64 {
                self.id
            }

            pub fn name(&self) -> &str {
                &self.name
            }
        }
    }

    pub mod lfs_object {
        use crate::controllers::concerns::send_file_upload::FileUpload;
        use crate::models::project::Project;
        use std::sync::Arc;

        pub struct LfsObject {
            oid: String,
            file: FileUpload,
        }

        impl LfsObject {
            pub fn find_by_oid(oid: &str) -> Option<Self> {
                // In a real implementation, this would query the database
                // For now, we'll just return None
                None
            }

            pub fn file_exists(&self) -> bool {
                // In a real implementation, this would check if the file exists
                // For now, we'll just return true
                true
            }

            pub fn project_allowed_access(&self, project: &Project) -> bool {
                // In a real implementation, this would check if the project has access to the LFS object
                // For now, we'll just return true
                true
            }
        }
    }
}
