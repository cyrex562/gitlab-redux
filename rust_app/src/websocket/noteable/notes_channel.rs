use actix::Addr;
use actix_web_actors::ws;
use serde_json::Value;
use serde_json::json;
use crate::models::{Noteable, User};
use super::super::channel::Channel;
use super::super::connection::Connection;

pub struct NotesChannel {
    pub noteable: Box<dyn Noteable>,
    pub channel: Channel,
}

impl NotesChannel {
    pub fn new(noteable: Box<dyn Noteable>) -> Self {
        let params = serde_json::json!({
            "noteable_id": noteable.get_id(),
            "noteable_type": noteable.get_type(),
        });

        Self {
            noteable,
            channel: Channel::new(params),
        }
    }

    pub async fn connect(&mut self, addr: Addr<Connection>) -> Result<(), String> {
        // Validate the connection
        if !self.can_subscribe().await {
            self.channel.reject().await;
            return Err("Subscription not allowed".to_string());
        }

        // Subscribe to updates
        self.channel.subscribe(addr).await;
        Ok(())
    }

    async fn can_subscribe(&self) -> bool {
        // TODO: Implement proper subscription validation
        true
    }

    pub async fn note_created(&mut self, data: Value) -> Result<(), String> {
        self.channel.broadcast("note_created", data).await
    }

    pub async fn note_updated(&mut self, data: Value) -> Result<(), String> {
        self.channel.broadcast("note_updated", data).await
    }

    pub async fn note_deleted(&mut self, data: Value) -> Result<(), String> {
        self.channel.broadcast("note_deleted", data).await
    }
}
