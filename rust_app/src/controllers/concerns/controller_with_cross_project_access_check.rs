use actix_web::{HttpRequest, HttpResponse};
use std::sync::OnceLock;

pub struct CrossProjectAccess {
    pub should_run: bool,
}

impl CrossProjectAccess {
    pub fn find_check(controller: &dyn Any) -> Option<&'static CrossProjectAccess> {
        static INSTANCE: OnceLock<CrossProjectAccess> = OnceLock::new();
        Some(INSTANCE.get_or_init(|| CrossProjectAccess { should_run: true }))
    }
}

pub trait ControllerWithCrossProjectAccessCheck {
    fn cross_project_check(&self, req: &HttpRequest) -> HttpResponse {
        if let Some(check) = CrossProjectAccess::find_check(self) {
            if check.should_run {
                return self.authorize_cross_project_page(req);
            }
        }
        HttpResponse::Ok().finish()
    }

    fn authorize_cross_project_page(&self, req: &HttpRequest) -> HttpResponse {
        if self.can_read_cross_project(self.get_current_user()) {
            return HttpResponse::Ok().finish();
        }

        let rejection_message = "This page is unavailable because you are not allowed to read information across multiple projects.";
        self.access_denied(rejection_message)
    }

    // Required methods to be implemented by concrete types
    fn get_current_user(&self) -> Option<&User>;
    fn can_read_cross_project(&self, user: Option<&User>) -> bool;
    fn access_denied(&self, message: &str) -> HttpResponse;
}

pub struct User {
    // Add user fields as needed
}
