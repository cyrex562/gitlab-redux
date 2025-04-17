use actix_web::{HttpRequest, HttpResponse};
use std::sync::OnceLock;

pub struct Settings {
    pub invisible_captcha_enabled: bool,
}

impl Settings {
    pub fn get() -> &'static Settings {
        static INSTANCE: OnceLock<Settings> = OnceLock::new();
        INSTANCE.get_or_init(|| Settings {
            invisible_captcha_enabled: true,
        })
    }
}

pub struct Metrics {
    honeypot_counter: Counter,
    timestamp_counter: Counter,
}

struct Counter {
    name: String,
    description: String,
    value: std::sync::atomic::AtomicU64,
}

impl Counter {
    fn new(name: &str, description: &str) -> Self {
        Counter {
            name: name.to_string(),
            description: description.to_string(),
            value: std::sync::atomic::AtomicU64::new(0),
        }
    }

    fn increment(&self) {
        self.value.fetch_add(1, std::sync::atomic::Ordering::SeqCst);
    }
}

pub trait InvisibleCaptchaOnSignup {
    fn on_honeypot_spam_callback(&self, req: &HttpRequest) -> HttpResponse {
        if !Settings::get().invisible_captcha_enabled {
            return HttpResponse::Ok().finish();
        }

        self.invisible_captcha_honeypot_counter().increment();
        self.log_request("Invisible_Captcha_Honeypot_Request", req);

        HttpResponse::Ok().finish()
    }

    fn on_timestamp_spam_callback(&self, req: &HttpRequest) -> HttpResponse {
        if !Settings::get().invisible_captcha_enabled {
            return HttpResponse::Ok().finish();
        }

        self.invisible_captcha_timestamp_counter().increment();
        self.log_request("Invisible_Captcha_Timestamp_Request", req);

        HttpResponse::Found()
            .header("Location", "/users/sign_in")
            .header("X-Flash-Message", "Invalid timestamp")
            .finish()
    }

    fn invisible_captcha_honeypot_counter(&self) -> &'static Counter {
        static COUNTER: OnceLock<Counter> = OnceLock::new();
        COUNTER.get_or_init(|| {
            Counter::new(
                "bot_blocked_by_invisible_captcha_honeypot",
                "Counter of blocked sign up attempts with filled honeypot",
            )
        })
    }

    fn invisible_captcha_timestamp_counter(&self) -> &'static Counter {
        static COUNTER: OnceLock<Counter> = OnceLock::new();
        COUNTER.get_or_init(|| {
            Counter::new(
                "bot_blocked_by_invisible_captcha_timestamp",
                "Counter of blocked sign up attempts with invalid timestamp",
            )
        })
    }

    fn log_request(&self, message: &str, req: &HttpRequest) {
        let request_information = serde_json::json!({
            "message": message,
            "env": "invisible_captcha_signup_bot_detected",
            "remote_ip": req.connection_info().peer_addr().unwrap_or("unknown"),
            "request_method": req.method().as_str(),
            "path": req.uri().path()
        });

        // In a real implementation, this would log to your logging system
        println!("{}", request_information);
    }
}
