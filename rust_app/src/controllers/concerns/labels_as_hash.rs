use serde::{Deserialize, Serialize};
use crate::models::label::Label;
use crate::models::user::User;
use crate::services::labels_finder::LabelsFinder;
use std::collections::HashMap;

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct LabelHash {
    pub title: String,
    pub color: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub set: Option<bool>,
}

pub trait LabelsAsHash {
    fn labels_as_hash(&self, target: Option<&dyn Labeled>, params: &HashMap<String, String>) -> Vec<LabelHash>;
}

pub trait Labeled {
    fn labels(&self) -> Vec<Label>;
}

pub struct LabelsAsHashImpl {
    current_user: User,
}

impl LabelsAsHashImpl {
    pub fn new(current_user: User) -> Self {
        Self { current_user }
    }
}

impl LabelsAsHash for LabelsAsHashImpl {
    fn labels_as_hash(&self, target: Option<&dyn Labeled>, params: &HashMap<String, String>) -> Vec<LabelHash> {
        // Find available labels
        let labels_finder = LabelsFinder::new(&self.current_user, params);
        let available_labels = labels_finder.execute();

        // Convert labels to hash format
        let mut label_hashes: Vec<LabelHash> = available_labels
            .iter()
            .map(|label| LabelHash {
                title: label.title.clone(),
                color: label.color.clone(),
                set: None,
            })
            .collect();

        // If target is provided and has labels, mark the ones that are already set
        if let Some(target) = target {
            let target_labels = target.labels();
            let already_set_labels: Vec<&Label> = available_labels
                .iter()
                .filter(|label| target_labels.iter().any(|tl| tl.title == label.title))
                .collect();

            if !already_set_labels.is_empty() {
                let titles: Vec<String> = already_set_labels.iter().map(|l| l.title.clone()).collect();
                
                for hash in &mut label_hashes {
                    if titles.contains(&hash.title) {
                        hash.set = Some(true);
                    }
                }
            }
        }

        label_hashes
    }
} 