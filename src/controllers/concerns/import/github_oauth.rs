use actix_web::{web, HttpResponse, Responder};
use oauth2::{basic::BasicClient, AuthUrl, ClientId, ClientSecret, RedirectUrl, TokenUrl};
use serde::{Deserialize, Serialize};
use std::sync::Arc;
use url::Url;

use crate::config::AppConfig;
use crate::models::import::GithubImportState;
use crate::utils::oauth::OAuthError;

#[derive(Debug, Serialize, Deserialize)]
pub struct OAuthConfig {
    client_id: String,
    client_secret: String,
    callback_url: String,
}

pub trait GitHubOAuth {
    fn app_config(&self) -> Arc<AppConfig>;
    fn github_import_state(&self) -> Arc<GithubImportState>;
    fn set_oauth_error(&self, error: OAuthError);
    fn set_oauth_state(&self, state: String);
    fn oauth_state(&self) -> Option<String>;
    fn set_oauth_code(&self, code: String);
    fn oauth_code(&self) -> Option<String>;

    fn ensure_oauth_config(&self) -> Result<OAuthConfig, OAuthError> {
        let config = self.app_config();

        if config.github_client_id.is_empty() || config.github_client_secret.is_empty() {
            return Err(OAuthError::ConfigurationMissing);
        }

        Ok(OAuthConfig {
            client_id: config.github_client_id.clone(),
            client_secret: config.github_client_secret.clone(),
            callback_url: config.github_callback_url.clone(),
        })
    }

    fn oauth_client(&self) -> Result<BasicClient, OAuthError> {
        let config = self.ensure_oauth_config()?;

        let auth_url = AuthUrl::new("https://github.com/login/oauth/authorize".to_string())
            .map_err(|_| OAuthError::InvalidUrl)?;
        let token_url = TokenUrl::new("https://github.com/login/oauth/access_token".to_string())
            .map_err(|_| OAuthError::InvalidUrl)?;

        let client_id = ClientId::new(config.client_id);
        let client_secret = ClientSecret::new(config.client_secret);
        let redirect_url =
            RedirectUrl::new(config.callback_url).map_err(|_| OAuthError::InvalidUrl)?;

        Ok(
            BasicClient::new(client_id, Some(client_secret), auth_url, Some(token_url))
                .set_redirect_uri(redirect_url),
        )
    }

    fn authorize_url(&self) -> Result<String, OAuthError> {
        let client = self.oauth_client()?;
        let state = self.generate_oauth_state();
        self.set_oauth_state(state.clone());

        let (auth_url, _) = client
            .authorize_url(|| state)
            .add_scope(oauth2::Scope::new("repo".to_string()))
            .add_scope(oauth2::Scope::new("user".to_string()))
            .url();

        Ok(auth_url.to_string())
    }

    fn generate_oauth_state(&self) -> String {
        use rand::{thread_rng, Rng};
        let mut rng = thread_rng();
        let state: u64 = rng.gen();
        format!("{:x}", state)
    }

    fn handle_oauth_callback(&self) -> impl Responder {
        match self.oauth_code() {
            Some(code) => {
                self.set_oauth_code(code.clone());
                match self.exchange_code_for_token(&code) {
                    Ok(_) => HttpResponse::Ok().finish(),
                    Err(e) => {
                        self.set_oauth_error(e);
                        HttpResponse::InternalServerError().finish()
                    }
                }
            }
            None => HttpResponse::BadRequest().finish(),
        }
    }

    fn exchange_code_for_token(&self, code: &str) -> Result<(), OAuthError> {
        let client = self.oauth_client()?;

        let token = client
            .exchange_code(oauth2::AuthorizationCode::new(code.to_string()))
            .request(oauth2::reqwest::async_http_client)
            .await
            .map_err(|_| OAuthError::TokenExchangeFailed)?;

        self.github_import_state()
            .set_access_token(token.access_token().secret().clone());
        Ok(())
    }
}
