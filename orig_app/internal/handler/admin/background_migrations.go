package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

type BackgroundMigrationsController struct {
	migrationService *service.BackgroundMigrationService
}

func NewBackgroundMigrationsController(migrationService *service.BackgroundMigrationService) *BackgroundMigrationsController {
	return &BackgroundMigrationsController{
		migrationService: migrationService,
	}
}

// Index handles listing background migrations with filtering by status
func (c *BackgroundMigrationsController) Index(ctx *gin.Context) {
	tab := ctx.DefaultQuery("tab", "queued")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	database := ctx.DefaultQuery("database", "main")

	migrations, err := c.migrationService.ListMigrations(ctx, tab, page, database)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	databases, err := c.migrationService.ListDatabases(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"migrations": migrations,
		"databases": databases,
		"current_tab": tab,
	})
}

// Show displays details of a specific migration
func (c *BackgroundMigrationsController) Show(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid migration ID"})
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))

	migration, failedJobs, err := c.migrationService.GetMigration(ctx, id, page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"migration": migration,
		"failed_jobs": failedJobs,
	})
}

// Pause pauses a specific migration
func (c *BackgroundMigrationsController) Pause(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid migration ID"})
		return
	}

	if err := c.migrationService.PauseMigration(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Redirect(http.StatusFound, "/admin/background_migrations")
}

// Resume resumes a paused migration
func (c *BackgroundMigrationsController) Resume(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid migration ID"})
		return
	}

	if err := c.migrationService.ResumeMigration(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Redirect(http.StatusFound, "/admin/background_migrations")
}

// Retry retries failed jobs in a migration
func (c *BackgroundMigrationsController) Retry(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid migration ID"})
		return
	}

	if err := c.migrationService.RetryFailedJobs(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Redirect(http.StatusFound, "/admin/background_migrations")
}
