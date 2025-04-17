use actix_web::{web, HttpRequest, HttpResponse, Result};
use chrono::{DateTime, Duration, Utc};
use serde::{Deserialize, Serialize};
use std::sync::Arc;
use std::time::SystemTime;

// Define the response for verification
#[derive(Debug, Serialize, Deserialize)]
pub struct VerificationResponse {
    pub status: String,
    pub message: Option<String>,
    pub redirect_path: Option<String>,
}

// Define the User trait
pub trait User: Send + Sync {
    fn id(&self) -> i64;
    fn username(&self) -> &str;
    fn email(&self) -> &str;
    fn is_active(&self) -> bool;
    fn is_locked(&self) -> bool;
    fn is_confirmed(&self) -> bool;
    fn valid_password(&self, password: &str) -> bool;
    fn unlock_token(&self) -> Option<&str>;
    fn set_unlock_token(&mut self, token: Option<String>);
    fn lock_access(&mut self, send_instructions: bool, reason: Option<&str>);
    fn unlock_access(&mut self);
    fn confirm(&mut self);
    fn locked_at(&self) -> Option<DateTime<Utc>>;
    fn email_reset_offered_at(&self) -> Option<DateTime<Utc>>;
    fn set_email_reset_offered_at(&mut self, time: Option<DateTime<Utc>>);
    fn secondary_emails(&self) -> Vec<String>;
    fn find_confirmed_email(&self, email: &str) -> Option<String>;
}

// Define the EmailVerificationService trait
pub trait EmailVerificationService: Send + Sync {
    fn generate_token(&self, attr: &str, user: &dyn User) -> (String, String);
    fn validate_token(&self, attr: &str, user: &dyn User, token: &str) -> VerificationResult;
}

// Define the VerificationResult struct
#[derive(Debug, Serialize, Deserialize)]
pub struct VerificationResult {
    pub status: String,
    pub reason: Option<String>,
    pub message: Option<String>,
}

// Define the AuthenticationEvent trait
pub trait AuthenticationEvent: Send + Sync {
    fn initial_login_or_known_ip_address(user: &dyn User, ip: &str) -> bool;
}

// Define the Feature trait
pub trait Feature: Send + Sync {
    fn is_enabled(feature: &str, user: &dyn User) -> bool;
    fn is_disabled(feature: &str, user: &dyn User, feature_type: &str) -> bool;
}

// Define the ApplicationRateLimiter trait
pub trait ApplicationRateLimiter: Send + Sync {
    fn is_throttled(rate_limit: &str, scope: &dyn User) -> bool;
    fn rate_limits(rate_limit: &str) -> RateLimit;
}

// Define the RateLimit struct
pub struct RateLimit {
    pub interval: i64,
}

// Define the VerifiesWithEmail trait
pub trait VerifiesWithEmail {
    fn verify_with_email(&self, req: &HttpRequest) -> Result<HttpResponse>;
    fn resend_verification_code(&self, req: &HttpRequest) -> Result<HttpResponse>;
    fn update_email(&self, req: &HttpRequest) -> Result<HttpResponse>;
    fn successful_verification(&self, req: &HttpRequest) -> Result<HttpResponse>;
    fn find_user(&self, req: &HttpRequest) -> Option<Arc<dyn User>>;
    fn verification_email(&self, user: &dyn User) -> String;
    fn unconfirmed_verification_email(&self, user: &dyn User) -> bool;
    fn two_factor_enabled(&self) -> bool;
    fn is_qa_request(&self, user_agent: Option<&str>) -> bool;
    fn authentication_method(&self) -> &str;
    fn sign_in(&self, user: &dyn User);
    fn log_audit_event(&self, current_user: &dyn User, user: &dyn User, auth_method: &str);
    fn log_user_activity(&self, user: &dyn User);
    fn verify_known_sign_in(&self);
    fn after_sign_in_path(&self, user: &dyn User) -> String;
    fn get_session(&self) -> &dyn Session;
    fn set_session_value(&self, key: &str, value: i64);
    fn remove_session_value(&self, key: &str);
    fn get_request_ip(&self, req: &HttpRequest) -> String;
    fn get_user_agent(&self, req: &HttpRequest) -> Option<String>;
    fn get_user_params(&self, req: &HttpRequest) -> UserParams;
    fn get_verification_params(&self, req: &HttpRequest) -> VerificationParams;
    fn get_email_params(&self, req: &HttpRequest) -> EmailParams;
    fn render_json(&self, data: impl Serialize) -> Result<HttpResponse>;
    fn render_template(&self, template: &str, data: impl Serialize) -> Result<HttpResponse>;
    fn redirect_to(&self, path: &str, message: Option<&str>) -> Result<HttpResponse>;
    fn format_time_interval(&self, seconds: i64) -> String;
}

// Define the Session trait
pub trait Session {
    fn get(&self, key: &str) -> Option<i64>;
    fn set(&self, key: &str, value: i64);
    fn remove(&self, key: &str);
}

// Define the UserParams struct
#[derive(Debug, Deserialize)]
pub struct UserParams {
    pub password: Option<String>,
}

// Define the VerificationParams struct
#[derive(Debug, Deserialize)]
pub struct VerificationParams {
    pub verification_token: Option<String>,
}

// Define the EmailParams struct
#[derive(Debug, Deserialize)]
pub struct EmailParams {
    pub email: Option<String>,
}

// Define the VerifiesWithEmailHandler struct
pub struct VerifiesWithEmailHandler {
    email_verification_service: Arc<dyn EmailVerificationService>,
    authentication_event: Arc<dyn AuthenticationEvent>,
    feature: Arc<dyn Feature>,
    application_rate_limiter: Arc<dyn ApplicationRateLimiter>,
    token_valid_for_minutes: i64,
}

impl VerifiesWithEmailHandler {
    pub fn new(
        email_verification_service: Arc<dyn EmailVerificationService>,
        authentication_event: Arc<dyn AuthenticationEvent>,
        feature: Arc<dyn Feature>,
        application_rate_limiter: Arc<dyn ApplicationRateLimiter>,
        token_valid_for_minutes: i64,
    ) -> Self {
        VerifiesWithEmailHandler {
            email_verification_service,
            authentication_event,
            feature,
            application_rate_limiter,
            token_valid_for_minutes,
        }
    }

    fn skip_verify_with_email(&self, two_factor_enabled: bool, user_agent: Option<&str>) -> bool {
        two_factor_enabled || self.is_qa_request(user_agent)
    }

    fn find_verification_user(&self, session: &dyn Session) -> Option<Arc<dyn User>> {
        // In a real implementation, this would find the user by ID from the session
        // For now, we'll return None
        None
    }

    fn send_verification_instructions(
        &self,
        user: &mut dyn User,
        secondary_email: Option<&str>,
        reason: Option<&str>,
    ) -> Result<()> {
        let (raw_token, encrypted_token) = self
            .email_verification_service
            .generate_token("unlock_token", user);

        user.set_unlock_token(Some(encrypted_token));
        user.lock_access(false, reason);

        let email = secondary_email.unwrap_or_else(|| self.verification_email(user));
        self.send_verification_instructions_email(user, &raw_token, secondary_email);

        Ok(())
    }

    fn send_verification_instructions_email(
        &self,
        user: &dyn User,
        token: &str,
        secondary_email: Option<&str>,
    ) {
        // In a real implementation, this would send an email
        // For now, we'll just log it
        println!(
            "Sending verification instructions to {} with token {}",
            secondary_email.unwrap_or_else(|| self.verification_email(user)),
            token
        );

        self.log_verification(user, "instructions_sent", None);
    }

    fn verify_email(&self, user: &mut dyn User, password: Option<&str>) -> Result<()> {
        if let Some(token) = user.unlock_token() {
            // Prompt for the token if it already has been set. If the token has expired, send a new one.
            if self.unlock_token_expired(user) {
                self.send_verification_instructions(user, None, None)?;
            }
            self.prompt_for_email_verification(user)?;
        } else if user.is_locked() || !self.trusted_ip_address(user) {
            // require email verification if:
            // - their account has been locked because of too many failed login attempts, or
            // - they have logged in before, but never from the current ip address
            let reason = if !user.is_locked() {
                Some("sign in from untrusted IP address")
            } else {
                None
            };

            if !self.send_rate_limited(user) {
                self.send_verification_instructions(user, None, reason)?;
            }

            self.prompt_for_email_verification(user)?;
        }

        Ok(())
    }

    fn verify_token(&self, user: &mut dyn User, token: &str) -> Result<VerificationResponse> {
        let result = self
            .email_verification_service
            .validate_token("unlock_token", user, token);

        if result.status == "success" {
            self.handle_verification_success(user)?;

            Ok(VerificationResponse {
                status: "success".to_string(),
                message: None,
                redirect_path: Some("/users/successful_verification".to_string()),
            })
        } else {
            self.handle_verification_failure(
                user,
                result.reason.as_deref(),
                result.message.as_deref(),
            )?;

            Ok(result)
        }
    }

    fn render_sign_in_rate_limited(&self) -> Result<HttpResponse> {
        let interval = self.format_time_interval(
            self.application_rate_limiter
                .rate_limits("user_sign_in")
                .interval,
        );

        let message = format!(
            "Maximum login attempts exceeded. Wait {} and try again.",
            interval
        );

        self.redirect_to("/users/sign_in", Some(&message))
    }

    fn send_rate_limited(&self, user: &dyn User) -> bool {
        self.application_rate_limiter
            .is_throttled("email_verification_code_send", user)
    }

    fn handle_verification_failure(
        &self,
        user: &dyn User,
        reason: Option<&str>,
        message: Option<&str>,
    ) -> Result<()> {
        // In a real implementation, this would add an error to the user
        // For now, we'll just log it
        println!(
            "Verification failure for user {}: {} - {}",
            user.username(),
            reason.unwrap_or("unknown"),
            message.unwrap_or("No message")
        );

        self.log_verification(user, "failed_attempt", reason);

        Ok(())
    }

    fn handle_verification_success(&self, user: &mut dyn User) -> Result<()> {
        if self.unconfirmed_verification_email(user) {
            user.confirm();
        }

        if user.email_reset_offered_at().is_none() {
            user.set_email_reset_offered_at(Some(Utc::now()));
        }

        user.unlock_access();
        self.log_verification(user, "successful", None);

        // In a real implementation, this would sign in the user
        // For now, we'll just log it
        println!("User {} signed in successfully", user.username());

        // These would be implemented in the concrete class
        // self.sign_in(user);
        // self.log_audit_event(current_user, user, self.authentication_method());
        // self.log_user_activity(user);
        // self.verify_known_sign_in();

        Ok(())
    }

    fn trusted_ip_address(&self, user: &dyn User) -> bool {
        // In a real implementation, this would check if the IP is trusted
        // For now, we'll just return true
        true
    }

    fn prompt_for_email_verification(&self, user: &dyn User) -> Result<()> {
        // In a real implementation, this would set the session and render a template
        // For now, we'll just log it
        println!(
            "Prompting for email verification for user {}",
            user.username()
        );

        Ok(())
    }

    fn user_secondary_email(&self, user: &dyn User, email: Option<&str>) -> Option<String> {
        if let Some(email) = email {
            user.find_confirmed_email(email)
        } else {
            None
        }
    }

    fn log_verification(&self, user: &dyn User, event: &str, reason: Option<&str>) {
        // In a real implementation, this would log the verification
        // For now, we'll just print it
        println!(
            "Email Verification: {} for user {} from IP {} - Reason: {}",
            event,
            user.username(),
            "127.0.0.1", // This would be the actual IP in a real implementation
            reason.unwrap_or("None")
        );
    }

    fn require_email_verification_enabled(&self, user: &dyn User) -> bool {
        self.feature.is_enabled("require_email_verification", user)
            && self
                .feature
                .is_disabled("skip_require_email_verification", user, "ops")
    }

    fn unlock_token_expired(&self, user: &dyn User) -> bool {
        if let Some(locked_at) = user.locked_at() {
            let now = Utc::now();
            let token_valid_until = locked_at + Duration::minutes(self.token_valid_for_minutes);

            now > token_valid_until
        } else {
            false
        }
    }

    fn is_qa_request(&self, user_agent: Option<&str>) -> bool {
        // In a real implementation, this would check if the request is from QA
        // For now, we'll just return false
        false
    }
}

// Implement the VerifiesWithEmail trait for VerifiesWithEmailHandler
impl<T: VerifiesWithEmail> VerifiesWithEmail for T {
    fn verify_with_email(&self, req: &HttpRequest) -> Result<HttpResponse> {
        let user = self
            .find_user(req)
            .or_else(|| self.find_verification_user(self.get_session()));

        if let Some(mut user) = user {
            if user.is_active() {
                if let Some(session) = self.get_session().get("verification_user_id") {
                    if let Some(token) = self.get_verification_params(req).verification_token {
                        // The verification token is submitted, verify it
                        return self.verify_token(&mut *user, &token);
                    }
                } else if self.require_email_verification_enabled(&user) {
                    // Limit the amount of password guesses, since we now display the email verification page
                    // when the password is correct, which could be a giveaway when brute-forced.
                    if self
                        .application_rate_limiter
                        .is_throttled("user_sign_in", &*user)
                    {
                        return self.render_sign_in_rate_limited();
                    }

                    // Verify the email if the user has logged in successfully.
                    if let Some(password) = &self.get_user_params(req).password {
                        if user.valid_password(password) {
                            self.verify_email(&mut *user, Some(password))?;
                        }
                    }
                }
            }
        }

        Ok(HttpResponse::Ok().finish())
    }

    fn resend_verification_code(&self, req: &HttpRequest) -> Result<HttpResponse> {
        if let Some(user) = self.find_verification_user(self.get_session()) {
            if self.send_rate_limited(&user) {
                let interval = self.format_time_interval(
                    self.application_rate_limiter
                        .rate_limits("email_verification_code_send")
                        .interval,
                );

                let message = format!(
                    "You've reached the maximum amount of resends. Wait {} and try again.",
                    interval
                );

                return self.render_json(VerificationResponse {
                    status: "failure".to_string(),
                    message: Some(message),
                    redirect_path: None,
                });
            }

            let secondary_email =
                self.user_secondary_email(&user, self.get_email_params(req).email.as_deref());

            if let Some(email) = &self.get_email_params(req).email {
                if let Some(secondary_email) = secondary_email {
                    self.send_verification_instructions(&mut *user, Some(&secondary_email), None)?;
                }
            } else if self.get_email_params(req).email.is_none() {
                self.send_verification_instructions(&mut *user, None, None)?;
            }

            return self.render_json(VerificationResponse {
                status: "success".to_string(),
                message: None,
                redirect_path: None,
            });
        }

        Ok(HttpResponse::Ok().finish())
    }

    fn update_email(&self, req: &HttpRequest) -> Result<HttpResponse> {
        if let Some(mut user) = self.find_verification_user(self.get_session()) {
            self.log_verification(&user, "email_update_requested", None);

            // In a real implementation, this would update the email
            // For now, we'll just return a success response
            let result = VerificationResult {
                status: "success".to_string(),
                reason: None,
                message: None,
            };

            if result.status == "success" {
                if let Some(email) = &self.get_email_params(req).email {
                    self.send_verification_instructions(&mut *user, Some(email), None)?;
                }
            } else {
                self.handle_verification_failure(
                    &user,
                    result.reason.as_deref(),
                    result.message.as_deref(),
                )?;
            }

            return self.render_json(result);
        }

        Ok(HttpResponse::Ok().finish())
    }

    fn successful_verification(&self, req: &HttpRequest) -> Result<HttpResponse> {
        self.get_session().remove("verification_user_id");

        // In a real implementation, this would get the current user and redirect
        // For now, we'll just return a success response
        Ok(HttpResponse::Ok().finish())
    }

    fn find_user(&self, req: &HttpRequest) -> Option<Arc<dyn User>> {
        // This would be implemented by the concrete class
        None
    }

    fn verification_email(&self, user: &dyn User) -> String {
        user.email().to_string()
    }

    fn unconfirmed_verification_email(&self, user: &dyn User) -> bool {
        !user.is_confirmed()
    }

    fn two_factor_enabled(&self) -> bool {
        // This would be implemented by the concrete class
        false
    }

    fn is_qa_request(&self, user_agent: Option<&str>) -> bool {
        // This would be implemented by the concrete class
        false
    }

    fn authentication_method(&self) -> &str {
        // This would be implemented by the concrete class
        "email"
    }

    fn sign_in(&self, user: &dyn User) {
        // This would be implemented by the concrete class
    }

    fn log_audit_event(&self, current_user: &dyn User, user: &dyn User, auth_method: &str) {
        // This would be implemented by the concrete class
    }

    fn log_user_activity(&self, user: &dyn User) {
        // This would be implemented by the concrete class
    }

    fn verify_known_sign_in(&self) {
        // This would be implemented by the concrete class
    }

    fn after_sign_in_path(&self, user: &dyn User) -> String {
        // This would be implemented by the concrete class
        "/".to_string()
    }

    fn get_session(&self) -> &dyn Session {
        // This would be implemented by the concrete class
        unimplemented!("get_session must be implemented")
    }

    fn set_session_value(&self, key: &str, value: i64) {
        // This would be implemented by the concrete class
    }

    fn remove_session_value(&self, key: &str) {
        // This would be implemented by the concrete class
    }

    fn get_request_ip(&self, req: &HttpRequest) -> String {
        // This would be implemented by the concrete class
        "127.0.0.1".to_string()
    }

    fn get_user_agent(&self, req: &HttpRequest) -> Option<String> {
        // This would be implemented by the concrete class
        None
    }

    fn get_user_params(&self, req: &HttpRequest) -> UserParams {
        // This would be implemented by the concrete class
        UserParams { password: None }
    }

    fn get_verification_params(&self, req: &HttpRequest) -> VerificationParams {
        // This would be implemented by the concrete class
        VerificationParams {
            verification_token: None,
        }
    }

    fn get_email_params(&self, req: &HttpRequest) -> EmailParams {
        // This would be implemented by the concrete class
        EmailParams { email: None }
    }

    fn render_json(&self, data: impl Serialize) -> Result<HttpResponse> {
        // This would be implemented by the concrete class
        Ok(HttpResponse::Ok().json(data))
    }

    fn render_template(&self, template: &str, data: impl Serialize) -> Result<HttpResponse> {
        // This would be implemented by the concrete class
        Ok(HttpResponse::Ok().finish())
    }

    fn redirect_to(&self, path: &str, message: Option<&str>) -> Result<HttpResponse> {
        // This would be implemented by the concrete class
        Ok(HttpResponse::Found().header("Location", path).finish())
    }

    fn format_time_interval(&self, seconds: i64) -> String {
        // This would be implemented by the concrete class
        format!("{} seconds", seconds)
    }
}
