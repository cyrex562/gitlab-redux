use crate::controllers::concerns::ProductAnalyticsTracking;
use crate::controllers::explore::ExploreApplicationController;
use crate::services::ci::catalog::Listing;
use actix_web::{web, HttpResponse};

pub struct CatalogController {
    base: ExploreApplicationController,
}

impl CatalogController {
    pub fn new(base: ExploreApplicationController) -> Self {
        Self { base }
    }

    pub async fn show(&self, params: web::Path<String>) -> HttpResponse {
        if !self.check_resource_access(&params).await {
            return HttpResponse::NotFound().finish();
        }

        // TODO: Implement HTML template rendering
        HttpResponse::Ok().finish()
    }

    pub async fn index(&self) -> HttpResponse {
        self.show(web::Path::from("".to_string())).await
    }

    async fn check_resource_access(&self, params: &web::Path<String>) -> bool {
        self.catalog_resource(params).await.is_some()
    }

    async fn catalog_resource(&self, params: &web::Path<String>) -> Option<CatalogResource> {
        Listing::new(self.base.user.clone())
            .find_resource(params.to_string())
            .await
    }
}

impl ProductAnalyticsTracking for CatalogController {
    fn tracking_namespace_source(&self) -> Option<String> {
        self.base.user.namespace()
    }

    fn tracking_project_source(&self) -> Option<String> {
        None
    }
}
