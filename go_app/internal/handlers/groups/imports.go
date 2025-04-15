package groups

import (
	"fmt"
	"net/http"

	"github.com/cyrex562/gitlab-redux/internal/middleware"
)

// ImportsHandler handles group import operations
type ImportsHandler struct {
	*BaseHandler
}

// NewImportsHandler creates a new imports handler
func NewImportsHandler(baseHandler *BaseHandler) *ImportsHandler {
	return &ImportsHandler{
		BaseHandler: baseHandler,
	}
}

// RegisterRoutes registers the routes for the imports handler
func (h *ImportsHandler) RegisterRoutes(router *middleware.Router) {
	router.HandleFunc("/groups/:group_id/import", h.Show).Methods("GET")
}

// Show handles showing the import status
func (h *ImportsHandler) Show(w http.ResponseWriter, r *http.Request) {
	group := h.GetGroupFromContext(r)
	if group == nil {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	// Get continue parameters
	continueParams := h.getContinueParams(r)

	// Check import state
	if group.ImportState == nil || group.ImportState.IsFinished() {
		// Import is finished or not started
		if continueParams.To != "" {
			http.Redirect(w, r, continueParams.To, http.StatusFound)
			h.setFlashNotice(w, continueParams.Notice)
		} else {
			http.Redirect(w, r, "/groups/"+group.Path, http.StatusFound)
			h.setFlashNotice(w, "The group was successfully imported.")
		}
		return
	}

	if group.ImportState.IsFailed() {
		// Import failed
		http.Redirect(w, r, "/groups/"+group.Path+"/new", http.StatusFound)
		h.setFlashAlert(w, fmt.Sprintf("Failed to import group: %s", group.ImportState.LastError))
		return
	}

	// Import is in progress
	if continueParams.NoticeNow != "" {
		h.setFlashNoticeNow(w, continueParams.NoticeNow)
	}
}

// ContinueParams represents parameters for continuing after import
type ContinueParams struct {
	To        string
	Notice    string
	NoticeNow string
}

// getContinueParams gets the continue parameters from the request
func (h *ImportsHandler) getContinueParams(r *http.Request) ContinueParams {
	return ContinueParams{
		To:        r.URL.Query().Get("to"),
		Notice:    r.URL.Query().Get("notice"),
		NoticeNow: r.URL.Query().Get("notice_now"),
	}
}

// setFlashNotice sets a flash notice
func (h *ImportsHandler) setFlashNotice(w http.ResponseWriter, notice string) {
	if notice != "" {
		// TODO: Implement flash notice setting
	}
}

// setFlashNoticeNow sets a flash notice for the current request
func (h *ImportsHandler) setFlashNoticeNow(w http.ResponseWriter, notice string) {
	if notice != "" {
		// TODO: Implement flash notice now setting
	}
}

// setFlashAlert sets a flash alert
func (h *ImportsHandler) setFlashAlert(w http.ResponseWriter, alert string) {
	if alert != "" {
		// TODO: Implement flash alert setting
	}
} 