//! Jira Connect App Descriptor endpoint
//! Ported from Ruby: app/controllers/jira_connect/app_descriptor_controller.rb

use actix_web::{get, web, HttpResponse, Responder};
use serde::Serialize;

const HOME_URL: &str = "https://gitlab.com";
const DOC_URL: &str = "https://docs.gitlab.com/ee/integration/jira/";

#[derive(Serialize)]
struct AppDescriptor {
    name: String,
    description: String,
    key: String,
    #[serde(rename = "baseUrl")]
    base_url: String,
    lifecycle: Lifecycle,
    vendor: Vendor,
    links: Links,
    authentication: Authentication,
    modules: Modules,
    scopes: Vec<&'static str>,
    #[serde(rename = "apiVersion")]
    api_version: u8,
    #[serde(rename = "apiMigrations")]
    api_migrations: ApiMigrations,
}

#[derive(Serialize)]
struct Lifecycle {
    installed: String,
    uninstalled: String,
}

#[derive(Serialize)]
struct Vendor {
    name: &'static str,
    url: &'static str,
}

#[derive(Serialize)]
struct Links {
    documentation: &'static str,
}

#[derive(Serialize)]
struct Authentication {
    r#type: &'static str,
}

#[derive(Serialize)]
struct ApiMigrations {
    #[serde(rename = "context-qsh")]
    context_qsh: bool,
    #[serde(rename = "signed-install")]
    signed_install: bool,
    gdpr: bool,
}

// Modules struct can be expanded as needed
#[derive(Serialize)]
struct Modules {
    #[serde(rename = "postInstallPage")]
    post_install_page: PostInstallPage,
    #[serde(rename = "jiraDevelopmentTool")]
    jira_development_tool: JiraDevelopmentTool,
    #[serde(rename = "jiraBuildInfoProvider")]
    jira_build_info_provider: CommonModule,
    #[serde(rename = "jiraDeploymentInfoProvider")]
    jira_deployment_info_provider: CommonModule,
    #[serde(rename = "jiraFeatureFlagInfoProvider")]
    jira_feature_flag_info_provider: CommonModule,
}

#[derive(Serialize)]
struct PostInstallPage {
    key: &'static str,
    name: ValueField,
    url: String,
    conditions: Vec<Condition>,
}

#[derive(Serialize)]
struct ValueField {
    value: &'static str,
}

#[derive(Serialize)]
struct Condition {
    condition: &'static str,
    invert: bool,
}

#[derive(Serialize)]
struct JiraDevelopmentTool {
    actions: Actions,
    key: &'static str,
    application: ValueField,
    name: ValueField,
    url: &'static str,
    #[serde(rename = "logoUrl")]
    logo_url: &'static str,
    capabilities: Vec<&'static str>,
}

#[derive(Serialize)]
struct Actions {
    createBranch: TemplateUrl,
    searchConnectedWorkspaces: TemplateUrl,
    searchRepositories: TemplateUrl,
    associateRepository: TemplateUrl,
}

#[derive(Serialize)]
struct TemplateUrl {
    #[serde(rename = "templateUrl")]
    template_url: String,
}

#[derive(Serialize)]
struct CommonModule {
    actions: ActionsEmpty,
    name: ValueField,
    key: &'static str,
    #[serde(flatten)]
    common: CommonModuleProperties,
}

#[derive(Serialize)]
struct ActionsEmpty {}

#[derive(Serialize)]
struct CommonModuleProperties {
    #[serde(rename = "homeUrl")]
    home_url: &'static str,
    #[serde(rename = "logoUrl")]
    logo_url: &'static str,
    #[serde(rename = "documentationUrl")]
    documentation_url: &'static str,
}

#[get("/jira_connect/app_descriptor.json")]
pub async fn app_descriptor() -> impl Responder {
    let descriptor = AppDescriptor {
        name: "GitLab for Jira (gitlab.com)".to_string(),
        description: "Integrate commits, branches and merge requests from GitLab into Jira".to_string(),
        key: "gitlab-jira-connect-gitlab.com".to_string(),
        base_url: format!("{}/-/jira_connect", HOME_URL),
        lifecycle: Lifecycle {
            installed: "/events/installed".to_string(),
            uninstalled: "/events/uninstalled".to_string(),
        },
        vendor: Vendor {
            name: "GitLab",
            url: "https://gitlab.com",
        },
        links: Links {
            documentation: DOC_URL,
        },
        authentication: Authentication { r#type: "jwt" },
        modules: Modules {
            post_install_page: PostInstallPage {
                key: "gitlab-configuration",
                name: ValueField { value: "GitLab Configuration" },
                url: "/subscriptions".to_string(),
                conditions: vec![Condition { condition: "user_is_admin", invert: false }],
            },
            jira_development_tool: JiraDevelopmentTool {
                actions: Actions {
                    createBranch: TemplateUrl { template_url: "/branches/route?issue_key={issue.key}&issue_summary={issue.summary}&jwt={jwt}&addonkey=gitlab-jira-connect-gitlab.com".to_string() },
                    searchConnectedWorkspaces: TemplateUrl { template_url: "/workspaces/search".to_string() },
                    searchRepositories: TemplateUrl { template_url: "/repositories/search".to_string() },
                    associateRepository: TemplateUrl { template_url: "/repositories/associate".to_string() },
                },
                key: "gitlab-development-tool",
                application: ValueField { value: "GitLab (gitlab.com)" },
                name: ValueField { value: "GitLab (gitlab.com)" },
                url: HOME_URL,
                logo_url: "https://gitlab.com/assets/gitlab_logo.png",
                capabilities: vec!["branch", "commit", "pull_request"],
            },
            jira_build_info_provider: CommonModule {
                actions: ActionsEmpty {},
                name: ValueField { value: "GitLab CI" },
                key: "gitlab-ci",
                common: CommonModuleProperties {
                    home_url: HOME_URL,
                    logo_url: "https://gitlab.com/assets/gitlab_logo.png",
                    documentation_url: DOC_URL,
                },
            },
            jira_deployment_info_provider: CommonModule {
                actions: ActionsEmpty {},
                name: ValueField { value: "GitLab Deployments" },
                key: "gitlab-deployments",
                common: CommonModuleProperties {
                    home_url: HOME_URL,
                    logo_url: "https://gitlab.com/assets/gitlab_logo.png",
                    documentation_url: DOC_URL,
                },
            },
            jira_feature_flag_info_provider: CommonModule {
                actions: ActionsEmpty {},
                name: ValueField { value: "GitLab Feature Flags" },
                key: "gitlab-feature-flags",
                common: CommonModuleProperties {
                    home_url: HOME_URL,
                    logo_url: "https://gitlab.com/assets/gitlab_logo.png",
                    documentation_url: DOC_URL,
                },
            },
        },
        scopes: vec!["READ", "WRITE", "DELETE"],
        api_version: 1,
        api_migrations: ApiMigrations {
            context_qsh: true,
            signed_install: true,
            gdpr: true,
        },
    };
    HttpResponse::Ok().json(descriptor)
}
