use crate::models::project::Project;
use crate::models::snippet::Snippet;
use crate::models::user::User;
use actix_web::{web, HttpResponse};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

/// Module for handling snippet authorizations
pub trait SnippetAuthorizations {
    /// Get the current user
    fn current_user(&self) -> Option<&User>;

    /// Get the current snippet
    fn current_snippet(&self) -> Option<&Snippet>;

    /// Get the current project
    fn current_project(&self) -> Option<&Project>;

    /// Get the current user ID
    fn user_id(&self) -> Option<i32>;

    /// Get the snippet ID
    fn snippet_id(&self) -> Option<i32>;

    /// Get the snippet author ID
    fn snippet_author_id(&self) -> Option<i32>;

    /// Get the project ID
    fn project_id(&self) -> Option<i32>;

    /// Check if the user can read the snippet
    fn can_read_snippet(&self) -> bool {
        let user = match self.current_user() {
            Some(user) => user,
            None => return false,
        };

        let snippet = match self.current_snippet() {
            Some(snippet) => snippet,
            None => return false,
        };

        // Admin can read all snippets
        if user.is_admin() {
            return true;
        }

        // Project members can read project snippets
        if let Some(project) = self.current_project() {
            if project.is_member(user.id) {
                return true;
            }
        }

        // Public snippets can be read by anyone
        if snippet.is_public() {
            return true;
        }

        // Snippet author can read their own snippets
        if snippet.author_id == user.id {
            return true;
        }

        false
    }

    /// Check if the user can write to the snippet
    fn can_write_snippet(&self) -> bool {
        let user = match self.current_user() {
            Some(user) => user,
            None => return false,
        };

        let snippet = match self.current_snippet() {
            Some(snippet) => snippet,
            None => return false,
        };

        // Admin can write to all snippets
        if user.is_admin() {
            return true;
        }

        // Project members can write to project snippets
        if let Some(project) = self.current_project() {
            if project.is_member(user.id) {
                return true;
            }
        }

        // Snippet author can write to their own snippets
        if snippet.author_id == user.id {
            return true;
        }

        false
    }

    /// Check if the user can view the snippet
    fn can_view_snippet(&self) -> bool {
        // TODO: Implement actual authorization check
        // This would typically involve:
        // 1. Checking if the snippet is public
        // 2. Checking if the user has project access
        // 3. Checking if the user is the author
        true
    }

    /// Check if the user can edit the snippet
    fn can_edit_snippet(&self) -> bool {
        // TODO: Implement actual authorization check
        // This would typically involve:
        // 1. Checking if the user is the author
        // 2. Checking if the user has project write access
        // 3. Checking if the snippet is editable
        self.user_id() == self.snippet_author_id()
    }

    /// Check if the user can delete the snippet
    fn can_delete_snippet(&self) -> bool {
        // TODO: Implement actual authorization check
        // This would typically involve:
        // 1. Checking if the user is the author
        // 2. Checking if the user has project admin access
        // 3. Checking if the snippet is deletable
        self.user_id() == self.snippet_author_id()
    }

    /// Get snippet permissions
    fn get_snippet_permissions(&self) -> HashMap<String, bool> {
        let mut permissions = HashMap::new();

        permissions.insert("can_view".to_string(), self.can_view_snippet());
        permissions.insert("can_edit".to_string(), self.can_edit_snippet());
        permissions.insert("can_delete".to_string(), self.can_delete_snippet());

        permissions
    }

    /// Enforce snippet read authorization
    fn enforce_snippet_read_auth(&self) -> Result<(), HttpResponse> {
        if !self.can_read_snippet() {
            return Err(HttpResponse::Forbidden().json(json!({
                "error": "Forbidden",
                "message": "You don't have permission to read this snippet"
            })));
        }

        Ok(())
    }

    /// Enforce snippet write authorization
    fn enforce_snippet_write_auth(&self) -> Result<(), HttpResponse> {
        if !self.can_write_snippet() {
            return Err(HttpResponse::Forbidden().json(json!({
                "error": "Forbidden",
                "message": "You don't have permission to modify this snippet"
            })));
        }

        Ok(())
    }

    /// Enforce snippet authorization
    fn enforce_snippet_auth(&self) -> Result<(), HttpResponse> {
        if !self.can_view_snippet() {
            return Err(HttpResponse::Forbidden().json(serde_json::json!({
                "error": "You don't have permission to access this snippet"
            })));
        }
        Ok(())
    }
}
