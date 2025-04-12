package pagination

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// PaginationParams represents pagination parameters
type PaginationParams struct {
	Page      int
	PerPage   int
	Limit     int
	Sort      string
	OrderBy   string
	Pagination bool
}

// NewPaginationParams creates a new instance of PaginationParams
func NewPaginationParams() *PaginationParams {
	return &PaginationParams{
		Page:      1,
		PerPage:   20,
		Limit:     0,
		Sort:      "",
		OrderBy:   "",
		Pagination: true,
	}
}

// FromContext extracts pagination parameters from the context
func FromContext(c *gin.Context) *PaginationParams {
	params := NewPaginationParams()

	// Extract page
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			params.Page = page
		}
	}

	// Extract per_page
	if perPageStr := c.Query("per_page"); perPageStr != "" {
		if perPage, err := strconv.Atoi(perPageStr); err == nil && perPage > 0 {
			params.PerPage = perPage
		}
	}

	// Extract limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			params.Limit = limit
		}
	}

	// Extract sort
	if sort := c.Query("sort"); sort != "" {
		params.Sort = sort
	}

	// Extract order_by
	if orderBy := c.Query("order_by"); orderBy != "" {
		params.OrderBy = orderBy
	}

	// Extract pagination
	if paginationStr := c.Query("pagination"); paginationStr != "" {
		if pagination, err := strconv.ParseBool(paginationStr); err == nil {
			params.Pagination = pagination
		}
	}

	return params
}

// ToMap converts pagination parameters to a map
func (p *PaginationParams) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"page":       p.Page,
		"per_page":   p.PerPage,
		"limit":      p.Limit,
		"sort":       p.Sort,
		"order_by":   p.OrderBy,
		"pagination": p.Pagination,
	}
}
