package crm

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ContactsController handles CRM contacts within groups
type ContactsController struct {
	// Add any dependencies here
}

// NewContactsController creates a new ContactsController
func NewContactsController() *ContactsController {
	return &ContactsController{}
}

// RegisterRoutes registers the routes for the ContactsController
func (c *ContactsController) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/new", c.New)
	router.GET("/edit", c.Edit)
}

// New handles the new action
func (c *ContactsController) New(ctx *gin.Context) {
	// Render the index action
	ctx.HTML(http.StatusOK, "groups/crm/contacts/index", gin.H{
		"action": "new",
	})
}

// Edit handles the edit action
func (c *ContactsController) Edit(ctx *gin.Context) {
	// Render the index action
	ctx.HTML(http.StatusOK, "groups/crm/contacts/index", gin.H{
		"action": "edit",
	})
}

// ValidateCrmGroup middleware validates that the group has CRM enabled
func (c *ContactsController) ValidateCrmGroup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the group from the context
		group, exists := ctx.Get("group")
		if !exists {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
			ctx.Abort()
			return
		}

		// Check if the group has CRM enabled
		// This is a placeholder for the actual implementation
		// In the actual implementation, you would check if the group has CRM enabled
		hasCrmEnabled := true // Replace with actual check

		if !hasCrmEnabled {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "CRM is not enabled for this group"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// AuthorizeReadCrmContact middleware checks if the user has permission to read CRM contacts
func (c *ContactsController) AuthorizeReadCrmContact() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the current user from the context
		currentUser, exists := ctx.Get("current_user")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			ctx.Abort()
			return
		}

		// Get the group from the context
		group, exists := ctx.Get("group")
		if !exists {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
			ctx.Abort()
			return
		}

		// Check if the user has permission to read CRM contacts
		// This is a placeholder for the actual implementation
		// In the actual implementation, you would check if the user has the read_crm_contact permission
		canReadCrmContact := true // Replace with actual check

		if !canReadCrmContact {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "You don't have permission to read CRM contacts"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
