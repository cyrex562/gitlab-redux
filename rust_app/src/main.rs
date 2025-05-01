use axum::{
    routing::{get, post},
    Router,
};
use hyper::Server;
use std::net::SocketAddr;
use tower_http::{cors::CorsLayer, services::ServeDir, trace::TraceLayer};
use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt};

mod config;
mod controllers;
mod gitlab;
mod graphql;
mod handlers;
mod models;
mod services;
mod templates;
mod utils;
mod websocket;
use crate::controllers::admin::{
    AbuseReportsController, ApplicationsController, BackgroundJobsController,
    BackgroundMigrationsController, BatchedJobsController, BroadcastMessagesController,
    ClustersController, CohortsController, DashboardController, HealthCheckController,
    SystemInfoController, VersionCheckController,
};
use crate::controllers::groups::settings::repository_controller;
use crate::controllers::jira_connect::{app_descriptor_controller, oauth_callbacks_controller};
use crate::controllers::{
    application, health, jira_connect_app_descriptor, sessions, users, well_known,
};

use anyhow::Result;
use config::Config;
use gitlab::GitlabClient;
use graphql::GitlabClient as GraphQLClient;

#[tokio::main]
async fn main() -> Result<()> {
    // Initialize tracing
    tracing_subscriber::registry()
        .with(tracing_subscriber::EnvFilter::new(
            std::env::var("RUST_LOG").unwrap_or_else(|_| "info".into()),
        ))
        .with(tracing_subscriber::fmt::layer())
        .init();

    // Load environment variables
    dotenvy::dotenv().ok();

    // Initialize database pool
    let database_url = std::env::var("DATABASE_URL").expect("DATABASE_URL must be set");
    let pool = sqlx::postgres::PgPoolOptions::new()
        .max_connections(5)
        .connect(&database_url)
        .await
        .expect("Failed to connect to Postgres");

    // Initialize Redis connection
    let redis_url = std::env::var("REDIS_URL").expect("REDIS_URL must be set");
    let redis = redis::Client::open(redis_url).expect("Failed to connect to Redis");

    // Initialize auth service
    let auth_service = std::sync::Arc::new(services::AuthService::new());

    // Initialize Gitlab client
    let config = Config::load()?;
    let client = GitlabClient::new(&config)?;
    let graphql_client = GraphQLClient::new(&config)?;

    // Example: Get project details
    let project = client.get_project("your-project-path").await?;
    println!("Project: {:?}", project);

    // Example: Get merge requests
    let mrs = client.get_merge_requests("your-project-path").await?;
    println!("Merge Requests: {:?}", mrs);

    // Example: Get pipeline status
    let pipeline = client.get_pipeline("your-project-path", 123).await?;
    println!("Pipeline: {:?}", pipeline);

    // Example: Get project details via GraphQL
    let project_query = r#"
        query ProjectDetails($fullPath: ID!) {
            project(fullPath: $fullPath) {
                id
                name
                description
                webUrl
                repository {
                    rootRef
                }
            }
        }
    "#;

    let variables = serde_json::json!({
        "fullPath": "your-project-path"
    });

    let result = graphql_client.query(project_query, variables).await?;
    println!("GraphQL Result: {:?}", result);

    // Build our application with a route
    let app = Router::new()
        .route("/", get(handlers::home::index))
        .route("/health", get(handlers::health::check))
        .route(
            "/jira_connect/app_descriptor.json",
            get(app_descriptor_controller::show),
        )
        .route(
            "/jira_connect/app_descriptor.json",
            get(jira_connect_app_descriptor::app_descriptor),
        )
        .route(
            "/groups/:group_id/settings/repository/deploy_token",
            post(repository_controller::create_deploy_token),
        )
        .route(
            "/jira_connect/oauth/callback",
            get(oauth_callbacks_controller::index),
        )
        .nest_service("/assets", ServeDir::new("static/assets"))
        .layer(TraceLayer::new_for_http())
        .layer(CorsLayer::permissive())
        .with_state(AppState {
            db: pool,
            redis,
            auth_service,
        });

    // Register controllers
    app.service(health::HealthController::new().configure(health::configure))
        .service(application::ApplicationController::new().configure(application::configure))
        .service(users::UsersController::new().configure(users::configure))
        .service(sessions::SessionsController::new().configure(sessions::configure))
        .service(well_known::WellKnownController::new().configure(well_known::configure))
        // Register admin controllers
        .service(AbuseReportsController::new().configure(AbuseReportsController::configure))
        .service(ApplicationsController::new().configure(ApplicationsController::configure))
        .service(BackgroundJobsController::new().configure(BackgroundJobsController::configure))
        .service(
            BackgroundMigrationsController::new()
                .configure(BackgroundMigrationsController::configure),
        )
        .service(BatchedJobsController::new().configure(BatchedJobsController::configure))
        .service(
            BroadcastMessagesController::new().configure(BroadcastMessagesController::configure),
        )
        .service(ClustersController::new().configure(ClustersController::configure))
        .service(CohortsController::new().configure(CohortsController::configure))
        .service(DashboardController::new().configure(DashboardController::configure))
        .service(HealthCheckController::new().configure(HealthCheckController::configure))
        .service(SystemInfoController::new().configure(SystemInfoController::configure))
        .service(VersionCheckController::new().configure(VersionCheckController::configure));

    // Run it with hyper on localhost:3000
    let addr = SocketAddr::from(([127, 0, 0, 1], 3000));
    tracing::info!("listening on {}", addr);
    Server::bind(&addr)
        .serve(app.into_make_service())
        .await
        .unwrap();

    Ok(())
}

// Application state
#[derive(Clone)]
pub struct AppState {
    pub db: sqlx::PgPool,
    pub redis: redis::Client,
    pub auth_service: std::sync::Arc<services::AuthService>,
}
