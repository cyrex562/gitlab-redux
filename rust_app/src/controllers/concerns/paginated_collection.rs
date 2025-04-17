use actix_web::{web, HttpResponse, Responder};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

#[derive(Debug, Serialize, Deserialize)]
pub struct PaginatedCollection<T> {
    pub items: Vec<T>,
    pub current_page: i32,
    pub total_pages: i32,
    pub total_count: i64,
    pub per_page: i32,
}

impl<T> PaginatedCollection<T> {
    pub fn new(
        items: Vec<T>,
        current_page: i32,
        total_pages: i32,
        total_count: i64,
        per_page: i32,
    ) -> Self {
        PaginatedCollection {
            items,
            current_page,
            total_pages,
            total_count,
            per_page,
        }
    }
}

pub trait PaginatedCollectionHandler {
    fn redirect_out_of_range(
        &self,
        collection: &PaginatedCollection<impl std::fmt::Debug>,
        total_pages: Option<i32>,
    ) -> impl Responder {
        let total_pages = total_pages.unwrap_or(collection.total_pages);

        if total_pages == 0 {
            return HttpResponse::Ok().finish();
        }

        let out_of_range = collection.current_page > total_pages;

        if out_of_range {
            let mut params = self.safe_params();
            params.insert("page".to_string(), total_pages.to_string());

            let url = self.build_url(&params);
            return HttpResponse::Found()
                .append_header(("Location", url))
                .finish();
        }

        HttpResponse::Ok().finish()
    }

    // Required methods to be implemented by concrete types
    fn safe_params(&self) -> HashMap<String, String>;
    fn build_url(&self, params: &HashMap<String, String>) -> String;
}
