// Ported from: orig_app/app/controllers/projects/harbor/tags_controller.rb
// Ported on: 2025-05-19
// This file implements the Projects::Harbor::TagsController from the Ruby codebase.

use crate::controllers::projects::harbor::application_controller::ProjectsHarborApplicationController;

pub struct ProjectsHarborTagsController {
    base: ProjectsHarborApplicationController,
}

impl ProjectsHarborTagsController {
    pub fn new(base: ProjectsHarborApplicationController) -> Self {
        Self { base }
    }

    fn container(&self) -> &crate::models::project::Project {
        self.base.base.project()
    }
}
