use actix::Addr;
use actix_web_actors::ws;
use bytes::Bytes;
use serde_json::Value;
use std::sync::Arc;
use tokio::sync::RwLock;

use super::connection::Connection;

type ConnectionAddr = Addr<Connection>;

#[derive(Debug)]
pub struct Channel {
    pub params: Value,
    pub connection: Option<Arc<RwLock<ConnectionAddr>>>,
}

impl Channel {
    pub fn new(params: Value) -> Self {
        Self {
            params,
            connection: None,
        }
    }

    pub async fn subscribe(&mut self, addr: ConnectionAddr) {
        self.connection = Some(Arc::new(RwLock::new(addr)));
    }

    pub async fn unsubscribe(&mut self) {
        self.connection = None;
    }

    pub async fn broadcast(&self, event: &str, data: Value) -> Result<(), String> {
        if let Some(connection) = &self.connection {
            let payload = self.notification_payload(event);
            let text = serde_json::to_string(&payload).unwrap();
            let msg = ws::Message::Text(text.into());

            let conn = connection.read().await;
            conn.do_send(msg).map_err(|e| e.to_string())?;
            Ok(())
        } else {
            Err("No active connection".to_string())
        }
    }

    pub fn notification_payload(&self, _event: &str) -> Value {
        // TODO: Implement actual payload construction
        Value::Null
    }

    pub async fn reject(&mut self) {
        if let Some(connection) = &self.connection {
            let conn = connection.read().await;
            let _ = conn.do_send(ws::Message::Close(None));
        }
        self.unsubscribe().await;
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_new_channel() {
        let params = serde_json::json!({});
        let channel = Channel::new(params.clone());
        assert_eq!(channel.params, params);
        assert!(channel.connection.is_none());
    }
}
