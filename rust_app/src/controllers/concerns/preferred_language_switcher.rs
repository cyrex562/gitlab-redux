use actix_web::{dev::ServiceRequest, web::Data};
use std::sync::Arc;
use crate::{
    utils::{
        feature_flags::Feature,
        i18n::I18n,
        strong_memoize::StrongMemoize,
    },
    config::Settings,
};

pub trait PreferredLanguageSwitcher {
    fn init_preferred_language(&self, req: &ServiceRequest) -> Result<(), Error>;
    fn preferred_language(&self, req: &ServiceRequest) -> String;
}

pub struct PreferredLanguageSwitcherImpl {
    settings: Data<Settings>,
    i18n: Arc<I18n>,
}

impl PreferredLanguageSwitcherImpl {
    pub fn new(settings: Data<Settings>, i18n: Arc<I18n>) -> Self {
        Self { settings, i18n }
    }

    fn selectable_language(&self, language_options: &[String]) -> Option<String> {
        language_options.iter()
            .find(|lan| self.ordered_selectable_locales_codes().contains(&lan.to_string()))
            .cloned()
    }

    fn ordered_selectable_locales_codes(&self) -> Vec<String> {
        self.i18n.ordered_selectable_locales()
            .iter()
            .map(|locale| locale.value().to_string())
            .collect()
    }

    fn browser_languages(&self, req: &ServiceRequest) -> Vec<String> {
        req.headers()
            .get("accept-language")
            .and_then(|h| h.to_str().ok())
            .map(|header| {
                header.replace('-', "_")
                    .split(|c| c == ',' || c == ';')
                    .filter(|s| !s.starts_with('q'))
                    .map(|s| s.trim().to_string())
                    .collect()
            })
            .unwrap_or_default()
    }

    fn language_from_params(&self, _req: &ServiceRequest) -> Vec<String> {
        vec![] // To be overridden in EE
    }
}

impl PreferredLanguageSwitcher for PreferredLanguageSwitcherImpl {
    fn init_preferred_language(&self, req: &ServiceRequest) -> Result<(), Error> {
        if !Feature::enabled("disable_preferred_language_cookie", None, None) {
            let language = self.preferred_language(req);
            req.cookies_mut().add("preferred_language", language);
        }
        Ok(())
    }

    fn preferred_language(&self, req: &ServiceRequest) -> String {
        req.cookies()
            .get("preferred_language")
            .and_then(|c| c.value().to_string().into())
            .filter(|lang| self.i18n.available_locales().contains(lang))
            .or_else(|| self.selectable_language(&self.language_from_params(req)))
            .or_else(|| self.selectable_language(&self.browser_languages(req)))
            .unwrap_or_else(|| self.settings.default_preferred_language.clone())
    }
} 