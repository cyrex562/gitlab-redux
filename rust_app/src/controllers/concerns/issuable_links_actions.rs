// Ported from: orig_app/app/controllers/concerns/issuable_links.rb
// Ported on: 2025-04-25

use actix_web::{web, HttpResponse, Responder};
use serde_json::json;

/// Trait for controller actions related to issuable links
#[allow(unused_variables)]
#[async_trait::async_trait]
pub trait IssuableLinksActions {
    /// List issuables (GET /...)
    async fn index(&self) -> impl Responder {
        let issuables = self.issuables().await;
        HttpResponse::Ok().json(issuables)
    }

    /// Create issuable links (POST /...)
    async fn create(&self) -> impl Responder {
        let result = self.create_service().await.execute().await;
        let issuables = self.issuables().await;
        HttpResponse::build(result.http_status).json(json!({
            "message": result.message,
            "issuables": issuables
        }))
    }

    /// Destroy issuable links (DELETE /...)
    async fn destroy(&self) -> impl Responder {
        let result = self.destroy_service().await.execute().await;
        let issuables = self.issuables().await;
        HttpResponse::build(result.http_status).json(json!({
            "issuables": issuables
        }))
    }

    /// List issuables (service call)
    async fn issuables(&self) -> serde_json::Value {
        self.list_service().await.execute().await
    }

    /// Service for listing issuables
    async fn list_service(&self) -> Box<dyn ListService> {
        unimplemented!("list_service must be implemented by the controller");
    }

    /// Params for create (simulate strong params)
    fn create_params(&self, params: &serde_json::Value) -> serde_json::Value {
        // In a real implementation, filter params here
        params.clone()
    }

    /// Service for creating links
    async fn create_service(&self) -> Box<dyn CreateService> {
        unimplemented!("create_service must be implemented by the controller");
    }

    /// Service for destroying links
    async fn destroy_service(&self) -> Box<dyn DestroyService> {
        unimplemented!("destroy_service must be implemented by the controller");
    }
}

// Traits for service objects (to be implemented elsewhere)
#[async_trait::async_trait]
pub trait ListService: Send + Sync {
    async fn execute(&self) -> serde_json::Value;
}

#[async_trait::async_trait]
pub trait CreateService: Send + Sync {
    async fn execute(&self) -> ServiceResult;
}

#[async_trait::async_trait]
pub trait DestroyService: Send + Sync {
    async fn execute(&self) -> ServiceResult;
}

pub struct ServiceResult {
    pub message: String,
    pub http_status: actix_web::http::StatusCode,
}
