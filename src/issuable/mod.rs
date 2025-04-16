use serde::{Deserialize, Serialize};
use std::collections::HashMap;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Issuable {
    pub id: i64,
    pub title: String,
    pub description: Option<String>,
    pub state: IssuableState,
    pub confidential: bool,
    pub author_id: i64,
    pub assignee_ids: Vec<i64>,
    pub label_ids: Vec<i64>,
    pub milestone_id: Option<i64>,
    pub project_id: i64,
    pub created_at: chrono::DateTime<chrono::Utc>,
    pub updated_at: chrono::DateTime<chrono::Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum IssuableState {
    Opened,
    Closed,
    Merged,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct IssuableMetadata {
    pub upvotes: i32,
    pub downvotes: i32,
    pub user_notes_count: i32,
    pub subscribed: bool,
}

pub trait IssuableActions {
    fn show(&self) -> Result<(), Box<dyn std::error::Error>>;
    fn update(&self) -> Result<(), Box<dyn std::error::Error>>;
    fn destroy(&self) -> Result<(), Box<dyn std::error::Error>>;
    fn bulk_update(&self) -> Result<(), Box<dyn std::error::Error>>;
}

pub trait IssuableCollections {
    fn set_issuables_index(&mut self) -> Result<(), Box<dyn std::error::Error>>;
    fn set_pagination(&mut self) -> Result<(), Box<dyn std::error::Error>>;
    fn issuables_collection(&self) -> Result<Vec<Issuable>, Box<dyn std::error::Error>>;
}

pub trait IssuableLinks {
    fn index(&self) -> Result<Vec<Issuable>, Box<dyn std::error::Error>>;
    fn create(&self) -> Result<(), Box<dyn std::error::Error>>;
    fn destroy(&self) -> Result<(), Box<dyn std::error::Error>>;
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct IssuableFinderOptions {
    pub scope: Option<String>,
    pub state: Option<IssuableState>,
    pub confidential: Option<bool>,
    pub sort: Option<String>,
    pub project_id: Option<i64>,
    pub group_id: Option<i64>,
    pub include_subgroups: Option<bool>,
    pub search: Option<String>,
    pub iids: Option<Vec<i64>>,
}

impl Default for IssuableFinderOptions {
    fn default() -> Self {
        Self {
            scope: None,
            state: Some(IssuableState::Opened),
            confidential: None,
            sort: None,
            project_id: None,
            group_id: None,
            include_subgroups: None,
            search: None,
            iids: None,
        }
    }
}
