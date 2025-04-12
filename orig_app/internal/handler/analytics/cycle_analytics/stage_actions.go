package cycle_analytics

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// StageActions provides common functionality for cycle analytics stage actions
type StageActions struct {
	stageService    *service.StageService
	dataCollector   *service.DataCollector
	valueStream     *model.ValueStream
	paramsValidator *service.ParamsValidator
}

// NewStageActions creates a new instance of StageActions
func NewStageActions(
	stageService *service.StageService,
	dataCollector *service.DataCollector,
	valueStream *model.ValueStream,
	paramsValidator *service.ParamsValidator,
) *StageActions {
	return &StageActions{
		stageService:    stageService,
		dataCollector:   dataCollector,
		valueStream:     valueStream,
		paramsValidator: paramsValidator,
	}
}

// RegisterRoutes registers the routes for stage actions
func (s *StageActions) RegisterRoutes(r *gin.RouterGroup) {
	stages := r.Group("/stages")
	{
		// Apply middleware
		stages.Use(s.validateParams)
		stages.Use(s.authorizeStage)

		// Register routes
		stages.GET("/", s.index)
		stages.GET("/:id/median", s.median)
		stages.GET("/:id/average", s.average)
		stages.GET("/:id/records", s.records)
		stages.GET("/:id/count", s.count)
	}
}

// index handles the GET /stages endpoint
func (s *StageActions) index(ctx *gin.Context) {
	// Get namespace from context
	namespace, exists := ctx.Get("namespace")
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Namespace not found"})
		return
	}

	// List stages
	result, err := s.stageService.ListStages(ctx, namespace.(model.Namespace), s.valueStream)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list stages"})
		return
	}

	// Configure stages
	configuration, err := s.configureStages(result.Stages)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to configure stages"})
		return
	}

	ctx.JSON(http.StatusOK, configuration)
}

// median handles the GET /stages/:id/median endpoint
func (s *StageActions) median(ctx *gin.Context) {
	stage, exists := ctx.Get("stage")
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Stage not found"})
		return
	}

	// Get median value
	median, err := s.dataCollector.GetMedian(ctx, stage.(*model.Stage))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get median"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"value": median.Seconds()})
}

// average handles the GET /stages/:id/average endpoint
func (s *StageActions) average(ctx *gin.Context) {
	stage, exists := ctx.Get("stage")
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Stage not found"})
		return
	}

	// Get average value
	average, err := s.dataCollector.GetAverage(ctx, stage.(*model.Stage))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get average"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"value": average.Seconds()})
}

// records handles the GET /stages/:id/records endpoint
func (s *StageActions) records(ctx *gin.Context) {
	stage, exists := ctx.Get("stage")
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Stage not found"})
		return
	}

	// Get records with pagination
	records, err := s.dataCollector.GetRecords(ctx, stage.(*model.Stage))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get records"})
		return
	}

	// Add pagination headers
	s.addPaginationHeaders(ctx, records)

	ctx.JSON(http.StatusOK, records)
}

// count handles the GET /stages/:id/count endpoint
func (s *StageActions) count(ctx *gin.Context) {
	stage, exists := ctx.Get("stage")
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Stage not found"})
		return
	}

	// Get count
	count, err := s.dataCollector.GetCount(ctx, stage.(*model.Stage))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get count"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"count": count})
}

// validateParams middleware validates request parameters
func (s *StageActions) validateParams(ctx *gin.Context) {
	if err := s.paramsValidator.Validate(ctx); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}
}

// authorizeStage middleware checks stage authorization
func (s *StageActions) authorizeStage(ctx *gin.Context) {
	// This should be implemented by child handlers
	// It will check if the user has permission to access the stage
}

// addPaginationHeaders adds pagination headers to the response
func (s *StageActions) addPaginationHeaders(ctx *gin.Context, records *model.PaginatedRecords) {
	// Add pagination headers
	ctx.Header("X-Total", records.Total)
	ctx.Header("X-Page", records.Page)
	ctx.Header("X-Per-Page", records.PerPage)
	ctx.Header("X-Next-Page", records.NextPage)
	ctx.Header("X-Prev-Page", records.PrevPage)
}

// configureStages configures stages for the response
func (s *StageActions) configureStages(stages []*model.Stage) (*model.Configuration, error) {
	// Create stage presenters
	presenters := make([]*model.StagePresenter, len(stages))
	for i, stage := range stages {
		presenters[i] = model.NewStagePresenter(stage)
	}

	// Create configuration
	return model.NewConfiguration(presenters), nil
}
