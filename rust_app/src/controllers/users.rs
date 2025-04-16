use crate::controllers::ApplicationController;
use actix_web::http::header;
use actix_web::{delete, get, post, put, web, HttpResponse, Responder};
use chrono::{Date, NaiveDate, Utc};
use serde::{Deserialize, Serialize};
use std::sync::Arc;
use tokio::sync::RwLock;

const FOLLOWERS_FOLLOWING_USERS_PER_PAGE: usize = 21;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct User {
    id: i64,
    username: String,
    name: String,
    email: String,
    // TODO: Add more user fields
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct UserQueryParams {
    page: Option<i32>,
    per_page: Option<i32>,
    skip_pagination: Option<bool>,
    skip_namespace: Option<bool>,
    compact_mode: Option<bool>,
    card_mode: Option<bool>,
    scope: Option<String>,
    date: Option<String>,
}

#[derive(Debug)]
pub struct UsersController {
    app_controller: Arc<ApplicationController>,
}

impl UsersController {
    pub fn new(app_controller: Arc<ApplicationController>) -> Self {
        Self { app_controller }
    }

    #[get("/users/{username}")]
    pub async fn show(&self, username: web::Path<String>) -> impl Responder {
        // TODO: Implement user lookup and profile display
        HttpResponse::Ok().body("User profile")
    }

    #[get("/users/{username}/ssh_keys")]
    pub async fn ssh_keys(&self, username: web::Path<String>) -> impl Responder {
        // TODO: Implement SSH keys retrieval
        HttpResponse::Ok().body("SSH keys")
    }

    #[get("/users/{username}/activity")]
    pub async fn activity(
        &self,
        username: web::Path<String>,
        query: web::Query<UserQueryParams>,
    ) -> impl Responder {
        // TODO: Implement user activity
        HttpResponse::Ok().json(serde_json::json!({
            "message": "User activity"
        }))
    }

    #[get("/users/{username}/gpg_keys")]
    pub async fn gpg_keys(&self, username: web::Path<String>) -> impl Responder {
        // TODO: Implement GPG keys retrieval
        HttpResponse::Ok().body("GPG keys")
    }

    #[get("/users/{username}/groups")]
    pub async fn groups(
        &self,
        username: web::Path<String>,
        query: web::Query<UserQueryParams>,
    ) -> impl Responder {
        // TODO: Implement user groups
        HttpResponse::Ok().json(serde_json::json!({
            "html": "Groups list"
        }))
    }

    #[get("/users/{username}/projects")]
    pub async fn projects(
        &self,
        username: web::Path<String>,
        query: web::Query<UserQueryParams>,
    ) -> impl Responder {
        // TODO: Implement user projects
        self.present_projects(username, query, |_| {
            // TODO: Implement project loading
            vec![]
        })
        .await
    }

    #[get("/users/{username}/contributed")]
    pub async fn contributed(
        &self,
        username: web::Path<String>,
        query: web::Query<UserQueryParams>,
    ) -> impl Responder {
        // TODO: Implement contributed projects
        self.present_projects(username, query, |_| {
            // TODO: Implement contributed projects loading
            vec![]
        })
        .await
    }

    #[get("/users/{username}/starred")]
    pub async fn starred(
        &self,
        username: web::Path<String>,
        query: web::Query<UserQueryParams>,
    ) -> impl Responder {
        // TODO: Implement starred projects
        self.present_projects(username, query, |_| {
            // TODO: Implement starred projects loading
            vec![]
        })
        .await
    }

    #[get("/users/{username}/followers")]
    pub async fn followers(
        &self,
        username: web::Path<String>,
        query: web::Query<UserQueryParams>,
    ) -> impl Responder {
        // TODO: Implement followers
        self.present_users(username, query, |_| {
            // TODO: Implement followers loading
            vec![]
        })
        .await
    }

    #[get("/users/{username}/following")]
    pub async fn following(
        &self,
        username: web::Path<String>,
        query: web::Query<UserQueryParams>,
    ) -> impl Responder {
        // TODO: Implement following
        self.present_users(username, query, |_| {
            // TODO: Implement following loading
            vec![]
        })
        .await
    }

    #[get("/users/{username}/snippets")]
    pub async fn snippets(
        &self,
        username: web::Path<String>,
        query: web::Query<UserQueryParams>,
    ) -> impl Responder {
        // TODO: Implement snippets
        HttpResponse::Ok().json(serde_json::json!({
            "html": "Snippets list"
        }))
    }

    #[get("/users/{username}/calendar")]
    pub async fn calendar(&self, username: web::Path<String>) -> impl Responder {
        // TODO: Implement contributions calendar
        HttpResponse::Ok().json(serde_json::json!({
            "activity_dates": {}
        }))
    }

    #[get("/users/{username}/calendar_activities")]
    pub async fn calendar_activities(
        &self,
        username: web::Path<String>,
        query: web::Query<UserQueryParams>,
    ) -> impl Responder {
        // TODO: Implement calendar activities
        let calendar_date = query
            .date
            .as_ref()
            .and_then(|date_str| NaiveDate::parse_from_str(date_str, "%Y-%m-%d").ok())
            .unwrap_or_else(|| Utc::now().date_naive());

        HttpResponse::Ok().body("Calendar activities")
    }

    #[get("/users/exists")]
    pub async fn exists(&self, query: web::Query<UserQueryParams>) -> impl Responder {
        // TODO: Implement username existence check
        HttpResponse::Ok().json(serde_json::json!({
            "exists": false
        }))
    }

    #[post("/users/{username}/follow")]
    pub async fn follow(&self, username: web::Path<String>) -> impl Responder {
        // TODO: Implement follow functionality
        HttpResponse::Ok().body("Followed user")
    }

    #[delete("/users/{username}/unfollow")]
    pub async fn unfollow(&self, username: web::Path<String>) -> impl Responder {
        // TODO: Implement unfollow functionality
        HttpResponse::Ok().body("Unfollowed user")
    }

    async fn present_projects<F>(
        &self,
        username: web::Path<String>,
        query: web::Query<UserQueryParams>,
        loader: F,
    ) -> impl Responder
    where
        F: FnOnce(&str) -> Vec<String>,
    {
        let projects = loader(&username);

        HttpResponse::Ok().json(serde_json::json!({
            "html": "Projects list",
            "projects": projects,
            "skip_pagination": query.skip_pagination.unwrap_or(false),
            "skip_namespace": query.skip_namespace.unwrap_or(false),
            "compact_mode": query.compact_mode.unwrap_or(false),
            "card_mode": query.card_mode.unwrap_or(false)
        }))
    }

    async fn present_users<F>(
        &self,
        username: web::Path<String>,
        query: web::Query<UserQueryParams>,
        loader: F,
    ) -> impl Responder
    where
        F: FnOnce(&str) -> Vec<String>,
    {
        let users = loader(&username);

        HttpResponse::Ok().json(serde_json::json!({
            "html": "Users list",
            "users": users
        }))
    }
}

// Helper functions
impl UsersController {
    fn find_user(&self, username: &str) -> Option<User> {
        // TODO: Implement user lookup
        None
    }

    fn authorize_read_user_profile(&self, user: &User) -> bool {
        // TODO: Implement authorization check
        true
    }
}
