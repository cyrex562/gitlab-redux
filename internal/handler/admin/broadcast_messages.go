package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

type BroadcastMessagesController struct {
	broadcastService *service.BroadcastMessageService
}

func NewBroadcastMessagesController(broadcastService *service.BroadcastMessageService) *BroadcastMessagesController {
	return &BroadcastMessagesController{
		broadcastService: broadcastService,
	}
}

// Index displays the list of broadcast messages and a form to create new ones
func (c *BroadcastMessagesController) Index(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	messages, err := c.broadcastService.ListMessages(ctx, page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"messages": messages,
		"new_message": service.BroadcastMessage{}, // Empty message for the form
	})
}

// Show displays a specific broadcast message for editing
func (c *BroadcastMessagesController) Show(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	message, err := c.broadcastService.GetMessage(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": message})
}

// Create creates a new broadcast message
func (c *BroadcastMessagesController) Create(ctx *gin.Context) {
	var message service.BroadcastMessage
	if err := ctx.ShouldBindJSON(&message); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.broadcastService.CreateMessage(ctx, &message); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, message)
}

// Update updates an existing broadcast message
func (c *BroadcastMessagesController) Update(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	var message service.BroadcastMessage
	if err := ctx.ShouldBindJSON(&message); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message.ID = id
	if err := c.broadcastService.UpdateMessage(ctx, &message); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, message)
}

// Delete removes a broadcast message
func (c *BroadcastMessagesController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	if err := c.broadcastService.DeleteMessage(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

// Preview generates a preview of how the broadcast message will look
func (c *BroadcastMessagesController) Preview(ctx *gin.Context) {
	var message service.BroadcastMessage
	if err := ctx.ShouldBindJSON(&message); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	preview, err := c.broadcastService.GeneratePreview(ctx, &message)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.String(http.StatusOK, preview)
}
