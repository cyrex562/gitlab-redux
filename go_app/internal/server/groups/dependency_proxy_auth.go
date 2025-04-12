package groups

import (
	"net/http"

	"gitlab.com/gitlab-org/gitlab/internal/middleware"
	"gitlab.com/gitlab-org/gitlab/internal/models"
)

// DependencyProxyAuthController handles authentication for group dependency proxies
type DependencyProxyAuthController struct {
	group *models.Group
}

// NewDependencyProxyAuthController creates a new dependency proxy auth controller
func NewDependencyProxyAuthController(group *models.Group) *DependencyProxyAuthController {
	return &DependencyProxyAuthController{
		group: group,
	}
}

// RegisterRoutes registers the dependency proxy auth routes
func (c *DependencyProxyAuthController) RegisterRoutes(r *http.ServeMux) {
	handler := middleware.Chain(
		c.verifyDependencyProxyEnabled,
	)(c.authenticate)
	r.HandleFunc("/groups/"+c.group.Path+"/dependency_proxy/auth", handler)
}

// verifyDependencyProxyEnabled middleware checks if dependency proxy is enabled
func (c *DependencyProxyAuthController) verifyDependencyProxyEnabled(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setting := c.getDependencyProxySetting()
		if setting == nil || !setting.Enabled {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		next(w, r)
	}
}

// getDependencyProxySetting gets or creates dependency proxy setting
func (c *DependencyProxyAuthController) getDependencyProxySetting() *models.DependencyProxySetting {
	// TODO: Implement actual database interaction
	setting := c.group.DependencyProxySetting
	if setting == nil {
		setting = &models.DependencyProxySetting{
			GroupID: c.group.ID,
			Enabled: false,
		}
		c.group.DependencyProxySetting = setting
	}
	return setting
}

// authenticate handles the authentication request
func (c *DependencyProxyAuthController) authenticate(w http.ResponseWriter, r *http.Request) {
	// Simply return 200 OK with empty body, matching the Ruby implementation
	w.WriteHeader(http.StatusOK)
} 