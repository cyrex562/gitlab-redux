use actix_web::{HttpRequest, HttpResponse, Responder};

pub trait HotlinkInterceptor {
    fn intercept_hotlinking(&self, req: &HttpRequest) -> impl Responder {
        if self.is_hotlinking(req) {
            self.render_406()
        } else {
            HttpResponse::Ok().finish()
        }
    }

    fn render_406(&self) -> impl Responder {
        HttpResponse::NotAcceptable().finish()
    }

    // Required method to be implemented by concrete types
    fn is_hotlinking(&self, req: &HttpRequest) -> bool;
}
