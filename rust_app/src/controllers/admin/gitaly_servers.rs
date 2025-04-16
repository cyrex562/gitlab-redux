use actix_web::{get, web, HttpResponse, Responder};
use crate::models::gitaly::Server;
use crate::auth::authorize;

#[get("/admin/gitaly/servers")]
pub async fn index() -> impl Responder {
    // TODO: Implement authorization check
    // authorize!("read_admin_gitaly_servers")?;

    // TODO: Implement Gitaly::Server.all equivalent
    let gitaly_servers = Vec::new(); // Placeholder until we implement the model

    HttpResponse::Ok().json(gitaly_servers)
} 