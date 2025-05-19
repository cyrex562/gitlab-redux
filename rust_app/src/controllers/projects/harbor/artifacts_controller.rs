// Ported from: orig_app/app/controllers/projects/harbor/artifacts_controller.rb
// Ported on: 2025-05-19
// This file implements the Projects::Harbor::ArtifactsController from the Ruby codebase.

use crate::controllers::projects::harbor::application_controller::ProjectsHarborApplicationController;

pub struct ProjectsHarborArtifactsController {
    base: ProjectsHarborApplicationController,
}

impl ProjectsHarborArtifactsController {
    pub fn new(base: ProjectsHarborApplicationController) -> Self {
        Self { base }
    }

    fn container(&self) -> &crate::models::project::Project {
        self.base.base.project()
    }
}
