use std::error::Error;

/// Trait for Harbor access authorization
pub trait HarborAccess {
    /// Authorize read access to Harbor registry
    ///
    /// This method should be implemented by controllers that need to authorize
    /// access to Harbor registry resources.
    fn authorize_read_harbor_registry(&self) -> Result<(), Box<dyn Error>> {
        Err("Not implemented".into())
    }
}

/// Default implementation for HarborAccess
pub struct DefaultHarborAccess;

impl HarborAccess for DefaultHarborAccess {
    fn authorize_read_harbor_registry(&self) -> Result<(), Box<dyn Error>> {
        Err("Not implemented".into())
    }
}
