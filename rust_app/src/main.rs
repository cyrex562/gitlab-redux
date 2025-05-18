use actix_web::middleware::{Compress, Logger};
use actix_web::{web, App, HttpServer};
use tower_http::cors::{Any, CorsLayer};

mod config;
mod controllers;
mod handlers;
mod models;
mod services;
mod websocket;

use crate::controllers::{oauth, organizations};
use crate::handlers::health;

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    // Initialize tracing
    tracing_subscriber::fmt::init();

    // Initialize CORS with proper type annotations
    let cors = CorsLayer::new()
        .allow_origin(Any)
        .allow_methods(vec!["GET", "POST", "PUT", "DELETE"].into_iter());

    // Create the HTTP server
    HttpServer::new(move || {
        App::new()
            .wrap(Logger::default())
            .wrap(Compress::default())
            .app_data(web::Data::new(cors.clone()))
            .service(web::scope("/api").configure(configure_api))
            .service(web::scope("/oauth").configure(configure_oauth))
            .service(
                web::scope("/organizations")
                    .configure(organizations::projects_controller::configure),
            )
            .service(web::scope("").configure(oauth::applications_controller::configure))
    })
    .bind(("127.0.0.1", 3000))?
    .run()
    .await
}

fn configure_api(cfg: &mut web::ServiceConfig) {
    cfg.service(web::scope("/health").route("", web::get().to(health::health_check)))
        .service(
            web::scope("/organizations").configure(organizations::projects_controller::configure),
        );
}

fn configure_oauth(cfg: &mut web::ServiceConfig) {
    cfg.service(
        web::scope("/device")
            .route("/confirm", web::post().to(oauth::device_authorizations_controller::DeviceAuthorizationsController::confirm))
            .route("", web::get().to(oauth::device_authorizations_controller::DeviceAuthorizationsController::index))
    );
}
