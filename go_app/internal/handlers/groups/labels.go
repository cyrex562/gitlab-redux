package groups

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cyrex562/gitlab-redux/internal/middleware"
	"github.com/cyrex562/gitlab-redux/internal/models"
	"github.com/cyrex562/gitlab-redux/internal/services"
)

// LabelsHandler handles group label operations
type LabelsHandler struct {
	*BaseHandler
	labelService *services.LabelService
}

// NewLabelsHandler creates a new labels handler
func NewLabelsHandler(baseHandler *BaseHandler, labelService *services.LabelService) *LabelsHandler {
	return &LabelsHandler{
		BaseHandler:  baseHandler,
		labelService: labelService,
	}
}

// RegisterRoutes registers the routes for the labels handler
func (h *LabelsHandler) RegisterRoutes(router *middleware.Router) {
	router.HandleFunc("/groups/{group_id}/labels", h.Index).Methods("GET")
	router.HandleFunc("/groups/{group_id}/labels/new", h.New).Methods("GET")
	router.HandleFunc("/groups/{group_id}/labels", h.Create).Methods("POST")
	router.HandleFunc("/groups/{group_id}/labels/{id}/edit", h.Edit).Methods("GET")
	router.HandleFunc("/groups/{group_id}/labels/{id}", h.Update).Methods("PUT")
	router.HandleFunc("/groups/{group_id}/labels/{id}", h.Destroy).Methods("DELETE")
}

// Index handles listing labels
func (h *LabelsHandler) Index(w http.ResponseWriter, r *http.Request) {
	group := h.GetGroupFromContext(r)
	if group == nil {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	if !h.canReadLabels(r) {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	onlyGroupLabels := true // At group level we only want group labels
	search := r.URL.Query().Get("search")
	sort := r.URL.Query().Get("sort")
	if sort == "" {
		sort = "name_asc"
	}

	// Get labels
	labels, err := h.labelService.FindLabels(r.Context(), &services.LabelFindOptions{
		GroupID:         group.ID,
		OnlyGroupLabels: onlyGroupLabels,
		IncludeAncestors: true,
		Sort:            sort,
		Search:          search,
		Page:            page,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Handle different response formats
	switch r.Header.Get("Accept") {
	case "application/json":
		h.respondWithJSON(w, labels)
	default:
		h.respondWithHTML(w, "labels/index", map[string]interface{}{
			"labels": labels,
		})
	}
}

// New handles the new label form
func (h *LabelsHandler) New(w http.ResponseWriter, r *http.Request) {
	group := h.GetGroupFromContext(r)
	if group == nil {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	if !h.canAdminLabels(r) {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	h.respondWithHTML(w, "labels/new", map[string]interface{}{
		"label": &models.Label{},
		"previous_labels_path": h.getPreviousLabelsPath(r),
	})
}

// Create handles label creation
func (h *LabelsHandler) Create(w http.ResponseWriter, r *http.Request) {
	group := h.GetGroupFromContext(r)
	if group == nil {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	if !h.canAdminLabels(r) {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	// Parse label parameters
	var params struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Color       string `json:"color"`
		LockOnMerge bool   `json:"lock_on_merge"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create label
	label, err := h.labelService.CreateLabel(r.Context(), &services.LabelCreateOptions{
		GroupID:     group.ID,
		Title:       params.Title,
		Description: params.Description,
		Color:       params.Color,
		LockOnMerge: params.LockOnMerge,
	})
	if err != nil {
		if r.Header.Get("Accept") == "application/json" {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		h.respondWithHTML(w, "labels/new", map[string]interface{}{
			"label": label,
			"error": err.Error(),
		})
		return
	}

	// Handle different response formats
	switch r.Header.Get("Accept") {
	case "application/json":
		h.respondWithJSON(w, label)
	default:
		http.Redirect(w, r, "/groups/"+strconv.FormatInt(group.ID, 10)+"/labels", http.StatusSeeOther)
	}
}

// Edit handles the edit label form
func (h *LabelsHandler) Edit(w http.ResponseWriter, r *http.Request) {
	group := h.GetGroupFromContext(r)
	if group == nil {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	labelID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid label ID", http.StatusBadRequest)
		return
	}

	label, err := h.labelService.GetLabel(r.Context(), labelID)
	if err != nil {
		http.Error(w, "Label not found", http.StatusNotFound)
		return
	}

	if !h.canAdminLabel(r, label) {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	// Save previous path for redirect after update
	h.savePreviousLabelPath(r)

	h.respondWithHTML(w, "labels/edit", map[string]interface{}{
		"label": label,
		"previous_labels_path": h.getPreviousLabelsPath(r),
	})
}

// Update handles label updates
func (h *LabelsHandler) Update(w http.ResponseWriter, r *http.Request) {
	group := h.GetGroupFromContext(r)
	if group == nil {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	labelID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid label ID", http.StatusBadRequest)
		return
	}

	label, err := h.labelService.GetLabel(r.Context(), labelID)
	if err != nil {
		http.Error(w, "Label not found", http.StatusNotFound)
		return
	}

	if !h.canAdminLabel(r, label) {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	// Parse label parameters
	var params struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Color       string `json:"color"`
		LockOnMerge bool   `json:"lock_on_merge"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update label
	updatedLabel, err := h.labelService.UpdateLabel(r.Context(), labelID, &services.LabelUpdateOptions{
		Title:       params.Title,
		Description: params.Description,
		Color:       params.Color,
		LockOnMerge: params.LockOnMerge,
	})
	if err != nil {
		h.respondWithHTML(w, "labels/edit", map[string]interface{}{
			"label": updatedLabel,
			"error": err.Error(),
		})
		return
	}

	http.Redirect(w, r, h.getPreviousLabelsPath(r), http.StatusSeeOther)
}

// Destroy handles label deletion
func (h *LabelsHandler) Destroy(w http.ResponseWriter, r *http.Request) {
	group := h.GetGroupFromContext(r)
	if group == nil {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	labelID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid label ID", http.StatusBadRequest)
		return
	}

	label, err := h.labelService.GetLabel(r.Context(), labelID)
	if err != nil {
		http.Error(w, "Label not found", http.StatusNotFound)
		return
	}

	if !h.canAdminLabel(r, label) {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if err := h.labelService.DeleteLabel(r.Context(), labelID); err != nil {
		http.Redirect(w, r, "/groups/"+strconv.FormatInt(group.ID, 10)+"/labels", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/groups/"+strconv.FormatInt(group.ID, 10)+"/labels", http.StatusSeeOther)
}

// Helper methods

func (h *LabelsHandler) canReadLabels(r *http.Request) bool {
	user := h.GetUserFromContext(r)
	return h.authorizer.CanReadLabel(user, h.GetGroupFromContext(r))
}

func (h *LabelsHandler) canAdminLabels(r *http.Request) bool {
	user := h.GetUserFromContext(r)
	return h.authorizer.CanAdminLabel(user, h.GetGroupFromContext(r))
}

func (h *LabelsHandler) canAdminLabel(r *http.Request, label *models.Label) bool {
	user := h.GetUserFromContext(r)
	return h.authorizer.CanAdminLabel(user, label)
}

func (h *LabelsHandler) getPreviousLabelsPath(r *http.Request) string {
	if path := r.URL.Query().Get("previous_labels_path"); path != "" {
		return path
	}
	return "/groups/" + strconv.FormatInt(h.GetGroupFromContext(r).ID, 10) + "/labels"
}

func (h *LabelsHandler) savePreviousLabelPath(r *http.Request) {
	if referer := r.Header.Get("Referer"); referer != "" {
		// Store in session or context for later use
		// Implementation depends on your session management
	}
} 