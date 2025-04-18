use actix_web::{web, HttpRequest, HttpResponse, Result};
use std::sync::Arc;

// Define the Project trait
pub trait Project: Send + Sync {
    fn id(&self) -> i64;
    fn name(&self) -> &str;
}

// Define the User trait
pub trait User: Send + Sync {
    fn id(&self) -> i64;
    fn name(&self) -> &str;
}

// Define the Subscribable trait
pub trait Subscribable: Send + Sync {
    fn id(&self) -> i64;
    fn type_name(&self) -> &str;
    fn toggle_subscription(&self, user: Arc<dyn User>, project: Arc<dyn Project>);
}

pub trait ToggleSubscriptionAction {
    fn toggle_subscription(&self, req: &HttpRequest) -> Result<HttpResponse>;
    fn subscribable_project(&self) -> Arc<dyn Project>;
    fn subscribable_resource(&self) -> Arc<dyn Subscribable>;
}

pub struct ToggleSubscriptionActionHandler {
    current_user: Option<Arc<dyn User>>,
}

impl ToggleSubscriptionActionHandler {
    pub fn new(current_user: Option<Arc<dyn User>>) -> Self {
        ToggleSubscriptionActionHandler { current_user }
    }
}

impl ToggleSubscriptionAction for ToggleSubscriptionActionHandler {
    fn toggle_subscription(&self, req: &HttpRequest) -> Result<HttpResponse> {
        // Check if the user is authenticated
        if let Some(user) = &self.current_user {
            // Get the subscribable resource and project
            let resource = self.subscribable_resource();
            let project = self.subscribable_project();

            // Toggle the subscription
            resource.toggle_subscription(user.clone(), project);

            // Return a success response
            Ok(HttpResponse::Ok().finish())
        } else {
            // Return a success response even if the user is not authenticated
            // This matches the Ruby implementation
            Ok(HttpResponse::Ok().finish())
        }
    }

    fn subscribable_project(&self) -> Arc<dyn Project> {
        // This would be implemented by the concrete class
        unimplemented!("subscribable_project must be implemented")
    }

    fn subscribable_resource(&self) -> Arc<dyn Subscribable> {
        // This would be implemented by the concrete class
        unimplemented!("subscribable_resource must be implemented")
    }
}
