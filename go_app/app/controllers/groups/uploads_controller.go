package groups

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/models"
	"gitlab.com/gitlab-org/gitlab-redux/internal/services/uploads"
	"gitlab.com/gitlab-org/gitlab-redux/internal/services/workhorse"
)

// UploadsController handles requests for group uploads
type UploadsController struct {
	uploadsService  *uploads.Service
	workhorseService *workhorse.Service
}

// NewUploadsController creates a new uploads controller
func NewUploadsController(
	uploadsService *uploads.Service,
	workhorseService *workhorse.Service,
) *UploadsController {
	return &UploadsController{
		uploadsService:  uploadsService,
		workhorseService: workhorseService,
	}
}

// RegisterRoutes registers the routes for the uploads controller
func (c *UploadsController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/uploads", c.authorizeUploadFile(), c.Create)
	router.POST("/uploads/authorize", c.authorizeUploadFile(), c.verifyWorkhorseAPI(), c.Authorize)
	router.GET("/uploads/:id", c.loadGroupIfNeeded(), c.disallowNewUploads(), c.Show)
}

// Create handles POST requests for creating uploads
func (c *UploadsController) Create(ctx *gin.Context) {
	// Implementation will be added when needed
	ctx.Status(http.StatusOK)
}

// Authorize handles POST requests for authorizing uploads
func (c *UploadsController) Authorize(ctx *gin.Context) {
	// Implementation will be added when needed
	ctx.Status(http.StatusOK)
}

// Show handles GET requests for showing uploads
func (c *UploadsController) Show(ctx *gin.Context) {
	// Implementation will be added when needed
	ctx.Status(http.StatusOK)
}

// Helper methods for middleware and authorization

func (c *UploadsController) loadGroupIfNeeded() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Skip loading group if the action is 'show' and the upload is embeddable
		if ctx.Request.Method == "GET" && c.isEmbeddable(ctx) {
			ctx.Next()
			return
		}

		// Load the group
		group := ctx.MustGet("group").(*models.Group)
		ctx.Set("group", group)
		ctx.Next()
	}
}

func (c *UploadsController) authorizeUploadFile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("current_user").(*models.User)
		group := ctx.MustGet("group").(*models.Group)

		if !user.CanUploadFile(group) {
			ctx.Status(http.StatusNotFound)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func (c *UploadsController) verifyWorkhorseAPI() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !c.workhorseService.VerifyAPI(ctx) {
			ctx.Status(http.StatusForbidden)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func (c *UploadsController) disallowNewUploads() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if c.uploadsService.IsVersionAtLeast(uploads.IDBasedUploadPathVersion) {
			ctx.Status(http.StatusNotFound)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func (c *UploadsController) isEmbeddable(ctx *gin.Context) bool {
	// TODO: Implement the actual logic
	return false
} 