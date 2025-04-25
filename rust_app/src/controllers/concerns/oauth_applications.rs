// Ported from orig_app/app/controllers/concerns/oauth_applications.rb
// Handles OAuth application-related controller logic

use std::collections::HashMap;

pub const CREATED_SESSION_KEY: &str = "oauth_applications_created";

pub trait OauthApplications {
    fn prepare_scopes(&mut self, params: &mut HashMap<String, serde_json::Value>);
    fn set_created_session(&mut self, session: &mut HashMap<String, serde_json::Value>);
    fn get_created_session(&mut self, session: &mut HashMap<String, serde_json::Value>) -> bool;
    fn load_scopes(&self) -> Vec<String>;
    fn permitted_params(&self) -> Vec<&'static str>;
    fn application_params(
        &self,
        params: &HashMap<String, serde_json::Value>,
    ) -> HashMap<String, serde_json::Value>;
}

pub struct OauthApplicationsHandler;

impl OauthApplications for OauthApplicationsHandler {
    fn prepare_scopes(&mut self, params: &mut HashMap<String, serde_json::Value>) {
        if let Some(app) = params.get_mut("doorkeeper_application") {
            if let Some(scopes) = app.get_mut("scopes") {
                if let Some(arr) = scopes.as_array_mut() {
                    let joined = arr
                        .iter()
                        .filter_map(|v| v.as_str())
                        .collect::<Vec<_>>()
                        .join(" ");
                    *scopes = serde_json::Value::String(joined);
                }
            }
        }
    }

    fn set_created_session(&mut self, session: &mut HashMap<String, serde_json::Value>) {
        session.insert(
            CREATED_SESSION_KEY.to_string(),
            serde_json::Value::Bool(true),
        );
    }

    fn get_created_session(&mut self, session: &mut HashMap<String, serde_json::Value>) -> bool {
        session
            .remove(CREATED_SESSION_KEY)
            .and_then(|v| v.as_bool())
            .unwrap_or(false)
    }

    fn load_scopes(&self) -> Vec<String> {
        // Placeholder: Replace with actual Doorkeeper config integration
        let all_scopes = vec![
            "read",
            "write",
            "admin",
            "ai_workflow",
            "dynamic_user",
            "self_rotate_scope",
        ];
        all_scopes
            .into_iter()
            .filter(|s| s != "ai_workflow" && s != "dynamic_user" && s != "self_rotate_scope")
            .map(|s| s.to_string())
            .collect()
    }

    fn permitted_params(&self) -> Vec<&'static str> {
        vec!["name", "redirect_uri", "scopes", "confidential"]
    }

    fn application_params(
        &self,
        params: &HashMap<String, serde_json::Value>,
    ) -> HashMap<String, serde_json::Value> {
        if let Some(app) = params.get("doorkeeper_application") {
            if let Some(obj) = app.as_object() {
                let permitted = self.permitted_params();
                obj.iter()
                    .filter(|(k, _)| permitted.contains(&k.as_str()))
                    .map(|(k, v)| (k.clone(), v.clone()))
                    .collect()
            } else {
                HashMap::new()
            }
        } else {
            HashMap::new()
        }
    }
}
