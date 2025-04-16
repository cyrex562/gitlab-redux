use actix_web::{web, HttpResponse};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::time::{Duration, Instant};

/// Module for handling performance bar functionality
pub trait WithPerformanceBar {
    /// Get the performance bar enabled status
    fn performance_bar_enabled(&self) -> bool {
        false // Implement based on your needs
    }

    /// Get the performance bar threshold
    fn performance_bar_threshold(&self) -> Duration {
        Duration::from_millis(1000) // Default 1 second threshold
    }

    /// Get the performance bar request ID
    fn performance_bar_request_id(&self) -> Option<String> {
        None // Implement based on your needs
    }

    /// Get the performance bar user ID
    fn performance_bar_user_id(&self) -> Option<i32> {
        None // Implement based on your needs
    }

    /// Get the performance bar data
    fn performance_bar_data(&self) -> HashMap<String, serde_json::Value> {
        HashMap::new() // Implement based on your needs
    }

    /// Check if performance bar should be shown
    fn should_show_performance_bar(&self, duration: Duration) -> bool {
        self.performance_bar_enabled() && duration >= self.performance_bar_threshold()
    }

    /// Add performance bar headers to response
    fn add_performance_bar_headers(
        &self,
        response: &mut HttpResponse,
        duration: Duration,
    ) -> Result<(), HttpResponse> {
        if !self.should_show_performance_bar(duration) {
            return Ok(());
        }

        let data = self.performance_bar_data();
        let request_id = self.performance_bar_request_id();
        let user_id = self.performance_bar_user_id();

        // Add performance bar headers
        response.headers_mut().insert(
            "X-Performance-Bar-Enabled",
            "true".parse().map_err(|e| {
                HttpResponse::InternalServerError().json(serde_json::json!({
                    "error": format!("Failed to parse performance bar header: {}", e)
                }))
            })?,
        );

        if let Some(request_id) = request_id {
            response.headers_mut().insert(
                "X-Performance-Bar-Request-ID",
                request_id.parse().map_err(|e| {
                    HttpResponse::InternalServerError().json(serde_json::json!({
                        "error": format!("Failed to parse request ID header: {}", e)
                    }))
                })?,
            );
        }

        if let Some(user_id) = user_id {
            response.headers_mut().insert(
                "X-Performance-Bar-User-ID",
                user_id.to_string().parse().map_err(|e| {
                    HttpResponse::InternalServerError().json(serde_json::json!({
                        "error": format!("Failed to parse user ID header: {}", e)
                    }))
                })?,
            );
        }

        // Add performance data as JSON
        response.headers_mut().insert(
            "X-Performance-Bar-Data",
            serde_json::to_string(&data)
                .map_err(|e| {
                    HttpResponse::InternalServerError().json(serde_json::json!({
                        "error": format!("Failed to serialize performance data: {}", e)
                    }))
                })?
                .parse()
                .map_err(|e| {
                    HttpResponse::InternalServerError().json(serde_json::json!({
                        "error": format!("Failed to parse performance data header: {}", e)
                    }))
                })?,
        );

        Ok(())
    }

    /// Measure execution time of a block
    fn measure_execution_time<F, T>(&self, f: F) -> (T, Duration)
    where
        F: FnOnce() -> T,
    {
        let start = Instant::now();
        let result = f();
        let duration = start.elapsed();
        (result, duration)
    }

    /// Execute a block with performance bar
    fn with_performance_bar<F, T>(&self, f: F) -> Result<T, HttpResponse>
    where
        F: FnOnce() -> T,
    {
        let (result, duration) = self.measure_execution_time(f);

        // Create a dummy response to add headers
        let mut response = HttpResponse::Ok();
        self.add_performance_bar_headers(&mut response, duration)?;

        Ok(result)
    }
}
