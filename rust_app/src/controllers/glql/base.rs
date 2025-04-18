use crate::controllers::graphql::GraphqlController;
use crate::error::Error;
use crate::metrics::{ApplicationRateLimiter, GlqlSlis, System};
use actix_web::{web, HttpResponse};
use sha2::{Digest, Sha256};
use std::time::Instant;

#[derive(Debug)]
pub struct GlqlQueryLockedError {
    message: String,
}

impl std::fmt::Display for GlqlQueryLockedError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{}", self.message)
    }
}

impl std::error::Error for GlqlQueryLockedError {}

pub struct BaseController {
    base: GraphqlController,
}

impl BaseController {
    pub fn new(base: GraphqlController) -> Self {
        Self { base }
    }

    pub async fn execute(&self, params: web::Json<GraphqlParams>) -> Result<HttpResponse, Error> {
        self.check_rate_limit(&params).await?;

        let start_time = System::monotonic_time();
        let result = self.base.execute(params.clone()).await;

        match result {
            Ok(response) => {
                self.increment_glql_sli(System::monotonic_time() - start_time, None)
                    .await;
                Ok(response)
            }
            Err(error) => {
                if error.is_query_aborted() {
                    self.increment_rate_limit_counter(&params).await;
                }

                self.increment_glql_sli(
                    System::monotonic_time() - start_time,
                    Some(self.error_type_from(&error)),
                )
                .await;

                Err(error)
            }
        }
    }

    async fn check_rate_limit(&self, params: &GraphqlParams) -> Result<(), Error> {
        if ApplicationRateLimiter::peek("glql", &self.query_sha(params)).await {
            Err(Error::from(GlqlQueryLockedError {
                message: "Query execution is locked due to repeated failures.".to_string(),
            }))
        } else {
            Ok(())
        }
    }

    async fn increment_rate_limit_counter(&self, params: &GraphqlParams) {
        ApplicationRateLimiter::throttled("glql", &self.query_sha(params)).await;
    }

    fn query_sha(&self, params: &GraphqlParams) -> String {
        let mut hasher = Sha256::new();
        hasher.update(params.query.as_bytes());
        format!("{:x}", hasher.finalize())
    }

    async fn increment_glql_sli(&self, duration: f64, error_type: Option<ErrorType>) {
        let query_urgency = EndpointAttributes::Config::REQUEST_URGENCIES
            .get("low")
            .unwrap();

        let labels = GlqlLabels {
            endpoint_id: ApplicationContext::current_context_attribute("caller_id"),
            feature_category: ApplicationContext::current_context_attribute("feature_category"),
            query_urgency: query_urgency.name.clone(),
        };

        GlqlSlis::record_error(
            labels.clone().merge(ErrorLabels { error_type }),
            error_type.is_some(),
        )
        .await;

        if error_type.is_none() {
            GlqlSlis::record_apdex(labels, duration <= query_urgency.duration).await;
        }
    }

    fn error_type_from(&self, error: &Error) -> ErrorType {
        match error {
            Error::QueryAborted(_) => ErrorType::QueryAborted,
            _ => ErrorType::Other,
        }
    }
}

#[derive(serde::Deserialize)]
pub struct GraphqlParams {
    query: String,
}

#[derive(Clone)]
struct GlqlLabels {
    endpoint_id: Option<String>,
    feature_category: Option<String>,
    query_urgency: String,
}

impl GlqlLabels {
    fn merge(self, error_labels: ErrorLabels) -> Self {
        Self {
            endpoint_id: self.endpoint_id,
            feature_category: self.feature_category,
            query_urgency: self.query_urgency,
            error_type: error_labels.error_type,
        }
    }
}

#[derive(Clone)]
struct ErrorLabels {
    error_type: ErrorType,
}

#[derive(Clone, Debug)]
enum ErrorType {
    QueryAborted,
    Other,
}
