package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/models"
	"gitlab.com/gitlab-org/gitlab-redux/internal/routing"
)

// AvatarsController handles topic avatar operations in the admin panel
type AvatarsController struct {
	*routing.BaseController
}

// NewAvatarsController creates a new instance of AvatarsController
func NewAvatarsController() *AvatarsController {
	return &AvatarsController{
		BaseController: routing.NewBaseController(),
	}
}

// SetupRoutes configures the routes for the avatars controller
func (c *AvatarsController) SetupRoutes(router *gin.Engine) {
	admin := router.Group("/admin/topics")
	{
		admin.DELETE("/:topic_id/avatar", c.Destroy)
	}
}

// Destroy removes the avatar from a topic
func (c *AvatarsController) Destroy(ctx *gin.Context) {
	topicID, err := strconv.ParseUint(ctx.Param("topic_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
		return
	}

	topic, err := models.GetTopicByID(uint(topicID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Topic not found"})
		return
	}

	if err := topic.RemoveAvatar(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := topic.Save(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Redirect(http.StatusFound, "/admin/topics/"+strconv.FormatUint(topicID, 10)+"/edit")
}
