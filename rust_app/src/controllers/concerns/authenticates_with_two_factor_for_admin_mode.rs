// Ported from orig_app/app/controllers/concerns/authenticates_with_two_factor_for_admin_mode.rb
// Handles two-factor authentication for admin mode

use std::collections::HashMap;
use std::sync::Arc;

#[derive(Debug, Clone)]
pub struct AdminTwoFactorContext {
    pub session: HashMap<String, String>,
    pub user: Arc<User>,
    pub current_user_mode: Arc<UserMode>,
    pub request_ip: String,
}

pub trait AuthenticatesWithTwoFactorForAdminMode {
    fn admin_mode_prompt_for_two_factor(&mut self, ctx: &mut AdminTwoFactorContext) -> Response;
    fn admin_mode_authenticate_with_two_factor(
        &mut self,
        ctx: &mut AdminTwoFactorContext,
        user_params: &HashMap<String, String>,
    ) -> Response;
    fn admin_mode_authenticate_with_two_factor_via_otp(
        &mut self,
        ctx: &mut AdminTwoFactorContext,
        otp_attempt: &str,
    ) -> Response;
    fn admin_mode_authenticate_with_two_factor_via_webauthn(
        &mut self,
        ctx: &mut AdminTwoFactorContext,
        device_response: &str,
    ) -> Response;
    fn enable_admin_mode(&mut self, ctx: &mut AdminTwoFactorContext) -> Response;
    fn invalid_login_redirect(&self) -> Response;
    fn admin_handle_two_factor_success(&mut self, ctx: &mut AdminTwoFactorContext) -> Response;
    fn admin_handle_two_factor_failure(
        &mut self,
        ctx: &mut AdminTwoFactorContext,
        method: &str,
        message: &str,
    ) -> Response;
}

pub struct AdminTwoFactorHandler;

impl AuthenticatesWithTwoFactorForAdminMode for AdminTwoFactorHandler {
    fn admin_mode_prompt_for_two_factor(&mut self, ctx: &mut AdminTwoFactorContext) -> Response {
        // Set @user for admin views (not needed in Rust, just context)
        if !ctx.user.can_log_in() {
            return handle_locked_user(&ctx.user);
        }
        ctx.session
            .insert("otp_user_id".to_string(), ctx.user.id.to_string());
        setup_webauthn_authentication(&ctx.user, &mut ctx.session);
        Response::render_template("admin/sessions/two_factor", Some("application"))
    }

    fn admin_mode_authenticate_with_two_factor(
        &mut self,
        ctx: &mut AdminTwoFactorContext,
        user_params: &HashMap<String, String>,
    ) -> Response {
        if !ctx.user.can_log_in() {
            return handle_locked_user(&ctx.user);
        }
        if let (Some(otp), Some(_)) = (
            user_params.get("otp_attempt"),
            ctx.session.get("otp_user_id"),
        ) {
            self.admin_mode_authenticate_with_two_factor_via_otp(ctx, otp)
        } else if let (Some(device_response), Some(_)) = (
            user_params.get("device_response"),
            ctx.session.get("otp_user_id"),
        ) {
            self.admin_mode_authenticate_with_two_factor_via_webauthn(ctx, device_response)
        } else if let Some(password) = user_params.get("password") {
            if ctx.user.valid_password(password) {
                self.admin_mode_prompt_for_two_factor(ctx)
            } else {
                self.invalid_login_redirect()
            }
        } else {
            self.invalid_login_redirect()
        }
    }

    fn admin_mode_authenticate_with_two_factor_via_otp(
        &mut self,
        ctx: &mut AdminTwoFactorContext,
        otp_attempt: &str,
    ) -> Response {
        if valid_otp_attempt(&ctx.user, otp_attempt) {
            ctx.session.remove("otp_user_id");
            // user.save! unless read_only
            self.enable_admin_mode(ctx)
        } else {
            self.admin_handle_two_factor_failure(ctx, "OTP", "Invalid two-factor code.")
        }
    }

    fn admin_mode_authenticate_with_two_factor_via_webauthn(
        &mut self,
        ctx: &mut AdminTwoFactorContext,
        device_response: &str,
    ) -> Response {
        if webauthn_authenticate(&ctx.user, device_response, ctx.session.get("challenge")) {
            self.admin_handle_two_factor_success(ctx)
        } else {
            self.admin_handle_two_factor_failure(
                ctx,
                "WebAuthn",
                "Authentication via WebAuthn device failed.",
            )
        }
    }

    fn enable_admin_mode(&mut self, ctx: &mut AdminTwoFactorContext) -> Response {
        if ctx.current_user_mode.enable_admin_mode(true) {
            Response::redirect_with_notice("/admin", "Admin mode enabled")
        } else {
            self.invalid_login_redirect()
        }
    }

    fn invalid_login_redirect(&self) -> Response {
        Response::render_with_alert("new", "Invalid login or password")
    }

    fn admin_handle_two_factor_success(&mut self, ctx: &mut AdminTwoFactorContext) -> Response {
        ctx.session.remove("otp_user_id");
        ctx.session.remove("challenge");
        self.enable_admin_mode(ctx)
    }

    fn admin_handle_two_factor_failure(
        &mut self,
        ctx: &mut AdminTwoFactorContext,
        method: &str,
        message: &str,
    ) -> Response {
        ctx.user.increment_failed_attempts();
        log_failed_two_factor(&ctx.user, method, &ctx.request_ip);
        // flash.now[:alert] = message
        self.admin_mode_prompt_for_two_factor(ctx)
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
}

#[derive(Debug, Clone)]
pub struct UserMode {
    pub admin_mode: bool,
}

impl UserMode {
    pub fn enable_admin_mode(&self, _skip_password: bool) -> bool {
        true
    }
}

pub struct Response;

impl Response {
    pub fn render_template(_template: &str, _layout: Option<&str>) -> Self {
        Response
    }
    pub fn redirect_with_notice(_path: &str, _notice: &str) -> Self {
        Response
    }
    pub fn render_with_alert(_view: &str, _alert: &str) -> Self {
        Response
    }
}

fn handle_locked_user(_user: &User) -> Response {
    Response
}
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
fn log_failed_two_factor(_user: &User, _method: &str, _ip: &str) {}
