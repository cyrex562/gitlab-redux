package storage

import (
	"crypto/subtle"
	"net/http"

	"gitlab.com/gitlab-org/gitlab-redux/internal/handler/application"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/utils/errors"
)

// StaticObjectExternalStorage provides functionality for handling external storage
// for static objects
type StaticObjectExternalStorage struct {
	settingsService    *model.SettingsService
	applicationHelper *application.Helper
}

// NewStaticObjectExternalStorage creates a new instance of StaticObjectExternalStorage
func NewStaticObjectExternalStorage(
	settingsService *model.SettingsService,
	applicationHelper *application.Helper,
) *StaticObjectExternalStorage {
	return &StaticObjectExternalStorage{
		settingsService:    settingsService,
		applicationHelper: applicationHelper,
	}
}

// RedirectToExternalStorage redirects to external storage if the request is not
// already from external storage
func (s *StaticObjectExternalStorage) RedirectToExternalStorage(w http.ResponseWriter, r *http.Request, project *model.Project) error {
	if s.IsExternalStorageRequest(r) {
		return nil
	}

	externalURL := s.applicationHelper.ExternalStorageURLOrPath(r.URL.Path, project)
	http.Redirect(w, r, externalURL, http.StatusTemporaryRedirect)
	return nil
}

// IsExternalStorageRequest checks if the request is from external storage by validating
// the X-Gitlab-External-Storage-Token header
func (s *StaticObjectExternalStorage) IsExternalStorageRequest(r *http.Request) bool {
	headerToken := r.Header.Get("X-Gitlab-External-Storage-Token")
	if headerToken == "" {
		return false
	}

	externalStorageToken := s.settingsService.StaticObjectsExternalStorageAuthToken()
	if externalStorageToken == "" {
		return false
	}

	// Use constant-time comparison to prevent timing attacks
	if subtle.ConstantTimeCompare([]byte(headerToken), []byte(externalStorageToken)) == 1 {
		return true
	}

	// If token is present but invalid, raise access denied error
	panic(errors.NewAccessDeniedError("invalid external storage token"))
}
