use actix_web::web;
use chrono::{DateTime, Duration, Utc};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

#[derive(Debug, Serialize, Deserialize)]
pub struct CycleAnalyticsProjectParams {
    pub start_date: Option<String>,
    pub created_after: Option<String>,
    pub created_before: Option<String>,
    pub branch_name: Option<String>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct CycleAnalyticsGroupParams {
    pub group_id: Option<i32>,
    pub start_date: Option<String>,
    pub created_after: Option<String>,
    pub created_before: Option<String>,
    pub project_ids: Option<Vec<i32>>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct CycleAnalyticsOptions {
    pub current_user_id: Option<i32>,
    pub projects: Option<Vec<i32>>,
    pub from: DateTime<Utc>,
    pub to: Option<DateTime<Utc>>,
    pub end_event_filter: Option<String>,
    pub use_aggregated_data_collector: Option<bool>,
}

pub trait CycleAnalyticsParams {
    fn cycle_analytics_project_params(
        &self,
        params: &web::Json<HashMap<String, serde_json::Value>>,
    ) -> CycleAnalyticsProjectParams {
        if let Some(cycle_analytics) = params.get("cycle_analytics") {
            serde_json::from_value(cycle_analytics.clone()).unwrap_or_default()
        } else {
            CycleAnalyticsProjectParams {
                start_date: None,
                created_after: None,
                created_before: None,
                branch_name: None,
            }
        }
    }

    fn cycle_analytics_group_params(
        &self,
        params: &web::Json<HashMap<String, serde_json::Value>>,
    ) -> CycleAnalyticsGroupParams {
        serde_json::from_value(params.0.clone()).unwrap_or_default()
    }

    fn options(
        &self,
        params: &web::Json<HashMap<String, serde_json::Value>>,
    ) -> CycleAnalyticsOptions {
        let mut opts = CycleAnalyticsOptions {
            current_user_id: self.current_user_id(),
            projects: params
                .get("project_ids")
                .and_then(|v| serde_json::from_value(v.clone()).ok()),
            from: params
                .get("from")
                .and_then(|v| serde_json::from_value(v.clone()).ok())
                .unwrap_or_else(|| self.start_date(params)),
            to: params
                .get("to")
                .and_then(|v| serde_json::from_value(v.clone()).ok()),
            end_event_filter: params
                .get("end_event_filter")
                .and_then(|v| serde_json::from_value(v.clone()).ok()),
            use_aggregated_data_collector: params
                .get("use_aggregated_data_collector")
                .and_then(|v| serde_json::from_value(v.clone()).ok()),
        };

        if let Some(date_range) = self.date_range(params) {
            opts.from = date_range.from;
            if let Some(to) = date_range.to {
                opts.to = Some(to);
            }
        }

        opts
    }

    fn start_date(&self, params: &web::Json<HashMap<String, serde_json::Value>>) -> DateTime<Utc> {
        match params.get("start_date").and_then(|v| v.as_str()) {
            Some("7") => Utc::now() - Duration::days(7),
            Some("30") => Utc::now() - Duration::days(30),
            _ => Utc::now() - Duration::days(90),
        }
    }

    fn date_range(
        &self,
        params: &web::Json<HashMap<String, serde_json::Value>>,
    ) -> Option<DateRange> {
        let created_after = params
            .get("created_after")
            .and_then(|v| v.as_str())
            .map(|s| self.to_utc_time(s).beginning_of_day());
        let created_before = params
            .get("created_before")
            .and_then(|v| v.as_str())
            .map(|s| self.to_utc_time(s).end_of_day());

        if created_after.is_some() || created_before.is_some() {
            Some(DateRange {
                from: created_after.unwrap_or_else(|| Utc::now() - Duration::days(90)),
                to: created_before,
            })
        } else {
            None
        }
    }

    fn to_utc_time(&self, field: &str) -> DateTime<Utc> {
        chrono::NaiveDate::parse_from_str(field, "%Y-%m-%d")
            .unwrap_or_else(|_| chrono::NaiveDate::from_ymd_opt(2000, 1, 1).unwrap())
            .and_hms_opt(0, 0, 0)
            .unwrap()
            .and_utc()
    }

    fn validate_params(
        &self,
        params: &web::Json<HashMap<String, serde_json::Value>>,
    ) -> Result<(), actix_web::Error> {
        if !self.is_valid_params(params) {
            return Err(actix_web::error::ErrorBadRequest("Invalid parameters"));
        }
        Ok(())
    }

    // Required methods to be implemented by concrete types
    fn current_user_id(&self) -> Option<i32>;
    fn is_valid_params(&self, params: &web::Json<HashMap<String, serde_json::Value>>) -> bool;
}

#[derive(Debug)]
pub struct DateRange {
    pub from: DateTime<Utc>,
    pub to: Option<DateTime<Utc>>,
}
