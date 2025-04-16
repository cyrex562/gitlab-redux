use serde_json::json;
use std::sync::Arc;
use tokio::sync::RwLock;

use crate::models::{Noteable, Project};
use crate::services::notes::NotesFinder;
use crate::websocket::Channel;

pub struct NotesChannel {
    pub channel: Channel,
    pub noteable: Option<Arc<Noteable>>,
}

impl NotesChannel {
    pub fn new(channel: Channel) -> Self {
        Self {
            channel,
            noteable: None,
        }
    }

    pub async fn subscribe(&mut self) -> Result<(), String> {
        // First call the parent subscribe method to validate token scope
        self.channel.subscribe().await?;

        // Get parameters from the channel
        let params = &self.channel.params;

        // Find project if project_id is present
        let project = if let Some(project_id) = params.get("project_id").and_then(|id| id.as_str())
        {
            // TODO: Implement project lookup
            None
        } else {
            None
        };

        // Find noteable using NotesFinder
        let group_id = params.get("group_id").and_then(|id| id.as_str());
        let noteable_type = params.get("noteable_type").and_then(|t| t.as_str());
        let noteable_id = params.get("noteable_id").and_then(|id| id.as_str());

        // Get current user from connection
        let connection = self.channel.connection.read().await;
        let current_user = connection.current_user.clone();

        if let Some(user) = current_user {
            let notes_finder =
                NotesFinder::new(user, project, group_id, noteable_type, noteable_id);

            self.noteable = notes_finder.find_target().await;

            if self.noteable.is_none() {
                self.channel.reject();
                return Err("Noteable not found".to_string());
            }

            // Stream for the noteable
            // TODO: Implement streaming
        } else {
            self.channel.reject();
            return Err("User not found".to_string());
        }

        Ok(())
    }
}
