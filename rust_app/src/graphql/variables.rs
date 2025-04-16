use anyhow::{anyhow, Result};
use serde_json::Value;
use std::collections::HashMap;
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct Variables {
    pub query: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub variables: Option<serde_json::Value>,
}

impl Variables {
    pub fn new(variables: Value) -> Self {
        Self { variables }
    }

    pub fn to_map(&self) -> Result<HashMap<String, Value>> {
        match &self.variables {
            Value::Object(map) => {
                let mut result = HashMap::new();
                for (key, value) in map {
                    result.insert(key.clone(), value.clone());
                }
                Ok(result)
            }
            Value::Null => Ok(HashMap::new()),
            _ => Err(anyhow!("Variables must be an object or null")),
        }
    }
}

#[derive(Debug)]
pub struct InvalidVariablesError {
    pub message: String,
}

impl std::fmt::Display for InvalidVariablesError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{}", self.message)
    }
}

impl std::error::Error for InvalidVariablesError {} 
