use serde::{Deserialize, Serialize};
use std::env;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct OAuthConfig {
    pub client_id: String,
    pub client_secret: String,
    pub redirect_uri: String,
    pub authorize_url: String,
    pub token_url: String,
    pub scope: String,
}

impl OAuthConfig {
    pub fn new(
        client_id: String,
        client_secret: String,
        redirect_uri: String,
        authorize_url: String,
        token_url: String,
        scope: String,
    ) -> Self {
        Self {
            client_id,
            client_secret,
            redirect_uri,
            authorize_url,
            token_url,
            scope,
        }
    }

    pub fn from_env() -> Option<Self> {
        let client_id = env::var("GITHUB_CLIENT_ID").ok()?;
        let client_secret = env::var("GITHUB_CLIENT_SECRET").ok()?;
        let redirect_uri = env::var("GITHUB_REDIRECT_URI").ok()?;
        let authorize_url = env::var("GITHUB_AUTHORIZE_URL")
            .unwrap_or_else(|_| "https://github.com/login/oauth/authorize".to_string());
        let token_url = env::var("GITHUB_TOKEN_URL")
            .unwrap_or_else(|_| "https://github.com/login/oauth/access_token".to_string());
        let scope = env::var("GITHUB_SCOPE")
            .unwrap_or_else(|_| "read:user user:email".to_string());

        Some(Self::new(
            client_id,
            client_secret,
            redirect_uri,
            authorize_url,
            token_url,
            scope,
        ))
    }
}

impl Default for OAuthConfig {
    fn default() -> Self {
        Self::new(
            String::new(),
            String::new(),
            String::new(),
            "https://github.com/login/oauth/authorize".to_string(),
            "https://github.com/login/oauth/access_token".to_string(),
            "read:user user:email".to_string(),
        )
    }
} 