use std::collections::HashMap;

/// Trait for components that need to format input options
pub trait CheckboxRadioOptions {
    /// Get the input options
    fn input_options(&self) -> HashMap<String, String>;

    /// Format the input options with CSS classes
    fn formatted_input_options(&self) -> HashMap<String, String> {
        let options = self.input_options();
        let css_classes = vec!["custom-control-input"];

        crate::components::pajamas::Component::format_options(options, &css_classes, HashMap::new())
    }
}
