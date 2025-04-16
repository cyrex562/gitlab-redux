use actix_web::{web, HttpResponse, Responder};
use serde::Deserialize;
use crate::controllers::base::BaseController;
use crate::models::{Agent, User};

#[derive(Debug, Deserialize)]
pub struct AgentParams {
    agent_id: i32,
}

pub struct DashboardController;

impl DashboardController {
    #[get("/clusters/agents")]
    pub async fn index(
        current_user: Option<web::ReqData<User>>,
    ) -> impl Responder {
        // TODO: Implement agents dashboard index
        // This should return the agents dashboard
        HttpResponse::Ok().json(serde_json::json!({
            "status": "ok",
            "message": "Agents Dashboard"
        }))
    }

    #[get("/clusters/agents/{agent_id}")]
    pub async fn show(
        params: web::Path<AgentParams>,
        current_user: Option<web::ReqData<User>>,
    ) -> impl Responder {
        // TODO: Implement agent details
        // This should return the details of the specified agent
        HttpResponse::Ok().json(serde_json::json!({
            "id": params.agent_id,
            "name": "Unknown Agent"
        }))
    }

    async fn find_agent(params: &AgentParams) -> Result<Agent, actix_web::Error> {
        Agent::find(params.agent_id)
            .await
            .map_err(|e| actix_web::error::ErrorNotFound(e))
    }

    async fn check_feature_flag(current_user: Option<&User>) -> bool {
        // TODO: Implement feature flag check
        // This should check if the k8s_dashboard feature is enabled
        true
    }

    async fn authorize_read_cluster_agent(
        current_user: Option<&User>,
        agent: &Agent,
    ) -> bool {
        // TODO: Implement authorization check
        // This should check if the user has read_cluster_agent permission
        true
    }
} 