package organizations

import (
	"net/http"

	"github.com/cyrex562/gitlab-redux/app/handlers"
	"github.com/cyrex562/gitlab-redux/app/services"
	"github.com/gin-gonic/gin"
)

// BaseHandler provides common functionality for organization handlers
type BaseHandler struct {
	*handlers.BaseHandler
	organizationService *services.OrganizationService
	featureService     *services.FeatureService
}

// NewBaseHandler creates a new BaseHandler
func NewBaseHandler(baseHandler *handlers.BaseHandler, organizationService *services.OrganizationService, featureService *services.FeatureService) *BaseHandler {
	return &BaseHandler{
		BaseHandler:         baseHandler,
		organizationService: organizationService,
		featureService:     featureService,
	}
}

// CheckFeatureFlag checks if a feature is enabled for the current user
func (h *BaseHandler) CheckFeatureFlag(c *gin.Context, featureName string) bool {
	user := h.GetCurrentUser(c)
	if user == nil {
		return false
	}
	return h.featureService.IsEnabled(featureName, user.ID)
}

// GetOrganization gets the organization from the path parameter
func (h *BaseHandler) GetOrganization(c *gin.Context) (*services.Organization, error) {
	organizationPath := c.Param("organization_path")
	if organizationPath == "" {
		return nil, nil
	}
	return h.organizationService.FindByPath(organizationPath)
}

// RequireFeatureFlag middleware checks if a feature is enabled
func (h *BaseHandler) RequireFeatureFlag(featureName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !h.CheckFeatureFlag(c, featureName) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Feature not enabled",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireOrganizationFeature middleware checks if organizations feature is enabled
func (h *BaseHandler) RequireOrganizationFeature() gin.HandlerFunc {
	return h.RequireFeatureFlag("ui_for_organizations")
}

// RequireOrganizationCreation middleware checks if organization creation is allowed
func (h *BaseHandler) RequireOrganizationCreation() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := h.GetCurrentUser(c)
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}

		if !h.CheckFeatureFlag(c, "allow_organization_creation") {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Organization creation not allowed",
			})
			c.Abort()
			return
		}

		if !h.organizationService.CanCreateOrganization(user.ID) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Cannot create organization",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireOrganizationRead middleware checks if user can read organization
func (h *BaseHandler) RequireOrganizationRead() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := h.GetCurrentUser(c)
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}

		organization, err := h.GetOrganization(c)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Organization not found",
			})
			c.Abort()
			return
		}

		if !h.organizationService.CanReadOrganization(user.ID, organization.ID) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Cannot read organization",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireOrganizationAdmin middleware checks if user can admin organization
func (h *BaseHandler) RequireOrganizationAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := h.GetCurrentUser(c)
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}

		organization, err := h.GetOrganization(c)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Organization not found",
			})
			c.Abort()
			return
		}

		if !h.organizationService.CanAdminOrganization(user.ID, organization.ID) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Cannot admin organization",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireGroupCreation middleware checks if user can create groups in organization
func (h *BaseHandler) RequireGroupCreation() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := h.GetCurrentUser(c)
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}

		organization, err := h.GetOrganization(c)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Organization not found",
			})
			c.Abort()
			return
		}

		if !h.organizationService.CanCreateGroup(user.ID, organization.ID) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Cannot create group",
			})
			c.Abort()
			return
		}

		c.Next()
	}
} 