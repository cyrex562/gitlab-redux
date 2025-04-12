package cycle_analytics

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// ValueStreamActions provides common functionality for cycle analytics value stream actions
type ValueStreamActions struct {
	valueStreamService *service.ValueStreamService
	authService       *service.AuthService
	namespace         model.Namespace
}

// NewValueStreamActions creates a new instance of ValueStreamActions
func NewValueStreamActions(
	valueStreamService *service.ValueStreamService,
	authService *service.AuthService,
	namespace model.Namespace,
) *ValueStreamActions {
	return &ValueStreamActions{
		valueStreamService: valueStreamService,
		authService:       authService,
		namespace:         namespace,
	}
}

// RegisterRoutes registers the routes for value stream actions
func (v *ValueStreamActions) RegisterRoutes(r *gin.RouterGroup) {
	// Apply middleware
	r.Use(v.authorize)
	r.Use(v.authorizeModification)

	// Register routes
	r.GET("/", v.index)
}

// index handles the GET / endpoint
func (v *ValueStreamActions) index(ctx *gin.Context) {
	// In FOSS, users can only see the default value stream
	valueStreams := []*model.ValueStream{
		model.NewDefaultValueStream(v.namespace),
	}

	// Serialize the response
	serializer := model.NewValueStreamSerializer()
	response, err := serializer.Serialize(valueStreams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize value streams"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// authorize middleware checks if the user has permission to read cycle analytics
func (v *ValueStreamActions) authorize(ctx *gin.Context) {
	if err := v.authService.AuthorizeReadCycleAnalytics(ctx); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to read cycle analytics"})
		ctx.Abort()
		return
	}
}

// authorizeModification middleware checks if the user has permission to modify value streams
// This is a no-op in FOSS, but can be overridden in EE
func (v *ValueStreamActions) authorizeModification(ctx *gin.Context) {
	// No-op in FOSS
}
