// Ported from orig_app/app/controllers/concerns/authenticates_with_two_factor.rb
// Handles two-factor authentication for general user sessions

use std::collections::HashMap;
use std::sync::Arc;

#[derive(Debug, Clone)]
pub struct TwoFactorContext {
    pub session: HashMap<String, String>,
    pub user: Arc<User>,
    pub request_ip: String,
}

pub trait AuthenticatesWithTwoFactor {
    fn prompt_for_two_factor(&mut self, ctx: &mut TwoFactorContext) -> Response;
    fn authenticate_with_two_factor(
        &mut self,
        ctx: &mut TwoFactorContext,
        user_params: &HashMap<String, String>,
    ) -> Response;
    fn authenticate_with_two_factor_via_otp(
        &mut self,
        ctx: &mut TwoFactorContext,
        otp_attempt: &str,
    ) -> Response;
    fn authenticate_with_two_factor_via_webauthn(
        &mut self,
        ctx: &mut TwoFactorContext,
        device_response: &str,
    ) -> Response;
    fn handle_locked_user(&mut self, ctx: &mut TwoFactorContext) -> Response;
    fn handle_changed_user(&mut self, ctx: &mut TwoFactorContext) -> Response;
    fn handle_two_factor_success(&mut self, ctx: &mut TwoFactorContext) -> Response;
    fn handle_two_factor_failure(
        &mut self,
        ctx: &mut TwoFactorContext,
        method: &str,
        message: &str,
    ) -> Response;
}

pub struct TwoFactorHandler;

impl AuthenticatesWithTwoFactor for TwoFactorHandler {
    fn prompt_for_two_factor(&mut self, ctx: &mut TwoFactorContext) -> Response {
        if !ctx.user.can_log_in() {
            return self.handle_locked_user(ctx);
        }
        ctx.session
            .insert("otp_user_id".to_string(), ctx.user.id.to_string());
        ctx.session.insert(
            "user_password_hash".to_string(),
            hash_password(&ctx.user.encrypted_password),
        );
        add_gon_variables();
        setup_webauthn_authentication(&ctx.user, &mut ctx.session);
        Response::render_template("devise/sessions/two_factor", None)
    }

    fn authenticate_with_two_factor(
        &mut self,
        ctx: &mut TwoFactorContext,
        user_params: &HashMap<String, String>,
    ) -> Response {
        if !ctx.user.can_log_in() {
            return self.handle_locked_user(ctx);
        }
        if user_password_changed(&ctx.user, &ctx.session) {
            return self.handle_changed_user(ctx);
        }
        if let (Some(otp), Some(_)) = (
            user_params.get("otp_attempt"),
            ctx.session.get("otp_user_id"),
        ) {
            self.authenticate_with_two_factor_via_otp(ctx, otp)
        } else if let (Some(device_response), Some(_)) = (
            user_params.get("device_response"),
            ctx.session.get("otp_user_id"),
        ) {
            self.authenticate_with_two_factor_via_webauthn(ctx, device_response)
        } else if let Some(password) = user_params.get("password") {
            if ctx.user.valid_password(password) {
                self.prompt_for_two_factor(ctx)
            } else {
                Response::render_with_alert("new", "Invalid login or password")
            }
        } else {
            Response::render_with_alert("new", "Invalid login or password")
        }
    }

    fn authenticate_with_two_factor_via_otp(
        &mut self,
        ctx: &mut TwoFactorContext,
        otp_attempt: &str,
    ) -> Response {
        if valid_otp_attempt(&ctx.user, otp_attempt) {
            clear_two_factor_attempt(&mut ctx.session);
            remember_me(&ctx.user, &ctx.session);
            ctx.user.save();
            sign_in(&ctx.user, "two_factor_authenticated", "authentication");
            self.handle_two_factor_success(ctx)
        } else {
            send_two_factor_otp_attempt_failed_email(&ctx.user, &ctx.request_ip);
            self.handle_two_factor_failure(ctx, "OTP", "Invalid two-factor code.")
        }
    }

    fn authenticate_with_two_factor_via_webauthn(
        &mut self,
        ctx: &mut TwoFactorContext,
        device_response: &str,
    ) -> Response {
        if webauthn_authenticate(&ctx.user, device_response, ctx.session.get("challenge")) {
            self.handle_two_factor_success(ctx)
        } else {
            self.handle_two_factor_failure(
                ctx,
                "WebAuthn",
                "Authentication via WebAuthn device failed.",
            )
        }
    }

    fn handle_locked_user(&mut self, ctx: &mut TwoFactorContext) -> Response {
        clear_two_factor_attempt(&mut ctx.session);
        Response::redirect_with_alert("/users/sign_in", locked_user_redirect_alert(&ctx.user))
    }

    fn handle_changed_user(&mut self, ctx: &mut TwoFactorContext) -> Response {
        clear_two_factor_attempt(&mut ctx.session);
        Response::redirect_with_alert("/users/sign_in", "An error occurred. Please sign in again.")
    }

    fn handle_two_factor_success(&mut self, ctx: &mut TwoFactorContext) -> Response {
        clear_two_factor_attempt(&mut ctx.session);
        remember_me(&ctx.user, &ctx.session);
        sign_in(&ctx.user, "two_factor_authenticated", "authentication");
        Response::redirect_with_notice("/", "Signed in with two-factor authentication.")
    }

    fn handle_two_factor_failure(
        &mut self,
        ctx: &mut TwoFactorContext,
        method: &str,
        message: &str,
    ) -> Response {
        ctx.user.increment_failed_attempts();
        log_failed_two_factor(&ctx.user, method, &ctx.request_ip);
        Response::render_with_alert("devise/sessions/two_factor", message)
    }
}

// Helper stubs and types for integration
#[derive(Debug, Clone)]
pub struct User {
    pub id: i64,
    pub username: String,
    pub encrypted_password: String,
    pub locked: bool,
}

impl User {
    pub fn can_log_in(&self) -> bool {
        !self.locked
    }
    pub fn valid_password(&self, _pw: &str) -> bool {
        true
    } // stub
    pub fn increment_failed_attempts(&self) {}
    pub fn save(&self) {}
}

pub struct Response;

impl Response {
    pub fn render_template(_template: &str, _layout: Option<&str>) -> Self {
        Response
    }
    pub fn render_with_alert(_view: &str, _alert: &str) -> Self {
        Response
    }
    pub fn redirect_with_alert(_path: &str, _alert: &str) -> Self {
        Response
    }
    pub fn redirect_with_notice(_path: &str, _notice: &str) -> Self {
        Response
    }
}

fn add_gon_variables() {}
fn setup_webauthn_authentication(_user: &User, _session: &mut HashMap<String, String>) {}
fn valid_otp_attempt(_user: &User, _otp: &str) -> bool {
    true
}
fn webauthn_authenticate(
    _user: &User,
    _device_response: &str,
    _challenge: Option<&String>,
) -> bool {
    true
}
fn clear_two_factor_attempt(session: &mut HashMap<String, String>) {
    session.remove("otp_user_id");
    session.remove("user_password_hash");
    session.remove("challenge");
}
fn remember_me(_user: &User, _session: &HashMap<String, String>) {}
fn sign_in(_user: &User, _message: &str, _event: &str) {}
fn send_two_factor_otp_attempt_failed_email(_user: &User, _ip: &str) {}
fn log_failed_two_factor(_user: &User, _method: &str, _ip: &str) {}
fn hash_password(_pw: &str) -> String {
    "hashed".to_string()
}
fn user_password_changed(user: &User, session: &HashMap<String, String>) -> bool {
    if let Some(hash) = session.get("user_password_hash") {
        hash_password(&user.encrypted_password) != *hash
    } else {
        false
    }
}
fn locked_user_redirect_alert(_user: &User) -> &'static str {
    "Your account is locked."
}
