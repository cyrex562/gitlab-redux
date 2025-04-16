use serde::{Deserialize, Serialize};
use uuid::Uuid;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct MergeRequest {
    pub id: Uuid,
    pub iid: i32,
    pub title: String,
    pub description: Option<String>,
    pub state: String,
    pub created_at: chrono::DateTime<chrono::Utc>,
    pub updated_at: chrono::DateTime<chrono::Utc>,
    pub project_id: Uuid,
    // Add other merge request fields as needed
}

impl MergeRequest {
    pub fn is_persisted(&self) -> bool {
        // In a real implementation, this would check if the merge request is saved to the database
        true
    }
}
