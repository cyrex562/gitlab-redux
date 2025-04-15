package routes

import (
	"time"

	"github.com/cyrex562/gitlab-redux/app/handlers"
	"github.com/cyrex562/gitlab-redux/app/middleware"
	"github.com/cyrex562/gitlab-redux/app/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupOAuthRoutes sets up all OAuth-related routes
func SetupOAuthRoutes(router *gin.Engine, db *gorm.DB) {
	baseHandler := handlers.NewBaseHandler()
	oauthService := services.NewOAuthService(db)
	authService := services.NewAuthService(db, "your-jwt-secret", time.Hour*24)
	authMiddleware := middleware.NewAuthMiddleware(authService)
	
	appsHandler := handlers.NewOAuthApplicationsHandler(baseHandler, oauthService)
	authHandler := handlers.NewOAuthAuthorizationsHandler(baseHandler, oauthService)
	authorizedAppsHandler := handlers.NewOAuthAuthorizedApplicationsHandler(baseHandler, oauthService)

	// OAuth Applications routes (require authentication)
	oauthGroup := router.Group("/oauth", authMiddleware.RequireAuth())
	{
		// OAuth Applications routes
		oauthGroup.GET("/applications", appsHandler.Index)
		oauthGroup.GET("/applications/:id", appsHandler.Show)
		oauthGroup.POST("/applications", appsHandler.Create)
		oauthGroup.POST("/applications/:id/renew_secret", appsHandler.RenewSecret)

		// OAuth Authorizations routes
		oauthGroup.GET("/authorizations/new", authHandler.New)
		oauthGroup.POST("/authorizations", authHandler.Create)

		// OAuth Authorized Applications routes
		oauthGroup.GET("/authorized_applications", authorizedAppsHandler.Index)
		oauthGroup.DELETE("/authorized_applications/:id", authorizedAppsHandler.Destroy)
	}

	// Token revocation
	tokenRevocationHandler := handlers.NewOAuthTokenRevocationHandler(baseHandler, oauthService)
	router.POST("/oauth/revoke", tokenRevocationHandler.Revoke)
} 