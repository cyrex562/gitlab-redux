package analytics

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
	"gitlab.com/gitlab-org/gitlab-redux/internal/util"
)

// CycleAnalyticsParams handles parameters for cycle analytics
type CycleAnalyticsParams struct {
	analyticsService *service.AnalyticsService
	logger           *util.Logger
}

// NewCycleAnalyticsParams creates a new instance of CycleAnalyticsParams
func NewCycleAnalyticsParams(
	analyticsService *service.AnalyticsService,
	logger *util.Logger,
) *CycleAnalyticsParams {
	return &CycleAnalyticsParams{
		analyticsService: analyticsService,
		logger:           logger,
	}
}

// GetProjectParams gets the cycle analytics project parameters
func (c *CycleAnalyticsParams) GetProjectParams(ctx *gin.Context) map[string]interface{} {
	// Check if the cycle_analytics parameter is present
	cycleAnalytics, exists := ctx.GetQuery("cycle_analytics")
	if !exists || cycleAnalytics == "" {
		return make(map[string]interface{})
	}

	// Create a map for the parameters
	params := make(map[string]interface{})

	// Get the start date
	if startDate, exists := ctx.GetQuery("cycle_analytics[start_date]"); exists {
		params["start_date"] = startDate
	}

	// Get the created after date
	if createdAfter, exists := ctx.GetQuery("cycle_analytics[created_after]"); exists {
		params["created_after"] = createdAfter
	}

	// Get the created before date
	if createdBefore, exists := ctx.GetQuery("cycle_analytics[created_before]"); exists {
		params["created_before"] = createdBefore
	}

	// Get the branch name
	if branchName, exists := ctx.GetQuery("cycle_analytics[branch_name]"); exists {
		params["branch_name"] = branchName
	}

	return params
}

// GetGroupParams gets the cycle analytics group parameters
func (c *CycleAnalyticsParams) GetGroupParams(ctx *gin.Context) map[string]interface{} {
	// Check if the params are present
	if ctx.Request.URL.RawQuery == "" {
		return make(map[string]interface{})
	}

	// Create a map for the parameters
	params := make(map[string]interface{})

	// Get the group ID
	if groupID, exists := ctx.GetQuery("group_id"); exists {
		params["group_id"] = groupID
	}

	// Get the start date
	if startDate, exists := ctx.GetQuery("start_date"); exists {
		params["start_date"] = startDate
	}

	// Get the created after date
	if createdAfter, exists := ctx.GetQuery("created_after"); exists {
		params["created_after"] = createdAfter
	}

	// Get the created before date
	if createdBefore, exists := ctx.GetQuery("created_before"); exists {
		params["created_before"] = createdBefore
	}

	// Get the project IDs
	if projectIDs, exists := ctx.GetQueryArray("project_ids[]"); exists {
		params["project_ids"] = projectIDs
	}

	return params
}

// GetOptions gets the cycle analytics options
func (c *CycleAnalyticsParams) GetOptions(ctx *gin.Context, params map[string]interface{}) map[string]interface{} {
	// Get the current user from the context
	currentUser, exists := ctx.Get("current_user")
	if !exists {
		return make(map[string]interface{})
	}
	user := currentUser.(*model.User)

	// Create a map for the options
	opts := make(map[string]interface{})

	// Set the current user
	opts["current_user"] = user

	// Set the projects if present
	if projectIDs, ok := params["project_ids"]; ok {
		opts["projects"] = projectIDs
	}

	// Set the from date
	if from, ok := params["from"]; ok {
		opts["from"] = from
	} else {
		opts["from"] = c.getStartDate(params)
	}

	// Set the to date if present
	if to, ok := params["to"]; ok {
		opts["to"] = to
	}

	// Set the end event filter if present
	if endEventFilter, ok := params["end_event_filter"]; ok {
		opts["end_event_filter"] = endEventFilter
	}

	// Set the use aggregated data collector if present
	if useAggregatedDataCollector, ok := params["use_aggregated_data_collector"]; ok {
		opts["use_aggregated_data_collector"] = useAggregatedDataCollector
	}

	// Merge the finder parameters
	for _, paramName := range c.analyticsService.GetFinderParamNames() {
		if value, ok := params[paramName]; ok {
			opts[paramName] = value
		}
	}

	// Merge the date range parameters
	dateRange := c.getDateRange(params)
	for key, value := range dateRange {
		opts[key] = value
	}

	return opts
}

// getStartDate gets the start date based on the parameters
func (c *CycleAnalyticsParams) getStartDate(params map[string]interface{}) time.Time {
	// Get the start date from the parameters
	startDate, ok := params["start_date"]
	if !ok {
		return time.Now().AddDate(0, 0, -90)
	}

	// Convert the start date to a string
	startDateStr, ok := startDate.(string)
	if !ok {
		return time.Now().AddDate(0, 0, -90)
	}

	// Parse the start date
	switch startDateStr {
	case "7":
		return time.Now().AddDate(0, 0, -7)
	case "30":
		return time.Now().AddDate(0, 0, -30)
	default:
		return time.Now().AddDate(0, 0, -90)
	}
}

// getDateRange gets the date range based on the parameters
func (c *CycleAnalyticsParams) getDateRange(params map[string]interface{}) map[string]interface{} {
	// Create a map for the date range parameters
	dateRangeParams := make(map[string]interface{})

	// Get the created after date
	if createdAfter, ok := params["created_after"]; ok {
		dateRangeParams["from"] = c.toUTCTime(createdAfter).Truncate(24 * time.Hour)
	}

	// Get the created before date
	if createdBefore, ok := params["created_before"]; ok {
		dateRangeParams["to"] = c.toUTCTime(createdBefore).Add(24*time.Hour - time.Second)
	}

	return dateRangeParams
}

// toUTCTime converts a field to UTC time
func (c *CycleAnalyticsParams) toUTCTime(field interface{}) time.Time {
	// Handle different types
	switch v := field.(type) {
	case time.Time:
		return v.UTC()
	case string:
		// Parse the date
		date, err := time.Parse("2006-01-02", v)
		if err != nil {
			c.logger.Error("failed to parse date", "date", v, "error", err)
			return time.Now().UTC()
		}
		return date.UTC()
	default:
		return time.Now().UTC()
	}
}

// GetPermittedParams gets the permitted cycle analytics parameters
func (c *CycleAnalyticsParams) GetPermittedParams(ctx *gin.Context) map[string]interface{} {
	// Get the permitted parameters from the analytics service
	permittedParams := c.analyticsService.GetPermittedParams(ctx)

	// Create a map for the parameters
	params := make(map[string]interface{})

	// Copy the permitted parameters
	for key, value := range permittedParams {
		params[key] = value
	}

	return params
}

// GetAllParams gets all the cycle analytics parameters
func (c *CycleAnalyticsParams) GetAllParams(ctx *gin.Context) map[string]interface{} {
	// Get the permitted parameters
	params := c.GetPermittedParams(ctx)

	// Get the current user from the context
	currentUser, exists := ctx.Get("current_user")
	if !exists {
		return params
	}
	user := currentUser.(*model.User)

	// Get the namespace from the context
	namespace, exists := ctx.Get("namespace")
	if !exists {
		return params
	}
	namespaceObj := namespace.(*model.Namespace)

	// Set the current user and namespace
	params["current_user"] = user
	params["namespace"] = namespaceObj

	return params
}

// ValidateParams validates the cycle analytics parameters
func (c *CycleAnalyticsParams) ValidateParams(ctx *gin.Context) error {
	// Get all the parameters
	params := c.GetAllParams(ctx)

	// Create a request params object
	requestParams := c.analyticsService.NewRequestParams(params)

	// Validate the parameters
	if !requestParams.IsValid() {
		// Return an error response
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Invalid parameters",
			"errors":  requestParams.GetErrors(),
		})
		return util.NewValidationError("invalid parameters")
	}

	return nil
}
