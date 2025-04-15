package organizations

import (
	"net/http"

	"github.com/cyrex562/gitlab-redux/app/services"
	"github.com/gin-gonic/gin"
)

// GroupsHandler handles organization group operations
type GroupsHandler struct {
	*BaseHandler
	groupService *services.GroupService
}

// NewGroupsHandler creates a new GroupsHandler
func NewGroupsHandler(baseHandler *BaseHandler, groupService *services.GroupService) *GroupsHandler {
	return &GroupsHandler{
		BaseHandler:  baseHandler,
		groupService: groupService,
	}
}

// New renders the new group page
func (h *GroupsHandler) New(c *gin.Context) {
	// This would typically render a template
	// For now, just return a success response
	c.JSON(http.StatusOK, gin.H{
		"message": "New group page",
	})
}

// Edit renders the edit group page
func (h *GroupsHandler) Edit(c *gin.Context) {
	// Get the group
	group, err := h.getGroup(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Group not found",
		})
		return
	}

	// Check if user can view edit page
	if !h.groupService.CanViewEditPage(h.GetCurrentUser(c).ID, group.ID) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Cannot view edit page",
		})
		return
	}

	// This would typically render a template
	// For now, just return the group data
	c.JSON(http.StatusOK, group)
}

// Create creates a new group in the organization
func (h *GroupsHandler) Create(c *gin.Context) {
	// Get the organization
	organization, err := h.GetOrganization(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Organization not found",
		})
		return
	}

	// Parse group parameters
	var params services.GroupParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid group parameters",
		})
		return
	}

	// Set organization ID
	params.OrganizationID = organization.ID

	// Create the group
	group, err := h.groupService.CreateGroup(h.GetCurrentUser(c).ID, params)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": err.Error(),
		})
		return
	}

	// Return the created group
	c.JSON(http.StatusOK, group)
}

// Destroy deletes a group
func (h *GroupsHandler) Destroy(c *gin.Context) {
	// Get the group
	group, err := h.getGroup(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Group not found",
		})
		return
	}

	// Check if user can remove group
	if !h.groupService.CanRemoveGroup(h.GetCurrentUser(c).ID, group.ID) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Cannot remove group",
		})
		return
	}

	// Delete the group
	err = h.groupService.DestroyGroup(h.GetCurrentUser(c).ID, group.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": err.Error(),
		})
		return
	}

	// Return success message
	c.JSON(http.StatusOK, gin.H{
		"message": "Group is being deleted",
	})
}

// getGroup gets the group from the path parameter
func (h *GroupsHandler) getGroup(c *gin.Context) (*services.Group, error) {
	// Get the organization
	organization, err := h.GetOrganization(c)
	if err != nil {
		return nil, err
	}

	// Get the group ID from the path
	groupID := c.Param("id")
	if groupID == "" {
		return nil, nil
	}

	// Find the group in the organization
	return h.groupService.FindByID(organization.ID, groupID)
} 