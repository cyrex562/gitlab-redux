use actix::{Actor, Handler, StreamHandler};
use actix_web::{
    web,
    HttpRequest,
    HttpResponse,
    Error,
    error::ErrorUnauthorized,
};
use actix_web_actors::ws;
use serde_json::Value;
use std::sync::Arc;
use crate::models::user::User;
use super::logging::{Logger, Logging};
use super::channel::Channel;

#[derive(Debug)]
pub struct WsMessage(pub String);

pub struct Connection {
    request: HttpRequest,
    channel: Channel,
    logger: Arc<Logger>,
}

impl Actor for Connection {
    type Context = ws::WebsocketContext<Self>;
}

impl Handler<ws::Message> for Connection {
    type Result = ();

    fn handle(&mut self, msg: ws::Message, ctx: &mut Self::Context) {
        match msg {
            ws::Message::Text(text) => ctx.text(text),
            ws::Message::Close(reason) => {
                self.logger.log_disconnect();
                ctx.close(reason);
                ctx.stop();
            }
            _ => (),
        }
    }
}

impl StreamHandler<Result<ws::Message, ws::ProtocolError>> for Connection {
    fn handle(&mut self, msg: Result<ws::Message, ws::ProtocolError>, ctx: &mut Self::Context) {
        match msg {
            Ok(ws::Message::Text(text)) => {
                if let Ok(value) = serde_json::from_str(&text) {
                    self.logger.log_event("message", &value);
                    ctx.text(text);
                }
            }
            Ok(ws::Message::Close(reason)) => {
                self.logger.log_disconnect();
                ctx.close(reason);
                ctx.stop();
            }
            _ => {}
        }
    }
}

impl Connection {
    pub fn new(request: HttpRequest, user: Option<Arc<User>>, params: Value) -> Self {
        let logger = Arc::new(Logger::new(request.clone(), user));
        let channel = Channel::new(params);

        Self {
            request,
            channel,
            logger,
        }
    }

    pub async fn start(mut self) -> Result<HttpResponse, Error> {
        // Authenticate connection request
        if !self.authenticate().await {
            return Err(ErrorUnauthorized("Unauthorized"));
        }

        self.logger.log_connect();

        // Start WebSocket connection
        ws::start(self, &self.request)
    }

    async fn authenticate(&self) -> bool {
        // TODO: Implement actual authentication
        true
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use actix_web::test::TestRequest;

    #[actix_web::test]
    async fn test_connection_new() {
        let req = TestRequest::default().to_http_request();
        let params = serde_json::json!({});
        let conn = Connection::new(req, None, params);
        assert!(conn.authenticate().await);
    }
}
