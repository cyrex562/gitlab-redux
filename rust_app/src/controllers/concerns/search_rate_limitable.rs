use crate::config::settings::Settings;
use actix_web::{web, HttpResponse};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::sync::Arc;
use std::time::{Duration, Instant, SystemTime};
use tokio::sync::RwLock;

/// Module for implementing rate limiting for search operations
pub trait SearchRateLimitable {
    /// Get the current user ID
    fn user_id(&self) -> i32;

    /// Get the rate limit window in seconds
    fn rate_limit_window(&self) -> u64 {
        60 // Default 1 minute window
    }

    /// Get the maximum number of requests allowed in the window
    fn max_requests_per_window(&self) -> u32 {
        30 // Default 30 requests per minute
    }

    /// Get the rate limit key for the current user
    fn rate_limit_key(&self) -> String {
        format!("search_rate_limit:user_{}", self.user_id())
    }

    /// Check if the current request is rate limited
    async fn is_rate_limited(
        &self,
        storage: Arc<RwLock<HashMap<String, Vec<Instant>>>>,
    ) -> Result<bool, HttpResponse> {
        let key = self.rate_limit_key();
        let window = Duration::from_secs(self.rate_limit_window());
        let now = Instant::now();

        let mut storage = storage.write().await;
        let requests = storage.entry(key).or_insert_with(Vec::new);

        // Remove expired timestamps
        requests.retain(|&timestamp| now.duration_since(timestamp) < window);

        // Check if we're over the limit
        let is_limited = requests.len() >= self.max_requests_per_window() as usize;

        // Add current request timestamp
        requests.push(now);

        Ok(is_limited)
    }

    /// Get the remaining requests for the current window
    async fn remaining_requests(
        &self,
        storage: Arc<RwLock<HashMap<String, Vec<Instant>>>>,
    ) -> Result<u32, HttpResponse> {
        let key = self.rate_limit_key();
        let window = Duration::from_secs(self.rate_limit_window());
        let now = Instant::now();

        let storage = storage.read().await;
        let requests = storage
            .get(&key)
            .map(|reqs| {
                reqs.iter()
                    .filter(|&timestamp| now.duration_since(*timestamp) < window)
                    .count()
            })
            .unwrap_or(0);

        Ok(self
            .max_requests_per_window()
            .saturating_sub(requests as u32))
    }

    /// Get rate limit headers
    async fn rate_limit_headers(
        &self,
        storage: Arc<RwLock<HashMap<String, Vec<Instant>>>>,
    ) -> Result<HashMap<String, String>, HttpResponse> {
        let mut headers = HashMap::new();
        let remaining = self.remaining_requests(storage).await?;
        let reset = Instant::now() + Duration::from_secs(self.rate_limit_window());

        headers.insert(
            "X-RateLimit-Limit".to_string(),
            self.max_requests_per_window().to_string(),
        );
        headers.insert("X-RateLimit-Remaining".to_string(), remaining.to_string());
        headers.insert(
            "X-RateLimit-Reset".to_string(),
            reset.elapsed().as_secs().to_string(),
        );

        Ok(headers)
    }

    /// Enforce rate limiting
    async fn enforce_rate_limit(
        &self,
        storage: Arc<RwLock<HashMap<String, Vec<Instant>>>>,
    ) -> Result<(), HttpResponse> {
        if self.is_rate_limited(storage).await? {
            return Err(HttpResponse::TooManyRequests().json(serde_json::json!({
                "error": "Rate limit exceeded. Please try again later."
            })));
        }
        Ok(())
    }
}
