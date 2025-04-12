package sorting

import (
	"github.com/gin-gonic/gin"
)

// GetSortValueName gets the sort value name from the context
func GetSortValueName(ctx *gin.Context) string {
	sort := ctx.DefaultQuery("sort", "")
	if sort == "" {
		sort = ctx.DefaultQuery("sort_by", "")
	}
	return sort
}

// SetSortOrder sets the sort order in the context
func SetSortOrder(ctx *gin.Context, defaultSort string) string {
	sort := GetSortValueName(ctx)
	if sort == "" {
		sort = defaultSort
	}
	ctx.Set("sort", sort)
	return sort
}

// GetSortDirection gets the sort direction from the context
func GetSortDirection(ctx *gin.Context) string {
	direction := ctx.DefaultQuery("direction", "")
	if direction == "" {
		direction = ctx.DefaultQuery("sort_direction", "")
	}
	return direction
}

// SetSortDirection sets the sort direction in the context
func SetSortDirection(ctx *gin.Context, defaultDirection string) string {
	direction := GetSortDirection(ctx)
	if direction == "" {
		direction = defaultDirection
	}
	ctx.Set("direction", direction)
	return direction
}
