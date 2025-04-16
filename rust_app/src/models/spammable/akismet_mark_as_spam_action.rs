use crate::models::user::User;
use crate::spammable::Spammable;
use actix_web::{web, HttpResponse};

pub trait AkismetMarkAsSpamAction {
    // TODO: Implement spammable getter
    fn spammable(&self) -> &dyn Spammable;

    // TODO: Implement spammable_path getter
    fn spammable_path(&self) -> String;

    fn mark_as_spam(&self, user: &User) -> HttpResponse {
        // TODO: Implement authorization check
        if !user.can_admin_all_resources() {
            return HttpResponse::Forbidden().finish();
        }

        // TODO: Implement AkismetMarkAsSpamService
        let success = true; // Placeholder for actual service call

        if success {
            HttpResponse::Found()
                .header("Location", self.spammable_path())
                .json(serde_json::json!({
                    "notice": format!("{} was submitted to Akismet successfully.",
                        self.spammable().spammable_entity_type())
                }))
        } else {
            HttpResponse::Found()
                .header("Location", self.spammable_path())
                .json(serde_json::json!({
                    "alert": "Error with Akismet. Please check the logs for more info."
                }))
        }
    }
}
