use actix_web::{get, web, HttpResponse, Responder};
use crate::controllers::base::BaseController;

pub struct GraphqlExplorerController;

impl GraphqlExplorerController {
    #[get("/api/graphql/explorer")]
    pub async fn show() -> impl Responder {
        // TODO: Implement gon variables setup for Apollo client
        // This will need to be adapted to work with the Rust frontend setup
        HttpResponse::Ok().json(serde_json::json!({
            "status": "ok",
            "message": "GraphQL Explorer endpoint"
        }))
    }
} 