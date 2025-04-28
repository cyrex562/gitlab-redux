// Ported from: orig_app/app/controllers/concerns/project_unauthorized.rb
// Provides controller actions for handling unauthorized project access due to external authorization

use crate::models::Project;
use crate::utils::external_authorization;
use actix_web::error::Error;

pub mod controller_actions {
    use super::*;
    use std::sync::Arc;

    pub fn on_routable_not_found(
        current_user: &crate::models::User,
    ) -> Arc<dyn Fn(&dyn crate::models::Routable, &str) -> Result<(), Error> + Send + Sync> {
        Arc::new(move |routable, _full_path| {
            if let Some(project) = routable.as_any().downcast_ref::<Project>() {
                let label = project.external_authorization_classification_label();
                let mut rejection_reason = None;

                if !external_authorization::access_allowed(current_user, &label) {
                    rejection_reason =
                        external_authorization::rejection_reason(current_user, &label).or(Some(
                            "External authorization denied access to this project".to_string(),
                        ));
                }

                if let Some(reason) = rejection_reason {
                    return Err(actix_web::error::ErrorForbidden(reason));
                }
            }
            Ok(())
        })
    }
}
