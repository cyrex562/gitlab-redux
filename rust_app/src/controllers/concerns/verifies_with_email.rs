use crate::models::user::User;
use crate::services::email_verification_service::EmailVerificationService;
use crate::services::rate_limit_service::RateLimitService;
use actix_web::{web, HttpResponse};
use chrono::{DateTime, Utc};

/// Module for handling email verification
pub trait VerifiesWithEmail {
    /// Verify user with email
    fn verify_with_email(&self, user: Option<&User>) -> HttpResponse {
        if let Some(user) = user {
            if !user.is_active() {
                return HttpResponse::Unauthorized().finish();
            }

            if let Some(token) = self.verification_token() {
                // Verify token if submitted
                self.verify_token(user, token)
            } else if self.require_email_verification_enabled(user) {
                // Check rate limit for sign in
                if self.check_rate_limit("user_sign_in", user) {
                    return self.render_sign_in_rate_limited();
                }

                // Verify email if password is valid
                if user.valid_password(self.user_password()) {
                    self.verify_email(user)
                } else {
                    HttpResponse::Unauthorized().finish()
                }
            } else {
                HttpResponse::Ok().finish()
            }
        } else {
            HttpResponse::Unauthorized().finish()
        }
    }

    /// Resend verification code
    fn resend_verification_code(&self, user: Option<&User>) -> HttpResponse {
        if let Some(user) = user {
            if self.is_send_rate_limited(user) {
                let interval = self.rate_limit_interval("email_verification_code_send");
                return HttpResponse::TooManyRequests().json(json!({
                    "status": "failure",
                    "message": format!(
                        "You've reached the maximum amount of resends. Wait {} and try again.",
                        interval
                    )
                }));
            }

            if let Some(email) = self.email_param() {
                if let Some(secondary_email) = self.user_secondary_email(user, email) {
                    self.send_verification_instructions(user, Some(secondary_email))
                } else {
                    HttpResponse::BadRequest().finish()
                }
            } else {
                self.send_verification_instructions(user, None)
            }

            HttpResponse::Ok().json(json!({ "status": "success" }))
        } else {
            HttpResponse::Unauthorized().finish()
        }
    }

    /// Handle verification success
    fn handle_verification_success(&self, user: &User) {
        if self.is_unconfirmed_verification_email(user) {
            user.confirm();
        }

        if user.email_reset_offered_at.is_none() {
            user.email_reset_offered_at = Some(Utc::now());
        }

        user.unlock_access();
        self.log_verification(user, "successful", None);
        self.sign_in(user);
        self.log_audit_event(user);
        self.log_user_activity(user);
        self.verify_known_sign_in();
    }

    /// Get verification token from params
    fn verification_token(&self) -> Option<String>;

    /// Get user password from params
    fn user_password(&self) -> Option<String>;

    /// Get email from params
    fn email_param(&self) -> Option<String>;

    /// Check if email verification is required
    fn require_email_verification_enabled(&self, user: &User) -> bool;

    /// Check rate limit
    fn check_rate_limit(&self, scope: &str, user: &User) -> bool;

    /// Get rate limit interval
    fn rate_limit_interval(&self, scope: &str) -> String;

    /// Check if send is rate limited
    fn is_send_rate_limited(&self, user: &User) -> bool;

    /// Get user secondary email
    fn user_secondary_email(&self, user: &User, email: &str) -> Option<String>;

    /// Send verification instructions
    fn send_verification_instructions(&self, user: &User, secondary_email: Option<String>);

    /// Verify token
    fn verify_token(&self, user: &User, token: String) -> HttpResponse;

    /// Verify email
    fn verify_email(&self, user: &User) -> HttpResponse;

    /// Check if verification email is unconfirmed
    fn is_unconfirmed_verification_email(&self, user: &User) -> bool;

    /// Sign in user
    fn sign_in(&self, user: &User);

    /// Log verification
    fn log_verification(&self, user: &User, event: &str, reason: Option<String>);

    /// Log audit event
    fn log_audit_event(&self, user: &User);

    /// Log user activity
    fn log_user_activity(&self, user: &User);

    /// Verify known sign in
    fn verify_known_sign_in(&self);

    /// Render sign in rate limited
    fn render_sign_in_rate_limited(&self) -> HttpResponse;
}
