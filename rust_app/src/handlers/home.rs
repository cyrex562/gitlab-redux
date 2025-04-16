use axum::{
    response::IntoResponse,
    http::StatusCode,
};
use tera::Tera;

pub async fn index() -> impl IntoResponse {
    let mut tera = Tera::new("src/templates/**/*.html").unwrap();
    
    let context = tera::Context::new();
    
    match tera.render("index.html", &context) {
        Ok(html) => (StatusCode::OK, html).into_response(),
        Err(_) => (StatusCode::INTERNAL_SERVER_ERROR, "Template error").into_response(),
    }
} 