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

/// Module for handling issuable metadata
pub trait IssuableMetadata {
    /// Set issuable metadata
    fn set_issuable_metadata(&self) -> Result<(), AppError> {
        let current_user = self.current_user();
        let project = self.project();
        let finder_params = self.finder_params();

        // Create finder
        let finder = IssuableFinder::new(current_user, project, finder_params);

        // Execute finder
        let issuables = finder.execute();

        // Set metadata
        self.set_metadata(issuables);

        Ok(())
    }

    /// Set metadata
    fn set_metadata(&self, issuables: Vec<Box<dyn Issuable>>) {
        // Set total count
        self.set_header("X-Total-Count", issuables.len().to_string());

        // Set state counts
        let opened_count = issuables.iter().filter(|i| i.state() == "opened").count();
        let closed_count = issuables.iter().filter(|i| i.state() == "closed").count();

        self.set_header("X-Opened-Count", opened_count.to_string());
        self.set_header("X-Closed-Count", closed_count.to_string());

        // Set label counts
        let label_counts = self.calculate_label_counts(&issuables);
        for (label, count) in label_counts {
            self.set_header(&format!("X-Label-{}", label), count.to_string());
        }

        // Set milestone counts
        let milestone_counts = self.calculate_milestone_counts(&issuables);
        for (milestone, count) in milestone_counts {
            self.set_header(&format!("X-Milestone-{}", milestone), count.to_string());
        }
    }

    /// Calculate label counts
    fn calculate_label_counts(&self, issuables: &[Box<dyn Issuable>]) -> Vec<(String, usize)> {
        let mut counts = std::collections::HashMap::new();

        for issuable in issuables {
            for label in issuable.labels() {
                *counts.entry(label.to_string()).or_insert(0) += 1;
            }
        }

        counts.into_iter().collect()
    }

    /// Calculate milestone counts
    fn calculate_milestone_counts(&self, issuables: &[Box<dyn Issuable>]) -> Vec<(String, usize)> {
        let mut counts = std::collections::HashMap::new();

        for issuable in issuables {
            if let Some(milestone) = issuable.milestone() {
                *counts.entry(milestone.to_string()).or_insert(0) += 1;
            }
        }

        counts.into_iter().collect()
    }

    // Required trait methods that need to be implemented by the controller
    fn current_user(&self) -> Option<&User>;
    fn project(&self) -> &Project;
    fn finder_params(&self) -> IssuableFinderParams;
    fn set_header(&self, name: &str, value: &str);
}
