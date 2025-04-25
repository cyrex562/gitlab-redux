// Ported from: orig_app/app/controllers/concerns/default_branch_protection.rb
// Provides normalization for default branch protection parameters.
use std::collections::HashMap;

pub struct BranchProtectionDefaults {
    pub allowed_to_push: Vec<HashMap<String, i32>>,
    pub allowed_to_merge: Vec<HashMap<String, i32>>,
    pub allow_force_push: Option<bool>,
    pub code_owner_approval_required: Option<bool>,
    pub developer_can_initial_push: Option<bool>,
}

pub fn normalize_default_branch_params(
    params: &mut HashMap<String, serde_json::Value>,
    form_key: &str,
    protection_none: &serde_json::Value,
    protected_fully: &HashMap<String, bool>,
) {
    if let Some(entity_settings_params) = params.get_mut(form_key) {
        if let Some(default_branch_protected) =
            entity_settings_params.get("default_branch_protected")
        {
            if to_bool(default_branch_protected) == Some(false) {
                entity_settings_params["default_branch_protection_defaults"] =
                    protection_none.clone();
                return;
            }
        }
        if !entity_settings_params
            .get("default_branch_protection_defaults")
            .is_some()
        {
            return;
        }
        entity_settings_params.as_object_mut().map(|obj| {
            obj.remove("default_branch_protection_level");
        });
        if let Some(defaults) = entity_settings_params.get_mut("default_branch_protection_defaults")
        {
            if let Some(allowed_to_push) = defaults.get_mut("allowed_to_push") {
                if let Some(arr) = allowed_to_push.as_array_mut() {
                    for entry in arr.iter_mut() {
                        if let Some(access_level) = entry.get_mut("access_level") {
                            if let Some(val) = access_level.as_i64() {
                                *access_level = serde_json::Value::from(val as i32);
                            } else if let Some(s) = access_level.as_str() {
                                if let Ok(val) = s.parse::<i32>() {
                                    *access_level = serde_json::Value::from(val);
                                }
                            }
                        }
                    }
                }
            }
            if let Some(allowed_to_merge) = defaults.get_mut("allowed_to_merge") {
                if let Some(arr) = allowed_to_merge.as_array_mut() {
                    for entry in arr.iter_mut() {
                        if let Some(access_level) = entry.get_mut("access_level") {
                            if let Some(val) = access_level.as_i64() {
                                *access_level = serde_json::Value::from(val as i32);
                            } else if let Some(s) = access_level.as_str() {
                                if let Ok(val) = s.parse::<i32>() {
                                    *access_level = serde_json::Value::from(val);
                                }
                            }
                        }
                    }
                }
            }
            for key in [
                "allow_force_push",
                "code_owner_approval_required",
                "developer_can_initial_push",
            ] {
                if let Some(val) = defaults.get_mut(key) {
                    let default = protected_fully.get(key).copied().unwrap_or(false);
                    *val = serde_json::Value::from(to_bool_with_default(val, default));
                }
            }
        }
    }
}

fn to_bool(val: &serde_json::Value) -> Option<bool> {
    match val {
        serde_json::Value::Bool(b) => Some(*b),
        serde_json::Value::Number(n) => Some(n.as_i64().unwrap_or(0) != 0),
        serde_json::Value::String(s) => match s.as_str() {
            "true" | "1" => Some(true),
            "false" | "0" => Some(false),
            _ => None,
        },
        _ => None,
    }
}

fn to_bool_with_default(val: &serde_json::Value, default: bool) -> bool {
    to_bool(val).unwrap_or(default)
}
