package sourcegraph

import (
	"net/http"

	"gitlab.com/gitlab-org/gitlab-redux/internal/handler/auth"
	"gitlab.com/gitlab-org/gitlab-redux/internal/handler/project"
	"gitlab.com/gitlab-org/gitlab-redux/internal/handler/security"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/utils"
)

// SourcegraphDecorator provides functionality for Sourcegraph integration
type SourcegraphDecorator struct {
	settingsService *model.SettingsService
	authService     *auth.AuthService
	projectService  *project.ProjectService
	securityService *security.SecurityService
	gonService      *model.GonService
}

// NewSourcegraphDecorator creates a new instance of SourcegraphDecorator
func NewSourcegraphDecorator(
	settingsService *model.SettingsService,
	authService *auth.AuthService,
	projectService *project.ProjectService,
	securityService *security.SecurityService,
	gonService *model.GonService,
) *SourcegraphDecorator {
	return &SourcegraphDecorator{
		settingsService: settingsService,
		authService:     authService,
		projectService:  projectService,
		securityService: securityService,
		gonService:      gonService,
	}
}

// RegisterMiddleware registers the middleware for Sourcegraph integration
func (s *SourcegraphDecorator) RegisterMiddleware(router http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only apply for HTML requests
		if !isHTMLRequest(r) {
			router.ServeHTTP(w, r)
			return
		}

		// Push Sourcegraph data to gon
		s.pushSourcegraphGon()

		// Update CSP if needed
		s.updateContentSecurityPolicy(w, r)

		router.ServeHTTP(w, r)
	})
}

// pushSourcegraphGon pushes Sourcegraph data to gon
func (s *SourcegraphDecorator) pushSourcegraphGon() {
	if !s.sourcegraphEnabled() {
		return
	}

	s.gonService.Push("sourcegraph", map[string]interface{}{
		"url": s.settingsService.SourcegraphURL(),
	})
}

// sourcegraphEnabled checks if Sourcegraph is enabled
func (s *SourcegraphDecorator) sourcegraphEnabled() bool {
	return s.settingsService.SourcegraphEnabled() &&
		s.sourcegraphEnabledForProject() &&
		s.authService.CurrentUser() != nil &&
		s.authService.CurrentUser().SourcegraphEnabled
}

// sourcegraphEnabledForProject checks if Sourcegraph is enabled for the current project
func (s *SourcegraphDecorator) sourcegraphEnabledForProject() bool {
	currentProject := s.projectService.CurrentProject()
	if currentProject == nil {
		return false
	}

	if s.settingsService.SourcegraphPublicOnly() {
		return currentProject.IsPublic()
	}

	return true
}

// updateContentSecurityPolicy updates the Content Security Policy for Sourcegraph
func (s *SourcegraphDecorator) updateContentSecurityPolicy(w http.ResponseWriter, r *http.Request) {
	if !s.settingsService.SourcegraphEnabled() {
		return
	}

	csp := s.securityService.GetContentSecurityPolicy(r)
	if csp == nil || len(csp.Directives) == 0 {
		return
	}

	defaultConnectSrc := csp.Directives["connect-src"]
	if defaultConnectSrc == "" {
		defaultConnectSrc = csp.Directives["default-src"]
	}

	// Add Sourcegraph API URL to connect-src
	sourcegraphAPIURL := utils.AppendPath(s.settingsService.SourcegraphURL(), ".api/")
	connectSrcValues := append([]string{defaultConnectSrc}, sourcegraphAPIURL)

	csp.Directives["connect-src"] = connectSrcValues
	s.securityService.SetContentSecurityPolicy(w, r, csp)
}

// isHTMLRequest checks if the request is for HTML content
func isHTMLRequest(r *http.Request) bool {
	accept := r.Header.Get("Accept")
	return accept == "" || accept == "*/*" || accept == "text/html" || accept == "text/html,application/xhtml+xml"
}
