use serde::{Deserialize, Serialize};
use uuid::Uuid;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Noteable {
    pub id: Uuid,
    pub noteable_type: String,
    pub project_id: Option<Uuid>,
    pub group_id: Option<Uuid>,
    // Add other noteable fields as needed
} 