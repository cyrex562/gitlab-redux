package todos

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmadden/gitlab-redux/internal/service"
)

// TodosActionsHandler handles HTTP requests for todos
type TodosActionsHandler struct {
	todoService *service.TodoService
}

// NewTodosActionsHandler creates a new handler instance
func NewTodosActionsHandler(todoService *service.TodoService) *TodosActionsHandler {
	return &TodosActionsHandler{
		todoService: todoService,
	}
}

// RegisterRoutes registers the handler routes
func (h *TodosActionsHandler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		todos := api.Group("/todos")
		{
			todos.POST("/:issuableType/:issuableID", h.create)
		}
	}
}

// create handles POST /api/todos/:issuableType/:issuableID
func (h *TodosActionsHandler) create(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	issuableType := c.Param("issuableType")
	issuableID, err := strconv.ParseInt(c.Param("issuableID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid issuable ID"})
		return
	}

	// Create the todo
	todo, err := h.todoService.MarkTodo(c.Request.Context(), issuableType, issuableID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create todo"})
		return
	}

	// Get count of pending todos
	count, err := h.todoService.GetPendingCount(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get todo count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count":      count,
		"delete_path": "/dashboard/todos/" + strconv.FormatInt(todo.ID, 10),
	})
}
