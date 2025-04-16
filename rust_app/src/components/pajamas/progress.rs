/// Renders a Pajamas compliant progress bar
#[derive(Debug, Clone)]
pub struct ProgressComponent {
    /// The value of the progress bar (0-100)
    value: u8,
    /// The variant of the progress bar
    variant: ProgressVariant,
}

/// The variant of the progress bar
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum ProgressVariant {
    /// Primary variant
    Primary,
    /// Success variant
    Success,
}

impl ProgressComponent {
    /// Create a new progress component
    pub fn new(value: u8, variant: ProgressVariant) -> Self {
        Self {
            value: value.min(100), // Ensure value is between 0 and 100
            variant,
        }
    }

    /// Get the variant class
    fn variant_class(&self) -> &'static str {
        match self.variant {
            ProgressVariant::Primary => "gl-progress-bar-primary",
            ProgressVariant::Success => "gl-progress-bar-success",
        }
    }

    /// Render the progress component
    pub fn render(&self) -> String {
        format!(
            r#"<div class="gl-progress-bar progress"><div class="gl-progress {}" style="width: {}%;"></div></div>"#,
            self.variant_class(),
            self.value
        )
    }
}
