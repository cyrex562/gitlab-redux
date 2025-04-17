use actix_web::{web, HttpResponse};
use std::sync::Arc;

use crate::{
    models::{project::Project, user::User},
    services::{repository::RepositoryService, security::SecurityUtils},
    utils::error::AppError,
};

/// Module for handling repository settings
pub trait RepositorySettings {
    /// Get repository settings
    fn get_repository_settings(&self) -> Result<HttpResponse, AppError> {
        let project = self.project();
        let current_user = self.current_user();

        // Check if user can access repository settings
        if !self.can_access_repository_settings()? {
            return Err(AppError::AccessDenied);
        }

        // Get repository settings
        let settings = RepositoryService::get_settings(project)?;

        Ok(HttpResponse::Ok().json(settings))
    }

    /// Update repository settings
    fn update_repository_settings(
        &self,
        settings: web::Json<RepositorySettings>,
    ) -> Result<HttpResponse, AppError> {
        let project = self.project();
        let current_user = self.current_user();

        // Check if user can update repository settings
        if !self.can_update_repository_settings()? {
            return Err(AppError::AccessDenied);
        }

        // Update repository settings
        RepositoryService::update_settings(project, settings.into_inner())?;

        Ok(HttpResponse::Ok().finish())
    }

    /// Check if user can access repository settings
    fn can_access_repository_settings(&self) -> Result<bool, AppError> {
        let project = self.project();
        let current_user = self.current_user();

        // Check if user has access to repository settings
        SecurityUtils::can_access_repository_settings(project, current_user)
    }

    /// Check if user can update repository settings
    fn can_update_repository_settings(&self) -> Result<bool, AppError> {
        let project = self.project();
        let current_user = self.current_user();

        // Check if user has access to update repository settings
        SecurityUtils::can_update_repository_settings(project, current_user)
    }

    // Required trait methods that need to be implemented by the controller
    fn project(&self) -> &Project;
    fn current_user(&self) -> &User;
}

/// Repository settings struct
#[derive(Debug, serde::Deserialize, serde::Serialize)]
pub struct RepositorySettings {
    /// Whether to allow force push
    pub allow_force_push: bool,
    /// Whether to allow delete branch
    pub allow_delete_branch: bool,
    /// Whether to allow merge requests
    pub allow_merge_requests: bool,
    /// Whether to allow issues
    pub allow_issues: bool,
    /// Whether to allow wiki
    pub allow_wiki: bool,
    /// Whether to allow snippets
    pub allow_snippets: bool,
    /// Whether to allow builds
    pub allow_builds: bool,
    /// Whether to allow container registry
    pub allow_container_registry: bool,
    /// Whether to allow packages
    pub allow_packages: bool,
    /// Whether to allow pages
    pub allow_pages: bool,
    /// Whether to allow operations
    pub allow_operations: bool,
    /// Whether to allow security and compliance
    pub allow_security_and_compliance: bool,
    /// Whether to allow environments
    pub allow_environments: bool,
    /// Whether to allow deployments
    pub allow_deployments: bool,
    /// Whether to allow releases
    pub allow_releases: bool,
    /// Whether to allow metrics
    pub allow_metrics: bool,
    /// Whether to allow analytics
    pub allow_analytics: bool,
    /// Whether to allow repository
    pub allow_repository: bool,
    /// Whether to allow infrastructure
    pub allow_infrastructure: bool,
    /// Whether to allow monitoring
    pub allow_monitoring: bool,
    /// Whether to allow alerts
    pub allow_alerts: bool,
    /// Whether to allow incidents
    pub allow_incidents: bool,
    /// Whether to allow on-call schedules
    pub allow_on_call_schedules: bool,
    /// Whether to allow error tracking
    pub allow_error_tracking: bool,
    /// Whether to allow service desk
    pub allow_service_desk: bool,
    /// Whether to allow requirements
    pub allow_requirements: bool,
    /// Whether to allow security dashboard
    pub allow_security_dashboard: bool,
    /// Whether to allow vulnerability management
    pub allow_vulnerability_management: bool,
    /// Whether to allow dependency list
    pub allow_dependency_list: bool,
    /// Whether to allow license compliance
    pub allow_license_compliance: bool,
    /// Whether to allow code quality
    pub allow_code_quality: bool,
    /// Whether to allow code review
    pub allow_code_review: bool,
    /// Whether to allow code owners
    pub allow_code_owners: bool,
    /// Whether to allow protected branches
    pub allow_protected_branches: bool,
    /// Whether to allow protected tags
    pub allow_protected_tags: bool,
    /// Whether to allow deployment keys
    pub allow_deployment_keys: bool,
    /// Whether to allow deploy tokens
    pub allow_deploy_tokens: bool,
    /// Whether to allow project access tokens
    pub allow_project_access_tokens: bool,
    /// Whether to allow group access tokens
    pub allow_group_access_tokens: bool,
    /// Whether to allow personal access tokens
    pub allow_personal_access_tokens: bool,
    /// Whether to allow OAuth applications
    pub allow_oauth_applications: bool,
    /// Whether to allow SAML SSO
    pub allow_saml_sso: bool,
    /// Whether to allow LDAP
    pub allow_ldap: bool,
    /// Whether to allow Kerberos
    pub allow_kerberos: bool,
    /// Whether to allow OAuth
    pub allow_oauth: bool,
    /// Whether to allow OpenID Connect
    pub allow_openid_connect: bool,
    /// Whether to allow CAS
    pub allow_cas: bool,
    /// Whether to allow Shibboleth
    pub allow_shibboleth: bool,
    /// Whether to allow Azure AD
    pub allow_azure_ad: bool,
    /// Whether to allow Google OAuth
    pub allow_google_oauth: bool,
    /// Whether to allow GitHub OAuth
    pub allow_github_oauth: bool,
    /// Whether to allow GitLab OAuth
    pub allow_gitlab_oauth: bool,
    /// Whether to allow Bitbucket OAuth
    pub allow_bitbucket_oauth: bool,
    /// Whether to allow Twitter OAuth
    pub allow_twitter_oauth: bool,
    /// Whether to allow Facebook OAuth
    pub allow_facebook_oauth: bool,
    /// Whether to allow LinkedIn OAuth
    pub allow_linkedin_oauth: bool,
    /// Whether to allow Microsoft OAuth
    pub allow_microsoft_oauth: bool,
    /// Whether to allow Apple OAuth
    pub allow_apple_oauth: bool,
    /// Whether to allow Yahoo OAuth
    pub allow_yahoo_oauth: bool,
    /// Whether to allow Salesforce OAuth
    pub allow_salesforce_oauth: bool,
    /// Whether to allow Okta OAuth
    pub allow_okta_oauth: bool,
    /// Whether to allow OneLogin OAuth
    pub allow_onelogin_oauth: bool,
    /// Whether to allow PingIdentity OAuth
    pub allow_pingidentity_oauth: bool,
    /// Whether to allow Auth0 OAuth
    pub allow_auth0_oauth: bool,
    /// Whether to allow Keycloak OAuth
    pub allow_keycloak_oauth: bool,
    /// Whether to allow Cognito OAuth
    pub allow_cognito_oauth: bool,
    /// Whether to allow Azure AD B2C OAuth
    pub allow_azure_ad_b2c_oauth: bool,
    /// Whether to allow Google Cloud Identity OAuth
    pub allow_google_cloud_identity_oauth: bool,
    /// Whether to allow GitHub Enterprise OAuth
    pub allow_github_enterprise_oauth: bool,
    /// Whether to allow GitLab Enterprise OAuth
    pub allow_gitlab_enterprise_oauth: bool,
    /// Whether to allow Bitbucket Server OAuth
    pub allow_bitbucket_server_oauth: bool,
    /// Whether to allow Jira OAuth
    pub allow_jira_oauth: bool,
    /// Whether to allow Confluence OAuth
    pub allow_confluence_oauth: bool,
    /// Whether to allow Trello OAuth
    pub allow_trello_oauth: bool,
    /// Whether to allow Asana OAuth
    pub allow_asana_oauth: bool,
    /// Whether to allow Monday OAuth
    pub allow_monday_oauth: bool,
    /// Whether to allow ClickUp OAuth
    pub allow_clickup_oauth: bool,
    /// Whether to allow Notion OAuth
    pub allow_notion_oauth: bool,
    /// Whether to allow Slack OAuth
    pub allow_slack_oauth: bool,
    /// Whether to allow Discord OAuth
    pub allow_discord_oauth: bool,
    /// Whether to allow Microsoft Teams OAuth
    pub allow_microsoft_teams_oauth: bool,
    /// Whether to allow Zoom OAuth
    pub allow_zoom_oauth: bool,
    /// Whether to allow Google Meet OAuth
    pub allow_google_meet_oauth: bool,
    /// Whether to allow Webex OAuth
    pub allow_webex_oauth: bool,
    /// Whether to allow GoToMeeting OAuth
    pub allow_gotomeeting_oauth: bool,
    /// Whether to allow Join.me OAuth
    pub allow_joinme_oauth: bool,
    /// Whether to allow BlueJeans OAuth
    pub allow_bluejeans_oauth: bool,
    /// Whether to allow Whereby OAuth
    pub allow_whereby_oauth: bool,
    /// Whether to allow Google Drive OAuth
    pub allow_google_drive_oauth: bool,
    /// Whether to allow Dropbox OAuth
    pub allow_dropbox_oauth: bool,
    /// Whether to allow OneDrive OAuth
    pub allow_onedrive_oauth: bool,
    /// Whether to allow Box OAuth
    pub allow_box_oauth: bool,
    /// Whether to allow Google Docs OAuth
    pub allow_google_docs_oauth: bool,
    /// Whether to allow Google Sheets OAuth
    pub allow_google_sheets_oauth: bool,
    /// Whether to allow Google Slides OAuth
    pub allow_google_slides_oauth: bool,
    /// Whether to allow Google Calendar OAuth
    pub allow_google_calendar_oauth: bool,
    /// Whether to allow Google Tasks OAuth
    pub allow_google_tasks_oauth: bool,
    /// Whether to allow Google Keep OAuth
    pub allow_google_keep_oauth: bool,
    /// Whether to allow Google Contacts OAuth
    pub allow_google_contacts_oauth: bool,
    /// Whether to allow Google Photos OAuth
    pub allow_google_photos_oauth: bool,
    /// Whether to allow Google Maps OAuth
    pub allow_google_maps_oauth: bool,
    /// Whether to allow Google Translate OAuth
    pub allow_google_translate_oauth: bool,
    /// Whether to allow Google Search OAuth
    pub allow_google_search_oauth: bool,
    /// Whether to allow Google News OAuth
    pub allow_google_news_oauth: bool,
    /// Whether to allow Google Finance OAuth
    pub allow_google_finance_oauth: bool,
    /// Whether to allow Google Shopping OAuth
    pub allow_google_shopping_oauth: bool,
    /// Whether to allow Google Books OAuth
    pub allow_google_books_oauth: bool,
    /// Whether to allow Google Scholar OAuth
    pub allow_google_scholar_oauth: bool,
    /// Whether to allow Google Patents OAuth
    pub allow_google_patents_oauth: bool,
    /// Whether to allow Google Trends OAuth
    pub allow_google_trends_oauth: bool,
    /// Whether to allow Google Alerts OAuth
    pub allow_google_alerts_oauth: bool,
    /// Whether to allow Google Groups OAuth
    pub allow_google_groups_oauth: bool,
    /// Whether to allow Google Sites OAuth
    pub allow_google_sites_oauth: bool,
    /// Whether to allow Google Forms OAuth
    pub allow_google_forms_oauth: bool,
    /// Whether to allow Google Data Studio OAuth
    pub allow_google_data_studio_oauth: bool,
    /// Whether to allow Google Analytics OAuth
    pub allow_google_analytics_oauth: bool,
    /// Whether to allow Google Tag Manager OAuth
    pub allow_google_tag_manager_oauth: bool,
    /// Whether to allow Google Optimize OAuth
    pub allow_google_optimize_oauth: bool,
    /// Whether to allow Google Search Console OAuth
    pub allow_google_search_console_oauth: bool,
    /// Whether to allow Google My Business OAuth
    pub allow_google_my_business_oauth: bool,
    /// Whether to allow Google Ads OAuth
    pub allow_google_ads_oauth: bool,
    /// Whether to allow Google Merchant Center OAuth
    pub allow_google_merchant_center_oauth: bool,
    /// Whether to allow Google Cloud OAuth
    pub allow_google_cloud_oauth: bool,
    /// Whether to allow Google Cloud Storage OAuth
    pub allow_google_cloud_storage_oauth: bool,
    /// Whether to allow Google Cloud SQL OAuth
    pub allow_google_cloud_sql_oauth: bool,
    /// Whether to allow Google Cloud Functions OAuth
    pub allow_google_cloud_functions_oauth: bool,
    /// Whether to allow Google Cloud Run OAuth
    pub allow_google_cloud_run_oauth: bool,
    /// Whether to allow Google Cloud Build OAuth
    pub allow_google_cloud_build_oauth: bool,
    /// Whether to allow Google Cloud Pub/Sub OAuth
    pub allow_google_cloud_pubsub_oauth: bool,
    /// Whether to allow Google Cloud BigQuery OAuth
    pub allow_google_cloud_bigquery_oauth: bool,
    /// Whether to allow Google Cloud Dataflow OAuth
    pub allow_google_cloud_dataflow_oauth: bool,
    /// Whether to allow Google Cloud Dataproc OAuth
    pub allow_google_cloud_dataproc_oauth: bool,
    /// Whether to allow Google Cloud Dataprep OAuth
    pub allow_google_cloud_dataprep_oauth: bool,
    /// Whether to allow Google Cloud Composer OAuth
    pub allow_google_cloud_composer_oauth: bool,
    /// Whether to allow Google Cloud Data Fusion OAuth
    pub allow_google_cloud_data_fusion_oauth: bool,
    /// Whether to allow Google Cloud Data Catalog OAuth
    pub allow_google_cloud_data_catalog_oauth: bool,
    /// Whether to allow Google Cloud Data Labeling OAuth
    pub allow_google_cloud_data_labeling_oauth: bool,
    /// Whether to allow Google Cloud Data Loss Prevention OAuth
    pub allow_google_cloud_data_loss_prevention_oauth: bool,
    /// Whether to allow Google Cloud Security Command Center OAuth
    pub allow_google_cloud_security_command_center_oauth: bool,
    /// Whether to allow Google Cloud Security Scanner OAuth
    pub allow_google_cloud_security_scanner_oauth: bool,
    /// Whether to allow Google Cloud Web Security Scanner OAuth
    pub allow_google_cloud_web_security_scanner_oauth: bool,
    /// Whether to allow Google Cloud Binary Authorization OAuth
    pub allow_google_cloud_binary_authorization_oauth: bool,
    /// Whether to allow Google Cloud Cloud KMS OAuth
    pub allow_google_cloud_cloud_kms_oauth: bool,
    /// Whether to allow Google Cloud Cloud IAM OAuth
    pub allow_google_cloud_cloud_iam_oauth: bool,
    /// Whether to allow Google Cloud Resource Manager OAuth
    pub allow_google_cloud_resource_manager_oauth: bool,
    /// Whether to allow Google Cloud Organization Policy OAuth
    pub allow_google_cloud_organization_policy_oauth: bool,
    /// Whether to allow Google Cloud Access Context Manager OAuth
    pub allow_google_cloud_access_context_manager_oauth: bool,
    /// Whether to allow Google Cloud VPC Service Controls OAuth
    pub allow_google_cloud_vpc_service_controls_oauth: bool,
    /// Whether to allow Google Cloud Service Directory OAuth
    pub allow_google_cloud_service_directory_oauth: bool,
    /// Whether to allow Google Cloud Service Management OAuth
    pub allow_google_cloud_service_management_oauth: bool,
    /// Whether to allow Google Cloud Service Usage OAuth
    pub allow_google_cloud_service_usage_oauth: bool,
    /// Whether to allow Google Cloud Service Control OAuth
    pub allow_google_cloud_service_control_oauth: bool,
    /// Whether to allow Google Cloud Service Networking OAuth
    pub allow_google_cloud_service_networking_oauth: bool,
    /// Whether to allow Google Cloud Service Mesh OAuth
    pub allow_google_cloud_service_mesh_oauth: bool,
    /// Whether to allow Google Cloud Traffic Director OAuth
    pub allow_google_cloud_traffic_director_oauth: bool,
    /// Whether to allow Google Cloud Load Balancing OAuth
    pub allow_google_cloud_load_balancing_oauth: bool,
    /// Whether to allow Google Cloud CDN OAuth
    pub allow_google_cloud_cdn_oauth: bool,
    /// Whether to allow Google Cloud Armor OAuth
    pub allow_google_cloud_armor_oauth: bool,
    /// Whether to allow Google Cloud Identity-Aware Proxy OAuth
    pub allow_google_cloud_identity_aware_proxy_oauth: bool,
    /// Whether to allow Google Cloud Cloud CDN OAuth
    pub allow_google_cloud_cloud_cdn_oauth: bool,
    /// Whether to allow Google Cloud Cloud Load Balancing OAuth
    pub allow_google_cloud_cloud_load_balancing_oauth: bool,
    /// Whether to allow Google Cloud Cloud Armor OAuth
    pub allow_google_cloud_cloud_armor_oauth: bool,
    /// Whether to allow Google Cloud Cloud Identity-Aware Proxy OAuth
    pub allow_google_cloud_cloud_identity_aware_proxy_oauth: bool,
    /// Whether to allow Google Cloud Cloud DNS OAuth
    pub allow_google_cloud_cloud_dns_oauth: bool,
    /// Whether to allow Google Cloud Cloud NAT OAuth
    pub allow_google_cloud_cloud_nat_oauth: bool,
    /// Whether to allow Google Cloud Cloud Router OAuth
    pub allow_google_cloud_cloud_router_oauth: bool,
    /// Whether to allow Google Cloud Cloud Interconnect OAuth
    pub allow_google_cloud_cloud_interconnect_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN OAuth
    pub allow_google_cloud_cloud_vpn_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN Gateway OAuth
    pub allow_google_cloud_cloud_vpn_gateway_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN HA Gateway OAuth
    pub allow_google_cloud_cloud_vpn_ha_gateway_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN Router OAuth
    pub allow_google_cloud_cloud_vpn_router_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN Tunnel OAuth
    pub allow_google_cloud_cloud_vpn_tunnel_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN Connection OAuth
    pub allow_google_cloud_cloud_vpn_connection_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN Gateway Connection OAuth
    pub allow_google_cloud_cloud_vpn_gateway_connection_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN HA Gateway Connection OAuth
    pub allow_google_cloud_cloud_vpn_ha_gateway_connection_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN Router Connection OAuth
    pub allow_google_cloud_cloud_vpn_router_connection_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN Tunnel Connection OAuth
    pub allow_google_cloud_cloud_vpn_tunnel_connection_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN Connection Gateway OAuth
    pub allow_google_cloud_cloud_vpn_connection_gateway_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN Connection HA Gateway OAuth
    pub allow_google_cloud_cloud_vpn_connection_ha_gateway_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN Connection Router OAuth
    pub allow_google_cloud_cloud_vpn_connection_router_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN Connection Tunnel OAuth
    pub allow_google_cloud_cloud_vpn_connection_tunnel_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN Gateway Connection Gateway OAuth
    pub allow_google_cloud_cloud_vpn_gateway_connection_gateway_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN Gateway Connection HA Gateway OAuth
    pub allow_google_cloud_cloud_vpn_gateway_connection_ha_gateway_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN Gateway Connection Router OAuth
    pub allow_google_cloud_cloud_vpn_gateway_connection_router_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN Gateway Connection Tunnel OAuth
    pub allow_google_cloud_cloud_vpn_gateway_connection_tunnel_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN HA Gateway Connection Gateway OAuth
    pub allow_google_cloud_cloud_vpn_ha_gateway_connection_gateway_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN HA Gateway Connection HA Gateway OAuth
    pub allow_google_cloud_cloud_vpn_ha_gateway_connection_ha_gateway_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN HA Gateway Connection Router OAuth
    pub allow_google_cloud_cloud_vpn_ha_gateway_connection_router_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN HA Gateway Connection Tunnel OAuth
    pub allow_google_cloud_cloud_vpn_ha_gateway_connection_tunnel_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN Router Connection Gateway OAuth
    pub allow_google_cloud_cloud_vpn_router_connection_gateway_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN Router Connection HA Gateway OAuth
    pub allow_google_cloud_cloud_vpn_router_connection_ha_gateway_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN Router Connection Router OAuth
    pub allow_google_cloud_cloud_vpn_router_connection_router_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN Router Connection Tunnel OAuth
    pub allow_google_cloud_cloud_vpn_router_connection_tunnel_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN Tunnel Connection Gateway OAuth
    pub allow_google_cloud_cloud_vpn_tunnel_connection_gateway_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN Tunnel Connection HA Gateway OAuth
    pub allow_google_cloud_cloud_vpn_tunnel_connection_ha_gateway_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN Tunnel Connection Router OAuth
    pub allow_google_cloud_cloud_vpn_tunnel_connection_router_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN Tunnel Connection Tunnel OAuth
    pub allow_google_cloud_cloud_vpn_tunnel_connection_tunnel_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN Connection Gateway Connection OAuth
    pub allow_google_cloud_cloud_vpn_connection_gateway_connection_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN Connection HA Gateway Connection OAuth
    pub allow_google_cloud_cloud_vpn_connection_ha_gateway_connection_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN Connection Router Connection OAuth
    pub allow_google_cloud_cloud_vpn_connection_router_connection_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN Connection Tunnel Connection OAuth
    pub allow_google_cloud_cloud_vpn_connection_tunnel_connection_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN Gateway Connection Gateway Connection OAuth
    pub allow_google_cloud_cloud_vpn_gateway_connection_gateway_connection_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN Gateway Connection HA Gateway Connection OAuth
    pub allow_google_cloud_cloud_vpn_gateway_connection_ha_gateway_connection_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN Gateway Connection Router Connection OAuth
    pub allow_google_cloud_cloud_vpn_gateway_connection_router_connection_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN Gateway Connection Tunnel Connection OAuth
    pub allow_google_cloud_cloud_vpn_gateway_connection_tunnel_connection_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN HA Gateway Connection Gateway Connection OAuth
    pub allow_google_cloud_cloud_vpn_ha_gateway_connection_gateway_connection_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN HA Gateway Connection HA Gateway Connection OAuth
    pub allow_google_cloud_cloud_vpn_ha_gateway_connection_ha_gateway_connection_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN HA Gateway Connection Router Connection OAuth
    pub allow_google_cloud_cloud_vpn_ha_gateway_connection_router_connection_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN HA Gateway Connection Tunnel Connection OAuth
    pub allow_google_cloud_cloud_vpn_ha_gateway_connection_tunnel_connection_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN Router Connection Gateway Connection OAuth
    pub allow_google_cloud_cloud_vpn_router_connection_gateway_connection_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN Router Connection HA Gateway Connection OAuth
    pub allow_google_cloud_cloud_vpn_router_connection_ha_gateway_connection_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN Router Connection Router Connection OAuth
    pub allow_google_cloud_cloud_vpn_router_connection_router_connection_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN Router Connection Tunnel Connection OAuth
    pub allow_google_cloud_cloud_vpn_router_connection_tunnel_connection_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN Tunnel Connection Gateway Connection OAuth
    pub allow_google_cloud_cloud_vpn_tunnel_connection_gateway_connection_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN Tunnel Connection HA Gateway Connection OAuth
    pub allow_google_cloud_cloud_vpn_tunnel_connection_ha_gateway_connection_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN Tunnel Connection Router Connection OAuth
    pub allow_google_cloud_cloud_vpn_tunnel_connection_router_connection_oauth: bool,
    /// Whether to allow Google Cloud Cloud VPN Tunnel Connection Tunnel Connection OAuth
    pub allow_google_cloud_cloud_vpn_tunnel_connection_tunnel_connection_oauth: bool,
}
