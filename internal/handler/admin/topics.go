package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// TopicsController handles topic management for GitLab projects
type TopicsController struct {
	topicService *service.TopicService
}

// NewTopicsController creates a new instance of TopicsController
func NewTopicsController(topicService *service.TopicService) *TopicsController {
	return &TopicsController{
		topicService: topicService,
	}
}

// RegisterRoutes registers the routes for the TopicsController
func (c *TopicsController) RegisterRoutes(r *gin.RouterGroup) {
	topics := r.Group("/admin/topics")
	{
		topics.Use(c.requireAdmin)
		topics.GET("/", c.index)
		topics.GET("/new", c.new)
		topics.GET("/:id/edit", c.edit)
		topics.POST("/", c.create)
		topics.PUT("/:id", c.update)
		topics.DELETE("/:id", c.destroy)
		topics.POST("/merge", c.merge)
	}
}

// requireAdmin middleware ensures that only admin users can access these endpoints
func (c *TopicsController) requireAdmin(ctx *gin.Context) {
	user := ctx.MustGet("user")
	if user == nil || !user.IsAdmin() {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		ctx.Abort()
		return
	}
	ctx.Next()
}

// index handles the GET /admin/topics endpoint
func (c *TopicsController) index(ctx *gin.Context) {
	search := ctx.Query("search")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	organizationID := ctx.MustGet("organization_id").(int64)

	topics, err := c.topicService.List(ctx, search, page, organizationID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch topics"})
		return
	}

	// TODO: Implement HTML rendering
	ctx.JSON(http.StatusOK, topics)
}

// new handles the GET /admin/topics/new endpoint
func (c *TopicsController) new(ctx *gin.Context) {
	// TODO: Implement HTML rendering for new topic form
	ctx.JSON(http.StatusOK, gin.H{"topic": &model.Topic{}})
}

// edit handles the GET /admin/topics/:id/edit endpoint
func (c *TopicsController) edit(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
		return
	}

	organizationID := ctx.MustGet("organization_id").(int64)
	topic, err := c.topicService.Get(ctx, id, organizationID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Topic not found"})
		return
	}

	// TODO: Implement HTML rendering for edit topic form
	ctx.JSON(http.StatusOK, gin.H{"topic": topic})
}

// create handles the POST /admin/topics endpoint
func (c *TopicsController) create(ctx *gin.Context) {
	var params model.TopicParams
	if err := ctx.ShouldBind(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic parameters"})
		return
	}

	organizationID := ctx.MustGet("organization_id").(int64)
	params.OrganizationID = organizationID

	topic, err := c.topicService.Create(ctx, &params)
	if err != nil {
		// TODO: Implement HTML rendering for new topic form with errors
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	ctx.Redirect(http.StatusFound, "/admin/topics")
}

// update handles the PUT /admin/topics/:id endpoint
func (c *TopicsController) update(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
		return
	}

	var params model.TopicParams
	if err := ctx.ShouldBind(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic parameters"})
		return
	}

	organizationID := ctx.MustGet("organization_id").(int64)
	topic, err := c.topicService.Update(ctx, id, organizationID, &params)
	if err != nil {
		// TODO: Implement HTML rendering for edit topic form with errors
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	ctx.Redirect(http.StatusFound, "/admin/topics/"+strconv.FormatInt(topic.ID, 10)+"/edit")
}

// destroy handles the DELETE /admin/topics/:id endpoint
func (c *TopicsController) destroy(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
		return
	}

	organizationID := ctx.MustGet("organization_id").(int64)
	if err := c.topicService.Destroy(ctx, id, organizationID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete topic"})
		return
	}

	ctx.Redirect(http.StatusFound, "/admin/topics")
}

// merge handles the POST /admin/topics/merge endpoint
func (c *TopicsController) merge(ctx *gin.Context) {
	var params struct {
		SourceTopicID int64 `json:"source_topic_id" binding:"required"`
		TargetTopicID int64 `json:"target_topic_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid merge parameters"})
		return
	}

	organizationID := ctx.MustGet("organization_id").(int64)
	if err := c.topicService.Merge(ctx, params.SourceTopicID, params.TargetTopicID, organizationID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.Redirect(http.StatusFound, "/admin/topics")
}
