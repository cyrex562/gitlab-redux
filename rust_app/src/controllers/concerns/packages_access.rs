// Ported from: app/controllers/concerns/packages_access.rb
use crate::{
    config::Settings,
    models::{Project, User},
    utils::authorization::Can,
};
use actix_web::{error::Error, HttpResponse};

pub trait PackagesAccess {
    fn verify_packages_enabled(&self) -> Result<(), Error>;
    fn verify_read_package(&self, user: &User, project: &Project) -> Result<(), Error>;
}

pub struct PackagesAccessImpl {
    settings: Settings,
}

impl PackagesAccessImpl {
    pub fn new(settings: Settings) -> Self {
        Self { settings }
    }
}

impl PackagesAccess for PackagesAccessImpl {
    fn verify_packages_enabled(&self) -> Result<(), Error> {
        if !self.settings.packages.enabled {
            return Err(actix_web::error::ErrorNotFound("Packages are not enabled"));
        }
        Ok(())
    }

    fn verify_read_package(&self, user: &User, project: &Project) -> Result<(), Error> {
        if let Some(packages_policy_subject) = project.packages_policy_subject() {
            if !Can::can(user, "read_package", packages_policy_subject) {
                return Err(actix_web::error::ErrorForbidden("Access denied"));
            }
        }
        Ok(())
    }
}
