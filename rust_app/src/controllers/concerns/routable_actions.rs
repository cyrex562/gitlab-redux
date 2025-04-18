use crate::{
    config::Settings,
    models::{Project, Routable},
    utils::authorization::Can,
};
use actix_web::{
    dev::ServiceRequest, error::Error, http::StatusCode, web::Data, HttpResponse, Responder,
};
use std::sync::Arc;

pub trait RoutableActions {
    fn find_routable<T: Routable>(
        &self,
        routable_klass: &str,
        routable_full_path: &str,
        full_path: &str,
        extra_authorization_proc: Option<Arc<dyn Fn(&T) -> bool + Send + Sync>>,
    ) -> Result<Option<T>, Error>;

    fn not_found_actions(
        &self,
    ) -> Vec<Arc<dyn Fn(&dyn Routable, &str) -> Result<(), Error> + Send + Sync>>;

    fn perform_not_found_actions(
        &self,
        routable: &dyn Routable,
        actions: &[Arc<dyn Fn(&dyn Routable, &str) -> Result<(), Error> + Send + Sync>],
        full_path: &str,
    ) -> Result<(), Error>;

    fn routable_authorized<T: Routable>(
        &self,
        routable: &T,
        extra_authorization_proc: Option<Arc<dyn Fn(&T) -> bool + Send + Sync>>,
    ) -> bool;

    fn ensure_canonical_path<T: Routable>(
        &self,
        routable: &T,
        routable_full_path: &str,
    ) -> Result<(), Error>;
}

pub struct RoutableActionsImpl {
    settings: Data<Settings>,
}

impl RoutableActionsImpl {
    pub fn new(settings: Data<Settings>) -> Self {
        Self { settings }
    }

    fn build_canonical_path<T: Routable>(&self, routable: &T) -> String {
        format!("/{}", routable.full_path())
    }

    fn route_not_found(&self) -> HttpResponse {
        HttpResponse::NotFound()
            .content_type("text/html")
            .body(include_str!("../../templates/errors/not_found.html"))
    }
}

impl RoutableActions for RoutableActionsImpl {
    fn find_routable<T: Routable>(
        &self,
        routable_klass: &str,
        routable_full_path: &str,
        full_path: &str,
        extra_authorization_proc: Option<Arc<dyn Fn(&T) -> bool + Send + Sync>>,
    ) -> Result<Option<T>, Error> {
        let routable = T::find_by_full_path(routable_full_path, true)?;

        if let Some(routable) = &routable {
            if self.routable_authorized(routable, extra_authorization_proc.clone()) {
                self.ensure_canonical_path(routable, routable_full_path)?;
                Ok(routable)
            } else {
                self.perform_not_found_actions(
                    routable.as_ref(),
                    &self.not_found_actions(),
                    full_path,
                )?;
                Ok(None)
            }
        } else {
            Ok(None)
        }
    }

    fn not_found_actions(
        &self,
    ) -> Vec<Arc<dyn Fn(&dyn Routable, &str) -> Result<(), Error> + Send + Sync>> {
        vec![Arc::new(|routable, _full_path| {
            if let Some(project) = routable.as_any().downcast_ref::<Project>() {
                let label = project.external_authorization_classification_label();
                if !crate::utils::external_authorization::access_allowed(label) {
                    let reason = crate::utils::external_authorization::rejection_reason(label)
                        .unwrap_or_else(|| {
                            "External authorization denied access to this project".to_string()
                        });
                    return Err(actix_web::error::ErrorForbidden(reason));
                }
            }
            Ok(())
        })]
    }

    fn perform_not_found_actions(
        &self,
        routable: &dyn Routable,
        actions: &[Arc<dyn Fn(&dyn Routable, &str) -> Result<(), Error> + Send + Sync>],
        full_path: &str,
    ) -> Result<(), Error> {
        for action in actions {
            action(routable, full_path)?;
        }
        Ok(())
    }

    fn routable_authorized<T: Routable>(
        &self,
        routable: &T,
        extra_authorization_proc: Option<Arc<dyn Fn(&T) -> bool + Send + Sync>>,
    ) -> bool {
        let action = format!("read_{}", routable.type_name().to_lowercase());
        if !Can::can(&action, routable) {
            return false;
        }

        if let Some(proc) = extra_authorization_proc {
            proc(routable)
        } else {
            true
        }
    }

    fn ensure_canonical_path<T: Routable>(
        &self,
        routable: &T,
        routable_full_path: &str,
    ) -> Result<(), Error> {
        let canonical_path = routable.full_path();
        if canonical_path != routable_full_path {
            let canonical_url = self.build_canonical_path(routable);
            return Ok(HttpResponse::MovedPermanently()
                .header("Location", canonical_url)
                .finish());
        }
        Ok(())
    }
}
