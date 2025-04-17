use actix_web::{HttpRequest, HttpResponse};
use std::sync::OnceLock;

pub struct GitlabConfig {
    pub packages: PackagesConfig,
}

pub struct PackagesConfig {
    pub enabled: bool,
}

impl GitlabConfig {
    pub fn get() -> &'static GitlabConfig {
        static INSTANCE: OnceLock<GitlabConfig> = OnceLock::new();
        INSTANCE.get_or_init(|| GitlabConfig {
            packages: PackagesConfig { enabled: true },
        })
    }
}

pub struct Project {
    // Add project fields as needed
}

impl Project {
    pub fn packages_policy_subject(&self) -> &dyn Any {
        // Implementation would depend on your policy system
        self
    }
}

pub trait PackagesAccess {
    fn verify_packages_enabled(&self) -> HttpResponse {
        if GitlabConfig::get().packages.enabled {
            HttpResponse::Ok().finish()
        } else {
            HttpResponse::NotFound().finish()
        }
    }

    fn verify_read_package(&self, req: &HttpRequest) -> HttpResponse {
        if self.can_read_package(
            self.get_current_user(),
            self.get_project().map(|p| p.packages_policy_subject()),
        ) {
            HttpResponse::Ok().finish()
        } else {
            HttpResponse::Forbidden().finish()
        }
    }

    // Required methods to be implemented by concrete types
    fn get_current_user(&self) -> &User;
    fn get_project(&self) -> Option<&Project>;
    fn can_read_package(&self, user: &User, subject: Option<&dyn Any>) -> bool;
}

pub struct User {
    // Add user fields as needed
}
