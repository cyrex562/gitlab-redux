package storage

import (
	"net/http"

	"gitlab.com/gitlab-org/gitlab-redux/internal/handler/security"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
)

// StaticObjectExternalStorageCSP provides functionality for Content Security Policy
// for static objects stored in external storage
type StaticObjectExternalStorageCSP struct {
	settingsService *model.SettingsService
	securityService *security.SecurityService
}

// NewStaticObjectExternalStorageCSP creates a new instance of StaticObjectExternalStorageCSP
func NewStaticObjectExternalStorageCSP(
	settingsService *model.SettingsService,
	securityService *security.SecurityService,
) *StaticObjectExternalStorageCSP {
	return &StaticObjectExternalStorageCSP{
		settingsService: settingsService,
		securityService: securityService,
	}
}

// UpdateContentSecurityPolicy updates the Content Security Policy for static objects
// stored in external storage
func (s *StaticObjectExternalStorageCSP) UpdateContentSecurityPolicy(w http.ResponseWriter, r *http.Request) {
	if !s.settingsService.StaticObjectsExternalStorageEnabled() {
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

	// Add external storage URL to connect-src
	externalStorageURL := s.settingsService.StaticObjectsExternalStorageURL()
	connectSrcValues := append([]string{defaultConnectSrc}, externalStorageURL)

	csp.Directives["connect-src"] = connectSrcValues
	s.securityService.SetContentSecurityPolicy(w, r, csp)
}
