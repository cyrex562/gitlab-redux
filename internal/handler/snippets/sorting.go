package snippets

import (
	"github.com/gin-gonic/gin"
)

// GetSortParam gets the sort parameter from the context
func GetSortParam(ctx *gin.Context) string {
	sort := ctx.DefaultQuery("sort", "")
	if sort == "" {
		sort = ctx.DefaultQuery("sort_by", "")
	}
	return sort
}

// SetSortParam sets the sort parameter in the context
func SetSortParam(ctx *gin.Context, defaultSort string) string {
	sort := GetSortParam(ctx)
	if sort == "" {
		sort = defaultSort
	}
	ctx.Set("sort", sort)
	return sort
}

// GetSortOptions gets the sort options for snippets
func GetSortOptions() map[string]string {
	return map[string]string{
		"created_desc": "Created date, newest first",
		"created_asc":  "Created date, oldest first",
		"updated_desc": "Updated date, newest first",
		"updated_asc":  "Updated date, oldest first",
		"title_asc":    "Title, A to Z",
		"title_desc":   "Title, Z to A",
	}
}
