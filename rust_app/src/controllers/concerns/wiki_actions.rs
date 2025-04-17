use actix_web::{web, HttpResponse, Responder};
use serde::{Deserialize, Serialize};
use std::collections::{HashMap, HashSet};
use std::sync::Arc;
use tokio::sync::RwLock;
use uuid::Uuid;

use crate::models::wiki::{Wiki, WikiDirectory, WikiPage};
use crate::services::wiki_pages::{CreateService, UpdateService};
use crate::utils::strong_memoize::StrongMemoize;

/// Module for handling wiki actions
pub trait WikiActions: StrongMemoize {
    /// Get the container (Project or Group)
    fn container(&self) -> Box<dyn std::any::Any>;

    /// Get the current user ID
    fn current_user_id(&self) -> Option<i32>;

    /// Get the wiki
    fn wiki(&self) -> Arc<RwLock<Wiki>> {
        self.strong_memoize("wiki", || {
            let container = self.container();
            // TODO: Implement wiki creation based on container
            Arc::new(RwLock::new(Wiki::new()))
        })
    }

    /// Get the current page
    fn page(&self) -> Option<Arc<RwLock<WikiPage>>> {
        self.strong_memoize("page", || {
            let wiki = self.wiki();
            let page_id = self.page_params().get("id").cloned();
            page_id.and_then(|id| {
                // TODO: Implement page lookup
                Some(Arc::new(RwLock::new(WikiPage::new())))
            })
        })
    }

    /// Get page parameters
    fn page_params(&self) -> HashMap<String, String> {
        let mut params = HashMap::new();
        if let Some(id) = self.request_param("id") {
            params.insert("id".to_string(), id);
        }
        params
    }

    /// Get a request parameter
    fn request_param(&self, name: &str) -> Option<String>;

    /// Create a new wiki page
    async fn create(&self) -> impl Responder {
        let wiki_params = self.wiki_params();
        let container = self.container();
        let current_user_id = self.current_user_id();

        let service = CreateService::new(container, current_user_id);
        match service.execute(wiki_params).await {
            Ok(response) => {
                if response.success {
                    self.handle_action_success("created", response.page).await
                } else {
                    let templates = self.templates_list().await;
                    // TODO: Render edit template with error
                    HttpResponse::Ok().finish()
                }
            }
            Err(_) => HttpResponse::InternalServerError().finish(),
        }
    }

    /// Update a wiki page
    async fn update(&self) -> impl Responder {
        if !self.can_create_wiki() {
            return HttpResponse::Ok().json(serde_json::json!({
                "error": "You are not authorized to create wiki pages"
            }));
        }

        let wiki_params = self.wiki_params();
        let container = self.container();
        let current_user_id = self.current_user_id();

        let service = UpdateService::new(container, current_user_id);
        match service.execute(self.page().unwrap(), wiki_params).await {
            Ok(response) => {
                if response.success {
                    self.handle_action_success("updated", response.page).await
                } else {
                    let templates = self.templates_list().await;
                    // TODO: Render edit template with error
                    HttpResponse::Ok().finish()
                }
            }
            Err(_) => HttpResponse::InternalServerError().finish(),
        }
    }

    /// Show a wiki page
    async fn show(&self) -> impl Responder {
        if let Some(page) = self.page() {
            if !self.valid_encoding() {
                self.set_encoding_error();
            }

            let version_id = self.request_param("version_id");
            let path = page.read().await.path();
            let templates = self.templates_list().await;

            // TODO: Render show template
            HttpResponse::Ok().finish()
        } else if let Some(file_blob) = self.file_blob().await {
            self.send_wiki_file_blob(self.wiki().clone(), file_blob)
                .await
        } else {
            self.handle_redirection().await
        }
    }

    /// Get wiki parameters
    fn wiki_params(&self) -> HashMap<String, String> {
        let mut params = HashMap::new();
        if let Some(title) = self.request_param("title") {
            params.insert("title".to_string(), title);
        }
        if let Some(content) = self.request_param("content") {
            params.insert("content".to_string(), content);
        }
        if let Some(format) = self.request_param("format") {
            params.insert("format".to_string(), format);
        }
        if let Some(message) = self.request_param("message") {
            params.insert("message".to_string(), message);
        }
        if let Some(last_commit_sha) = self.request_param("last_commit_sha") {
            params.insert("last_commit_sha".to_string(), last_commit_sha);
        }
        params
    }

    /// Check if the current user can create wiki pages
    fn can_create_wiki(&self) -> bool {
        // TODO: Implement authorization check
        true
    }

    /// Handle action success
    async fn handle_action_success(
        &self,
        action: &str,
        page: Arc<RwLock<WikiPage>>,
    ) -> impl Responder {
        let page_title = page.read().await.title();
        if page_title == "sidebar" {
            // TODO: Handle sidebar update
            HttpResponse::Ok().finish()
        } else {
            // TODO: Redirect to the updated page
            HttpResponse::Ok().finish()
        }
    }

    /// Get the file blob
    async fn file_blob(&self) -> Option<Vec<u8>> {
        self.strong_memoize("file_blob", || {
            // TODO: Implement file blob lookup
            None
        })
    }

    /// Send wiki file blob
    async fn send_wiki_file_blob(
        &self,
        wiki: Arc<RwLock<Wiki>>,
        file_blob: Vec<u8>,
    ) -> impl Responder {
        // TODO: Implement file blob sending
        HttpResponse::Ok().finish()
    }

    /// Check if the page encoding is valid
    fn valid_encoding(&self) -> bool {
        // TODO: Implement encoding validation
        true
    }

    /// Set encoding error
    fn set_encoding_error(&self) {
        // TODO: Implement encoding error handling
    }

    /// Handle redirection
    async fn handle_redirection(&self) -> impl Responder {
        if self.show_create_form() {
            self.handle_create_form().await
        } else if self.wiki().read().await.exists() {
            // TODO: Render 404 template
            HttpResponse::NotFound().finish()
        } else {
            // TODO: Render empty template
            HttpResponse::Ok().finish()
        }
    }

    /// Handle create form
    async fn handle_create_form(&self) -> impl Responder {
        let title = self.request_param("id").unwrap_or_default();
        let page = self.build_page(title);
        let templates = self.templates_list().await;

        // TODO: Render edit template
        HttpResponse::Ok().finish()
    }

    /// Build a new page
    fn build_page(&self, title: String) -> Arc<RwLock<WikiPage>> {
        Arc::new(RwLock::new(WikiPage::new()))
    }

    /// Check if we should show the create form
    fn show_create_form(&self) -> bool {
        self.can_create_wiki()
    }

    /// Get templates list
    async fn templates_list(&self) -> Vec<Arc<RwLock<WikiPage>>> {
        self.strong_memoize("templates_list", || {
            let wiki = self.wiki();
            // TODO: Implement templates list
            Vec::new()
        })
    }

    /// Find redirection
    fn find_redirection(&self, path: &str, redirect_limit: usize) -> Option<String> {
        let mut seen = HashSet::new();
        let mut current = path.to_string();

        for _ in 0..redirect_limit {
            if seen.contains(&current) {
                return Some(current);
            }
            seen.insert(current.clone());

            if let Some(redir) = self.find_single_redirection(&current) {
                current = redir;
            } else {
                return None;
            }
        }

        None
    }

    /// Find single redirection
    fn find_single_redirection(&self, path: &str) -> Option<String> {
        // TODO: Implement single redirection lookup
        None
    }
}
