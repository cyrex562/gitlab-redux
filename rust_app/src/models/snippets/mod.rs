pub mod blobs_actions;
pub mod send_blob;

use actix_web::{web, HttpResponse};
use serde::{Deserialize, Serialize};
use crate::models::user::User;

// TODO: Implement Snippet model
pub struct Snippet {
    pub id: i32,
    // Add other necessary fields
}

// TODO: Implement Blob model
pub struct Blob {
    pub path: String,
    pub content: Vec<u8>,
    // Add other necessary fields
}

// TODO: Implement Repository model
pub struct Repository {
    // Add necessary fields
}

#[derive(Debug, Serialize, Deserialize)]
pub struct RefExtractor {
    pub id: Option<String>,
    pub ref_name: Option<String>,
    pub path: Option<String>,
    pub ref_type: Option<String>,
} 