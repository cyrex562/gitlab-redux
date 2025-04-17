use actix_web::{web, HttpResponse, Responder};
use serde::Deserialize;
use std::collections::HashMap;

use crate::models::web_hook::WebHook;
use crate::models::web_hook_log::WebHookLog;
use crate::services::web_hooks::{CreateService, DestroyService};

const SECRET_MASK: &str = "********";

#[derive(Debug, Deserialize)]
pub struct HookParams {
    enable_ssl_verification: Option<bool>,
    name: Option<String>,
    description: Option<String>,
    token: Option<String>,
    url: Option<String>,
    push_events_branch_filter: Option<String>,
    branch_filter_strategy: Option<String>,
    custom_webhook_template: Option<String>,
    url_variables: Option<Vec<UrlVariable>>,
    custom_headers: Option<Vec<CustomHeader>>,
}

#[derive(Debug, Deserialize)]
pub struct UrlVariable {
    key: String,
    value: Option<String>,
}

#[derive(Debug, Deserialize)]
pub struct CustomHeader {
    key: String,
    value: Option<String>,
}

pub trait HookActions {
    fn hooks(&self) -> &Vec<WebHook>;
    fn set_hooks(&mut self, hooks: Vec<WebHook>);
    fn hook(&self) -> &WebHook;
    fn set_hook(&mut self, hook: WebHook);
    fn current_user(&self) -> &User;
    fn relation(&self) -> &dyn WebHookRelation;

    fn index(&self) -> impl Responder {
        let hooks = self.relation().select_persisted();
        // In a real implementation, we'd need to handle the view rendering
        HttpResponse::Ok().finish()
    }

    fn create(&mut self, params: web::Json<HookParams>) -> impl Responder {
        let service = CreateService::new(self.current_user());
        let result = service.execute(params.into_inner(), self.relation());

        if result.is_success() {
            // In a real implementation, we'd need to handle flash messages
            HttpResponse::Found()
                .header("Location", "/web_hooks")
                .finish()
        } else {
            self.set_hooks(self.relation().select_persisted());
            // In a real implementation, we'd need to handle flash messages
            HttpResponse::Found()
                .header("Location", "/web_hooks")
                .finish()
        }
    }

    fn update(&mut self, params: web::Json<HookParams>) -> impl Responder {
        if self.hook().update(params.into_inner()) {
            // In a real implementation, we'd need to handle flash messages
            HttpResponse::Found()
                .header("Location", "/web_hooks/edit")
                .finish()
        } else {
            // In a real implementation, we'd need to handle the view rendering
            HttpResponse::Ok().finish()
        }
    }

    fn destroy(&self) -> impl Responder {
        self.destroy_hook(self.hook());
        HttpResponse::Found()
            .header("Location", "/web_hooks")
            .finish()
    }

    fn edit(&self) -> impl Responder {
        if self.hook().is_some() {
            // In a real implementation, we'd need to handle the view rendering
            HttpResponse::Ok().finish()
        } else {
            HttpResponse::Found()
                .header("Location", "/web_hooks")
                .finish()
        }
    }

    fn hook_logs(&self, page: i32) -> Vec<WebHookLog> {
        self.hook()
            .web_hook_logs()
            .recent()
            .page(page)
            .without_count()
    }

    fn destroy_hook(&self, hook: &WebHook) {
        let service = DestroyService::new(self.current_user());
        let result = service.execute(hook);

        match result.status {
            DestroyStatus::Success { async: true } => {
                // In a real implementation, we'd need to handle flash messages
            }
            DestroyStatus::Success { async: false } => {
                // In a real implementation, we'd need to handle flash messages
            }
            DestroyStatus::Error { message } => {
                // In a real implementation, we'd need to handle flash messages
            }
        }
    }
}

#[derive(Debug)]
pub enum DestroyStatus {
    Success { async: bool },
    Error { message: String },
}

pub trait WebHookRelation {
    fn select_persisted(&self) -> Vec<WebHook>;
}
