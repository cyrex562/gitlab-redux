use actix_web::{web, HttpResponse};
use std::sync::Arc;

use crate::{
    models::{issuable::Issuable, project::Project, user::User},
    services::{
        issuable::{
            IssuableBulkUpdateService, IssuableDestroyService, IssuableDiscussionsListService,
            IssuableService,
        },
        security::SecurityUtils,
    },
    utils::{error::AppError, flash::Flash, polling_interval::PollingInterval},
};

/// Module for handling issuable actions
pub trait IssuableActions {
    /// Show an issuable
    fn show(&self) -> Result<HttpResponse, AppError> {
        let issuable = self.issuable();
        let serializer = self.serializer();

        // Handle different response formats
        match self.format() {
            "html" => {
                let issuable_sidebar = serializer.represent(issuable, "sidebar");
                Ok(self.render_view("show", issuable_sidebar))
            }
            "json" => {
                let serializer_name = self.params().get("serializer").cloned().unwrap_or_default();
                Ok(HttpResponse::Ok().json(serializer.represent(issuable, &serializer_name)))
            }
            _ => Err(AppError::UnsupportedFormat),
        }
    }

    /// Update an issuable
    fn update(&self) -> Result<HttpResponse, AppError> {
        let issuable = self.issuable();
        let update_service = self.update_service();

        // Execute the update service
        let updated_issuable = update_service.execute(issuable)?;

        // Handle different response formats
        match self.format() {
            "html" => {
                if updated_issuable.is_valid() {
                    // For spammable issuables, redirect to spammable path
                    if updated_issuable.is_spammable() {
                        Ok(HttpResponse::Found()
                            .header("Location", self.spammable_path(&updated_issuable))
                            .finish())
                    } else {
                        // For non-spammable issuables, render edit view
                        Ok(self.render_view("edit", updated_issuable))
                    }
                } else {
                    // Handle validation errors with captcha check if needed
                    if updated_issuable.is_spammable() {
                        self.with_captcha_check_html_format(&updated_issuable, || {
                            self.render_view("edit", updated_issuable)
                        })
                    } else {
                        Ok(self.render_view("edit", updated_issuable))
                    }
                }
            }
            "json" => {
                if updated_issuable.is_spammable() {
                    self.with_captcha_check_json_format(&updated_issuable, || {
                        self.render_entity_json(&updated_issuable)
                    })
                } else {
                    self.render_entity_json(&updated_issuable)
                }
            }
            _ => Err(AppError::UnsupportedFormat),
        }
    }

    /// Get realtime changes for an issuable
    fn realtime_changes(&self) -> Result<HttpResponse, AppError> {
        let issuable = self.issuable();

        // Set polling interval header
        PollingInterval::set_header(3000);

        // Build response data
        let mut response = serde_json::json!({
            "title": self.markdown_field(issuable, "title"),
            "title_text": issuable.title(),
            "description": self.markdown_field(issuable, "description"),
            "description_text": issuable.description(),
            "task_completion_status": issuable.task_completion_status(),
            "lock_version": issuable.lock_version(),
        });

        // Add edited information if available
        if issuable.is_edited() {
            response["updated_at"] =
                serde_json::Value::String(issuable.last_edited_at().to_iso8601());
            response["updated_by_name"] =
                serde_json::Value::String(issuable.last_edited_by().name());
            response["updated_by_path"] =
                serde_json::Value::String(self.user_path(issuable.last_edited_by()));
        }

        Ok(HttpResponse::Ok().json(response))
    }

    /// Destroy an issuable
    fn destroy(&self) -> Result<HttpResponse, AppError> {
        let issuable = self.issuable();
        let project = issuable.project();
        let current_user = self.current_user();

        // Execute destroy service
        IssuableDestroyService::new(project, current_user).execute(issuable)?;

        // Get human-readable name for the issuable
        let name = issuable.human_class_name();

        // Set flash notice
        self.flash()
            .set_notice(format!("The {} was successfully deleted.", name));

        // Get index path
        let index_path = self.polymorphic_path(&[self.parent(), issuable.class_name()]);

        // Handle different response formats
        match self.format() {
            "html" => Ok(HttpResponse::Found()
                .header("Location", index_path)
                .finish()),
            "json" => Ok(HttpResponse::Ok().json(serde_json::json!({
                "web_url": index_path
            }))),
            _ => Err(AppError::UnsupportedFormat),
        }
    }

    /// Check destroy confirmation
    fn check_destroy_confirmation(&self) -> Result<bool, AppError> {
        if self.params().get("destroy_confirm").is_some() {
            return Ok(true);
        }

        let issuable = self.issuable();
        let error_message = format!(
            "Destroy confirmation not provided for {}",
            issuable.human_class_name()
        );

        // Log the error
        self.log_error(
            &error_message,
            &[
                ("project_path", issuable.project().full_path()),
                ("issuable_type", issuable.class_name()),
                ("issuable_id", issuable.id().to_string()),
            ],
        );

        // Get index path
        let index_path = self.polymorphic_path(&[self.parent(), issuable.class_name()]);

        // Handle different response formats
        match self.format() {
            "html" => {
                self.flash().set_notice(&error_message);
                Ok(HttpResponse::Found()
                    .header("Location", index_path)
                    .finish())
            }
            "json" => Ok(HttpResponse::UnprocessableEntity().json(serde_json::json!({
                "errors": error_message
            }))),
            _ => Err(AppError::UnsupportedFormat),
        }
    }

    /// Bulk update issuables
    fn bulk_update(&self) -> Result<HttpResponse, AppError> {
        let parent = self.parent();
        let current_user = self.current_user();
        let bulk_update_params = self.bulk_update_params();
        let resource_name = self.resource_name();

        // Execute bulk update service
        let result = IssuableBulkUpdateService::new(parent, current_user, bulk_update_params)
            .execute(&resource_name)?;

        if result.is_success() {
            let quantity = result.payload().get("count").unwrap_or(&0);
            Ok(HttpResponse::Ok().json(serde_json::json!({
                "notice": format!("{} {} updated", quantity, resource_name.pluralize(*quantity))
            })))
        } else {
            Ok(HttpResponse::UnprocessableEntity().json(serde_json::json!({
                "errors": result.message()
            })))
        }
    }

    /// Get discussions for an issuable
    fn discussions(&self) -> Result<HttpResponse, AppError> {
        let current_user = self.current_user();
        let issuable = self.issuable();
        let finder_params = self.finder_params_for_issuable();

        // Create discussions finder
        let finder = IssuableDiscussionsListService::new(current_user, issuable, finder_params);

        // Execute finder
        let discussion_notes = finder.execute();

        // Process discussion notes if needed
        if let Some(block) = self.discussion_block() {
            block(discussion_notes.clone());
        }

        // Set pagination header if needed
        if let Some(paginator) = finder.paginator() {
            if paginator.has_next_page() {
                self.set_header("X-Next-Page-Cursor", paginator.cursor_for_next_page());
            }
        }

        // Check if response is stale
        let is_stale = match issuable.class_name().as_str() {
            "MergeRequest" | "Issue" => {
                self.is_stale(&[self.discussion_cache_context(), discussion_notes.clone()])
            }
            _ => false,
        };

        // Return appropriate response
        if is_stale {
            Ok(HttpResponse::Ok().json(
                self.discussion_serializer()
                    .represent(discussion_notes, self),
            ))
        } else {
            Ok(HttpResponse::Ok().json(
                self.discussion_serializer()
                    .represent(discussion_notes, self),
            ))
        }
    }

    /// Get notes filter
    fn notes_filter(&self) -> Option<i32> {
        let notes_filter_param = self
            .params()
            .get("notes_filter")
            .and_then(|v| v.parse::<i32>().ok());

        // Check if database is read-only or filter should not be persisted
        if self.is_database_read_only()
            || self.params().get("persist_filter") == Some(&"false".to_string())
        {
            notes_filter_param.or_else(|| {
                self.current_user()
                    .and_then(|user| user.notes_filter_for(self.issuable()))
            })
        } else {
            // Set notes filter for user
            self.current_user()
                .and_then(|user| user.set_notes_filter(notes_filter_param, self.issuable()))
                .or(notes_filter_param)
        }
    }

    /// Get discussion cache context
    fn discussion_cache_context(&self) -> String {
        let current_user = self.current_user();
        let project = self.project();

        format!(
            "{}:{}:v2",
            current_user.map(|u| u.cache_key()).unwrap_or_default(),
            project
                .team()
                .human_max_access(current_user.map(|u| u.id())),
        )
    }

    /// Get discussion serializer
    fn discussion_serializer(&self) -> Box<dyn DiscussionSerializer> {
        Box::new(DiscussionSerializer::new(
            self.project(),
            self.issuable(),
            self.current_user(),
            Box::new(ProjectNoteEntity::new()),
        ))
    }

    /// Render conflict response
    fn render_conflict_response(&self) -> Result<HttpResponse, AppError> {
        let issuable = self.issuable();

        match self.format() {
            "html" => {
                Ok(self.render_view("edit", serde_json::json!({
                    "conflict": true
                })))
            },
            "json" => {
                Ok(HttpResponse::Conflict().json(serde_json::json!({
                    "errors": [
                        format!(
                            "Someone edited this {} at the same time you did. \
                            Please refresh your browser and make sure your changes will not unintentionally remove theirs.",
                            issuable.human_class_name()
                        )
                    ]
                })))
            },
            _ => Err(AppError::UnsupportedFormat),
        }
    }

    /// Authorize destroy issuable
    fn authorize_destroy_issuable(&self) -> Result<(), AppError> {
        let current_user = self.current_user();
        let issuable = self.issuable();

        if !SecurityUtils::can(
            current_user,
            &format!("destroy_{}", issuable.to_ability_name()),
            issuable,
        ) {
            return Err(AppError::AccessDenied);
        }

        Ok(())
    }

    /// Authorize admin issuable
    fn authorize_admin_issuable(&self) -> Result<(), AppError> {
        let current_user = self.current_user();
        let parent = self.parent();
        let resource_name = self.resource_name();

        if !SecurityUtils::can(current_user, &format!("admin_{}", resource_name), parent) {
            return Err(AppError::AccessDenied);
        }

        Ok(())
    }

    /// Authorize update issuable
    fn authorize_update_issuable(&self) -> Result<(), AppError> {
        let current_user = self.current_user();
        let issuable = self.issuable();
        let resource_name = self.resource_name();

        if !SecurityUtils::can(current_user, &format!("update_{}", resource_name), issuable) {
            return Err(AppError::NotFound);
        }

        Ok(())
    }

    /// Set application context
    fn set_application_context(&self) -> Result<(), AppError> {
        // No-op in CE, logic is defined in EE module
        Ok(())
    }

    /// Get bulk update params
    fn bulk_update_params(&self) -> std::collections::HashMap<String, String> {
        let update_params = self
            .params()
            .get("update")
            .and_then(|v| v.as_object())
            .map(|obj| {
                obj.iter()
                    .filter_map(|(k, v)| {
                        if let Some(v_str) = v.as_str() {
                            Some((k.clone(), v_str.to_string()))
                        } else {
                            None
                        }
                    })
                    .collect::<std::collections::HashMap<String, String>>()
            })
            .unwrap_or_default();

        self.clean_bulk_update_params(update_params)
    }

    /// Clean bulk update params
    fn clean_bulk_update_params(
        &self,
        mut permitted_params: std::collections::HashMap<String, String>,
    ) -> std::collections::HashMap<String, String> {
        permitted_params.retain(|k, v| {
            if k == "issuable_ids" {
                return true;
            }

            if v.is_empty() {
                return false;
            }

            true
        });

        permitted_params
    }

    /// Get bulk update permitted keys
    fn bulk_update_permitted_keys(&self) -> Vec<String> {
        vec![
            "issuable_ids".to_string(),
            "assignee_id".to_string(),
            "milestone_id".to_string(),
            "state_event".to_string(),
            "subscription_event".to_string(),
            "add_label_ids".to_string(),
            "remove_label_ids".to_string(),
            "add_labels".to_string(),
            "remove_labels".to_string(),
            "due_date".to_string(),
            "confidential".to_string(),
            "discussion_locked".to_string(),
            "todo_event".to_string(),
            "weight".to_string(),
            "health_status".to_string(),
            "epic_id".to_string(),
            "epic_action".to_string(),
            "epic_date".to_string(),
            "epic_date_type".to_string(),
            "epic_date_source".to_string(),
            "epic_date_fixed_date".to_string(),
            "epic_date_is_fixed".to_string(),
            "epic_date_start_date_fixed".to_string(),
            "epic_date_start_date_is_fixed".to_string(),
            "epic_date_due_date_fixed".to_string(),
            "epic_date_due_date_is_fixed".to_string(),
            "epic_date_start_date_sourced_from_milestones".to_string(),
            "epic_date_due_date_sourced_from_milestones".to_string(),
            "epic_date_start_date_inherited_from_milestones".to_string(),
            "epic_date_due_date_inherited_from_milestones".to_string(),
            "epic_date_start_date_from_milestones".to_string(),
            "epic_date_due_date_from_milestones".to_string(),
            "epic_date_start_date_from_inherited_source".to_string(),
            "epic_date_due_date_from_inherited_source".to_string(),
            "epic_date_start_date_from_fixed_date".to_string(),
            "epic_date_due_date_from_fixed_date".to_string(),
            "epic_date_start_date_from_milestones".to_string(),
            "epic_date_due_date_from_milestones".to_string(),
            "epic_date_start_date_from_inherited_source".to_string(),
            "epic_date_due_date_from_inherited_source".to_string(),
            "epic_date_start_date_from_fixed_date".to_string(),
            "epic_date_due_date_from_fixed_date".to_string(),
        ]
    }

    /// Get resource name
    fn resource_name(&self) -> String {
        self.controller_name().singularize()
    }

    /// Render entity JSON
    fn render_entity_json(&self, issuable: &dyn Issuable) -> Result<HttpResponse, AppError> {
        if issuable.is_valid() {
            Ok(HttpResponse::Ok().json(self.serializer().represent(issuable, "full")))
        } else {
            Ok(HttpResponse::UnprocessableEntity().json(serde_json::json!({
                "errors": issuable.errors()
            })))
        }
    }

    // Required trait methods that need to be implemented by the controller
    fn issuable(&self) -> &dyn Issuable;
    fn serializer(&self) -> Box<dyn IssuableSerializer>;
    fn update_service(&self) -> Box<dyn IssuableService>;
    fn parent(&self) -> &dyn Parent;
    fn project(&self) -> &Project;
    fn current_user(&self) -> Option<&User>;
    fn flash(&self) -> &dyn Flash;
    fn params(&self) -> &serde_json::Value;
    fn format(&self) -> &str;
    fn controller_name(&self) -> String;
    fn render_view(&self, view: &str, data: impl serde::Serialize) -> HttpResponse;
    fn markdown_field(&self, issuable: &dyn Issuable, field: &str) -> String;
    fn user_path(&self, user: &User) -> String;
    fn polymorphic_path(&self, parts: &[&dyn std::any::Any]) -> String;
    fn is_database_read_only(&self) -> bool;
    fn is_stale(&self, etag: &[impl std::any::Any]) -> bool;
    fn set_header(&self, name: &str, value: &str);
    fn log_error(&self, message: &str, context: &[(&str, String)]);
    fn discussion_block(&self) -> Option<Box<dyn Fn(Vec<Box<dyn Note>>)>>;
    fn finder_params_for_issuable(&self) -> std::collections::HashMap<String, String>;
    fn spammable_path(&self, issuable: &dyn Issuable) -> String;
    fn with_captcha_check_html_format<F>(
        &self,
        spammable: &dyn Spammable,
        f: F,
    ) -> Result<HttpResponse, AppError>
    where
        F: FnOnce() -> HttpResponse;
    fn with_captcha_check_json_format<F>(
        &self,
        spammable: &dyn Spammable,
        f: F,
    ) -> Result<HttpResponse, AppError>
    where
        F: FnOnce() -> Result<HttpResponse, AppError>;
}

/// Trait for issuable serializers
pub trait IssuableSerializer {
    fn represent(&self, issuable: &dyn Issuable, serializer: &str) -> serde_json::Value;
}

/// Trait for discussion serializers
pub trait DiscussionSerializer {
    fn represent(
        &self,
        discussion: Box<dyn Discussion>,
        context: &dyn std::any::Any,
    ) -> serde_json::Value;
}

/// Trait for parent objects
pub trait Parent {
    fn id(&self) -> i32;
}

/// Trait for spammable objects
pub trait Spammable {
    fn is_spammable(&self) -> bool;
}

/// Trait for notes
pub trait Note {
    fn id(&self) -> i32;
}

/// Trait for discussions
pub trait Discussion {
    fn id(&self) -> i32;
}

/// Trait for note entities
pub trait NoteEntity {
    fn new() -> Self
    where
        Self: Sized;
}

/// Implementation for ProjectNoteEntity
pub struct ProjectNoteEntity;

impl NoteEntity for ProjectNoteEntity {
    fn new() -> Self {
        ProjectNoteEntity
    }
}

/// Implementation for DiscussionSerializer
pub struct DiscussionSerializerImpl {
    project: Arc<Project>,
    noteable: Box<dyn Issuable>,
    current_user: Option<Arc<User>>,
    note_entity: Box<dyn NoteEntity>,
}

impl DiscussionSerializerImpl {
    pub fn new(
        project: &Project,
        noteable: &dyn Issuable,
        current_user: Option<&User>,
        note_entity: Box<dyn NoteEntity>,
    ) -> Self {
        Self {
            project: Arc::new(project.clone()),
            noteable: Box::new(noteable.clone()),
            current_user: current_user.map(|u| Arc::new(u.clone())),
            note_entity,
        }
    }
}

impl DiscussionSerializer for DiscussionSerializerImpl {
    fn represent(
        &self,
        discussion: Box<dyn Discussion>,
        _context: &dyn std::any::Any,
    ) -> serde_json::Value {
        // Implementation would go here
        serde_json::json!({})
    }
}
