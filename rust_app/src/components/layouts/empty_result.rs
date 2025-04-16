use crate::components::Component;
use std::collections::HashMap;

/// Component for displaying empty search or filter results
pub struct EmptyResultComponent {
    component_type: EmptyResultType,
    html_options: HashMap<String, String>,
}

/// Type of empty result
#[derive(Debug, Clone, Copy, PartialEq)]
pub enum EmptyResultType {
    Search,
    Filter,
}

impl EmptyResultComponent {
    /// Create a new empty result component
    ///
    /// # Arguments
    ///
    /// * `component_type` - The type of empty result (search or filter)
    /// * `html_options` - Optional HTML attributes
    pub fn new(component_type: EmptyResultType, html_options: HashMap<String, String>) -> Self {
        Self {
            component_type,
            html_options,
        }
    }

    /// Check if this is a filter empty result
    pub fn is_filter(&self) -> bool {
        self.component_type == EmptyResultType::Filter
    }

    /// Get formatted HTML options
    fn format_options(&self) -> String {
        self.html_options
            .iter()
            .map(|(k, v)| format!("{}={}", k, v))
            .collect::<Vec<String>>()
            .join(" ")
    }
}

impl Component for EmptyResultComponent {
    fn render(&self) -> String {
        let title = "No results found";
        let description = if self.is_filter() {
            "To widen your search, change or remove filters above."
        } else {
            "Edit your search and try again."
        };
        let svg_path = "illustrations/empty-state/empty-search-md.svg";

        // Format HTML options
        let options_str = self.format_options();

        // Build the HTML using the EmptyStateComponent
        format!(
            r#"<div class="gl-empty-state" {}>
                <div class="gl-empty-state-content">
                    <img src="{}" alt="" class="gl-empty-state-illustration" />
                    <h3 class="gl-empty-state-title">{}</h3>
                    <p class="gl-empty-state-description">{}</p>
                </div>
            </div>"#,
            options_str, svg_path, title, description
        )
    }
}
