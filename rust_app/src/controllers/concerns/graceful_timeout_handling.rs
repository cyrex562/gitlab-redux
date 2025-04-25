// Ported from: orig_app/app/controllers/concerns/graceful_timeout_handling.rb
// Provides graceful handling for database timeout errors in JSON requests.

use actix_web::{HttpResponse, ResponseError};
use log::error;
use thiserror::Error;

#[derive(Debug, Error)]
#[error("There is too much data to calculate. Please change your selection.")]
pub struct QueryCanceledError;

impl ResponseError for QueryCanceledError {
    fn error_response(&self) -> HttpResponse {
        HttpResponse::InternalServerError().json(serde_json::json!({
            "error": "There is too much data to calculate. Please change your selection."
        }))
    }
}

pub fn log_query_canceled<E: std::fmt::Debug>(err: &E) {
    error!("QueryCanceledError: {:?}", err);
}
