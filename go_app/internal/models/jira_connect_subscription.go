package models

// JiraConnectSubscription represents a subscription between a GitLab namespace and a Jira installation
type JiraConnectSubscription struct {
	ID             string `json:"id"`
	InstallationID string `json:"installation_id"`
	NamespacePath  string `json:"namespace_path"`
	JiraUser       string `json:"jira_user"`
	NamespaceRoute string `json:"namespace_route"`
} 