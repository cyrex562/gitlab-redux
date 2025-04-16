use actix_web::{web, HttpResponse};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

/// Module for handling members presentation
pub trait MembersPresentation {
    /// Get the members
    fn members(&self) -> Vec<HashMap<String, serde_json::Value>>;

    /// Get the source type
    fn source_type(&self) -> String;

    /// Get the source ID
    fn source_id(&self) -> i32;

    /// Get the current user ID
    fn current_user_id(&self) -> Option<i32> {
        None
    }

    /// Get the member role
    fn member_role(&self) -> Option<String> {
        None
    }

    /// Get the member access level
    fn member_access_level(&self) -> Option<i32> {
        None
    }

    /// Get the member expires at
    fn member_expires_at(&self) -> Option<String> {
        None
    }

    /// Get the member invited by
    fn member_invited_by(&self) -> Option<i32> {
        None
    }

    /// Get the member invited at
    fn member_invited_at(&self) -> Option<String> {
        None
    }

    /// Get the member accepted at
    fn member_accepted_at(&self) -> Option<String> {
        None
    }

    /// Get the member requested at
    fn member_requested_at(&self) -> Option<String> {
        None
    }

    /// Get the member requested by
    fn member_requested_by(&self) -> Option<i32> {
        None
    }

    /// Get the member state
    fn member_state(&self) -> Option<String> {
        None
    }

    /// Present members
    fn present_members(&self) -> Result<Vec<HashMap<String, serde_json::Value>>, HttpResponse> {
        let mut presented_members = Vec::new();

        for member in &self.members() {
            let mut presented_member = member.clone();

            // Add source information
            presented_member.insert(
                "source_type".to_string(),
                serde_json::json!(self.source_type()),
            );
            presented_member.insert("source_id".to_string(), serde_json::json!(self.source_id()));

            // Add member information
            if let Some(role) = self.member_role() {
                presented_member.insert("role".to_string(), serde_json::json!(role));
            }

            if let Some(access_level) = self.member_access_level() {
                presented_member
                    .insert("access_level".to_string(), serde_json::json!(access_level));
            }

            if let Some(expires_at) = self.member_expires_at() {
                presented_member.insert("expires_at".to_string(), serde_json::json!(expires_at));
            }

            if let Some(invited_by) = self.member_invited_by() {
                presented_member.insert("invited_by".to_string(), serde_json::json!(invited_by));
            }

            if let Some(invited_at) = self.member_invited_at() {
                presented_member.insert("invited_at".to_string(), serde_json::json!(invited_at));
            }

            if let Some(accepted_at) = self.member_accepted_at() {
                presented_member.insert("accepted_at".to_string(), serde_json::json!(accepted_at));
            }

            if let Some(requested_at) = self.member_requested_at() {
                presented_member
                    .insert("requested_at".to_string(), serde_json::json!(requested_at));
            }

            if let Some(requested_by) = self.member_requested_by() {
                presented_member
                    .insert("requested_by".to_string(), serde_json::json!(requested_by));
            }

            if let Some(state) = self.member_state() {
                presented_member.insert("state".to_string(), serde_json::json!(state));
            }

            // Add current user information
            if let Some(current_user_id) = self.current_user_id() {
                presented_member.insert(
                    "current_user_id".to_string(),
                    serde_json::json!(current_user_id),
                );
            }

            presented_members.push(presented_member);
        }

        Ok(presented_members)
    }

    /// Get member metadata
    fn get_member_metadata(&self) -> Result<HashMap<String, serde_json::Value>, HttpResponse> {
        let mut metadata = HashMap::new();

        metadata.insert(
            "source_type".to_string(),
            serde_json::json!(self.source_type()),
        );
        metadata.insert("source_id".to_string(), serde_json::json!(self.source_id()));

        if let Some(current_user_id) = self.current_user_id() {
            metadata.insert(
                "current_user_id".to_string(),
                serde_json::json!(current_user_id),
            );
        }

        if let Some(role) = self.member_role() {
            metadata.insert("role".to_string(), serde_json::json!(role));
        }

        if let Some(access_level) = self.member_access_level() {
            metadata.insert("access_level".to_string(), serde_json::json!(access_level));
        }

        if let Some(state) = self.member_state() {
            metadata.insert("state".to_string(), serde_json::json!(state));
        }

        Ok(metadata)
    }

    /// Validate member data
    fn validate_member_data(&self) -> Result<(), HttpResponse> {
        if self.members().is_empty() {
            return Err(HttpResponse::BadRequest().json(serde_json::json!({
                "error": "Members cannot be empty"
            })));
        }

        let source_type = self.source_type();
        if !["Project", "Group"].contains(&source_type.as_str()) {
            return Err(HttpResponse::BadRequest().json(serde_json::json!({
                "error": format!("Invalid source type: {}", source_type)
            })));
        }

        if let Some(state) = &self.member_state() {
            if !["active", "pending", "expired", "blocked"].contains(&state.as_str()) {
                return Err(HttpResponse::BadRequest().json(serde_json::json!({
                    "error": format!("Invalid member state: {}", state)
                })));
            }
        }

        Ok(())
    }
}
