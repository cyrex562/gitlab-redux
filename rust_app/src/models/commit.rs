use serde::{Deserialize, Serialize};
use uuid::Uuid;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Commit {
    pub id: Uuid,
    pub sha: String,
    pub title: String,
    pub message: String,
    pub author_name: String,
    pub author_email: String,
    pub created_at: chrono::DateTime<chrono::Utc>,
    // Add other commit fields as needed
}
