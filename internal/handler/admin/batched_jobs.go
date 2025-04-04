package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

type BatchedJobsController struct {
	jobService *service.BatchedJobService
}

func NewBatchedJobsController(jobService *service.BatchedJobService) *BatchedJobsController {
	return &BatchedJobsController{
		jobService: jobService,
	}
}

// Show displays details of a specific batched job and its transition logs
func (c *BatchedJobsController) Show(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	database := ctx.DefaultQuery("database", "main")

	job, transitionLogs, err := c.jobService.GetJob(ctx, id, database)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"job": job,
		"transition_logs": transitionLogs,
	})
}
