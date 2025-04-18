use chrono::{DateTime, Duration, Utc};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::sync::Arc;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct User {
    pub id: i64,
    pub name: String,
    pub two_factor_enabled: bool,
    pub source_groups: Vec<Group>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Group {
    pub id: i64,
    pub name: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Session {
    pub skip_two_factor: Option<DateTime<Utc>>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Request {
    pub path: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Response {
    pub status: u16,
    pub body: String,
    pub headers: HashMap<String, String>,
    pub redirect: Option<String>,
}

pub trait EnforcesTwoFactorAuthentication {
    fn check_two_factor_requirement(
        &self,
        current_user: Option<&User>,
        session: &Session,
        request: &Request,
    ) -> Result<Response, String>;
    fn two_factor_authentication_required(&self) -> bool;
    fn current_user_requires_two_factor(&self, current_user: &User, session: &Session) -> bool;
    fn execute_action_for_2fa_reason(
        &self,
        current_user: &User,
        actions: &HashMap<String, Box<dyn Fn(&Vec<Group>) -> Result<Response, String>>>,
    ) -> Result<Response, String>;
    fn two_factor_grace_period(&self) -> Duration;
    fn two_factor_grace_period_expired(&self, current_user: &User) -> bool;
    fn two_factor_skippable(&self, current_user: &User) -> bool;
    fn skip_two_factor(&self, session: &Session) -> bool;
    fn mfa_help_page_url(&self) -> String;
}

pub struct TwoFactorAuthVerifier {
    current_user: Option<Arc<User>>,
    request: Arc<Request>,
}

impl TwoFactorAuthVerifier {
    pub fn new(current_user: Option<Arc<User>>, request: Arc<Request>) -> Self {
        Self {
            current_user,
            request,
        }
    }

    pub fn two_factor_authentication_required(&self) -> bool {
        // In a real implementation, this would check if 2FA is required
        true
    }

    pub fn current_user_needs_to_setup_two_factor(&self) -> bool {
        if let Some(user) = &self.current_user {
            !user.two_factor_enabled
        } else {
            false
        }
    }

    pub fn two_factor_authentication_reason(&self) -> String {
        // In a real implementation, this would return the reason for 2FA requirement
        "security".to_string()
    }

    pub fn two_factor_grace_period(&self) -> Duration {
        // In a real implementation, this would return the grace period
        Duration::days(7)
    }

    pub fn two_factor_grace_period_expired(&self) -> bool {
        // In a real implementation, this would check if the grace period has expired
        false
    }
}

pub struct EnforcesTwoFactorAuthenticationHandler {
    verifier: Arc<TwoFactorAuthVerifier>,
}

impl EnforcesTwoFactorAuthenticationHandler {
    pub fn new(current_user: Option<Arc<User>>, request: Arc<Request>) -> Self {
        Self {
            verifier: Arc::new(TwoFactorAuthVerifier::new(current_user, request)),
        }
    }

    fn render_error(&self, message: String, status: u16) -> Response {
        let mut headers = HashMap::new();
        headers.insert("Content-Type".to_string(), "application/json".to_string());

        Response {
            status,
            body: serde_json::json!({
                "error": message
            })
            .to_string(),
            headers,
            redirect: None,
        }
    }
}

impl EnforcesTwoFactorAuthentication for EnforcesTwoFactorAuthenticationHandler {
    fn check_two_factor_requirement(
        &self,
        current_user: Option<&User>,
        session: &Session,
        request: &Request,
    ) -> Result<Response, String> {
        if current_user.is_none() {
            return Ok(Response {
                status: 200,
                body: "".to_string(),
                headers: HashMap::new(),
                redirect: None,
            });
        }

        let current_user = current_user.unwrap();

        if self.two_factor_authentication_required()
            && self.current_user_requires_two_factor(current_user, session)
        {
            // Check if this is a GraphQL request
            if request.path.contains("/api/graphql") {
                return Ok(self.render_error(
                    format!(
                        "Authentication error: enable 2FA in your profile settings to continue using GitLab: {}",
                        self.mfa_help_page_url()
                    ),
                    401
                ));
            } else {
                let mut headers = HashMap::new();
                headers.insert("Content-Type".to_string(), "text/html".to_string());

                return Ok(Response {
                    status: 302,
                    body: "".to_string(),
                    headers,
                    redirect: Some("/profile/two_factor_auth".to_string()),
                });
            }
        }

        Ok(Response {
            status: 200,
            body: "".to_string(),
            headers: HashMap::new(),
            redirect: None,
        })
    }

    fn two_factor_authentication_required(&self) -> bool {
        self.verifier.two_factor_authentication_required()
    }

    fn current_user_requires_two_factor(&self, current_user: &User, session: &Session) -> bool {
        self.verifier.current_user_needs_to_setup_two_factor() && !self.skip_two_factor(session)
    }

    fn execute_action_for_2fa_reason(
        &self,
        current_user: &User,
        actions: &HashMap<String, Box<dyn Fn(&Vec<Group>) -> Result<Response, String>>>,
    ) -> Result<Response, String> {
        let reason = self.verifier.two_factor_authentication_reason();
        let groups_enforcing_two_factor = current_user.source_groups.clone();

        if let Some(action) = actions.get(&reason) {
            action(&groups_enforcing_two_factor)
        } else {
            Ok(Response {
                status: 400,
                body: "Invalid 2FA reason".to_string(),
                headers: HashMap::new(),
                redirect: None,
            })
        }
    }

    fn two_factor_grace_period(&self) -> Duration {
        self.verifier.two_factor_grace_period()
    }

    fn two_factor_grace_period_expired(&self, current_user: &User) -> bool {
        self.verifier.two_factor_grace_period_expired()
    }

    fn two_factor_skippable(&self, current_user: &User) -> bool {
        self.two_factor_authentication_required()
            && !current_user.two_factor_enabled
            && !self.two_factor_grace_period_expired(current_user)
    }

    fn skip_two_factor(&self, session: &Session) -> bool {
        if let Some(skip_time) = session.skip_two_factor {
            skip_time > Utc::now()
        } else {
            false
        }
    }

    fn mfa_help_page_url(&self) -> String {
        "https://docs.gitlab.com/ee/user/profile/account/two_factor_authentication.html#enable-two-factor-authentication".to_string()
    }
}
