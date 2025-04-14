package groups

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/models"
	"gitlab.com/gitlab-org/gitlab-redux/internal/services/ci"
	"gitlab.com/gitlab-org/gitlab-redux/internal/services/tracking"
)

// RunnersController handles requests for group runners
type RunnersController struct {
	runnersFinder      *ci.RunnersFinder
	updateRunnerService *ci.UpdateRunnerService
	trackingService    *tracking.Service
}

// NewRunnersController creates a new runners controller
func NewRunnersController(
	runnersFinder *ci.RunnersFinder,
	updateRunnerService *ci.UpdateRunnerService,
	trackingService *tracking.Service,
) *RunnersController {
	return &RunnersController{
		runnersFinder:      runnersFinder,
		updateRunnerService: updateRunnerService,
		trackingService:    trackingService,
	}
}

// RegisterRoutes registers the routes for the runners controller
func (c *RunnersController) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/runners", c.authorizeReadGroupRunners(), c.Index)
	router.GET("/runners/:id", c.authorizeReadGroupRunners(), c.Show)
	router.GET("/runners/new", c.authorizeCreateGroupRunners(), c.New)
	router.GET("/runners/:id/edit", c.authorizeUpdateRunner(), c.Edit)
	router.PUT("/runners/:id", c.authorizeUpdateRunner(), c.Update)
	router.GET("/runners/:id/register", c.authorizeCreateGroupRunners(), c.Register)
}

// Index handles GET requests for group runners
func (c *RunnersController) Index(ctx *gin.Context) {
	group := ctx.MustGet("group").(*models.Group)
	user := ctx.MustGet("current_user").(*models.User)

	allowRegistrationToken := group.AllowRunnerRegistrationToken()
	var groupRunnerRegistrationToken string
	if user.CanRegisterGroupRunners(group) {
		groupRunnerRegistrationToken = group.RunnersToken
	}

	var groupNewRunnerPath string
	if user.CanCreateRunner(group) {
		groupNewRunnerPath = "/groups/" + strconv.FormatInt(group.ID, 10) + "/runners/new"
	}

	// Track the event
	c.trackingService.TrackEvent(ctx, "RunnersController", "index", user, group)

	ctx.JSON(http.StatusOK, gin.H{
		"allow_registration_token":      allowRegistrationToken,
		"group_runner_registration_token": groupRunnerRegistrationToken,
		"group_new_runner_path":          groupNewRunnerPath,
	})
}

// Show handles GET requests for a specific runner
func (c *RunnersController) Show(ctx *gin.Context) {
	// Implementation will be added when needed
	ctx.Status(http.StatusOK)
}

// Edit handles GET requests for editing a runner
func (c *RunnersController) Edit(ctx *gin.Context) {
	// Implementation will be added when needed
	ctx.Status(http.StatusOK)
}

// Update handles PUT requests for updating a runner
func (c *RunnersController) Update(ctx *gin.Context) {
	runner := ctx.MustGet("runner").(*models.Runner)
	user := ctx.MustGet("current_user").(*models.User)
	group := ctx.MustGet("group").(*models.Group)

	var params models.RunnerUpdateParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.updateRunnerService.Execute(ctx, user, runner, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.Success {
		ctx.Redirect(http.StatusFound, "/groups/"+strconv.FormatInt(group.ID, 10)+"/runners/"+strconv.FormatInt(runner.ID, 10))
		ctx.Set("notice", "Runner was successfully updated.")
	} else {
		ctx.HTML(http.StatusOK, "runners/edit", gin.H{
			"runner": runner,
			"group":  group,
			"errors": result.Errors,
		})
	}
}

// New handles GET requests for creating a new runner
func (c *RunnersController) New(ctx *gin.Context) {
	// Implementation will be added when needed
	ctx.Status(http.StatusOK)
}

// Register handles GET requests for registering a runner
func (c *RunnersController) Register(ctx *gin.Context) {
	runner := ctx.MustGet("runner").(*models.Runner)
	
	if !runner.RegistrationAvailable() {
		ctx.Status(http.StatusNotFound)
		return
	}
	
	// Implementation will be added when needed
	ctx.Status(http.StatusOK)
}

// Helper methods for authorization

func (c *RunnersController) authorizeReadGroupRunners() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("current_user").(*models.User)
		group := ctx.MustGet("group").(*models.Group)
		
		if !user.CanReadGroupRunners(group) {
			ctx.Status(http.StatusNotFound)
			ctx.Abort()
			return
		}
		
		ctx.Next()
	}
}

func (c *RunnersController) authorizeCreateGroupRunners() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("current_user").(*models.User)
		group := ctx.MustGet("group").(*models.Group)
		
		if !user.CanCreateRunner(group) {
			ctx.Status(http.StatusNotFound)
			ctx.Abort()
			return
		}
		
		ctx.Next()
	}
}

func (c *RunnersController) authorizeUpdateRunner() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("current_user").(*models.User)
		runner := ctx.MustGet("runner").(*models.Runner)
		
		if !user.CanUpdateRunner(runner) {
			ctx.Status(http.StatusNotFound)
			ctx.Abort()
			return
		}
		
		ctx.Next()
	}
}

// Middleware to load the runner
func (c *RunnersController) loadRunner() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("current_user").(*models.User)
		group := ctx.MustGet("group").(*models.Group)
		
		runnerID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			ctx.Abort()
			return
		}
		
		params := ci.RunnersFinderParams{
			Group:      group,
			Membership: "all_available",
		}
		
		runner, err := c.runnersFinder.Execute(ctx, user, params, runnerID)
		if err != nil {
			ctx.Status(http.StatusNotFound)
			ctx.Abort()
			return
		}
		
		ctx.Set("runner", runner)
		ctx.Next()
	}
} 