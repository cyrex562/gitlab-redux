use actix_web::{web, HttpResponse, Responder};
use std::fmt::Display;

use crate::models::flash_message::FlashMessage;
use crate::models::spammable::Spammable;
use crate::models::user::User;
use crate::services::spam::akismet_mark_as_spam::AkismetMarkAsSpamService;

pub trait AkismetMarkAsSpamAction {
    fn current_user(&self) -> &User;
    fn spammable(&self) -> &dyn Spammable;
    fn spammable_path(&self) -> String;
    fn flash(&self) -> &FlashMessage;
    fn flash_mut(&mut self) -> &mut FlashMessage;
    fn access_denied(&self) -> HttpResponse;

    fn mark_as_spam(&mut self) -> impl Responder {
        if !self.authorize_submit_spammable() {
            return self.access_denied();
        }

        let service = AkismetMarkAsSpamService::new(self.spammable());

        if service.execute() {
            let spammable_titlecase = self.spammable().spammable_entity_type().to_titlecase();
            self.flash_mut().add_notice(format!(
                "{} was submitted to Akismet successfully.",
                spammable_titlecase
            ));

            HttpResponse::Found()
                .header("Location", self.spammable_path())
                .finish()
        } else {
            self.flash_mut()
                .add_alert("Error with Akismet. Please check the logs for more info.");

            HttpResponse::Found()
                .header("Location", self.spammable_path())
                .finish()
        }
    }

    fn authorize_submit_spammable(&self) -> bool {
        self.current_user().can_admin_all_resources()
    }
}

pub trait ToTitlecase {
    fn to_titlecase(&self) -> String;
}

impl<T: Display> ToTitlecase for T {
    fn to_titlecase(&self) -> String {
        let s = self.to_string();
        if s.is_empty() {
            return s;
        }

        let mut chars: Vec<char> = s.chars().collect();
        chars[0] = chars[0].to_uppercase().next().unwrap();

        chars.into_iter().collect()
    }
}
