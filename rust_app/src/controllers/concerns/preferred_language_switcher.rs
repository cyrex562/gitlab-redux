use actix_web::{HttpRequest, HttpResponse};
use once_cell::sync::Lazy;
use std::collections::HashMap;

static AVAILABLE_LOCALES: Lazy<Vec<String>> = Lazy::new(|| {
    vec![
        "en".to_string(),
        "es".to_string(),
        "fr".to_string(),
        "de".to_string(),
        "it".to_string(),
        "ja".to_string(),
        "ko".to_string(),
        "nl".to_string(),
        "pl".to_string(),
        "pt_BR".to_string(),
        "ru".to_string(),
        "zh_CN".to_string(),
        "zh_TW".to_string(),
    ]
});

pub trait PreferredLanguageSwitcher {
    fn init_preferred_language(&self, req: &HttpRequest) -> HttpResponse {
        if self.is_feature_enabled("disable_preferred_language_cookie") {
            return HttpResponse::Ok().finish();
        }

        let preferred_lang = self.preferred_language(req);
        self.set_cookie("preferred_language", &preferred_lang);

        HttpResponse::Ok().finish()
    }

    fn preferred_language(&self, req: &HttpRequest) -> String {
        // Check cookie first
        if let Some(cookie_lang) = self.get_cookie("preferred_language") {
            if AVAILABLE_LOCALES.contains(&cookie_lang) {
                return cookie_lang;
            }
        }

        // Check language from params
        if let Some(lang) = self.selectable_language(&self.language_from_params()) {
            return lang;
        }

        // Check browser languages
        if let Some(lang) = self.selectable_language(&self.browser_languages(req)) {
            return lang;
        }

        // Default to system setting
        self.default_preferred_language()
            .unwrap_or_else(|| "en".to_string())
    }

    fn selectable_language(&self, language_options: &[String]) -> Option<String> {
        let selectable_codes = self.ordered_selectable_locales_codes();
        language_options
            .iter()
            .find(|lang| selectable_codes.contains(lang))
            .cloned()
    }

    fn ordered_selectable_locales_codes(&self) -> Vec<String> {
        self.ordered_selectable_locales()
            .iter()
            .map(|locale| locale.value.clone())
            .collect()
    }

    fn browser_languages(&self, req: &HttpRequest) -> Vec<String> {
        if let Some(accept_language) = req.headers().get("accept-language") {
            if let Ok(accept_language_str) = accept_language.to_str() {
                let formatted = accept_language_str.replace('-', "_");
                return formatted
                    .split(|c| c == ',' || c == ';')
                    .filter(|s| !s.starts_with('q'))
                    .map(|s| s.trim().to_string())
                    .collect();
            }
        }

        Vec::new()
    }

    fn language_from_params(&self) -> Vec<String> {
        // Overridden in EE version
        Vec::new()
    }

    // Required methods to be implemented by concrete types
    fn is_feature_enabled(&self, feature_name: &str) -> bool;
    fn get_cookie(&self, name: &str) -> Option<String>;
    fn set_cookie(&self, name: &str, value: &str);
    fn default_preferred_language(&self) -> Option<String>;
    fn ordered_selectable_locales(&self) -> Vec<Locale>;
}

#[derive(Debug, Clone)]
pub struct Locale {
    pub value: String,
    pub name: String,
}
