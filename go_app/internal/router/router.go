package router

import (
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/jmadden/gitlab-redux/internal/handlers"
	"github.com/jmadden/gitlab-redux/internal/services"
)

// SetupRouter configures and returns the application router
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Load HTML templates
	r.LoadHTMLGlob(filepath.Join("internal", "templates", "**", "*.html"))

	// Initialize services
	jiraService := services.NewJiraService()

	// Initialize handlers
	installationsHandler := handlers.NewInstallationsHandler(jiraService)
	oauthApplicationIdsHandler := handlers.NewOauthApplicationIdsHandler(jiraService)
	oauthCallbacksHandler := handlers.NewOauthCallbacksHandler()
	publicKeysHandler := handlers.NewPublicKeysHandler(jiraService)
	repositoriesHandler := handlers.NewRepositoriesHandler(jiraService)
	subscriptionsHandler := handlers.NewSubscriptionsHandler(jiraService)

	// Jira Connect routes
	jiraGroup := r.Group("/api/v4/jira/connect")
	{
		// Installation routes
		jiraGroup.GET("/installations", installationsHandler.GetInstallation)
		jiraGroup.PUT("/installations", installationsHandler.UpdateInstallation)

		// OAuth application ID route - no JWT verification required
		jiraGroup.GET("/oauth_application_ids", oauthApplicationIdsHandler.ShowApplicationId)

		// OAuth callbacks route - no authentication required
		jiraGroup.GET("/oauth/callbacks", oauthCallbacksHandler.Index)

		// Public keys route - no authentication required
		jiraGroup.GET("/public_keys/:id", publicKeysHandler.ShowPublicKey)

		// Repository routes
		jiraGroup.GET("/repositories/search", repositoriesHandler.SearchRepositories)
		jiraGroup.GET("/repositories/associate", repositoriesHandler.AssociateRepository)

		// Subscription routes
		jiraGroup.GET("/subscriptions", subscriptionsHandler.Index)
		jiraGroup.POST("/subscriptions", subscriptionsHandler.Create)
		jiraGroup.DELETE("/subscriptions/:id", subscriptionsHandler.Destroy)
	}

	return r
} 