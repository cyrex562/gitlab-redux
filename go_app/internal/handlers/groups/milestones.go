package groups

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/cyrex562/gitlab-redux/internal/middleware"
	"github.com/cyrex562/gitlab-redux/internal/models"
	"github.com/cyrex562/gitlab-redux/internal/services"
)

// MilestonesHandler handles group milestone operations
type MilestonesHandler struct {
	*BaseHandler
	milestoneService *services.MilestoneService
}

// NewMilestonesHandler creates a new milestones handler
func NewMilestonesHandler(baseHandler *BaseHandler, milestoneService *services.MilestoneService) *MilestonesHandler {
	return &MilestonesHandler{
		BaseHandler:      baseHandler,
		milestoneService: milestoneService,
	}
}

// RegisterRoutes registers the routes for the milestones handler
func (h *MilestonesHandler) RegisterRoutes(router *middleware.Router) {
	router.HandleFunc("/groups/{group_id}/milestones", h.Index).Methods("GET")
	router.HandleFunc("/groups/{group_id}/milestones/new", h.New).Methods("GET")
	router.HandleFunc("/groups/{group_id}/milestones", h.Create).Methods("POST")
	router.HandleFunc("/groups/{group_id}/milestones/{id}", h.Show).Methods("GET")
	router.HandleFunc("/groups/{group_id}/milestones/{id}/edit", h.Edit).Methods("GET")
	router.HandleFunc("/groups/{group_id}/milestones/{id}", h.Update).Methods("PUT")
	router.HandleFunc("/groups/{group_id}/milestones/{id}", h.Destroy).Methods("DELETE")
}

// Index handles listing milestones
func (h *MilestonesHandler) Index(w http.ResponseWriter, r *http.Request) {
	group := h.GetGroupFromContext(r)
	if group == nil {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	state := r.URL.Query().Get("state")
	searchTitle := r.URL.Query().Get("search_title")
	sort := r.URL.Query().Get("sort")
	if sort == "" {
		sort = "due_date_asc"
	}

	// Get group IDs based on format
	includeAncestors := r.Header.Get("Accept") == "application/json"
	groupIDs, err := h.getGroupIDs(r, includeAncestors)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get project IDs with access
	projectIDs, err := h.getProjectsWithAccess(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get milestones
	milestones, err := h.milestoneService.FindMilestones(r.Context(), &services.MilestoneFindOptions{
		GroupIDs:   groupIDs,
		ProjectIDs: projectIDs,
		State:      state,
		Search:     searchTitle,
		Sort:       sort,
		Page:       page,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get milestone states count
	statesCount, err := h.milestoneService.GetMilestoneStatesCount(r.Context(), projectIDs, groupIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Handle different response formats
	switch r.Header.Get("Accept") {
	case "application/json":
		// For JSON, only return minimal milestone data
		type milestoneJSON struct {
			ID      int64     `json:"id"`
			Title   string    `json:"title"`
			DueDate time.Time `json:"due_date"`
			Name    string    `json:"name"`
		}
		jsonMilestones := make([]milestoneJSON, len(milestones))
		for i, m := range milestones {
			jsonMilestones[i] = milestoneJSON{
				ID:      m.ID,
				Title:   m.Title,
				DueDate: m.DueDate,
				Name:    m.Name,
			}
		}
		h.respondWithJSON(w, jsonMilestones)
	default:
		h.respondWithHTML(w, "milestones/index", map[string]interface{}{
			"milestones":     milestones,
			"milestone_states": statesCount,
		})
	}
}

// New handles the new milestone form
func (h *MilestonesHandler) New(w http.ResponseWriter, r *http.Request) {
	group := h.GetGroupFromContext(r)
	if group == nil {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	if !h.canAdminMilestones(r) {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	h.respondWithHTML(w, "milestones/new", map[string]interface{}{
		"milestone": &models.Milestone{},
	})
}

// Create handles milestone creation
func (h *MilestonesHandler) Create(w http.ResponseWriter, r *http.Request) {
	group := h.GetGroupFromContext(r)
	if group == nil {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	if !h.canAdminMilestones(r) {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	// Parse milestone parameters
	var params struct {
		Description string    `json:"description"`
		DueDate     time.Time `json:"due_date"`
		LockVersion int       `json:"lock_version"`
		StartDate   time.Time `json:"start_date"`
		StateEvent  string    `json:"state_event"`
		Title       string    `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create milestone
	milestone, err := h.milestoneService.CreateMilestone(r.Context(), &services.MilestoneCreateOptions{
		GroupID:     group.ID,
		Description: params.Description,
		DueDate:     params.DueDate,
		LockVersion: params.LockVersion,
		StartDate:   params.StartDate,
		StateEvent:  params.StateEvent,
		Title:       params.Title,
	})
	if err != nil {
		h.respondWithHTML(w, "milestones/new", map[string]interface{}{
			"milestone": milestone,
			"error":     err.Error(),
		})
		return
	}

	http.Redirect(w, r, "/milestones/"+strconv.FormatInt(milestone.ID, 10), http.StatusSeeOther)
}

// Show handles displaying a milestone
func (h *MilestonesHandler) Show(w http.ResponseWriter, r *http.Request) {
	milestone, err := h.getMilestone(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	h.respondWithHTML(w, "milestones/show", map[string]interface{}{
		"milestone": milestone,
	})
}

// Edit handles the edit milestone form
func (h *MilestonesHandler) Edit(w http.ResponseWriter, r *http.Request) {
	milestone, err := h.getMilestone(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if !h.canAdminMilestones(r) {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	h.respondWithHTML(w, "milestones/edit", map[string]interface{}{
		"milestone": milestone,
	})
}

// Update handles milestone updates
func (h *MilestonesHandler) Update(w http.ResponseWriter, r *http.Request) {
	milestone, err := h.getMilestone(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if !h.canAdminMilestones(r) {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	// Parse milestone parameters
	var params struct {
		Description string    `json:"description"`
		DueDate     time.Time `json:"due_date"`
		LockVersion int       `json:"lock_version"`
		StartDate   time.Time `json:"start_date"`
		StateEvent  string    `json:"state_event"`
		Title       string    `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update milestone
	updatedMilestone, err := h.milestoneService.UpdateMilestone(r.Context(), milestone.ID, &services.MilestoneUpdateOptions{
		Description: params.Description,
		DueDate:     params.DueDate,
		LockVersion: params.LockVersion,
		StartDate:   params.StartDate,
		StateEvent:  params.StateEvent,
		Title:       params.Title,
	})
	if err != nil {
		switch err.(type) {
		case *services.StaleObjectError:
			if r.Header.Get("Accept") == "application/json" {
				http.Error(w, "Someone edited this milestone at the same time you did. Please refresh your browser and make sure your changes will not unintentionally remove theirs.", http.StatusConflict)
				return
			}
			h.respondWithHTML(w, "milestones/edit", map[string]interface{}{
				"milestone": updatedMilestone,
				"conflict":  true,
			})
			return
		default:
			if r.Header.Get("Accept") == "application/json" {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
			h.respondWithHTML(w, "milestones/edit", map[string]interface{}{
				"milestone": updatedMilestone,
				"error":     err.Error(),
			})
			return
		}
	}

	// Handle different response formats
	switch r.Header.Get("Accept") {
	case "application/json":
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Redirect(w, r, "/milestones/"+strconv.FormatInt(updatedMilestone.ID, 10), http.StatusSeeOther)
	}
}

// Destroy handles milestone deletion
func (h *MilestonesHandler) Destroy(w http.ResponseWriter, r *http.Request) {
	milestone, err := h.getMilestone(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if !h.canAdminMilestones(r) {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if err := h.milestoneService.DeleteMilestone(r.Context(), milestone.ID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Handle different response formats
	switch r.Header.Get("Accept") {
	case "application/javascript":
		w.WriteHeader(http.StatusOK)
	default:
		http.Redirect(w, r, "/groups/"+strconv.FormatInt(h.GetGroupFromContext(r).ID, 10)+"/milestones", http.StatusSeeOther)
	}
}

// Helper methods

func (h *MilestonesHandler) canAdminMilestones(r *http.Request) bool {
	user := h.GetUserFromContext(r)
	return h.authorizer.CanAdminMilestone(user, h.GetGroupFromContext(r))
}

func (h *MilestonesHandler) getMilestone(r *http.Request) (*models.Milestone, error) {
	group := h.GetGroupFromContext(r)
	if group == nil {
		return nil, fmt.Errorf("group not found")
	}

	milestoneID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid milestone ID")
	}

	milestone, err := h.milestoneService.GetMilestone(r.Context(), milestoneID)
	if err != nil {
		return nil, fmt.Errorf("milestone not found")
	}

	return milestone, nil
}

func (h *MilestonesHandler) getGroupIDs(r *http.Request, includeAncestors bool) ([]int64, error) {
	group := h.GetGroupFromContext(r)
	if group == nil {
		return nil, fmt.Errorf("group not found")
	}

	if includeAncestors {
		return h.milestoneService.GetGroupHierarchyIDs(r.Context(), group.ID)
	}
	return h.milestoneService.GetGroupDescendantIDs(r.Context(), group.ID)
}

func (h *MilestonesHandler) getProjectsWithAccess(r *http.Request) ([]int64, error) {
	group := h.GetGroupFromContext(r)
	if group == nil {
		return nil, fmt.Errorf("group not found")
	}

	user := h.GetUserFromContext(r)
	return h.milestoneService.GetProjectsWithAccess(r.Context(), group.ID, user.ID)
} 