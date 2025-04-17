use std::any::Any;

pub struct Project {
    pub external_authorization_classification_label: Option<String>,
}

pub struct ExternalAuthorization;

impl ExternalAuthorization {
    pub fn access_allowed(user: &User, label: Option<&str>) -> bool {
        // Implementation would depend on your authorization system
        true
    }

    pub fn rejection_reason(user: &User, label: Option<&str>) -> Option<String> {
        // Implementation would depend on your authorization system
        None
    }
}

pub struct User {
    // Add user fields as needed
}

pub mod ControllerActions {
    use super::*;
    use actix_web::HttpResponse;

    pub fn on_routable_not_found() -> Box<dyn Fn(&dyn Any, &str) -> bool> {
        Box::new(|routable, _full_path| {
            if let Some(project) = routable.downcast_ref::<Project>() {
                let label = project
                    .external_authorization_classification_label
                    .as_deref();
                let mut rejection_reason = None;

                if !ExternalAuthorization::access_allowed(get_current_user(), label) {
                    rejection_reason = ExternalAuthorization::rejection_reason(
                        get_current_user(),
                        label,
                    )
                    .or_else(|| {
                        Some("External authorization denied access to this project".to_string())
                    });
                }

                if let Some(reason) = rejection_reason {
                    access_denied(&reason);
                    return false;
                }
            }
            true
        })
    }

    // These functions would be provided by the controller implementation
    fn get_current_user() -> &'static User {
        // This would be implemented by the controller
        unimplemented!("get_current_user must be implemented by the controller")
    }

    fn access_denied(reason: &str) {
        // This would be implemented by the controller
        // It would typically render a 403 response with the reason
        unimplemented!("access_denied must be implemented by the controller")
    }
}
