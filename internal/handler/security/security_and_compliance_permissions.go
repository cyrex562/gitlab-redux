package security

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// SecurityAndCompliancePermissions ensures security and compliance features are enabled for a project
type SecurityAndCompliancePermissions struct {
	authorizer *service.Authorizer
}

// NewSecurityAndCompliancePermissions creates a new instance of SecurityAndCompliancePermissions
func NewSecurityAndCompliancePermissions(authorizer *service.Authorizer) *SecurityAndCompliancePermissions {
	return &SecurityAndCompliancePermissions{
		authorizer: authorizer,
	}
}

// EnsureSecurityAndComplianceEnabled is a middleware that checks if the current user has permission
// to access security and compliance features for the project
func (s *SecurityAndCompliancePermissions) EnsureSecurityAndComplianceEnabled() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the current user from the context
		currentUser, exists := c.Get("current_user")
		if !exists {
			c.JSON(404, gin.H{"error": "Not Found"})
			c.Abort()
			return
		}

		// Get the project from the context
		project, exists := c.Get("project")
		if !exists {
			c.JSON(404, gin.H{"error": "Not Found"})
			c.Abort()
			return
		}

		// Check if the user has permission to access security and compliance features
		canAccess, err := s.authorizer.CanAccess(currentUser.(*model.User), "access_security_and_compliance", project.(*model.Project))
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal Server Error"})
			c.Abort()
			return
		}

		// If the user doesn't have permission, return a 404 error
		if !canAccess {
			c.JSON(404, gin.H{"error": "Not Found"})
			c.Abort()
			return
		}

		// Continue to the next handler
		c.Next()
	}
}
