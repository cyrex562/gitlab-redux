use actix_web::{HttpRequest, HttpResponse};
use std::sync::OnceLock;

pub struct Settings {
    pub email_confirmation_setting_soft: bool,
}

impl Settings {
    pub fn get() -> &'static Settings {
        static INSTANCE: OnceLock<Settings> = OnceLock::new();
        INSTANCE.get_or_init(|| Settings {
            email_confirmation_setting_soft: true,
        })
    }
}

pub struct User {
    pub email: String,
    pub unconfirmed_email: Option<String>,
    pub confirmed: bool,
}

pub trait ConfirmEmailWarning {
    fn show_confirm_warning(&self, req: &HttpRequest) -> bool {
        self.is_html_request(req)
            && req.method() == "GET"
            && Settings::get().email_confirmation_setting_soft
    }

    fn set_confirm_warning(&self, req: &HttpRequest) -> Option<String> {
        let user = self.get_current_user()?;
        if user.confirmed {
            return None;
        }

        let email = self.get_email(user);
        let warning_message = format!(
            "Please check your email ({}) to verify that you own this address and unlock the power of CI/CD. \
            Didn't receive it? <a href='/users/confirmation?email={}'>Resend it</a>. \
            Wrong email address? <a href='/users/settings/profile'>Update it</a>.",
            self.escape_html(&email),
            self.escape_html(&email)
        );

        Some(warning_message)
    }

    fn get_email(&self, user: &User) -> String {
        user.unconfirmed_email
            .clone()
            .unwrap_or_else(|| user.email.clone())
    }

    fn escape_html(&self, text: &str) -> String {
        text.replace("&", "&amp;")
            .replace("<", "&lt;")
            .replace(">", "&gt;")
            .replace("\"", "&quot;")
            .replace("'", "&#39;")
    }

    // Required methods to be implemented by concrete types
    fn is_html_request(&self, req: &HttpRequest) -> bool;
    fn get_current_user(&self) -> Option<&User>;
}
