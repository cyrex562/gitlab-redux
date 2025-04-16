use actix_web::{web, HttpResponse, Responder};
use serde::Deserialize;
use crate::controllers::base::BaseController;
use crate::models::{Cluster, User};

#[derive(Debug, Deserialize)]
pub struct ClusterParams {
    id: i32,
}

pub struct ClustersBaseController;

impl ClustersBaseController {
    pub async fn cluster(
        clusterable: &dyn crate::models::Clusterable,
        params: &ClusterParams,
        current_user: Option<&User>,
    ) -> Result<Cluster, actix_web::Error> {
        clusterable
            .find_cluster(params.id)
            .await
            .map_err(|e| actix_web::error::ErrorNotFound(e))
    }

    pub async fn authorize_update_cluster(
        current_user: Option<&User>,
        clusterable: &dyn crate::models::Clusterable,
    ) -> bool {
        // TODO: Implement proper authorization checks
        // This should check if the user has update_cluster permission
        true
    }

    pub async fn authorize_admin_cluster(
        current_user: Option<&User>,
        clusterable: &dyn crate::models::Clusterable,
    ) -> bool {
        // TODO: Implement proper authorization checks
        // This should check if the user has admin_cluster permission
        true
    }

    pub async fn authorize_read_cluster(
        current_user: Option<&User>,
        clusterable: &dyn crate::models::Clusterable,
    ) -> bool {
        // TODO: Implement proper authorization checks
        // This should check if the user has read_cluster permission
        true
    }

    pub async fn authorize_create_cluster(
        current_user: Option<&User>,
        clusterable: &dyn crate::models::Clusterable,
    ) -> bool {
        // TODO: Implement proper authorization checks
        // This should check if the user has create_cluster permission
        true
    }

    pub async fn authorize_read_prometheus(
        current_user: Option<&User>,
        clusterable: &dyn crate::models::Clusterable,
    ) -> bool {
        // TODO: Implement proper authorization checks
        // This should check if the user has read_prometheus permission
        true
    }
} 