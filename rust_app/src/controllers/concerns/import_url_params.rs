use actix_web::web;
use serde::{Deserialize, Serialize};
use url::Url;

#[derive(Debug, Serialize, Deserialize)]
pub struct ImportUrlCredentials {
    pub user: Option<String>,
    pub password: Option<String>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ImportUrlParams {
    pub import_url: String,
    pub import_type: String,
}

pub trait ImportUrlParams {
    fn import_url_params(&self, params: &web::Json<ProjectParams>) -> ImportUrlParams {
        if params.import_url.is_none() {
            return ImportUrlParams {
                import_url: String::new(),
                import_type: String::new(),
            };
        }

        ImportUrlParams {
            import_url: self.import_params_to_full_url(params),
            import_type: "git".to_string(),
        }
    }

    fn import_params_to_full_url(&self, params: &ProjectParams) -> String {
        let mut url = Url::parse(&params.import_url.as_ref().unwrap())
            .unwrap_or_else(|_| Url::parse("http://localhost").unwrap());

        if let Some(user) = &params.import_url_user {
            url.set_username(user).ok();
        }

        if let Some(password) = &params.import_url_password {
            url.set_password(Some(password)).ok();
        }

        url.to_string()
    }
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ProjectParams {
    pub import_url: Option<String>,
    pub import_url_user: Option<String>,
    pub import_url_password: Option<String>,
}
