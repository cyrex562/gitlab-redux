use actix_web::{web, HttpResponse};
use crate::controllers::ApplicationController;
use crate::auth::PublicVisibilityRestricted;

pub mod catalog;
pub mod groups;
pub mod projects;
pub mod snippets;

pub struct ExploreApplicationController {
    base: ApplicationController,
}

impl ExploreApplicationController {
    pub fn new(base: ApplicationController) -> Self {
        Self { base }
    }

    pub async fn authenticate_user(&self) -> HttpResponse {
        if !self.public_visibility_restricted() {
            return HttpResponse::Ok().finish();
        }
        
        self.base.authenticate_user().await
    }

    fn public_visibility_restricted(&self) -> bool {
        PublicVisibilityRestricted::check()
    }
} 