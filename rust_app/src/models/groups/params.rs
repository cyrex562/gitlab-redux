use serde::{Deserialize, Serialize};
use std::collections::HashMap;

#[derive(Debug, Serialize, Deserialize)]
pub struct GroupParams {
    pub avatar: Option<String>,
    pub description: Option<String>,
    pub emails_disabled: Option<bool>,
    pub emails_enabled: Option<bool>,
    pub show_diff_preview_in_email: Option<bool>,
    pub mentions_disabled: Option<bool>,
    pub lfs_enabled: Option<bool>,
    pub name: Option<String>,
    pub path: Option<String>,
    pub public: Option<bool>,
    pub request_access_enabled: Option<bool>,
    pub share_with_group_lock: Option<bool>,
    pub visibility_level: Option<i32>,
    pub parent_id: Option<i32>,
    pub create_chat_team: Option<bool>,
    pub chat_team_name: Option<String>,
    pub require_two_factor_authentication: Option<bool>,
    pub two_factor_grace_period: Option<i32>,
    pub enabled_git_access_protocol: Option<String>,
    pub project_creation_level: Option<String>,
    pub subgroup_creation_level: Option<String>,
    pub default_branch_protection: Option<bool>,
    pub default_branch_protection_defaults: Option<BranchProtectionDefaults>,
    pub default_branch_name: Option<String>,
    pub allow_mfa_for_subgroups: Option<bool>,
    pub resource_access_token_creation_allowed: Option<bool>,
    pub resource_access_token_notify_inherited: Option<bool>,
    pub lock_resource_access_token_notify_inherited: Option<bool>,
    pub prevent_sharing_groups_outside_hierarchy: Option<bool>,
    pub setup_for_company: Option<bool>,
    pub jobs_to_be_done: Option<String>,
    pub crm_enabled: Option<bool>,
    pub crm_source_group_id: Option<i32>,
    pub force_pages_access_control: Option<bool>,
    pub enable_namespace_descendants_cache: Option<bool>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct BranchProtectionDefaults {
    pub allow_force_push: Option<bool>,
    pub developer_can_initial_push: Option<bool>,
    pub code_owner_approval_required: Option<bool>,
    pub allowed_to_merge: Vec<AccessLevel>,
    pub allowed_to_push: Vec<AccessLevel>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct AccessLevel {
    pub access_level: i32,
}

pub struct Params;

impl Params {
    pub fn allowed_integration_params() -> Vec<&'static str> {
        vec![
            "active",
            "properties",
            "project_id",
            "group_id",
            "instance_level",
        ]
    }

    pub fn validate_params(params: &GroupParams) -> Result<(), String> {
        // TODO: Implement validation logic
        Ok(())
    }

    pub fn sanitize_params(params: &mut GroupParams) {
        // TODO: Implement sanitization logic
    }
}
