// Ported from: orig_app/app/controllers/explore/application_controller.rb
// Ported on: 2025-04-30
// This file implements the Explore::ApplicationController from the Ruby codebase.
//
// See porting_log.txt for details.

use crate::auth::PublicVisibilityRestricted;
use crate::controllers::ApplicationController;
use actix_web::{web, HttpResponse};

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
