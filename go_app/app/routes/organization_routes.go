package routes

import (
	"github.com/cyrex562/gitlab-redux/app/handlers"
	"github.com/cyrex562/gitlab-redux/app/handlers/organizations"
	"github.com/cyrex562/gitlab-redux/app/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupOrganizationRoutes sets up all organization-related routes
func SetupOrganizationRoutes(router *gin.Engine, db *gorm.DB) {
	baseHandler := handlers.NewBaseHandler()
	organizationService := services.NewOrganizationService(db)
	featureService := services.NewFeatureService(db)
	groupService := services.NewGroupService(db)

	orgBaseHandler := organizations.NewBaseHandler(baseHandler, organizationService, featureService)
	groupsHandler := organizations.NewGroupsHandler(orgBaseHandler, groupService)

	// Organization routes
	orgGroup := router.Group("/organizations/:organization_path")
	orgGroup.Use(orgBaseHandler.RequireOrganizationFeature())
	{
		// Groups routes
		orgGroup.GET("/groups/new", orgBaseHandler.RequireGroupCreation(), groupsHandler.New)
		orgGroup.GET("/groups/:id/edit", orgBaseHandler.RequireOrganizationRead(), groupsHandler.Edit)
		orgGroup.POST("/groups", orgBaseHandler.RequireGroupCreation(), groupsHandler.Create)
		orgGroup.DELETE("/groups/:id", orgBaseHandler.RequireOrganizationRead(), groupsHandler.Destroy)
	}
} 