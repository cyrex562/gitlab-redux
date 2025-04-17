use actix_web::web;
use std::collections::HashMap;

pub trait ParamsBackwardCompatibility {
    fn set_non_archived_param(&self, params: &mut web::Json<HashMap<String, serde_json::Value>>) {
        let archived = params
            .get("archived")
            .and_then(|v| v.as_bool())
            .unwrap_or(false);
        params.insert("non_archived".to_string(), serde_json::json!(!archived));
    }
}
