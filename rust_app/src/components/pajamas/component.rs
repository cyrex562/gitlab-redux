use std::collections::HashMap;

/// Base component for Pajamas UI components
pub struct Component;

impl Component {
    /// Filter a given value against a list of allowed values
    /// If no value is given or value is not allowed return default one
    ///
    /// # Arguments
    ///
    /// * `value` - The value to filter
    /// * `allowed_values` - The list of allowed values
    /// * `default` - The default value to return if the value is not allowed
    pub fn filter_attribute<T: PartialEq>(
        value: Option<&T>,
        allowed_values: &[T],
        default: Option<&T>,
    ) -> Option<&T> {
        match value {
            Some(val) => {
                if allowed_values.contains(val) {
                    Some(val)
                } else {
                    default
                }
            }
            None => default,
        }
    }

    /// Add CSS classes and additional options to an existing options hash
    ///
    /// # Arguments
    ///
    /// * `options` - The existing options hash
    /// * `css_classes` - The CSS classes to add
    /// * `additional_options` - Additional options to add
    pub fn format_options(
        mut options: HashMap<String, String>,
        css_classes: &[&str],
        additional_options: HashMap<String, String>,
    ) -> HashMap<String, String> {
        // Get existing class or create empty string
        let existing_class = options.remove("class").unwrap_or_default();

        // Combine classes
        let mut all_classes = css_classes.to_vec();
        if !existing_class.is_empty() {
            all_classes.push(&existing_class);
        }

        // Join classes with space
        let class_str = all_classes.join(" ");

        // Add class back to options
        options.insert("class".to_string(), class_str);

        // Add additional options
        for (key, value) in additional_options {
            options.insert(key, value);
        }

        options
    }
}
