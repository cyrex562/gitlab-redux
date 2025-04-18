use actix_web::{HttpResponse, Responder};
use chrono::{DateTime, Duration, Utc};
use std::collections::HashMap;
use std::sync::Arc;

pub trait SendsBlob {
    fn send_blob(
        &self,
        repository: Arc<dyn Repository>,
        blob: Option<Blob>,
        inline: bool,
        allow_caching: bool,
    ) -> impl Responder;

    fn cached_blob(&self, blob: &Blob, allow_caching: bool) -> bool;
    fn send_lfs_object(&self, blob: &Blob, project: Arc<dyn Project>) -> impl Responder;
    fn find_lfs_object(&self, blob: &Blob) -> Option<Arc<LfsObject>>;
}

pub trait Repository {
    fn project(&self) -> Arc<dyn Project>;
}

pub trait Project {
    fn id(&self) -> i64;
    fn name(&self) -> String;
}

#[derive(Debug, Clone)]
pub struct Blob {
    pub id: String,
    pub name: String,
    pub path: String,
    pub size: i64,
    pub mode: String,
    pub commit_id: String,
    pub lfs_oid: Option<String>,
    pub stored_externally: bool,
}

#[derive(Debug, Clone)]
pub struct LfsObject {
    pub id: i64,
    pub oid: String,
    pub size: i64,
    pub file: Arc<dyn FileUploader>,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

pub trait FileUploader {
    fn exists(&self) -> bool;
    fn path(&self) -> String;
    fn url(&self, params: HashMap<String, String>) -> String;
    fn file_storage(&self) -> bool;
    fn proxy_download_enabled(&self) -> bool;
}

impl<T> SendsBlob for T
where
    T: crate::uploads::SendFileUpload,
{
    fn send_blob(
        &self,
        repository: Arc<dyn Repository>,
        blob: Option<Blob>,
        inline: bool,
        allow_caching: bool,
    ) -> impl Responder {
        if let Some(blob) = blob {
            // Set security headers
            let mut response = HttpResponse::Ok();
            response.append_header(("X-Content-Type-Options", "nosniff"));

            // Check if blob is cached
            if self.cached_blob(&blob, allow_caching) {
                return response.finish();
            }

            // Set cache headers
            let max_age = if blob.commit_id == blob.id {
                // This is a link to a commit by its commit SHA. That means that the blob
                // is immutable. The only reason to invalidate the cache is if the commit
                // was deleted or if the user lost access to the repository.
                Blob::CACHE_TIME_IMMUTABLE
            } else {
                // A branch or tag points at this blob. That means that the expected blob
                // value may change over time.
                Blob::CACHE_TIME
            };

            response.append_header(("Cache-Control", format!(
                "max-age={}, {}must-revalidate, stale-if-error=300, stale-while-revalidate=60, s-maxage=60",
                max_age,
                if allow_caching { "" } else { "private, " }
            )));

            if blob.stored_externally {
                self.send_lfs_object(&blob, repository.project())
            } else {
                // In a real implementation, this would send the git blob
                response.finish()
            }
        } else {
            HttpResponse::NotFound().finish()
        }
    }

    fn cached_blob(&self, blob: &Blob, allow_caching: bool) -> bool {
        // In a real implementation, this would check if the blob is cached
        false
    }

    fn send_lfs_object(&self, blob: &Blob, project: Arc<dyn Project>) -> impl Responder {
        if let Some(lfs_object) = self.find_lfs_object(blob) {
            if lfs_object.file.exists() {
                // In a real implementation, this would check if the project has access to the LFS object
                let mut params = HashMap::new();
                params.insert("attachment".to_string(), blob.name.clone());

                self.send_upload(
                    Box::new(lfs_object.file.clone()),
                    HashMap::new(),
                    params,
                    Some(blob.name.clone()),
                    false,
                    "attachment".to_string(),
                )
            } else {
                HttpResponse::NotFound().finish()
            }
        } else {
            HttpResponse::NotFound().finish()
        }
    }

    fn find_lfs_object(&self, blob: &Blob) -> Option<Arc<LfsObject>> {
        // In a real implementation, this would find the LFS object in the database
        if let Some(oid) = &blob.lfs_oid {
            // This is a placeholder implementation
            Some(Arc::new(LfsObject {
                id: 1,
                oid: oid.clone(),
                size: blob.size,
                file: Arc::new(PlaceholderFileUploader {
                    path: format!("/lfs/objects/{}", oid),
                }),
                created_at: Utc::now(),
                updated_at: Utc::now(),
            }))
        } else {
            None
        }
    }
}

impl Blob {
    pub const CACHE_TIME: i32 = 3600; // 1 hour
    pub const CACHE_TIME_IMMUTABLE: i32 = 86400; // 24 hours
}

// Placeholder implementation for testing
struct PlaceholderFileUploader {
    path: String,
}

impl FileUploader for PlaceholderFileUploader {
    fn exists(&self) -> bool {
        true
    }

    fn path(&self) -> String {
        self.path.clone()
    }

    fn url(&self, _params: HashMap<String, String>) -> String {
        format!("https://example.com{}", self.path)
    }

    fn file_storage(&self) -> bool {
        true
    }

    fn proxy_download_enabled(&self) -> bool {
        false
    }
}
