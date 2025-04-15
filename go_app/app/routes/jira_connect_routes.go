package routes

import (
	"net/http"

	"gitlab.com/gitlab-org/gitlab-redux/app/handlers/jiraconnect"
	"gitlab.com/gitlab-org/gitlab-redux/app/services"
)

// RegisterJiraConnectRoutes registers all Jira Connect routes
func RegisterJiraConnectRoutes(mux *http.ServeMux, baseHandler *jiraconnect.BaseHandler, jiraService *services.JiraService) {
	// Create the application handler
	appHandler := jiraconnect.NewApplicationHandler(baseHandler, jiraService)
	
	// Create the branches handler
	branchesHandler := jiraconnect.NewBranchesHandler(appHandler)
	
	// Create the events handler
	eventsHandler := jiraconnect.NewEventsHandler(appHandler, jiraService)
	
	// Register routes
	
	// Branches routes
	mux.HandleFunc("/jira_connect/branches/new", branchesHandler.NewBranchHandler)
	mux.HandleFunc("/jira_connect/branches/route", branchesHandler.RouteHandler)
	
	// Events routes
	mux.HandleFunc("/jira_connect/events/installed", eventsHandler.InstalledHandler)
	mux.HandleFunc("/jira_connect/events/uninstalled", eventsHandler.UninstalledHandler)
	
	// Add more routes as needed
} 