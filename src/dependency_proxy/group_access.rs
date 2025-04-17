use std::error::Error;

/// Trait for DependencyProxy group access
pub trait GroupAccess {
    /// Verify that the dependency proxy is available
    fn verify_dependency_proxy_available(&self) -> Result<(), Box<dyn Error>>;

    /// Authorize read access to dependency proxy
    fn authorize_read_dependency_proxy(&self) -> Result<(), Box<dyn Error>>;

    /// Get the authenticated user or token
    fn auth_user_or_token(&self) -> Box<dyn std::any::Any>;

    /// Get the group
    fn group(&self) -> Option<&dyn Group>;

    /// Get the authenticated user
    fn auth_user(&self) -> &dyn std::any::Any;

    /// Get the personal access token
    fn personal_access_token(&self) -> Option<&dyn PersonalAccessToken>;
}

/// Trait for Group
pub trait Group {
    /// Check if the dependency proxy feature is available
    fn dependency_proxy_feature_available(&self) -> bool;

    /// Get the dependency proxy for containers policy subject
    fn dependency_proxy_for_containers_policy_subject(&self) -> &dyn std::any::Any;
}

/// Trait for PersonalAccessToken
pub trait PersonalAccessToken {
    /// Get the user associated with this token
    fn user(&self) -> &dyn User;
}

/// Trait for User
pub trait User {
    /// Check if the user is a project bot
    fn is_project_bot(&self) -> bool;

    /// Check if the user is a human
    fn is_human(&self) -> bool;

    /// Check if the user is a service account
    fn is_service_account(&self) -> bool;

    /// Get the resource bot resource
    fn resource_bot_resource(&self) -> &dyn std::any::Any;

    /// Check if the user has a specific permission on a resource
    fn can(&self, permission: &str, resource: &dyn std::any::Any) -> bool;
}

/// Default implementation for GroupAccess
pub struct DefaultGroupAccess {
    group: Option<Box<dyn Group>>,
    auth_user: Box<dyn std::any::Any>,
    personal_access_token: Option<Box<dyn PersonalAccessToken>>,
}

impl DefaultGroupAccess {
    pub fn new(
        group: Option<Box<dyn Group>>,
        auth_user: Box<dyn std::any::Any>,
        personal_access_token: Option<Box<dyn PersonalAccessToken>>,
    ) -> Self {
        Self {
            group,
            auth_user,
            personal_access_token,
        }
    }
}

impl GroupAccess for DefaultGroupAccess {
    fn verify_dependency_proxy_available(&self) -> Result<(), Box<dyn Error>> {
        if let Some(group) = &self.group {
            if !group.dependency_proxy_feature_available() {
                return Err("Dependency proxy feature is not available".into());
            }
        } else {
            return Err("Group not found".into());
        }
        Ok(())
    }

    fn authorize_read_dependency_proxy(&self) -> Result<(), Box<dyn Error>> {
        let auth = self.auth_user_or_token();

        if let Some(user) = auth.downcast_ref::<dyn User>() {
            self.authorize_read_dependency_proxy_for_users(user)?;
        } else {
            self.authorize_read_dependency_proxy_for_tokens(&auth)?;
        }

        Ok(())
    }

    fn auth_user_or_token(&self) -> Box<dyn std::any::Any> {
        if let Some(token) = &self.personal_access_token {
            if let Some(user) = self.auth_user.downcast_ref::<dyn User>() {
                if (user.is_project_bot()
                    && user
                        .resource_bot_resource()
                        .downcast_ref::<dyn Group>()
                        .is_some())
                    || user.is_human()
                    || user.is_service_account()
                {
                    return Box::new(token.user());
                }
            }
        }

        self.auth_user.clone()
    }

    fn group(&self) -> Option<&dyn Group> {
        self.group.as_ref().map(|g| g.as_ref())
    }

    fn auth_user(&self) -> &dyn std::any::Any {
        &*self.auth_user
    }

    fn personal_access_token(&self) -> Option<&dyn PersonalAccessToken> {
        self.personal_access_token.as_ref().map(|t| t.as_ref())
    }

    fn authorize_read_dependency_proxy_for_users(
        &self,
        user: &dyn User,
    ) -> Result<(), Box<dyn Error>> {
        if let Some(group) = &self.group {
            if !user.can("read_dependency_proxy", group) {
                return Err("Access denied".into());
            }
        } else {
            return Err("Group not found".into());
        }
        Ok(())
    }

    fn authorize_read_dependency_proxy_for_tokens(
        &self,
        token: &dyn std::any::Any,
    ) -> Result<(), Box<dyn Error>> {
        if let Some(group) = &self.group {
            if let Some(user) = token.downcast_ref::<dyn User>() {
                if !user.can(
                    "read_dependency_proxy",
                    group.dependency_proxy_for_containers_policy_subject(),
                ) {
                    return Err("Access denied".into());
                }
            } else {
                return Err("Invalid token".into());
            }
        } else {
            return Err("Group not found".into());
        }
        Ok(())
    }
}
