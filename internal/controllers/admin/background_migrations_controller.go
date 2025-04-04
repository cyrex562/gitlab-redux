package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gitlab-org/gitlab-redux/internal/controllers"
	"github.com/gitlab-org/gitlab-redux/internal/services"
)

// BackgroundMigrationsController handles background migration management in the admin interface
type BackgroundMigrationsController struct {
	controllers.BaseController
	migrationService *services.BackgroundMigrationService
}

// NewBackgroundMigrationsController creates a new instance of BackgroundMigrationsController
func NewBackgroundMigrationsController(migrationService *services.BackgroundMigrationService) *BackgroundMigrationsController {
	return &BackgroundMigrationsController{
		migrationService: migrationService,
	}
}

// RegisterRoutes registers the routes for this controller
func (c *BackgroundMigrationsController) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/background_migrations", c.RequireAdmin, c.Index)
	router.GET("/background_migrations/:id", c.RequireAdmin, c.Show)
	router.POST("/background_migrations/:id/pause", c.RequireAdmin, c.Pause)
	router.POST("/background_migrations/:id/resume", c.RequireAdmin, c.Resume)
	router.POST("/background_migrations/:id/retry", c.RequireAdmin, c.Retry)
}

// Index displays a list of background migrations
func (c *BackgroundMigrationsController) Index(ctx *gin.Context) {
	// Get query parameters
	tab := ctx.DefaultQuery("tab", "queued")
	database := ctx.DefaultQuery("database", "main")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))

	// Get migrations by status
	migrations, err := c.migrationService.GetMigrationsByStatus(tab, database, page)
	if err != nil {
		c.HandleError(ctx, err)
		return
	}

	// Get successful rows counts
	successfulRows, err := c.migrationService.GetSuccessfulRowsCounts(migrations)
	if err != nil {
		c.HandleError(ctx, err)
		return
	}

	// Get available databases
	databases, err := c.migrationService.GetAvailableDatabases()
	if err != nil {
		c.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"migrations":         migrations,
		"successful_rows":    successfulRows,
		"databases":         databases,
		"current_tab":       tab,
		"current_database":  database,
	})
}

// Show displays details of a specific background migration
func (c *BackgroundMigrationsController) Show(ctx *gin.Context) {
	// Get migration ID
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		c.HandleError(ctx, err)
		return
	}

	// Get migration details
	migration, err := c.migrationService.GetMigration(id)
	if err != nil {
		c.HandleError(ctx, err)
		return
	}

	// Get failed jobs
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	failedJobs, err := c.migrationService.GetFailedJobs(id, page)
	if err != nil {
		c.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"migration":   migration,
		"failed_jobs": failedJobs,
	})
}

// Pause pauses a background migration
func (c *BackgroundMigrationsController) Pause(ctx *gin.Context) {
	// Get migration ID
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		c.HandleError(ctx, err)
		return
	}

	// Pause the migration
	if err := c.migrationService.PauseMigration(id); err != nil {
		c.HandleError(ctx, err)
		return
	}

	ctx.Redirect(http.StatusFound, "/admin/background_migrations")
}

// Resume resumes a paused background migration
func (c *BackgroundMigrationsController) Resume(ctx *gin.Context) {
	// Get migration ID
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		c.HandleError(ctx, err)
		return
	}

	// Resume the migration
	if err := c.migrationService.ResumeMigration(id); err != nil {
		c.HandleError(ctx, err)
		return
	}

	ctx.Redirect(http.StatusFound, "/admin/background_migrations")
}

// Retry retries failed jobs in a background migration
func (c *BackgroundMigrationsController) Retry(ctx *gin.Context) {
	// Get migration ID
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		c.HandleError(ctx, err)
		return
	}

	// Retry failed jobs
	if err := c.migrationService.RetryFailedJobs(id); err != nil {
		c.HandleError(ctx, err)
		return
	}

	ctx.Redirect(http.StatusFound, "/admin/background_migrations")
}

// RequireAdmin is a middleware that ensures the user has admin privileges
func (c *BackgroundMigrationsController) RequireAdmin(ctx *gin.Context) {
	// TODO: Implement proper admin authorization check
	// This should check if the user has the :read_admin_background_migrations permission
	ctx.Next()
}
