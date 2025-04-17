use actix_web::{web, HttpResponse};
use std::sync::Arc;

use crate::models::project::Project;
use crate::models::user::User;
use crate::services::commit_service::CommitService;
use crate::utils::access_control::AccessControl;
use crate::utils::flash::Flash;
use crate::utils::merge_request::MergeRequestFinder;
use crate::utils::security::SecurityUtils;

/// Module for handling commit creation
pub trait CreatesCommit {
    /// Create a commit using the specified service
    fn create_commit(
        &self,
        service: Box<dyn CommitService>,
        success_path: String,
        failure_path: String,
        failure_view: Option<String>,
        success_notice: Option<String>,
        target_project: Option<Arc<Project>>,
    ) -> HttpResponse {
        let target_project = target_project.unwrap_or_else(|| self.project().clone());
        let branch_name_or_ref = self.branch_name_or_ref();

        let (project_to_commit_into, different_project, branch_name) = if self
            .user_access(&target_project)
            .can_push_to_branch(&branch_name_or_ref)
        {
            (target_project.clone(), false, branch_name_or_ref.clone())
        } else {
            let fork = self.current_user().fork_of(&target_project);
            let next_branch = fork.repository().next_branch("patch");
            (fork, true, next_branch)
        };

        let start_branch = self.ref_name().unwrap_or_else(|| branch_name.clone());

        let mut commit_params = self.commit_params().clone();
        commit_params.insert(
            "start_project".to_string(),
            project_to_commit_into.id().to_string(),
        );
        commit_params.insert("start_branch".to_string(), start_branch);
        commit_params.insert(
            "source_project".to_string(),
            self.project().id().to_string(),
        );
        commit_params.insert(
            "target_project".to_string(),
            target_project.id().to_string(),
        );
        commit_params.insert("branch_name".to_string(), branch_name);

        let result = service.execute(&project_to_commit_into, self.current_user(), commit_params);

        match result.status() {
            "success" => {
                let success_path = self.final_success_path(&success_path, &target_project);
                let success_notice = self.update_flash_notice(success_notice, &success_path);

                self.flash().set_notice(success_notice);

                HttpResponse::Found()
                    .header("Location", success_path)
                    .finish()
            }
            _ => {
                let error_message = self.format_flash_notice(result.message());
                self.flash().set_alert(error_message);

                let failure_path = if failure_path.contains("{{") {
                    // Handle dynamic paths with template variables
                    self.evaluate_path_template(&failure_path)
                } else {
                    failure_path
                };

                if let Some(view) = failure_view {
                    // Render the specified view
                    self.render_view(&view)
                } else {
                    // Redirect to failure path
                    HttpResponse::Found()
                        .header("Location", failure_path)
                        .finish()
                }
            }
        }
    }

    /// Authorize edit tree operation
    fn authorize_edit_tree(&self) -> Result<(), actix_web::Error> {
        if self.can_collaborate_with_project(self.project(), &self.branch_name_or_ref()) {
            Ok(())
        } else {
            Err(actix_web::error::ErrorForbidden("Access denied"))
        }
    }

    /// Format flash notice message
    fn format_flash_notice(&self, message: &str) -> String {
        let formatted_message = message.replace("\n", "<br>");
        self.sanitize_html(&formatted_message, &["br"])
    }

    /// Update flash notice with success message
    fn update_flash_notice(&self, success_notice: Option<String>, success_path: &str) -> String {
        let changes_link = format!("<a href=\"{}\" class=\"gl-link\">changes</a>", success_path);

        let default_message = format!("Your {} have been committed successfully.", changes_link);

        let mut notice = success_notice.unwrap_or(default_message);

        if self.create_merge_request() {
            if !self.merge_request_exists() {
                let mr_message = if self.different_project() {
                    "You can now submit a merge request to get this change into the original project."
                } else {
                    "You can now submit a merge request to get this change into the original branch."
                };

                notice.push_str(&format!(" {}", mr_message));
            }
        }

        notice
    }

    /// Get the final success path
    fn final_success_path(&self, success_path: &str, target_project: &Arc<Project>) -> String {
        if self.create_merge_request() {
            if self.merge_request_exists() {
                self.existing_merge_request_path()
            } else {
                self.new_merge_request_path(target_project)
            }
        } else {
            success_path.to_string()
        }
    }

    /// Get path for new merge request
    fn new_merge_request_path(&self, target_project: &Arc<Project>) -> String {
        format!("/projects/{}/merge_requests/new", target_project.id())
    }

    /// Get path for existing merge request
    fn existing_merge_request_path(&self) -> String {
        format!(
            "/projects/{}/merge_requests/{}",
            self.project().id(),
            self.merge_request().id()
        )
    }

    /// Check if merge request exists
    fn merge_request_exists(&self) -> bool {
        let finder = MergeRequestFinder::new(self.current_user());
        let params = vec![("project_id", self.project().id().to_string())];
        !finder.execute(params).is_empty()
    }

    /// Check if merge request should be created
    fn create_merge_request(&self) -> bool {
        // Even if the field is set, if we're checking the same branch
        // as the target branch in the same project,
        // we don't want to create a merge request.
        let create_mr = self
            .params()
            .get("create_merge_request")
            .and_then(|v| v.parse::<bool>().ok())
            .unwrap_or(false);

        create_mr && (self.different_project() || self.branch_name_or_ref() != self.target_branch())
    }

    /// Get branch name or ref
    fn branch_name_or_ref(&self) -> String {
        self.branch_name()
            .unwrap_or_else(|| self.ref_name().unwrap_or_default())
    }

    // Required trait methods that need to be implemented by the controller
    fn project(&self) -> Arc<Project>;
    fn current_user(&self) -> Arc<User>;
    fn user_access(&self, project: &Arc<Project>) -> Box<dyn AccessControl>;
    fn ref_name(&self) -> Option<String>;
    fn branch_name(&self) -> Option<String>;
    fn target_branch(&self) -> String;
    fn commit_params(&self) -> std::collections::HashMap<String, String>;
    fn flash(&self) -> &dyn Flash;
    fn params(&self) -> std::collections::HashMap<String, String>;
    fn different_project(&self) -> bool;
    fn merge_request(&self) -> Arc<crate::models::merge_request::MergeRequest>;
    fn can_collaborate_with_project(&self, project: &Arc<Project>, ref_name: &str) -> bool;
    fn sanitize_html(&self, html: &str, allowed_tags: &[&str]) -> String;
    fn render_view(&self, view: &str) -> HttpResponse;
    fn evaluate_path_template(&self, template: &str) -> String;
}
