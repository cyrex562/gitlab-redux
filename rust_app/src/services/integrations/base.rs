use serde::{Deserialize, Serialize};
use std::collections::HashMap;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Integration {
    pub id: i32,
    pub name: String,
    pub integration_type: String,
    pub active: bool,
    pub properties: HashMap<String, String>,
    pub project_id: Option<i32>,
    pub group_id: Option<i32>,
    pub instance_level: bool,
}

impl Integration {
    pub fn new(
        name: String,
        integration_type: String,
        properties: HashMap<String, String>,
        project_id: Option<i32>,
        group_id: Option<i32>,
        instance_level: bool,
    ) -> Self {
        Integration {
            id: 0, // This would be set by the database
            name,
            integration_type,
            active: true,
            properties,
            project_id,
            group_id,
            instance_level,
        }
    }

    pub fn activate(&mut self) {
        self.active = true;
    }

    pub fn deactivate(&mut self) {
        self.active = false;
    }

    pub fn update_properties(&mut self, properties: HashMap<String, String>) {
        self.properties = properties;
    }

    pub fn is_active(&self) -> bool {
        self.active
    }

    pub fn is_instance_level(&self) -> bool {
        self.instance_level
    }

    pub fn is_project_level(&self) -> bool {
        self.project_id.is_some()
    }

    pub fn is_group_level(&self) -> bool {
        self.group_id.is_some()
    }
}

pub trait IntegrationService {
    fn execute(&self, params: HashMap<String, String>) -> Result<String, String>;
    fn validate(&self, params: HashMap<String, String>) -> Result<(), String>;
    fn test(&self) -> Result<String, String>;
} 