package controller

import (
	"encoding/json"
	"net/http"

	"gitlab.com/gitlab-org/gitlab/workhorse/internal/api"
	"gitlab.com/gitlab-org/gitlab/workhorse/internal/helper"
	"gitlab.com/gitlab-org/gitlab/workhorse/internal/service"
)

// LabelsController handles admin label management
type LabelsController struct {
	baseController
	labelsService *service.LabelsService
}

// NewLabelsController creates a new labels controller
func NewLabelsController(apiClient *api.Client) *LabelsController {
	return &LabelsController{
		baseController: baseController{
			apiClient: apiClient,
		},
		labelsService: service.NewLabelsService(apiClient),
	}
}

// Index displays a list of template labels
func (c *LabelsController) Index(w http.ResponseWriter, r *http.Request) {
	page := helper.GetIntParam(r, "page", 1)
	labels, err := c.labelsService.GetTemplateLabels(page)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch labels")
		return
	}

	// TODO: Implement HTML template rendering for labels list
	helper.RespondWithError(w, http.StatusNotImplemented, "HTML response not yet implemented")
}

// Show displays a single label
func (c *LabelsController) Show(w http.ResponseWriter, r *http.Request) {
	labelID := r.URL.Query().Get("id")
	if labelID == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "Missing label ID")
		return
	}

	label, err := c.labelsService.GetLabel(labelID)
	if err != nil {
		helper.RespondWithError(w, http.StatusNotFound, "Label not found")
		return
	}

	// TODO: Implement HTML template rendering for label details
	helper.RespondWithError(w, http.StatusNotImplemented, "HTML response not yet implemented")
}

// New displays the form for creating a new label
func (c *LabelsController) New(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement HTML template rendering for new label form
	helper.RespondWithError(w, http.StatusNotImplemented, "HTML response not yet implemented")
}

// Create creates a new template label
func (c *LabelsController) Create(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Color       string `json:"color"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	label, err := c.labelsService.CreateTemplateLabel(params.Title, params.Description, params.Color)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnprocessableEntity, "Failed to create label")
		return
	}

	http.Redirect(w, r, "/admin/labels", http.StatusFound)
}

// Edit displays the form for editing a label
func (c *LabelsController) Edit(w http.ResponseWriter, r *http.Request) {
	labelID := r.URL.Query().Get("id")
	if labelID == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "Missing label ID")
		return
	}

	label, err := c.labelsService.GetLabel(labelID)
	if err != nil {
		helper.RespondWithError(w, http.StatusNotFound, "Label not found")
		return
	}

	// TODO: Implement HTML template rendering for edit label form
	helper.RespondWithError(w, http.StatusNotImplemented, "HTML response not yet implemented")
}

// Update updates an existing label
func (c *LabelsController) Update(w http.ResponseWriter, r *http.Request) {
	labelID := r.URL.Query().Get("id")
	if labelID == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "Missing label ID")
		return
	}

	var params struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Color       string `json:"color"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	label, err := c.labelsService.UpdateLabel(labelID, params.Title, params.Description, params.Color)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnprocessableEntity, "Failed to update label")
		return
	}

	http.Redirect(w, r, "/admin/labels", http.StatusFound)
}

// Destroy removes a label
func (c *LabelsController) Destroy(w http.ResponseWriter, r *http.Request) {
	labelID := r.URL.Query().Get("id")
	if labelID == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "Missing label ID")
		return
	}

	format := r.URL.Query().Get("format")
	if format == "js" {
		if err := c.labelsService.DeleteLabel(labelID); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	label, err := c.labelsService.GetLabel(labelID)
	if err != nil {
		helper.RespondWithError(w, http.StatusNotFound, "Label not found")
		return
	}

	if err := c.labelsService.DeleteLabel(labelID); err != nil {
		http.Redirect(w, r, "/admin/labels", http.StatusFound)
		return
	}

	http.Redirect(w, r, "/admin/labels", http.StatusFound)
}
