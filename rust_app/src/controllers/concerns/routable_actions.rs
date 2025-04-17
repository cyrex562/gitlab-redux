use actix_web::{HttpRequest, HttpResponse};
use std::any::Any;

pub trait Routable {
    fn find_by_full_path(&self, full_path: &str, follow_redirects: bool) -> Option<Self>
    where
        Self: Sized;
    fn full_path(&self) -> String;
    fn class_name(&self) -> String;
}

pub struct ProjectUnauthorized;

impl ProjectUnauthorized {
    pub fn on_routable_not_found() -> Box<dyn Fn(&dyn Any, &str) -> bool> {
        Box::new(|_routable, _full_path| {
            // Implementation for handling unauthorized access
            false
        })
    }
}

pub trait RoutableActions {
    fn find_routable<T: Routable + 'static>(
        &self,
        req: &HttpRequest,
        routable_klass: T,
        routable_full_path: &str,
        full_path: &str,
        extra_authorization_proc: Option<Box<dyn Fn(&dyn Any) -> bool>>,
    ) -> Option<T> {
        let routable = routable_klass.find_by_full_path(routable_full_path, req.method() == "GET");

        if self.routable_authorized(&routable, extra_authorization_proc.as_deref()) {
            self.ensure_canonical_path(&routable, routable_full_path);
            Some(routable)
        } else {
            self.perform_not_found_actions(&routable, self.not_found_actions(), full_path);
            None
        }
    }

    fn not_found_actions(&self) -> Vec<Box<dyn Fn(&dyn Any, &str) -> bool>> {
        vec![ProjectUnauthorized::on_routable_not_found()]
    }

    fn perform_not_found_actions(
        &self,
        routable: &dyn Any,
        actions: Vec<Box<dyn Fn(&dyn Any, &str) -> bool>>,
        full_path: &str,
    ) {
        for action in actions {
            if !action(routable, full_path) {
                break;
            }
        }
    }

    fn routable_authorized(
        &self,
        routable: &Option<impl Routable>,
        extra_authorization_proc: Option<&dyn Fn(&dyn Any) -> bool>,
    ) -> bool {
        if let Some(routable) = routable {
            let action = format!("read_{}", routable.class_name().to_lowercase());
            if !self.can(self.get_current_user(), &action, routable) {
                return false;
            }

            if let Some(proc) = extra_authorization_proc {
                proc(routable)
            } else {
                true
            }
        } else {
            false
        }
    }

    fn ensure_canonical_path(&self, routable: &impl Routable, routable_full_path: &str) {
        if self.is_get_request() {
            let canonical_path = routable.full_path();
            if canonical_path != routable_full_path {
                // Handle redirect to canonical path
            }
        }
    }

    // Required methods to be implemented by concrete types
    fn get_current_user(&self) -> &User;
    fn can(&self, user: &User, action: &str, subject: &impl Any) -> bool;
    fn is_get_request(&self) -> bool;
}

pub struct User {
    // Add user fields as needed
}
