package pagination

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmadden/gitlab-redux/internal/model"
)

// GetPaginationParams gets the pagination parameters from the context
func GetPaginationParams(ctx *gin.Context) (page, perPage int) {
	pageStr := ctx.DefaultQuery("page", "1")
	perPageStr := ctx.DefaultQuery("per_page", "20")

	page, _ = strconv.Atoi(pageStr)
	perPage, _ = strconv.Atoi(perPageStr)

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}

	return page, perPage
}

// ShouldRedirectOutOfRange checks if we need to redirect out of range
func ShouldRedirectOutOfRange(ctx *gin.Context, paginated interface{}) bool {
	var currentPage, totalPages int

	switch p := paginated.(type) {
	case *model.PaginatedSnippets:
		currentPage = p.CurrentPage
		totalPages = p.TotalPages
	default:
		return false
	}

	if currentPage > totalPages && totalPages > 0 {
		// Redirect to the last page
		ctx.Redirect(http.StatusFound, ctx.Request.URL.Path+"?page="+strconv.Itoa(totalPages))
		return true
	}

	return false
}
