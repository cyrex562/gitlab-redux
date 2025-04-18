use actix_web::web::Query;
use serde::Deserialize;

pub trait ParamsBackwardCompatibility {
    fn set_non_archived_param(&self, params: &mut Query<serde_json::Value>);
}

pub struct ParamsBackwardCompatibilityImpl;

impl ParamsBackwardCompatibilityImpl {
    pub fn new() -> Self {
        Self
    }
}

impl ParamsBackwardCompatibility for ParamsBackwardCompatibilityImpl {
    fn set_non_archived_param(&self, params: &mut Query<serde_json::Value>) {
        if !params.contains_key("archived") {
            params.insert("non_archived".to_string(), serde_json::Value::Bool(true));
        }
    }
}
