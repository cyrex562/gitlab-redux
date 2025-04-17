use actix_web::{web, HttpRequest, HttpResponse, Responder};
use chrono::{DateTime, Utc};
use std::sync::Arc;

use crate::auth::two_factor_auth_verifier::TwoFactorAuthVerifier;
use crate::config::settings::Settings;
use crate::models::user::User;
use crate::utils::session::Session;

/// Module for enforcing two-factor authentication requirements
pub trait EnforcesTwoFactorAuthentication {
    /// Check if two-factor authentication is required for the current request
    fn check_two_factor_requirement(&self, req: &HttpRequest) -> HttpResponse {
        if let Some(current_user) = self.current_user() {
            if self.two_factor_authentication_required() && self.current_user_requires_two_factor()
            {
                // Check if this is a GraphQL request
                if req.path().contains("/api/graphql") {
                    return HttpResponse::Unauthorized().json(json!({
                        "errors": [{
                            "message": format!(
                                "Authentication error: enable 2FA in your profile settings to continue using GitLab: {}",
                                self.mfa_help_page_url()
                            )
                        }]
                    }));
                } else {
                    // Redirect to 2FA setup page
                    return HttpResponse::Found()
                        .header("Location", "/profile/two_factor_auth")
                        .finish();
                }
            }
        }

        HttpResponse::Ok().finish()
    }

    /// Check if two-factor authentication is required
    fn two_factor_authentication_required(&self) -> bool {
        self.two_factor_verifier()
            .two_factor_authentication_required()
    }

    /// Check if the current user needs to set up two-factor authentication
    fn current_user_requires_two_factor(&self) -> bool {
        self.two_factor_verifier()
            .current_user_needs_to_setup_two_factor()
            && !self.skip_two_factor()
    }

    /// Execute an action based on the 2FA reason
    fn execute_action_for_2fa_reason<F>(&self, actions: &[(&str, F)]) -> HttpResponse
    where
        F: Fn(&[Arc<User>]) -> HttpResponse,
    {
        let reason = self
            .two_factor_verifier()
            .two_factor_authentication_reason();
        let groups_enforcing_two_factor = self
            .current_user()
            .map(|user| user.source_groups_of_two_factor_authentication_requirement())
            .unwrap_or_default();

        // Find the action for this reason
        if let Some((_, action)) = actions.iter().find(|(r, _)| *r == reason) {
            return action(&groups_enforcing_two_factor);
        }

        HttpResponse::Ok().finish()
    }

    /// Get the two-factor grace period
    fn two_factor_grace_period(&self) -> DateTime<Utc> {
        self.two_factor_verifier().two_factor_grace_period()
    }

    /// Check if the two-factor grace period has expired
    fn two_factor_grace_period_expired(&self) -> bool {
        self.two_factor_verifier().two_factor_grace_period_expired()
    }

    /// Check if two-factor authentication can be skipped
    fn two_factor_skippable(&self) -> bool {
        self.two_factor_authentication_required()
            && !self
                .current_user()
                .map_or(false, |user| user.two_factor_enabled())
            && !self.two_factor_grace_period_expired()
    }

    /// Check if two-factor authentication should be skipped
    fn skip_two_factor(&self) -> bool {
        if let Some(skip_time) = self.session().get::<DateTime<Utc>>("skip_two_factor") {
            return skip_time > Utc::now();
        }
        false
    }

    /// Get the two-factor verifier
    fn two_factor_verifier(&self) -> TwoFactorAuthVerifier {
        TwoFactorAuthVerifier::new(self.current_user(), self.request())
    }

    /// Get the MFA help page URL
    fn mfa_help_page_url(&self) -> String {
        "/help/user/profile/account/two_factor_authentication#enable-two-factor-authentication"
            .to_string()
    }

    // Required trait methods that need to be implemented by the controller
    fn current_user(&self) -> Option<Arc<User>>;
    fn session(&self) -> &Session;
    fn request(&self) -> &HttpRequest;
}
