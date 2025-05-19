// Ported from: orig_app/app/controllers/projects/design_management/designs/raw_images_controller.rb
// Ported on: 2025-05-19
// This file implements the Projects::DesignManagement::Designs::RawImagesController from the Ruby codebase.

use actix_web::{web, HttpResponse, Responder};

use crate::controllers::projects::design_management::designs_controller::ProjectsDesignManagementDesignsController;
use crate::services::design_management::design_repository::DesignRepository;
use crate::services::design_management::sends_blob::send_blob;

pub struct ProjectsDesignManagementDesignsRawImagesController {
    base: ProjectsDesignManagementDesignsController,
}

impl ProjectsDesignManagementDesignsRawImagesController {
    pub fn new(base: ProjectsDesignManagementDesignsController) -> Self {
        Self { base }
    }

    /// GET /projects/:project_id/design_management/designs/:id/raw_image
    pub async fn show(
        &self,
        ref_param: Option<String>,
        sha: Option<String>,
        design_full_path: String,
    ) -> impl Responder {
        let design_repository = self.design_repository();
        let ref_name = sha.unwrap_or_else(|| design_repository.root_ref());
        let blob = design_repository.blob_at(&ref_name, &design_full_path);
        let allow_caching = self.base.project().is_public();
        send_blob(&design_repository, &blob, false, allow_caching)
    }

    fn design_repository(&self) -> DesignRepository {
        self.base.project().design_repository()
    }
}
