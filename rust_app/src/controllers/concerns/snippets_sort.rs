// Ported from: orig_app/app/controllers/concerns/snippets_sort.rb
// Provides a trait for extracting the snippet sort parameter from pagination params.
use super::strong_pagination_params::PaginationParams;

pub trait SnippetsSort {
    fn sort_param(&self, pagination: &PaginationParams) -> String {
        pagination
            .sort
            .as_ref()
            .filter(|s| !s.is_empty())
            .cloned()
            .unwrap_or_else(|| "updated_desc".to_string())
    }
}
