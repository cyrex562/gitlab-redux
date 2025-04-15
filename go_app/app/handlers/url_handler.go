package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/gitlab-org/gitlab-redux/app/services"
)

// URLHandler handles URL validation for imports
type URLHandler struct {
	*BaseHandler
	urlService *services.URLService
}

// NewURLHandler creates a new URLHandler
func NewURLHandler(baseHandler *BaseHandler, urlService *services.URLService) *URLHandler {
	return &URLHandler{
		BaseHandler: baseHandler,
		urlService:  urlService,
	}
}

// RegisterRoutes registers the routes for the URL handler
func (h *URLHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/import/url/validate", h.Validate).Methods("POST")
}

// Validate validates a remote Git endpoint
func (h *URLHandler) Validate(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var params struct {
		User     string `json:"user"`
		Password string `json:"password"`
		URL      string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the remote Git endpoint
	result := h.urlService.ValidateRemoteGitEndpoint(params.URL, params.User, params.Password)

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	if result.Success {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
		})
	} else {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": result.Message,
		})
	}
} 