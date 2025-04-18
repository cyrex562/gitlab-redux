use actix_web::error::Error;
use crate::models::Project;
use crate::utils::external_authorization;

pub mod controller_actions {
    use super::*;
    use std::sync::Arc;

    pub fn on_routable_not_found() -> Arc<dyn Fn(&dyn crate::models::Routable, &str) -> Result<(), Error> + Send + Sync> {
        Arc::new(|routable, _full_path| {
            if let Some(project) = routable.as_any().downcast_ref::<Project>() {
                let label = project.external_authorization_classification_label();
                
                if !external_authorization::access_allowed(label) {
                    let rejection_reason = external_authorization::rejection_reason(label)
                        .unwrap_or_else(|| "External authorization denied access to this project".to_string());
                    
                    return Err(actix_web::error::ErrorForbidden(rejection_reason));
                }
            }
            Ok(())
        })
    }
} 