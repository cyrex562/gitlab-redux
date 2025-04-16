use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct GitlabSchema {
    pub query: String,
    pub mutation: Option<String>,
    pub subscription: Option<String>,
}

impl GitlabSchema {
    pub fn new(query: String) -> Self {
        Self {
            query,
            mutation: None,
            subscription: None,
        }
    }
}
