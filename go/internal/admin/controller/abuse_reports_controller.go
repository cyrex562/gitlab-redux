package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/models"
	"gitlab.com/gitlab-org/gitlab-redux/internal/routing"
	"gitlab.com/gitlab-org/gitlab-redux/internal/services"
)

// AbuseReportsController handles abuse report operations in the admin panel
type AbuseReportsController struct {
	*routing.BaseController
}

// NewAbuseReportsController creates a new instance of AbuseReportsController
func NewAbuseReportsController() *AbuseReportsController {
	return &AbuseReportsController{
		BaseController: routing.NewBaseController(),
	}
}

// SetupRoutes configures the routes for the abuse reports controller
func (c *AbuseReportsController) SetupRoutes(router *gin.Engine) {
	admin := router.Group("/admin/abuse_reports")
	{
		admin.GET("", c.Index)
		admin.GET("/:id", c.Show)
		admin.PUT("/:id", c.Update)
		admin.POST("/:id/moderate_user", c.ModerateUser)
		admin.DELETE("/:id", c.Destroy)
	}
}

// Index displays a list of abuse reports
func (c *AbuseReportsController) Index(ctx *gin.Context) {
	params := c.getIndexParams(ctx)
	reports, err := models.FindAbuseReports(params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.HTML(http.StatusOK, "admin/abuse_reports/index", gin.H{
		"abuse_reports": reports,
	})
}

// Show displays a single abuse report
func (c *AbuseReportsController) Show(ctx *gin.Context) {
	report, err := c.findAbuseReport(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Abuse report not found"})
		return
	}

	// Push feature flags to frontend
	ctx.Set("abuse_report_labels", true)
	ctx.Set("abuse_report_notes", true)

	ctx.HTML(http.StatusOK, "admin/abuse_reports/show", gin.H{
		"abuse_report": report,
	})
}

// Update updates an abuse report
func (c *AbuseReportsController) Update(ctx *gin.Context) {
	report, err := c.findAbuseReport(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Abuse report not found"})
		return
	}

	var params models.AbuseReportParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service := services.NewUpdateAbuseReportService(report, c.GetCurrentUser(ctx), params)
	response, err := service.Execute()
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

// ModerateUser moderates a user based on an abuse report
func (c *AbuseReportsController) ModerateUser(ctx *gin.Context) {
	report, err := c.findAbuseReport(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Abuse report not found"})
		return
	}

	var params models.AbuseReportParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service := services.NewModerateUserService(report, c.GetCurrentUser(ctx), params)
	response, err := service.Execute()
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": response.Message})
}

// Destroy removes an abuse report
func (c *AbuseReportsController) Destroy(ctx *gin.Context) {
	report, err := c.findAbuseReport(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Abuse report not found"})
		return
	}

	var params struct {
		RemoveUser bool `json:"remove_user"`
	}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if params.RemoveUser {
		if err := report.RemoveUser(c.GetCurrentUser(ctx)); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if err := report.Delete(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

// Helper methods

func (c *AbuseReportsController) findAbuseReport(ctx *gin.Context) (*models.AbuseReport, error) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		return nil, err
	}
	return models.GetAbuseReportByID(uint(id))
}

func (c *AbuseReportsController) getIndexParams(ctx *gin.Context) models.AbuseReportSearchParams {
	status := ctx.DefaultQuery("status", "open")
	return models.AbuseReportSearchParams{
		Page:     ctx.DefaultQuery("page", "1"),
		Status:   status,
		Category: ctx.Query("category"),
		User:     ctx.Query("user"),
		Reporter: ctx.Query("reporter"),
		Sort:     ctx.Query("sort"),
	}
}
