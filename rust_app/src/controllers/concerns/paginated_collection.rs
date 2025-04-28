// Ported from: orig_app/lib/gitlab/pagination/offset_pagination.rb
// and related pagination helpers in orig_app/app/helpers/pagination_helper.rb
// Ported on: 2025-04-28
//
// This file implements offset-based pagination logic similar to GitLab's Ruby backend.

use crate::models::Paginated;
use actix_web::{error::Error, HttpResponse};
use serde::Deserialize;

pub trait PaginatedCollection {
    fn redirect_out_of_range<T: Paginated>(
        &self,
        collection: &T,
        total_pages: Option<i32>,
    ) -> Result<bool, Error>;
}

pub struct PaginatedCollectionImpl;

impl PaginatedCollectionImpl {
    pub fn new() -> Self {
        Self
    }
}

impl PaginatedCollection for PaginatedCollectionImpl {
    fn redirect_out_of_range<T: Paginated>(
        &self,
        collection: &T,
        total_pages: Option<i32>,
    ) -> Result<bool, Error> {
        let total_pages = total_pages.unwrap_or_else(|| collection.total_pages());
        if total_pages == 0 {
            return Ok(false);
        }

        let out_of_range = collection.current_page() > total_pages;
        if out_of_range {
            let mut params = collection.params().clone();
            params.insert("page".to_string(), total_pages.to_string());

            return Ok(HttpResponse::Found()
                .header("Location", format!("/{}", params.to_query_string()))
                .finish());
        }

        Ok(false)
    }
}
