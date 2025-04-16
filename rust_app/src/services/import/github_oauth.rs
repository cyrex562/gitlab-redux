use crate::import::{OAuth2Client, OAuthProviderConfig};
use actix_web::{error::ResponseError, web, HttpResponse};
use rand::{thread_rng, Rng};
use std::collections::HashMap;

#[derive(Debug)]
pub struct OAuthConfigMissingError;

impl std::fmt::Display for OAuthConfigMissingError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "Missing OAuth configuration")
    }
}

impl std::error::Error for OAuthConfigMissingError {}

impl ResponseError for OAuthConfigMissingError {
    fn error_response(&self) -> HttpResponse {
        HttpResponse::Unauthorized().json(serde_json::json!({
            "errors": "Missing OAuth configuration for GitHub."
        }))
    }
}

pub trait GithubOauth {
    // TODO: Implement session handling
    fn get_session(&self) -> &HashMap<String, String>;
    fn set_session(&mut self, key: &str, value: String);

    // TODO: Implement access token key getter
    fn access_token_key(&self) -> &str;

    // TODO: Implement auth state key getter
    fn auth_state_key(&self) -> &str;

    // TODO: Implement provider name getter
    fn provider_name(&self) -> &str;

    // TODO: Implement new project path getter
    fn new_project_path(&self) -> String;

    // TODO: Implement new import URL getter
    fn new_import_url(&self) -> String;

    // TODO: Implement callback URL getter
    fn callback_url(&self, params: HashMap<String, String>) -> String;

    fn provider_auth(&mut self, ci_cd_only: bool) -> Option<HttpResponse> {
        if self.get_session().contains_key(self.access_token_key()) {
            return None;
        }

        if !ci_cd_only {
            return Some(self.go_to_provider_for_permissions());
        }

        None
    }

    fn ci_cd_only(&self, params: &web::Query<HashMap<String, String>>) -> bool {
        params
            .get("ci_cd_only")
            .map(|v| v == "1" || v == "true")
            .unwrap_or(false)
    }

    fn go_to_provider_for_permissions(&mut self) -> HttpResponse {
        HttpResponse::Found()
            .header("Location", self.authorize_url())
            .finish()
    }

    fn oauth_client(&self) -> Result<OAuth2Client, OAuthConfigMissingError> {
        let config = self.oauth_config().ok_or(OAuthConfigMissingError)?;
        Ok(self.oauth_client_from_config(config))
    }

    fn oauth_client_from_config(&self, config: &OAuthProviderConfig) -> OAuth2Client {
        let mut options = HashMap::new();
        if let Some(client_options) = config.args.get("client_options") {
            if let Some(obj) = client_options.as_object() {
                for (k, v) in obj {
                    options.insert(k.clone(), v.clone());
                }
            }
        }

        // Add SSL verification option
        options.insert(
            "verify_ssl".to_string(),
            serde_json::json!(config.verify_ssl),
        );

        OAuth2Client::new(config.app_id.clone(), config.app_secret.clone(), options)
    }

    fn oauth_config(&self) -> Option<OAuthProviderConfig> {
        OAuthProviderConfig::config_for("github")
    }

    fn authorize_url(&mut self) -> String {
        // Generate a random state
        let state = {
            let mut rng = thread_rng();
            let mut bytes = [0u8; 64];
            rng.fill(&mut bytes);
            base64::encode(&bytes)
        };

        // Store state in session
        self.set_session(self.auth_state_key(), state.clone());

        // Store failure path in session
        self.set_session(
            "auth_on_failure_path".to_string(),
            format!("{}#import_project", self.new_project_path()),
        );

        // Get OAuth client
        let client = self
            .oauth_client()
            .expect("OAuth client should be available");

        // Generate authorization URL
        client.auth_code().authorize_url(
            &self.callback_url(HashMap::new()),
            "repo, read:org",
            &state,
        )
    }

    fn get_token(&self, code: &str) -> Result<String, Box<dyn std::error::Error>> {
        let client = self.oauth_client()?;
        client.auth_code().get_token(code)
    }

    fn missing_oauth_config(&mut self) -> HttpResponse {
        // Clear access token from session
        self.set_session(self.access_token_key(), String::new());

        // Return error response
        HttpResponse::Found()
            .header("Location", self.new_import_url())
            .json(serde_json::json!({
                "alert": "Missing OAuth configuration for GitHub."
            }))
    }

    fn callback_import_url(&self) -> String {
        let mut params = HashMap::new();
        params.insert("namespace_id".to_string(), "dummy_namespace_id".to_string());

        // Add any extra import params
        for (k, v) in self.extra_import_params() {
            params.insert(k, v);
        }

        self.callback_url(params)
    }

    fn extra_import_params(&self) -> HashMap<String, String> {
        HashMap::new()
    }
}
