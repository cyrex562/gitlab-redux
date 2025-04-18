use crate::config::GitlabConfig;
use crate::models::deploy_token::DeployToken;
use crate::models::project::Project;
use crate::models::user::User;
use crate::utils::strong_memoize::StrongMemoize;
use actix_web::{web, HttpRequest, HttpResponse, Responder};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

pub const CONTENT_TYPE: &str = "application/vnd.git-lfs+json";

#[derive(Debug, Serialize, Deserialize)]
pub struct LfsObject {
    pub oid: String,
    #[serde(flatten)]
    pub extra: HashMap<String, String>,
}

pub trait LfsRequest {
    fn container(&self) -> Option<&dyn LfsContainer>;
    fn project(&self) -> &Project;
    fn user(&self) -> Option<&User>;
    fn deploy_token(&self) -> Option<&DeployToken>;
    fn authentication_result(&self) -> &AuthenticationResult;
    fn can(&self, object: &dyn std::any::Any, action: &str, subject: &dyn std::any::Any) -> bool;
    fn is_ci(&self) -> bool;
    fn is_download_request(&self) -> bool;
    fn is_upload_request(&self) -> bool;
    fn has_authentication_ability(&self, ability: &str) -> bool;
    fn objects(&self) -> Vec<LfsObject>;
    fn objects_oids(&self) -> Vec<String>;
    fn limit_exceeded(&self) -> bool;
}

pub trait LfsContainer {
    fn lfs_enabled(&self) -> bool;
}

pub struct AuthenticationResult {
    pub abilities: Vec<String>,
}

impl AuthenticationResult {
    pub fn lfs_deploy_token(&self, project: &Project) -> bool {
        // Implementation depends on your authentication logic
        false
    }
}

pub struct LfsRequestImpl {
    container: Option<Box<dyn LfsContainer>>,
    project: Project,
    user: Option<User>,
    deploy_token: Option<DeployToken>,
    authentication_result: AuthenticationResult,
    config: GitlabConfig,
    objects: Vec<LfsObject>,
}

impl LfsRequestImpl {
    pub fn new(
        container: Option<Box<dyn LfsContainer>>,
        project: Project,
        user: Option<User>,
        deploy_token: Option<DeployToken>,
        authentication_result: AuthenticationResult,
        config: GitlabConfig,
        objects: Vec<LfsObject>,
    ) -> Self {
        Self {
            container,
            project,
            user,
            deploy_token,
            authentication_result,
            config,
            objects,
        }
    }

    pub fn require_lfs_enabled(&self) -> impl Responder {
        if !self.config.lfs.enabled {
            return HttpResponse::NotImplemented()
                .content_type(CONTENT_TYPE)
                .json(serde_json::json!({
                    "message": "Git LFS is not enabled on this GitLab server, contact your admin.",
                    "documentation_url": self.help_url()
                }));
        }
        HttpResponse::Ok()
    }

    pub fn lfs_check_access(&self) -> impl Responder {
        if let Some(container) = &self.container {
            if !container.lfs_enabled() {
                return self.render_lfs_not_found();
            }
        }

        if self.is_download_request() && self.lfs_download_access() {
            return HttpResponse::Ok();
        }

        if self.is_upload_request() && self.lfs_upload_access() {
            return HttpResponse::Ok();
        }

        if self.lfs_download_access() {
            self.lfs_forbidden()
        } else {
            self.render_lfs_not_found()
        }
    }

    fn lfs_forbidden(&self) -> impl Responder {
        HttpResponse::Forbidden()
            .content_type(CONTENT_TYPE)
            .json(serde_json::json!({
                "message": "Access forbidden. Check your access level.",
                "documentation_url": self.help_url()
            }))
    }

    fn render_lfs_not_found(&self) -> impl Responder {
        HttpResponse::NotFound()
            .content_type(CONTENT_TYPE)
            .json(serde_json::json!({
                "message": "Not found.",
                "documentation_url": self.help_url()
            }))
    }

    fn help_url(&self) -> String {
        // Implementation depends on your help URL structure
        "https://docs.gitlab.com/ee/topics/git/lfs/".to_string()
    }
}

impl LfsRequest for LfsRequestImpl {
    fn container(&self) -> Option<&dyn LfsContainer> {
        self.container.as_ref().map(|c| c.as_ref())
    }

    fn project(&self) -> &Project {
        &self.project
    }

    fn user(&self) -> Option<&User> {
        self.user.as_ref()
    }

    fn deploy_token(&self) -> Option<&DeployToken> {
        self.deploy_token.as_ref()
    }

    fn authentication_result(&self) -> &AuthenticationResult {
        &self.authentication_result
    }

    fn can(&self, object: &dyn std::any::Any, action: &str, subject: &dyn std::any::Any) -> bool {
        // Implementation depends on your permission system
        false
    }

    fn is_ci(&self) -> bool {
        // Implementation depends on your CI detection logic
        false
    }

    fn is_download_request(&self) -> bool {
        // Implementation depends on your request type detection
        false
    }

    fn is_upload_request(&self) -> bool {
        // Implementation depends on your request type detection
        false
    }

    fn has_authentication_ability(&self, ability: &str) -> bool {
        self.authentication_result
            .abilities
            .contains(&ability.to_string())
    }

    fn objects(&self) -> Vec<LfsObject> {
        self.objects.clone()
    }

    fn objects_oids(&self) -> Vec<String> {
        self.objects.iter().map(|o| o.oid.clone()).collect()
    }

    fn limit_exceeded(&self) -> bool {
        // Default implementation, can be overridden in EE
        false
    }
}

impl LfsRequestImpl {
    fn lfs_download_access(&self) -> bool {
        self.is_ci()
            || self.lfs_deploy_token()
            || self.user_can_download_code()
            || self.build_can_download_code()
            || self.deploy_token_can_download_code()
    }

    fn deploy_token_can_download_code(&self) -> bool {
        if let Some(token) = &self.deploy_token {
            token.has_access_to(&self.project) && token.read_repository()
        } else {
            false
        }
    }

    fn lfs_upload_access(&self) -> bool {
        if !self.has_authentication_ability("push_code") {
            return false;
        }
        if self.limit_exceeded() {
            return false;
        }

        self.lfs_deploy_token()
            || self.can(
                self.user
                    .as_ref()
                    .map(|u| u as &dyn std::any::Any)
                    .unwrap_or(&()),
                "push_code",
                &self.project,
            )
            || self.can(
                self.deploy_token
                    .as_ref()
                    .map(|t| t as &dyn std::any::Any)
                    .unwrap_or(&()),
                "push_code",
                &self.project,
            )
            || self.any_branch_allows_collaboration()
    }

    fn any_branch_allows_collaboration(&self) -> bool {
        if let Some(user) = &self.user {
            !self
                .project
                .merge_requests_allowing_push_to_user(user)
                .is_empty()
        } else {
            false
        }
    }

    fn lfs_deploy_token(&self) -> bool {
        self.authentication_result.lfs_deploy_token(&self.project)
    }

    fn user_can_download_code(&self) -> bool {
        self.has_authentication_ability("download_code")
            && self.can(
                self.user
                    .as_ref()
                    .map(|u| u as &dyn std::any::Any)
                    .unwrap_or(&()),
                "download_code",
                &self.project,
            )
    }

    fn build_can_download_code(&self) -> bool {
        self.has_authentication_ability("build_download_code")
            && self.can(
                self.user
                    .as_ref()
                    .map(|u| u as &dyn std::any::Any)
                    .unwrap_or(&()),
                "build_download_code",
                &self.project,
            )
    }
}
