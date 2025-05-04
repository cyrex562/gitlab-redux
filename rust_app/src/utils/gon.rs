use serde::Serialize;
use std::collections::HashMap;

pub trait GonHelper {
    fn add_gon_variables(&self);
}

#[derive(Default)]
pub struct GonHelperImpl {
    gon: HashMap<String, serde_json::Value>,
}

impl GonHelperImpl {
    pub fn new() -> Self {
        Self {
            gon: HashMap::new(),
        }
    }

    pub fn push<T: Serialize>(&mut self, key: &str, value: T) {
        if let Ok(value) = serde_json::to_value(value) {
            self.gon.insert(key.to_string(), value);
        }
    }
}