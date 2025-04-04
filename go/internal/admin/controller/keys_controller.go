package controller

import (
	"net/http"

	"gitlab.com/gitlab-org/gitlab/workhorse/internal/api"
	"gitlab.com/gitlab-org/gitlab/workhorse/internal/helper"
	"gitlab.com/gitlab-org/gitlab/workhorse/internal/service"
)

// KeysController handles admin user key management
type KeysController struct {
	baseController
	keysService *service.KeysService
}

// NewKeysController creates a new keys controller
func NewKeysController(apiClient *api.Client) *KeysController {
	return &KeysController{
		baseController: baseController{
			apiClient: apiClient,
		},
		keysService: service.NewKeysService(apiClient),
	}
}

// Show displays a user's key details
func (c *KeysController) Show(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	keyID := r.URL.Query().Get("id")

	if userID == "" || keyID == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "Missing required parameters")
		return
	}

	key, err := c.keysService.GetUserKey(userID, keyID)
	if err != nil {
		helper.RespondWithError(w, http.StatusNotFound, "Key not found")
		return
	}

	// Handle different response formats
	format := r.URL.Query().Get("format")
	if format == "js" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Default to HTML response
	// TODO: Implement HTML template rendering for key details
	helper.RespondWithError(w, http.StatusNotImplemented, "HTML response not yet implemented")
}

// Destroy removes a user's key
func (c *KeysController) Destroy(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	keyID := r.URL.Query().Get("id")

	if userID == "" || keyID == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "Missing required parameters")
		return
	}

	if err := c.keysService.DeleteUserKey(userID, keyID); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to remove user key")
		return
	}

	// Redirect to user's keys page
	redirectURL := "/admin/users/" + userID + "/keys"
	http.Redirect(w, r, redirectURL, http.StatusFound)
}
