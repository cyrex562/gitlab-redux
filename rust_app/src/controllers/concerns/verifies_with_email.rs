use actix_web::{web, HttpRequest, HttpResponse, Responder};
use chrono::{DateTime, Duration, Utc};
use serde::{Deserialize, Serialize};
use std::sync::Arc;
use tokio::sync::RwLock;

use crate::models::user::User;
use crate::services::email_verification::{
    EmailVerificationService, GenerateTokenService, UpdateEmailService, ValidateTokenService,
};
use crate::utils::logger::AppLogger;
use crate::utils::rate_limiter::RateLimiter;

/// Module for handling email verification
pub trait VerifiesWithEmail {
    /// Verify with email
    async fn verify_with_email(&self, req: &HttpRequest) -> HttpResponse {
        let user = self
            .find_user()
            .await
            .or_else(|| self.find_verification_user().await);

        if let Some(user) = user {
            if !user.is_active() {
                return HttpResponse::Unauthorized().finish();
            }

            if let Some(session) = req.session() {
                if let Some(verification_user_id) = session.get::<i32>("verification_user_id") {
                    if let Some(token) = self.verification_params().verification_token {
                        // The verification token is submitted, verify it
                        return self.verify_token(&user, &token).await;
                    }
                }
            }

            if self.require_email_verification_enabled(&user).await {
                // Limit the amount of password guesses
                if self.check_rate_limit("user_sign_in", &user).await {
                    return self.render_sign_in_rate_limited();
                }

                // Verify the email if the user has logged in successfully
                if user.valid_password(&self.user_params().password) {
                    return self.verify_email(&user).await;
                }
            }
        }

        HttpResponse::Ok().finish()
    }

    /// Resend verification code
    async fn resend_verification_code(&self, req: &HttpRequest) -> impl Responder {
        let user = match self.find_verification_user().await {
            Some(user) => user,
            None => return HttpResponse::NotFound().finish(),
        };

        if self.send_rate_limited(&user).await {
            let interval = self.rate_limit_interval("email_verification_code_send");
            let message = format!(
                "You've reached the maximum amount of resends. Wait {} and try again.",
                interval
            );
            return HttpResponse::TooManyRequests().json(json!({
                "status": "failure",
                "message": message
            }));
        }

        let secondary_email = self
            .user_secondary_email(&user, &self.email_params().email)
            .await;

        if let Some(email) = self.email_params().email {
            if let Some(secondary_email) = secondary_email {
                self.send_verification_instructions(&user, Some(secondary_email), None)
                    .await;
            }
        } else {
            self.send_verification_instructions(&user, None, None).await;
        }

        HttpResponse::Ok().json(json!({ "status": "success" }))
    }

    /// Update email
    async fn update_email(&self, req: &HttpRequest) -> impl Responder {
        let user = match self.find_verification_user().await {
            Some(user) => user,
            None => return HttpResponse::NotFound().finish(),
        };

        self.log_verification(&user, "email_update_requested", None)
            .await;

        let service = UpdateEmailService::new(&user);
        let result = service.execute(&self.email_params().email).await;

        match result.status {
            "success" => {
                self.send_verification_instructions(&user, None, None).await;
                HttpResponse::Ok().json(result)
            }
            _ => {
                self.handle_verification_failure(&user, &result.reason, &result.message)
                    .await;
                HttpResponse::UnprocessableEntity().json(result)
            }
        }
    }

    /// Successful verification
    async fn successful_verification(&self, req: &HttpRequest) -> impl Responder {
        if let Some(session) = req.session() {
            session.remove("verification_user_id");
        }

        let redirect_url = self.after_sign_in_path().await;
        HttpResponse::Ok().content_type("text/html").body(format!(
            r#"<!DOCTYPE html>
                <html>
                <head>
                    <meta http-equiv="refresh" content="0;url={}">
                </head>
                <body>
                    <p>Redirecting to <a href="{}">{}</a>...</p>
                </body>
                </html>"#,
            redirect_url, redirect_url, redirect_url
        ))
    }

    /// Skip verify with email
    fn skip_verify_with_email(&self, req: &HttpRequest) -> bool {
        self.two_factor_enabled() || self.is_qa_request(req)
    }

    /// Find verification user
    async fn find_verification_user(&self) -> Option<Arc<User>> {
        if let Some(session) = self.request().session() {
            if let Some(user_id) = session.get::<i32>("verification_user_id") {
                return User::find_by_id(user_id).await;
            }
        }
        None
    }

    /// Send verification instructions
    async fn send_verification_instructions(
        &self,
        user: &User,
        secondary_email: Option<String>,
        reason: Option<String>,
    ) {
        let service = GenerateTokenService::new("unlock_token", user);
        let (raw_token, encrypted_token) = service.execute().await;

        user.set_unlock_token(&encrypted_token);
        user.lock_access(false, reason);

        self.send_verification_instructions_email(user, &raw_token, secondary_email)
            .await;
    }

    /// Send verification instructions email
    async fn send_verification_instructions_email(
        &self,
        user: &User,
        token: &str,
        secondary_email: Option<String>,
    ) {
        let email = secondary_email.unwrap_or_else(|| self.verification_email(user));
        self.notify_service()
            .send_verification_instructions_email(&email, token)
            .await;

        self.log_verification(user, "instructions_sent", None).await;
    }

    /// Verify email
    async fn verify_email(&self, user: &User) -> HttpResponse {
        if user.has_unlock_token() {
            if self.unlock_token_expired(user).await {
                self.send_verification_instructions(user, None, None).await;
            }
            self.prompt_for_email_verification(user).await
        } else if user.is_access_locked() || !self.trusted_ip_address(user).await {
            let reason = if !user.is_access_locked() {
                Some("sign in from untrusted IP address".to_string())
            } else {
                None
            };

            if !self.send_rate_limited(user).await {
                self.send_verification_instructions(user, None, reason)
                    .await;
            }
            self.prompt_for_email_verification(user).await
        } else {
            HttpResponse::Ok().finish()
        }
    }

    /// Verify token
    async fn verify_token(&self, user: &User, token: &str) -> HttpResponse {
        let service = ValidateTokenService::new("unlock_token", user, token);
        let result = service.execute().await;

        match result.status.as_str() {
            "success" => {
                self.handle_verification_success(user).await;
                HttpResponse::Ok().json(json!({
                    "status": "success",
                    "redirect_path": "/users/successful_verification"
                }))
            }
            _ => {
                self.handle_verification_failure(user, &result.reason, &result.message)
                    .await;
                HttpResponse::UnprocessableEntity().json(result)
            }
        }
    }

    /// Render sign in rate limited
    fn render_sign_in_rate_limited(&self) -> HttpResponse {
        let interval = self.rate_limit_interval("user_sign_in");
        let message = format!(
            "Maximum login attempts exceeded. Wait {} and try again.",
            interval
        );
        HttpResponse::TooManyRequests()
            .content_type("text/html")
            .body(format!(
                r#"<!DOCTYPE html>
                <html>
                <head>
                    <meta http-equiv="refresh" content="0;url=/users/sign_in">
                </head>
                <body>
                    <p>{}</p>
                    <p>Redirecting to <a href="/users/sign_in">login page</a>...</p>
                </body>
                </html>"#,
                message
            ))
    }

    /// Rate limit interval
    fn rate_limit_interval(&self, rate_limit: &str) -> String {
        let interval = self.rate_limiter().get_interval(rate_limit);
        self.format_time_interval(interval)
    }

    /// Send rate limited
    async fn send_rate_limited(&self, user: &User) -> bool {
        self.rate_limiter()
            .is_throttled("email_verification_code_send", user)
            .await
    }

    /// Handle verification failure
    async fn handle_verification_failure(&self, user: &User, reason: &str, message: &str) {
        user.add_error("base", message);
        self.log_verification(user, "failed_attempt", Some(reason))
            .await;
    }

    /// Handle verification success
    async fn handle_verification_success(&self, user: &User) {
        if self.unconfirmed_verification_email(user).await {
            user.confirm();
        }

        if user.email_reset_offered_at.is_none() {
            user.set_email_reset_offered_at(Utc::now());
        }

        user.unlock_access();
        self.log_verification(user, "successful", None).await;

        self.sign_in(user).await;

        self.log_audit_event(user, "email_verification");
        self.log_user_activity(user).await;
        self.verify_known_sign_in().await;
    }

    /// Trusted IP address
    async fn trusted_ip_address(&self, user: &User) -> bool {
        self.authentication_service()
            .is_initial_login_or_known_ip_address(user, self.request().peer_addr().unwrap().ip())
            .await
    }

    /// Prompt for email verification
    async fn prompt_for_email_verification(&self, user: &User) -> HttpResponse {
        if let Some(session) = self.request().session() {
            session.insert("verification_user_id", user.id);
        }

        HttpResponse::Ok()
            .content_type("text/html")
            .body(self.render_email_verification_template(user))
    }

    /// Verification params
    fn verification_params(&self) -> VerificationParams;

    /// Email params
    fn email_params(&self) -> EmailParams;

    /// User params
    fn user_params(&self) -> UserParams;

    /// User secondary email
    async fn user_secondary_email(&self, user: &User, email: &str) -> Option<String> {
        user.find_confirmed_email(email).await
    }

    /// Log verification
    async fn log_verification(&self, user: &User, event: &str, reason: Option<String>) {
        AppLogger::info(
            "Email Verification",
            event,
            &user.username,
            self.request().peer_addr().unwrap().ip(),
            reason,
        );
    }

    /// Require email verification enabled
    async fn require_email_verification_enabled(&self, user: &User) -> bool {
        self.feature_service()
            .is_enabled("require_email_verification", user)
            .await
            && !self
                .feature_service()
                .is_enabled("skip_require_email_verification", user)
                .await
    }

    /// Unlock token expired
    async fn unlock_token_expired(&self, user: &User) -> bool {
        if let Some(locked_at) = user.locked_at {
            let token_valid_for = Duration::minutes(ValidateTokenService::TOKEN_VALID_FOR_MINUTES);
            locked_at < Utc::now() - token_valid_for
        } else {
            false
        }
    }

    // Required trait methods that need to be implemented by the controller
    fn request(&self) -> &HttpRequest;
    fn two_factor_enabled(&self) -> bool;
    fn is_qa_request(&self, req: &HttpRequest) -> bool;
    fn find_user(&self) -> Option<Arc<User>>;
    fn after_sign_in_path(&self) -> String;
    fn verification_email(&self, user: &User) -> String;
    fn unconfirmed_verification_email(&self, user: &User) -> bool;
    fn sign_in(&self, user: &User);
    fn log_audit_event(&self, user: &User, auth_method: &str);
    fn log_user_activity(&self, user: &User);
    fn verify_known_sign_in(&self);
    fn rate_limiter(&self) -> Arc<RateLimiter>;
    fn feature_service(&self) -> Arc<FeatureService>;
    fn authentication_service(&self) -> Arc<AuthenticationService>;
    fn notify_service(&self) -> Arc<NotifyService>;
    fn render_email_verification_template(&self, user: &User) -> String;
    fn format_time_interval(&self, seconds: i64) -> String;
}

#[derive(Debug, Serialize, Deserialize)]
pub struct VerificationParams {
    pub verification_token: Option<String>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct EmailParams {
    pub email: Option<String>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct UserParams {
    pub password: String,
}
