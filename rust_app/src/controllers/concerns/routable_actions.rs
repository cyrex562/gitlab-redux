// Ported from: orig_app/app/controllers/concerns/routable_actions.rb
// This module provides RoutableActions trait for controller logic.
use actix_web::{HttpRequest, HttpResponse};

pub trait RoutableActions {
    // Find a routable entity by its full path, authorize, and ensure canonical path.
    fn find_routable(
        &self,
        routable_klass: &dyn RoutableClass,
        routable_full_path: &str,
        full_path: &str,
        req: &HttpRequest,
        extra_authorization: Option<&dyn Fn(&dyn Routable) -> bool>,
    ) -> Option<Box<dyn Routable>> {
        let routable = routable_klass.find_by_full_path(routable_full_path, req.method() == "GET");
        if self.routable_authorized(routable.as_deref(), extra_authorization) {
            self.ensure_canonical_path(routable.as_deref().unwrap(), routable_full_path, req);
            routable
        } else {
            self.perform_not_found_actions(
                routable.as_deref(),
                &self.not_found_actions(),
                full_path,
            );
            if !self.performed() {
                self.route_not_found();
            }
            None
        }
    }

    fn not_found_actions(&self) -> Vec<Box<dyn Fn(&dyn Routable, &str)>>;

    fn perform_not_found_actions(
        &self,
        routable: Option<&dyn Routable>,
        actions: &[Box<dyn Fn(&dyn Routable, &str)>],
        full_path: &str,
    );

    fn routable_authorized(
        &self,
        routable: Option<&dyn Routable>,
        extra_authorization: Option<&dyn Fn(&dyn Routable) -> bool>,
    ) -> bool;

    fn ensure_canonical_path(
        &self,
        routable: &dyn Routable,
        routable_full_path: &str,
        req: &HttpRequest,
    );

    fn performed(&self) -> bool;
    fn route_not_found(&self);
}

// Traits for Routable and RoutableClass to be implemented by your models.
pub trait Routable {
    fn full_path(&self) -> &str;
    fn class_name(&self) -> &str;
}

pub trait RoutableClass {
    fn find_by_full_path(
        &self,
        full_path: &str,
        follow_redirects: bool,
    ) -> Option<Box<dyn Routable>>;
}
