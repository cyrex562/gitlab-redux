package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// BaseController provides common functionality for all controllers
type BaseController struct{}

// HandleError handles errors in a consistent way across controllers
func (c *BaseController) HandleError(ctx *gin.Context, err error) {
	// TODO: Add more sophisticated error handling based on error types
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"error": err.Error(),
	})
}
