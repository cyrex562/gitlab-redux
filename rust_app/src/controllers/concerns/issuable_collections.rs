use actix_web::{web, HttpResponse};
use std::sync::Arc;

use crate::{
    models::{issuable::Issuable, project::Project, user::User},
    services::{
        issuable::{IssuableFinder, IssuableFinderParams},
        security::SecurityUtils,
    },
    utils::{error::AppError, pagination::Pagination},
};

/// Module for handling issuable collections
pub trait IssuableCollections {
    /// Set issuables index
    fn set_issuables_index(&self) -> Result<(), AppError> {
        let current_user = self.current_user();
        let project = self.project();
        let finder_params = self.finder_params();

        // Create finder
        let finder = IssuableFinder::new(current_user, project, finder_params);

        // Execute finder
        let issuables = finder.execute();

        // Set pagination
        self.set_pagination(issuables.paginator());

        Ok(())
    }

    /// Set pagination
    fn set_pagination(&self, paginator: Option<&dyn Pagination>) {
        if let Some(paginator) = paginator {
            if paginator.has_next_page() {
                self.set_header("X-Next-Page", paginator.next_page().to_string());
            }

            if paginator.has_prev_page() {
                self.set_header("X-Prev-Page", paginator.prev_page().to_string());
            }

            self.set_header("X-Total", paginator.total_count().to_string());
            self.set_header("X-Total-Pages", paginator.total_pages().to_string());
            self.set_header("X-Per-Page", paginator.per_page().to_string());
        }
    }

    /// Get issuables collection
    fn issuables_collection(&self) -> Result<Box<dyn IssuableCollection>, AppError> {
        let current_user = self.current_user();
        let project = self.project();
        let finder_params = self.finder_params();

        // Create finder
        let finder = IssuableFinder::new(current_user, project, finder_params);

        // Execute finder
        let issuables = finder.execute();

        // Create collection
        let collection =
            IssuableCollection::new(issuables, self.default_state(), self.default_sort());

        Ok(Box::new(collection))
    }

    /// Get default state
    fn default_state(&self) -> String {
        let current_user = self.current_user();
        let project = self.project();

        if let Some(user) = current_user {
            if let Some(state) = user.issuable_default_state_for(project) {
                return state;
            }
        }

        "opened".to_string()
    }

    /// Get default sort
    fn default_sort(&self) -> String {
        let current_user = self.current_user();
        let project = self.project();

        if let Some(user) = current_user {
            if let Some(sort) = user.issuable_default_sort_for(project) {
                return sort;
            }
        }

        "created_desc".to_string()
    }

    // Required trait methods that need to be implemented by the controller
    fn current_user(&self) -> Option<&User>;
    fn project(&self) -> &Project;
    fn finder_params(&self) -> IssuableFinderParams;
    fn set_header(&self, name: &str, value: &str);
}

/// Trait for issuable collections
pub trait IssuableCollection {
    fn new(issuables: Vec<Box<dyn Issuable>>, default_state: String, default_sort: String) -> Self
    where
        Self: Sized;
    fn issuables(&self) -> &[Box<dyn Issuable>];
    fn default_state(&self) -> &str;
    fn default_sort(&self) -> &str;
}

/// Implementation for IssuableCollection
pub struct IssuableCollectionImpl {
    issuables: Vec<Box<dyn Issuable>>,
    default_state: String,
    default_sort: String,
}

impl IssuableCollectionImpl {
    pub fn new(
        issuables: Vec<Box<dyn Issuable>>,
        default_state: String,
        default_sort: String,
    ) -> Self {
        Self {
            issuables,
            default_state,
            default_sort,
        }
    }
}

impl IssuableCollection for IssuableCollectionImpl {
    fn issuables(&self) -> &[Box<dyn Issuable>] {
        &self.issuables
    }

    fn default_state(&self) -> &str {
        &self.default_state
    }

    fn default_sort(&self) -> &str {
        &self.default_sort
    }
}
