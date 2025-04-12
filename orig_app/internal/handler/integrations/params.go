package integrations

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
)

// Params provides integration parameter validation and processing
type Params struct {
	integration *model.Integration
}

// NewParams creates a new instance of Params
func NewParams(integration *model.Integration) *Params {
	return &Params{
		integration: integration,
	}
}

// ParseParams parses and validates integration parameters
func (p *Params) ParseParams(ctx *gin.Context) (map[string]interface{}, error) {
	// Get allowed parameters
	allowed := p.getAllowedParams()

	// Parse request body
	var body struct {
		ID          string                 `json:"id"`
		Integration map[string]interface{} `json:"integration"`
		Service     map[string]interface{} `json:"service"`
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		return nil, fmt.Errorf("failed to parse request body: %w", err)
	}

	// Get integration parameters
	params := body.Integration
	if params == nil {
		params = body.Service
	}

	// Filter allowed parameters
	filtered := make(map[string]interface{})
	for key, value := range params {
		if allowed[key] {
			filtered[key] = value
		}
	}

	// Handle special cases for chat integrations
	if p.integration.IsChat() && (ctx.Request.Method == "PUT" || ctx.Request.Method == "POST") {
		// Remove webhook if it's the secret mask
		if webhook, ok := filtered["webhook"].(string); ok && webhook == model.ChatNotificationSecretMask {
			delete(filtered, "webhook")
		}

		// Handle configurable channels
		if p.integration.HasConfigurableChannels() {
			for _, channel := range p.integration.GetEventChannelNames() {
				if value, ok := filtered[channel].(string); ok && value == model.ChatNotificationSecretMask {
					delete(filtered, channel)
				}
			}
		}
	}

	// Remove empty secret fields
	for _, field := range p.integration.GetSecretFields() {
		if value, ok := filtered[field].(string); ok && value == "" {
			delete(filtered, field)
		}
	}

	return filtered, nil
}

// getAllowedParams returns a map of allowed parameter names
func (p *Params) getAllowedParams() map[string]bool {
	allowed := make(map[string]bool)

	// Add CE parameters
	for _, param := range allowedParamsCE {
		allowed[param] = true
	}

	// Add dynamic parameters from integration
	for _, param := range p.integration.GetEventChannelNames() {
		allowed[param] = true
	}
	for _, param := range p.integration.GetEventNames() {
		allowed[param] = true
	}

	return allowed
}

// allowedParamsCE contains the list of allowed parameters for CE integrations
var allowedParamsCE = []string{
	"app_store_issuer_id",
	"app_store_key_id",
	"app_store_private_key",
	"app_store_private_key_file_name",
	"app_store_protected_refs",
	"active",
	"alert_events",
	"api_key",
	"api_token",
	"api_url",
	"archive_trace_events",
	"bamboo_url",
	"branches_to_be_notified",
	"labels_to_be_notified",
	"labels_to_be_notified_behavior",
	"build_key",
	"build_type",
	"ca_pem",
	"channel",
	"channels",
	"color",
	"colorize_messages",
	"comment_on_event_enabled",
	"comment_detail",
	"confidential_issues_events",
	"confluence_url",
	"datadog_site",
	"datadog_env",
	"datadog_service",
	"datadog_tags",
	"datadog_ci_visibility",
	"default_irc_uri",
	"device",
	"disable_diffs",
	"diffblue_access_token_name",
	"diffblue_access_token_secret",
	"diffblue_license_key",
	"drone_url",
	"enable_ssl_verification",
	"exclude_service_accounts",
	"external_wiki_url",
	"google_iap_service_account_json",
	"google_iap_audience_client_id",
	"google_play_protected_refs",
	"group_confidential_mention_events",
	"group_mention_events",
	"hostname",
	"incident_events",
	"inherit_from_id",
	"issues_events",
	"issues_url",
	"jenkins_url",
	"jira_auth_type",
	"jira_issue_prefix",
	"jira_issue_regex",
	"jira_issue_transition_automatic",
	"jira_issue_transition_id",
	"jira_cloud_app_service_ids",
	"jira_cloud_app_enable_deployment_gating",
	"jira_cloud_app_deployment_gating_environments",
	"manual_configuration",
	"merge_requests_events",
	"mock_service_url",
	"namespace",
	"new_issue_url",
	"notify_only_broken_pipelines",
	"package_name",
	"password",
	"priority",
	"project_key",
	"project_name",
	"project_url",
	"recipients",
	"restrict_to_branch",
	"room",
	"send_from_committer_email",
	"server",
	"server_host",
	"server_port",
	"service_account_key",
	"service_account_key_file_name",
	"sound",
	"subdomain",
	"teamcity_url",
	"thread",
	"token",
	"type",
	"url",
	"user_key",
	"username",
	"webhook",
	"zentao_product_xid",
}
