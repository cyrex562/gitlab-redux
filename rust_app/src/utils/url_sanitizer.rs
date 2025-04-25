// Ported from: orig_app/lib/gitlab/url_sanitizer.rb
//
// Minimal Rust implementation for building a full URL with optional credentials.

pub struct UrlSanitizer<'a> {
    url: &'a str,
    credentials: Option<(String, String)>,
}

impl<'a> UrlSanitizer<'a> {
    pub fn new(url: &'a str, credentials: Option<(String, String)>) -> Self {
        Self { url, credentials }
    }

    pub fn full_url(&self) -> String {
        if let Some((ref user, ref password)) = self.credentials {
            if !user.is_empty() || !password.is_empty() {
                // Insert credentials into the URL (basic implementation)
                if let Ok(mut parsed) = url::Url::parse(self.url) {
                    if !user.is_empty() {
                        parsed.set_username(user).ok();
                    }
                    if !password.is_empty() {
                        parsed.set_password(Some(password)).ok();
                    }
                    return parsed.to_string();
                }
            }
        }
        self.url.to_string()
    }
}
