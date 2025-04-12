package groups

import (
	"net/http"

	"gitlab.com/gitlab-org/gitlab/internal/middleware"
	"gitlab.com/gitlab-org/gitlab/internal/models"
	"gitlab.com/gitlab-org/gitlab/internal/services"
)

// DeployTokensController handles group deploy token operations
type DeployTokensController struct {
	group *models.Group
	user  *models.User
}

// NewDeployTokensController creates a new deploy tokens controller
func NewDeployTokensController(group *models.Group, user *models.User) *DeployTokensController {
	return &DeployTokensController{
		group: group,
		user:  user,
	}
}

// RegisterRoutes registers the deploy token routes
func (c *DeployTokensController) RegisterRoutes(r *http.ServeMux) {
	r.HandleFunc("/groups/"+c.group.Path+"/deploy_tokens/revoke", middleware.Chain(
		c.authorizeDestroyDeployToken,
	)(c.revoke))
}

// authorizeDestroyDeployToken middleware checks if user can destroy deploy tokens
func (c *DeployTokensController) authorizeDestroyDeployToken(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement actual authorization check
		// For now, we'll assume the user has permission if they're authenticated
		if c.user == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

// revoke handles the token revocation request
func (c *DeployTokensController) revoke(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	params := make(map[string]interface{})
	// Add any necessary parameters from the request

	result := services.NewRevokeDeployTokenService(c.group, c.user, params).Execute()
	if !result.Success {
		http.Error(w, "Failed to revoke token", http.StatusInternalServerError)
		return
	}

	// Redirect to the group settings page
	http.Redirect(w, r, "/groups/"+c.group.Path+"/settings/repository#js-deploy-tokens", http.StatusSeeOther)
} 