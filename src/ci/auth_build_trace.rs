use std::error::Error;

/// Trait for CI build trace authorization
pub trait AuthBuildTrace {
    /// Authorize read access to build trace
    ///
    /// This method should be implemented by controllers that need to authorize
    /// access to CI build trace resources.
    fn authorize_read_build_trace(&self) -> Result<(), Box<dyn Error>>;

    /// Check if the current user can read the build trace
    fn can_read_build_trace(&self, user: &dyn User, build: &dyn Build) -> bool;

    /// Get the current build
    fn build(&self) -> &dyn Build;
}

/// Trait for User
pub trait User {
    /// Check if the user has a specific permission on a resource
    fn can(&self, permission: &str, resource: &dyn std::any::Any) -> bool;
}

/// Trait for Build
pub trait Build {
    /// Check if debug mode is enabled
    fn is_debug_mode(&self) -> bool;
}

/// Default implementation for AuthBuildTrace
pub struct DefaultAuthBuildTrace {
    build: Box<dyn Build>,
}

impl DefaultAuthBuildTrace {
    pub fn new(build: Box<dyn Build>) -> Self {
        Self { build }
    }
}

impl AuthBuildTrace for DefaultAuthBuildTrace {
    fn authorize_read_build_trace(&self) -> Result<(), Box<dyn Error>> {
        // This would be implemented to actually authorize access to build trace
        Ok(())
    }

    fn can_read_build_trace(&self, user: &dyn User, build: &dyn Build) -> bool {
        user.can("read_build_trace", build)
    }

    fn build(&self) -> &dyn Build {
        &*self.build
    }
}
