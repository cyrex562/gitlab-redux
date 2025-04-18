use actix_web::{dev::ServiceRequest, error::Error};
use crate::{
    models::{User, Group, Project},
    services::users::ActivityService,
    utils::{
        database::Database,
        event_store::EventStore,
        cookies::CookiesHelper,
    },
};

pub trait RecordUserLastActivity {
    fn set_user_last_activity(&self, req: &ServiceRequest) -> Result<(), Error>;
    fn set_member_last_activity(&self, req: &ServiceRequest) -> Result<(), Error>;
}

pub struct RecordUserLastActivityImpl {
    cookies_helper: Box<dyn CookiesHelper>,
}

impl RecordUserLastActivityImpl {
    pub fn new(cookies_helper: Box<dyn CookiesHelper>) -> Self {
        Self { cookies_helper }
    }
}

impl RecordUserLastActivity for RecordUserLastActivityImpl {
    fn set_user_last_activity(&self, req: &ServiceRequest) -> Result<(), Error> {
        if req.method() != actix_web::http::Method::GET {
            return Ok(());
        }

        if Database::is_read_only() {
            return Ok(());
        }

        if let Some(user) = req.extensions().get::<User>() {
            ActivityService::new(user.clone()).execute()?;
        }

        Ok(())
    }

    fn set_member_last_activity(&self, req: &ServiceRequest) -> Result<(), Error> {
        let context = req.extensions().get::<Group>()
            .or_else(|| req.extensions().get::<Project>());

        if let (Some(user), Some(context)) = (
            req.extensions().get::<User>(),
            context,
        ) {
            if context.is_persisted() {
                EventStore::publish(serde_json::json!({
                    "type": "Users::ActivityEvent",
                    "data": {
                        "user_id": user.id,
                        "namespace_id": context.root_ancestor().id
                    }
                }))?;
            }
        }

        Ok(())
    }
} 