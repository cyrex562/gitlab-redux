use crate::models::oauth::application::Application;
use crate::models::user::User;
use actix_web::HttpRequest;
use serde::{Deserialize, Serialize};
use std::sync::Arc;

pub struct CreateService {
    user: Arc<User>,
    params: ApplicationParams,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ApplicationParams {
    pub name: String,
    pub redirect_uri: String,
    pub scopes: String,
    pub confidential: bool,
}

impl CreateService {
    pub fn new(user: Arc<User>, params: ApplicationParams) -> Self {
        Self { user, params }
    }

    pub async fn execute(&self, _req: &HttpRequest) -> Result<Application, &'static str> {
        // Check if scopes are valid
        if !self.validate_scopes(&self.params.scopes) {
            return Err("Invalid scopes");
        }

        let application = Application::new(
            self.params.name.clone(),
            self.params.redirect_uri.clone(),
            self.params.scopes.clone(),
            self.params.confidential,
            Some(&self.user),
        );

        // TODO: Save to database
        // For now just return the created application
        Ok(application)
    }

    fn validate_scopes(&self, scopes: &str) -> bool {
        // TODO: Implement proper scope validation against allowed scopes
        !scopes.is_empty()
    }
}
