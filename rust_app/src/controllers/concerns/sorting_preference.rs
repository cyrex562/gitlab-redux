use crate::config::settings::Settings;
use crate::models::user::User;
use actix_web::{web, HttpResponse};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::sync::Arc;
use tokio::sync::RwLock;

/// Module for handling sorting preferences
pub trait SortingPreference {
    /// Get the current user
    fn current_user(&self) -> Option<&User>;

    /// Get the current sorting preference
    fn current_sorting_preference(&self) -> Option<String>;

    /// Get the sorting preference key
    fn sorting_preference_key(&self) -> String;

    /// Get the current user ID
    fn user_id(&self) -> Option<i32>;

    /// Get the preference key
    fn preference_key(&self) -> String;

    /// Get the default sort field
    fn default_sort_field(&self) -> String {
        "created_at".to_string()
    }

    /// Get the default sort direction
    fn default_sort_direction(&self) -> String {
        "desc".to_string()
    }

    /// Get the valid sort fields
    fn valid_sort_fields(&self) -> Vec<String> {
        vec![
            "created_at".to_string(),
            "updated_at".to_string(),
            "title".to_string(),
            "name".to_string(),
            "author".to_string(),
        ]
    }

    /// Get the valid sort directions
    fn valid_sort_directions(&self) -> Vec<String> {
        vec!["asc".to_string(), "desc".to_string()]
    }

    /// Get the user's sorting preference
    async fn get_sorting_preference(
        &self,
        storage: Arc<RwLock<HashMap<String, HashMap<String, String>>>>,
    ) -> Result<HashMap<String, String>, HttpResponse> {
        let user_id = match self.user_id() {
            Some(id) => id,
            None => return Ok(HashMap::new()),
        };

        let key = format!("user_{}_{}", user_id, self.preference_key());
        let storage = storage.read().await;

        Ok(storage.get(&key).cloned().unwrap_or_else(|| {
            let mut default_prefs = HashMap::new();
            default_prefs.insert("field".to_string(), self.default_sort_field());
            default_prefs.insert("direction".to_string(), self.default_sort_direction());
            default_prefs
        }))
    }

    /// Save the user's sorting preference
    async fn save_sorting_preference(
        &self,
        storage: Arc<RwLock<HashMap<String, HashMap<String, String>>>>,
        field: String,
        direction: String,
    ) -> Result<(), HttpResponse> {
        let user_id = match self.user_id() {
            Some(id) => id,
            None => return Ok(()),
        };

        if !self.valid_sort_fields().contains(&field) {
            return Err(HttpResponse::BadRequest().json(serde_json::json!({
                "error": format!("Invalid sort field: {}", field)
            })));
        }

        if !self.valid_sort_directions().contains(&direction) {
            return Err(HttpResponse::BadRequest().json(serde_json::json!({
                "error": format!("Invalid sort direction: {}", direction)
            })));
        }

        let key = format!("user_{}_{}", user_id, self.preference_key());
        let mut storage = storage.write().await;

        let mut prefs = HashMap::new();
        prefs.insert("field".to_string(), field);
        prefs.insert("direction".to_string(), direction);

        storage.insert(key, prefs);

        Ok(())
    }

    /// Clear the user's sorting preference
    async fn clear_sorting_preference(
        &self,
        storage: Arc<RwLock<HashMap<String, HashMap<String, String>>>>,
    ) -> Result<(), HttpResponse> {
        let user_id = match self.user_id() {
            Some(id) => id,
            None => return Ok(()),
        };

        let key = format!("user_{}_{}", user_id, self.preference_key());
        let mut storage = storage.write().await;

        storage.remove(&key);

        Ok(())
    }

    /// Get the sorting preference
    fn get_sorting_preference(&self) -> String {
        // First try to get from user preferences
        if let Some(user) = self.current_user() {
            if let Some(preference) = user.get_preference(&self.sorting_preference_key()) {
                return preference;
            }
        }

        // Then try to get from current request
        if let Some(preference) = self.current_sorting_preference() {
            return preference;
        }

        // Finally, fall back to default
        self.default_sorting_preference()
    }

    /// Get the default sorting preference
    fn default_sorting_preference(&self) -> String {
        let settings = Settings::current();
        settings.default_sorting_preference.clone()
    }

    /// Get all available sorting preferences
    fn available_sorting_preferences(&self) -> HashMap<String, String> {
        let mut preferences = HashMap::new();

        // Add default preferences
        preferences.insert(
            "created_at_desc".to_string(),
            "Created date (newest first)".to_string(),
        );
        preferences.insert(
            "created_at_asc".to_string(),
            "Created date (oldest first)".to_string(),
        );
        preferences.insert(
            "updated_at_desc".to_string(),
            "Last updated (newest first)".to_string(),
        );
        preferences.insert(
            "updated_at_asc".to_string(),
            "Last updated (oldest first)".to_string(),
        );
        preferences.insert("title_asc".to_string(), "Title (A-Z)".to_string());
        preferences.insert("title_desc".to_string(), "Title (Z-A)".to_string());

        preferences
    }

    /// Validate a sorting preference
    fn validate_sorting_preference(&self, preference: &str) -> bool {
        self.available_sorting_preferences()
            .contains_key(preference)
    }
}
