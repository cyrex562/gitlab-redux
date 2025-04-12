package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// GroupsHandler handles group management requests
type GroupsHandler struct {
	groupService *service.GroupService
}

// NewGroupsHandler creates a new GroupsHandler instance
func NewGroupsHandler(groupService *service.GroupService) *GroupsHandler {
	return &GroupsHandler{
		groupService: groupService,
	}
}

// Index handles the GET request to list all groups
func (h *GroupsHandler) Index(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	sort := c.DefaultQuery("sort", "name")
	name := c.Query("name")

	groups, err := h.groupService.GetGroups(c, page, sort, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch groups"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"groups": groups,
	})
}

// Show handles the GET request to display a specific group
func (h *GroupsHandler) Show(c *gin.Context) {
	groupID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	membersPage, _ := strconv.Atoi(c.DefaultQuery("members_page", "1"))
	projectsPage, _ := strconv.Atoi(c.DefaultQuery("projects_page", "1"))

	group, err := h.groupService.GetGroupWithDetails(c, groupID, membersPage, projectsPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch group details"})
		return
	}

	c.JSON(http.StatusOK, group)
}

// Create handles the POST request to create a new group
func (h *GroupsHandler) Create(c *gin.Context) {
	var group model.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group data"})
		return
	}

	createdGroup, err := h.groupService.CreateGroup(c, &group)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create group"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Group was successfully created",
		"group":   createdGroup,
	})
}

// Update handles the PUT request to update a group
func (h *GroupsHandler) Update(c *gin.Context) {
	groupID := c.Param("id")
	var group model.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group data"})
		return
	}

	updatedGroup, err := h.groupService.UpdateGroup(c, groupID, &group)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update group"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Group was successfully updated",
		"group":   updatedGroup,
	})
}

// Delete handles the DELETE request to remove a group
func (h *GroupsHandler) Delete(c *gin.Context) {
	groupID := c.Param("id")

	err := h.groupService.DeleteGroup(c, groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete group"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Group is being deleted",
	})
}
