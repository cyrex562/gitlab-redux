pub mod diffs;
pub mod layouts;
pub mod onboarding;
pub mod pajamas;
pub mod rapid_diffs;

pub use diffs::*;
pub use layouts::*;
pub use onboarding::*;
pub use pajamas::*;
pub use rapid_diffs::*;

use serde::Serialize;

/// Trait for all components
pub trait Component {
    /// Render the component to HTML
    fn render(&self) -> String;
}
