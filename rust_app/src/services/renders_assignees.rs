// Ported from: orig_app/app/controllers/concerns/renders_assignees.rb
// This module provides a function to preload assignees for rendering.

use std::collections::HashMap;
use uuid::Uuid;

use crate::models::{MergeRequest, Project};

/// Dummy trait for project team access, to be replaced with real logic
pub trait ProjectTeam {
    fn max_member_access_for_user_ids(&self, user_ids: &[Uuid]) -> HashMap<Uuid, i32>;
}

/// Preload assignees for rendering, similar to the Ruby concern
pub fn preload_assignees_for_render<T: ProjectTeam>(
    merge_request: &MergeRequest,
    project_team: &T,
    assignee_ids: &[Uuid],
) -> HashMap<Uuid, i32> {
    project_team.max_member_access_for_user_ids(assignee_ids)
}
