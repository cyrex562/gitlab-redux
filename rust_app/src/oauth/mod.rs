// Ported from gitlab-redux/orig_app/app/controllers/oauth/* Ruby controllers
pub mod applications;
pub mod authorizations;
pub mod authorized_applications;
pub mod device_authorizations;
pub mod token_info;
pub mod tokens;

use axum::{
    routing::{get, post},
    Router,
};

pub fn routes() -> Router {
    Router::new()
        .nest(
            "/oauth",
            Router::new()
                .route("/token", post(tokens::create))
                .route("/token/info", get(token_info::show))
                .route("/token/revoke", post(tokens::revoke))
                .route("/authorize", get(authorizations::new).post(authorizations::create))
                .route("/device_authorizations", get(device_authorizations::index))
                .route("/device_authorizations/confirm", get(device_authorizations::confirm))
                .route("/applications", get(applications::index).post(applications::create))
                .route("/applications/:id", get(applications::show))
                .route("/applications/:id/renew", post(applications::renew))
                .route("/authorized_applications", get(authorized_applications::index))
                .route("/authorized_applications/:id", post(authorized_applications::destroy))
        )
}