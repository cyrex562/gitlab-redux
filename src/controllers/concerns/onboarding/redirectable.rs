use actix_web::{HttpResponse, Responder};
use std::sync::Arc;

use crate::models::onboarding::OnboardingStatusPresenter;
use crate::utils::flash::FlashMessage;

pub trait Redirectable {
    fn onboarding_status_presenter(&self) -> Arc<OnboardingStatusPresenter>;
    fn stored_location_for(&self, key: &str) -> Option<String>;
    fn dashboard_projects_path(&self) -> String;
    fn polymorphic_path(&self, resource: &dyn std::any::Any) -> String;
    fn set_flash_notice(&self, message: String);

    fn after_sign_up_path(&self) -> impl Responder {
        if self.onboarding_status_presenter().single_invite() {
            let last_invited_member = self.onboarding_status_presenter().last_invited_member();
            let notice = format!("Invite accepted for {}", last_invited_member.name());
            self.set_flash_notice(notice);

            let source = self
                .onboarding_status_presenter()
                .last_invited_member_source();
            self.polymorphic_path(source.as_ref())
        } else {
            // Invites will come here if there is more than 1
            self.path_for_signed_in_user()
        }
    }

    fn path_for_signed_in_user(&self) -> impl Responder {
        self.stored_location_for("user")
            .unwrap_or_else(|| self.last_member_source_path())
    }

    fn last_member_source_path(&self) -> String {
        let source = self
            .onboarding_status_presenter()
            .last_invited_member_source();

        if source.is_none() {
            return self.dashboard_projects_path();
        }

        self.polymorphic_path(source.unwrap().as_ref())
    }
}
