// Ported from orig_app/app/controllers/concerns/spammable_actions/akismet_mark_as_spam_action.rb
// Date: 2025-04-24
// This trait provides the Akismet mark-as-spam action for spammable entities.
//
// To use, implement `spammable()` and `spammable_path()` for your controller struct.

use crate::models::user::User;
use crate::spammable::Spammable;
use actix_web::{web, HttpResponse};

pub trait AkismetMarkAsSpamAction {
    /// Returns a reference to the spammable entity.
    fn spammable(&self) -> &dyn Spammable;

    /// Returns the path to redirect to after the action.
    fn spammable_path(&self) -> String;

    /// Handles marking the entity as spam via Akismet.
    fn mark_as_spam(&self, user: &User) -> HttpResponse {
        // Authorization check
        if !user.can_admin_all_resources() {
            return HttpResponse::Forbidden().finish();
        }

        // TODO: Replace with actual AkismetMarkAsSpamService logic
        let success = true; // Placeholder for service call

        if success {
            HttpResponse::Found()
                .header("Location", self.spammable_path())
                .json(serde_json::json!({
                    "notice": format!(
                        "{} was submitted to Akismet successfully.",
                        self.spammable().spammable_entity_type().to_title_case()
                    )
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
