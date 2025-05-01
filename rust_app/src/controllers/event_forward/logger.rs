// Ported from orig_app/app/controllers/event_forward/logger.rb
use crate::logging::JsonLogger;

pub struct Logger {
    base: JsonLogger,
}

impl Logger {
    pub fn new() -> Self {
        Self {
            base: JsonLogger::new("event_collection"),
        }
    }

    pub fn build() -> Self {
        Self::new()
    }
}

impl std::ops::Deref for Logger {
    type Target = JsonLogger;

    fn deref(&self) -> &Self::Target {
        &self.base
    }
}
