use crate::{
    config::Settings,
    models::{DiffFile, DiffResource},
    utils::feature_flags::Feature,
};
use actix_web::{error::Error, web::Json, HttpResponse};
use serde_json::Value;

pub trait RapidDiffsResource {
    fn diffs_stream_url(
        &self,
        resource: &dyn DiffResource,
        offset: Option<usize>,
        diff_view: Option<&str>,
    ) -> Option<String>;
    fn diff_files_metadata(&self) -> Result<HttpResponse, Error>;
}

pub struct RapidDiffsResourceImpl {
    settings: Settings,
}

impl RapidDiffsResourceImpl {
    pub fn new(settings: Settings) -> Self {
        Self { settings }
    }

    fn rapid_diffs_enabled(&self) -> bool {
        Feature::enabled("rapid_diffs", None, Some("wip"))
    }
}

impl RapidDiffsResource for RapidDiffsResourceImpl {
    fn diffs_stream_url(
        &self,
        resource: &dyn DiffResource,
        offset: Option<usize>,
        diff_view: Option<&str>,
    ) -> Option<String> {
        if let Some(offset) = offset {
            if offset > resource.diffs_for_streaming().diff_files().count() {
                return None;
            }
        }
        self.diffs_stream_resource_url(resource, offset, diff_view)
    }

    fn diff_files_metadata(&self) -> Result<HttpResponse, Error> {
        if !self.rapid_diffs_enabled() {
            return Ok(HttpResponse::NotFound().finish());
        }

        let diffs_resource = self.diffs_resource()?;
        if diffs_resource.is_none() {
            return Ok(HttpResponse::NotFound().finish());
        }

        let diffs_resource = diffs_resource.unwrap();
        let diff_files = diffs_resource.raw_diff_files(true);

        Ok(HttpResponse::Ok().json(serde_json::json!({
            "diff_files": diff_files
        })))
    }
}

impl RapidDiffsResourceImpl {
    fn diffs_resource(&self) -> Result<Option<Box<dyn DiffResource>>, Error> {
        Err(actix_web::error::ErrorNotImplemented(
            "diffs_resource not implemented",
        ))
    }

    fn diffs_stream_resource_url(
        &self,
        resource: &dyn DiffResource,
        offset: Option<usize>,
        diff_view: Option<&str>,
    ) -> Option<String> {
        None // To be implemented by concrete implementations
    }
}
