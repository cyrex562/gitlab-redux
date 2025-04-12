package settings

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ApplicationsController handles group OAuth applications
type ApplicationsController struct {
	applicationsCreateService *ApplicationsCreateService
	doorkeeperApplicationModel *DoorkeeperApplicationModel
}

// NewApplicationsController creates a new ApplicationsController
func NewApplicationsController(
	applicationsCreateService *ApplicationsCreateService,
	doorkeeperApplicationModel *DoorkeeperApplicationModel,
) *ApplicationsController {
	return &ApplicationsController{
		applicationsCreateService: applicationsCreateService,
		doorkeeperApplicationModel: doorkeeperApplicationModel,
	}
}

// RegisterRoutes registers the routes for the ApplicationsController
func (c *ApplicationsController) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/", c.Index)
	router.GET("/:id", c.Show)
	router.GET("/:id/edit", c.Edit)
	router.POST("/", c.Create)
	router.PUT("/:id", c.Update)
	router.POST("/:id/renew", c.Renew)
	router.DELETE("/:id", c.Destroy)
}

// RegisterMiddleware registers the middleware for the ApplicationsController
func (c *ApplicationsController) RegisterMiddleware(router *gin.RouterGroup) {
	router.Use(c.AuthorizeAdminGroup())
	router.Use(c.LoadScopes())
}

// Index handles the index action
func (c *ApplicationsController) Index(ctx *gin.Context) {
	c.SetIndexVars(ctx)
	ctx.HTML(http.StatusOK, "groups/settings/applications/index", gin.H{
		"applications": ctx.MustGet("applications"),
		"applications_total_count": ctx.MustGet("applications_total_count"),
		"application": ctx.MustGet("application"),
	})
}

// Show handles the show action
func (c *ApplicationsController) Show(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "groups/settings/applications/show", gin.H{
		"application": ctx.MustGet("application"),
		"created": ctx.MustGet("created"),
	})
}

// Edit handles the edit action
func (c *ApplicationsController) Edit(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "groups/settings/applications/edit", gin.H{
		"application": ctx.MustGet("application"),
	})
}

// Create handles the create action
func (c *ApplicationsController) Create(ctx *gin.Context) {
	// Get the current user from the context
	currentUser, _ := ctx.Get("current_user")

	// Get the group from the context
	group, _ := ctx.Get("group")

	// Get the application parameters
	applicationParams := c.GetApplicationParams(ctx)
	applicationParams["owner"] = group

	// Create the application
	application, err := c.applicationsCreateService.Execute(currentUser, applicationParams, ctx.Request)
	if err != nil {
		ctx.Set("application", application)
		c.SetIndexVars(ctx)
		ctx.HTML(http.StatusOK, "groups/settings/applications/index", gin.H{
			"applications": ctx.MustGet("applications"),
			"applications_total_count": ctx.MustGet("applications_total_count"),
			"application": ctx.MustGet("application"),
		})
		return
	}

	// Set flash notice
	ctx.SetFlash("notice", "Application was successfully created.")

	// Set created flag
	ctx.Set("created", true)
	ctx.Set("application", application)

	// Render the show template
	ctx.HTML(http.StatusOK, "groups/settings/applications/show", gin.H{
		"application": ctx.MustGet("application"),
		"created": ctx.MustGet("created"),
	})
}

// Update handles the update action
func (c *ApplicationsController) Update(ctx *gin.Context) {
	// Get the application from the context
	application, _ := ctx.Get("application")

	// Get the application parameters
	applicationParams := c.GetApplicationParams(ctx)
	applicationParams["owner"] = ctx.MustGet("group")

	// Update the application
	err := application.(*DoorkeeperApplicationModel).Update(applicationParams)
	if err != nil {
		ctx.HTML(http.StatusOK, "groups/settings/applications/edit", gin.H{
			"application": ctx.MustGet("application"),
		})
		return
	}

	// Redirect to the show page
	ctx.Redirect(http.StatusFound, "/groups/:id/settings/applications/"+application.(*DoorkeeperApplicationModel).ID)
}

// Renew handles the renew action
func (c *ApplicationsController) Renew(ctx *gin.Context) {
	// Get the application from the context
	application, _ := ctx.Get("application")

	// Renew the secret
	err := application.(*DoorkeeperApplicationModel).RenewSecret()
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"errors": application.(*DoorkeeperApplicationModel).Errors,
		})
		return
	}

	// Save the application
	err = application.(*DoorkeeperApplicationModel).Save()
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"errors": application.(*DoorkeeperApplicationModel).Errors,
		})
		return
	}

	// Return the secret
	ctx.JSON(http.StatusOK, gin.H{
		"secret": application.(*DoorkeeperApplicationModel).PlaintextSecret,
	})
}

// Destroy handles the destroy action
func (c *ApplicationsController) Destroy(ctx *gin.Context) {
	// Get the application from the context
	application, _ := ctx.Get("application")

	// Destroy the application
	err := application.(*DoorkeeperApplicationModel).Destroy()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Redirect to the index page
	ctx.Redirect(http.StatusFound, "/groups/:id/settings/applications")
}

// SetIndexVars sets the index variables
func (c *ApplicationsController) SetIndexVars(ctx *gin.Context) {
	// Get the group from the context
	group, _ := ctx.Get("group")

	// Get the cursor from the query
	cursor := ctx.Query("cursor")

	// Get the applications
	applications, err := group.(*GroupModel).OAuthApplications.KeysetPaginate(map[string]string{"cursor": cursor})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Get the applications total count
	applicationsTotalCount, err := group.(*GroupModel).OAuthApplications.Count()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Set the applications
	ctx.Set("applications", applications)

	// Set the applications total count
	ctx.Set("applications_total_count", applicationsTotalCount)

	// Set the application if not already set
	if _, exists := ctx.Get("application"); !exists {
		ctx.Set("application", c.doorkeeperApplicationModel.New())
	}
}

// SetApplication sets the application
func (c *ApplicationsController) SetApplication(ctx *gin.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the group from the context
		group, _ := ctx.Get("group")

		// Get the application ID from the params
		applicationID := ctx.Param("id")

		// Get the application
		application, err := group.(*GroupModel).OAuthApplications.Find(applicationID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Application not found",
			})
			ctx.Abort()
			return
		}

		// Set the application
		ctx.Set("application", application)

		ctx.Next()
	}
}

// GetApplicationParams gets the application parameters
func (c *ApplicationsController) GetApplicationParams(ctx *gin.Context) map[string]interface{} {
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would get the application parameters
	return map[string]interface{}{
		"name": ctx.PostForm("name"),
		"redirect_uri": ctx.PostForm("redirect_uri"),
		"scopes": ctx.PostFormArray("scopes"),
	}
}

// AuthorizeAdminGroup middleware checks if the user has permission to admin the group
func (c *ApplicationsController) AuthorizeAdminGroup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the current user from the context
		currentUser, exists := ctx.Get("current_user")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			ctx.Abort()
			return
		}

		// Get the group from the context
		group, exists := ctx.Get("group")
		if !exists {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
			ctx.Abort()
			return
		}

		// Check if the user has permission to admin the group
		// This is a placeholder for the actual implementation
		// In the actual implementation, you would check if the user has the admin_group permission
		canAdminGroup := true // Replace with actual check

		if !canAdminGroup {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "You don't have permission to admin the group"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// LoadScopes middleware loads the scopes
func (c *ApplicationsController) LoadScopes() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// This is a placeholder for the actual implementation
		// In the actual implementation, you would load the scopes
		ctx.Next()
	}
}

// ApplicationsCreateService creates OAuth applications
type ApplicationsCreateService struct {
	// Add fields as needed
}

// Execute executes the service
func (s *ApplicationsCreateService) Execute(currentUser interface{}, params map[string]interface{}, request interface{}) (*DoorkeeperApplicationModel, error) {
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would create the application
	return nil, nil
}

// DoorkeeperApplicationModel represents a Doorkeeper application
type DoorkeeperApplicationModel struct {
	ID              string
	Name            string
	RedirectURI     string
	Scopes          []string
	PlaintextSecret string
	Errors          map[string][]string
}

// New creates a new DoorkeeperApplicationModel
func (m *DoorkeeperApplicationModel) New() *DoorkeeperApplicationModel {
	return &DoorkeeperApplicationModel{}
}

// Update updates the application
func (m *DoorkeeperApplicationModel) Update(params map[string]interface{}) error {
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would update the application
	return nil
}

// RenewSecret renews the secret
func (m *DoorkeeperApplicationModel) RenewSecret() error {
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would renew the secret
	return nil
}

// Save saves the application
func (m *DoorkeeperApplicationModel) Save() error {
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would save the application
	return nil
}

// Destroy destroys the application
func (m *DoorkeeperApplicationModel) Destroy() error {
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would destroy the application
	return nil
}

// GroupModel represents a GitLab group
type GroupModel struct {
	OAuthApplications *OAuthApplicationsModel
}

// OAuthApplicationsModel represents OAuth applications
type OAuthApplicationsModel struct {
	// Add fields as needed
}

// KeysetPaginate paginates the applications
func (m *OAuthApplicationsModel) KeysetPaginate(params map[string]string) (interface{}, error) {
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would paginate the applications
	return nil, nil
}

// Count counts the applications
func (m *OAuthApplicationsModel) Count() (int, error) {
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would count the applications
	return 0, nil
}

// Find finds an application
func (m *OAuthApplicationsModel) Find(id string) (*DoorkeeperApplicationModel, error) {
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would find an application
	return nil, nil
}
