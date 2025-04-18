use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::sync::Arc;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct User {
    pub id: i64,
    pub name: String,
    pub admin: bool,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct UserMode {
    pub admin_mode: bool,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Settings {
    pub admin_mode: bool,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Response {
    pub status: u16,
    pub body: String,
    pub headers: HashMap<String, String>,
    pub redirect: Option<String>,
}

pub trait EnforcesAdminAuthentication {
    fn authenticate_admin(
        &self,
        current_user: &User,
        current_user_mode: &UserMode,
        settings: &Settings,
        request_path: &str,
    ) -> Result<Response, String>;
    fn authorize_ability(&self, current_user: &User, ability: &str) -> Result<Response, String>;
}

pub struct EnforcesAdminAuthenticationHandler;

impl EnforcesAdminAuthenticationHandler {
    pub fn new() -> Self {
        Self
    }

    fn render_404(&self) -> Response {
        let mut headers = HashMap::new();
        headers.insert("Content-Type".to_string(), "text/plain".to_string());

        Response {
            status: 404,
            body: "Not Found".to_string(),
            headers,
            redirect: None,
        }
    }

    fn storable_location(&self, request_path: &str) -> bool {
        request_path != "/admin/session/new"
    }
}

impl EnforcesAdminAuthentication for EnforcesAdminAuthenticationHandler {
    fn authenticate_admin(
        &self,
        current_user: &User,
        current_user_mode: &UserMode,
        settings: &Settings,
        request_path: &str,
    ) -> Result<Response, String> {
        if !current_user.admin {
            return Ok(self.render_404());
        }

        if !settings.admin_mode {
            return Ok(Response {
                status: 200,
                body: "".to_string(),
                headers: HashMap::new(),
                redirect: None,
            });
        }

        if !current_user_mode.admin_mode {
            // In a real implementation, this would request admin mode
            let mut headers = HashMap::new();
            headers.insert("Content-Type".to_string(), "text/html".to_string());

            let redirect_path = if self.storable_location(request_path) {
                Some("/admin/session/new".to_string())
            } else {
                None
            };

            return Ok(Response {
                status: 302,
                body: "Re-authentication required".to_string(),
                headers,
                redirect: redirect_path,
            });
        }

        Ok(Response {
            status: 200,
            body: "".to_string(),
            headers: HashMap::new(),
            redirect: None,
        })
    }

    fn authorize_ability(&self, current_user: &User, ability: &str) -> Result<Response, String> {
        if current_user.admin {
            // In a real implementation, this would call authenticate_admin
            return Ok(Response {
                status: 200,
                body: "".to_string(),
                headers: HashMap::new(),
                redirect: None,
            });
        }

        // In a real implementation, this would check if the user can perform the ability
        // For now, we'll just return 404
        Ok(self.render_404())
    }
}

// Example of how to use this in a controller
pub struct AdminController {
    auth_handler: Arc<EnforcesAdminAuthenticationHandler>,
    current_user: Arc<User>,
    current_user_mode: Arc<UserMode>,
    settings: Arc<Settings>,
}

impl AdminController {
    pub fn new(
        auth_handler: Arc<EnforcesAdminAuthenticationHandler>,
        current_user: Arc<User>,
        current_user_mode: Arc<UserMode>,
        settings: Arc<Settings>,
    ) -> Self {
        Self {
            auth_handler,
            current_user,
            current_user_mode,
            settings,
        }
    }

    pub fn before_action(&self, request_path: &str) -> Result<Response, String> {
        self.auth_handler.authenticate_admin(
            &self.current_user,
            &self.current_user_mode,
            &self.settings,
            request_path,
        )
    }

    pub fn authorize(&self, ability: &str) -> Result<Response, String> {
        self.auth_handler
            .authorize_ability(&self.current_user, ability)
    }
}
