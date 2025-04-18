use actix_web::web::Query;
use actix_web::{web, HttpRequest, Result};
use serde::{Deserialize, Serialize};

// Define the pagination parameters
#[derive(Debug, Deserialize, Serialize)]
pub struct PaginationParams {
    pub page: Option<i32>,
    pub per_page: Option<i32>,
    pub limit: Option<i32>,
    pub sort: Option<String>,
    pub order_by: Option<String>,
    pub pagination: Option<bool>,
}

// Define the StrongPaginationParams trait
pub trait StrongPaginationParams {
    fn pagination_params(&self) -> PaginationParams;
}

// Define the StrongPaginationParamsHandler struct
pub struct StrongPaginationParamsHandler;

impl StrongPaginationParamsHandler {
    pub fn new() -> Self {
        StrongPaginationParamsHandler
    }
}

// Implement the StrongPaginationParams trait for StrongPaginationParamsHandler
impl StrongPaginationParams for StrongPaginationParamsHandler {
    fn pagination_params(&self) -> PaginationParams {
        // In a real implementation, this would extract parameters from the request
        // For now, we'll return default values
        PaginationParams {
            page: None,
            per_page: None,
            limit: None,
            sort: None,
            order_by: None,
            pagination: None,
        }
    }
}

impl<T> StrongPaginationParams for T
where
    T: Fn() -> Query<PaginationParams>,
{
    fn pagination_params(&self) -> PaginationParams {
        self().into_inner()
    }
}
