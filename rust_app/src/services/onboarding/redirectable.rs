use crate::onboarding::{InvitedMember, MemberSource, OnboardingStatusPresenter, User};
use actix_web::{web, HttpResponse};

pub trait Redirectable {
    // TODO: Implement onboarding_status_presenter getter
    fn onboarding_status_presenter(&self) -> &dyn OnboardingStatusPresenter;

    // TODO: Implement stored_location_for getter
    fn stored_location_for(&self, user: &User) -> Option<String>;

    fn after_sign_up_path(&self, user: &User) -> HttpResponse {
        if self.onboarding_status_presenter().single_invite() {
            if let Some(member) = self.onboarding_status_presenter().last_invited_member() {
                // TODO: Implement invite_accepted_notice helper
                let notice = format!("Invite accepted for member {}", member.id());

                if let Some(source) = self
                    .onboarding_status_presenter()
                    .last_invited_member_source()
                {
                    return HttpResponse::Found()
                        .header("Location", self.polymorphic_path(source))
                        .json(serde_json::json!({
                            "notice": notice
                        }));
                }
            }
        }

        // If we get here, either there are multiple invites or no invites
        HttpResponse::Found()
            .header("Location", self.path_for_signed_in_user(user))
            .finish()
    }

    fn path_for_signed_in_user(&self, user: &User) -> String {
        self.stored_location_for(user)
            .unwrap_or_else(|| self.last_member_source_path())
    }

    fn last_member_source_path(&self) -> String {
        if let Some(source) = self
            .onboarding_status_presenter()
            .last_invited_member_source()
        {
            self.polymorphic_path(source)
        } else {
            "/dashboard/projects".to_string()
        }
    }

    fn polymorphic_path(&self, source: &dyn MemberSource) -> String {
        format!("/{}/{}", source.type_name().to_lowercase(), source.id())
    }
}
