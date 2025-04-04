package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// UsersController handles user management for GitLab instance
type UsersController struct {
	userService *service.UserService
}

// NewUsersController creates a new instance of UsersController
func NewUsersController(userService *service.UserService) *UsersController {
	return &UsersController{
		userService: userService,
	}
}

// RegisterRoutes registers the routes for the UsersController
func (c *UsersController) RegisterRoutes(r *gin.RouterGroup) {
	users := r.Group("/admin/users")
	{
		users.Use(c.requireAdmin)
		users.GET("/", c.index)
		users.GET("/new", c.new)
		users.GET("/:id", c.show)
		users.GET("/:id/edit", c.edit)
		users.GET("/:id/projects", c.projects)
		users.GET("/:id/keys", c.keys)
		users.POST("/", c.create)
		users.PUT("/:id", c.update)
		users.DELETE("/:id", c.destroy)
		users.POST("/:id/impersonate", c.impersonate)
		users.POST("/:id/approve", c.approve)
		users.POST("/:id/reject", c.reject)
		users.POST("/:id/activate", c.activate)
		users.POST("/:id/deactivate", c.deactivate)
		users.POST("/:id/block", c.block)
		users.POST("/:id/unblock", c.unblock)
		users.POST("/:id/ban", c.ban)
		users.POST("/:id/unban", c.unban)
		users.POST("/:id/unlock", c.unlock)
		users.POST("/:id/trust", c.trust)
		users.POST("/:id/untrust", c.untrust)
		users.POST("/:id/confirm", c.confirm)
		users.POST("/:id/disable_two_factor", c.disableTwoFactor)
		users.DELETE("/:id/emails/:email_id", c.removeEmail)
	}
}

// requireAdmin middleware ensures that only admin users can access these endpoints
func (c *UsersController) requireAdmin(ctx *gin.Context) {
	user := ctx.MustGet("user")
	if user == nil || !user.IsAdmin() {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		ctx.Abort()
		return
	}
	ctx.Next()
}

// index handles the GET /admin/users endpoint
func (c *UsersController) index(ctx *gin.Context) {
	// Handle cohorts tab redirect
	if ctx.Query("tab") == "cohorts" {
		ctx.Redirect(http.StatusFound, "/admin/cohorts")
		return
	}

	filter := ctx.Query("filter")
	searchQuery := ctx.Query("search_query")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	sort := ctx.DefaultQuery("sort", "name_asc")

	users, err := c.userService.List(ctx, filter, searchQuery, page, sort)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	// TODO: Implement HTML rendering
	ctx.JSON(http.StatusOK, users)
}

// show handles the GET /admin/users/:id endpoint
func (c *UsersController) show(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := c.userService.Get(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// TODO: Implement HTML rendering
	ctx.JSON(http.StatusOK, user)
}

// projects handles the GET /admin/users/:id/projects endpoint
func (c *UsersController) projects(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	projects, err := c.userService.GetProjects(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user projects"})
		return
	}

	// TODO: Implement HTML rendering
	ctx.JSON(http.StatusOK, projects)
}

// keys handles the GET /admin/users/:id/keys endpoint
func (c *UsersController) keys(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	keys, err := c.userService.GetKeys(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user keys"})
		return
	}

	// TODO: Implement HTML rendering
	ctx.JSON(http.StatusOK, keys)
}

// create handles the POST /admin/users endpoint
func (c *UsersController) create(ctx *gin.Context) {
	var params model.UserParams
	if err := ctx.ShouldBind(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user parameters"})
		return
	}

	currentUser := ctx.MustGet("user").(*model.User)
	params.ResetPassword = true
	params.SkipConfirmation = true

	user, err := c.userService.Create(ctx, currentUser, &params)
	if err != nil {
		// TODO: Implement HTML rendering for new user form with errors
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	ctx.Redirect(http.StatusFound, "/admin/users/"+strconv.FormatInt(user.ID, 10))
}

// update handles the PUT /admin/users/:id endpoint
func (c *UsersController) update(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var params model.UserParams
	if err := ctx.ShouldBind(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user parameters"})
		return
	}

	currentUser := ctx.MustGet("user").(*model.User)
	user, err := c.userService.Update(ctx, currentUser, id, &params)
	if err != nil {
		// TODO: Implement HTML rendering for edit user form with errors
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	ctx.Redirect(http.StatusFound, "/admin/users/"+strconv.FormatInt(user.ID, 10))
}

// destroy handles the DELETE /admin/users/:id endpoint
func (c *UsersController) destroy(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	currentUser := ctx.MustGet("user").(*model.User)
	hardDelete := ctx.Query("hard_delete") == "true"

	if err := c.userService.Destroy(ctx, currentUser, id, hardDelete); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	ctx.Redirect(http.StatusFound, "/admin/users")
}

// impersonate handles the POST /admin/users/:id/impersonate endpoint
func (c *UsersController) impersonate(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	currentUser := ctx.MustGet("user").(*model.User)
	if err := c.userService.Impersonate(ctx, currentUser, id); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	ctx.Redirect(http.StatusFound, "/")
}

// Additional action handlers...
// TODO: Implement approve, reject, activate, deactivate, block, unblock, ban, unban, unlock, trust, untrust, confirm, disableTwoFactor, removeEmail
