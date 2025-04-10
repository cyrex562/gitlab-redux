package groups

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/serializer"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
	"gitlab.com/gitlab-org/gitlab-redux/internal/util"
)

// GroupTree handles rendering group trees with filtering, sorting, and pagination
type GroupTree struct {
	groupService      *service.GroupService
	serializerService *service.SerializerService
	logger            *util.Logger
}

// NewGroupTree creates a new instance of GroupTree
func NewGroupTree(
	groupService *service.GroupService,
	serializerService *service.SerializerService,
	logger *util.Logger,
) *GroupTree {
	return &GroupTree{
		groupService:      groupService,
		serializerService: serializerService,
		logger:            logger,
	}
}

// RenderGroupTree renders a group tree with filtering, sorting, and pagination
func (g *GroupTree) RenderGroupTree(ctx *gin.Context, groups []*model.Group) {
	// Get the sort parameter
	sort := ctx.Query("sort")
	if sort == "" {
		sort = "name_asc" // Default sort
	}

	// Sort the groups
	groups = g.groupService.SortGroupsByAttribute(groups, sort)

	// Apply filtering or parent filtering
	if filter := ctx.Query("filter"); filter != "" {
		// Filter groups and include ancestors
		groups = g.filteredGroupsWithAncestors(ctx, groups, filter)
	} else if parentID := ctx.Query("parent_id"); parentID != "" {
		// Filter by parent ID
		groups = g.groupService.GetGroupsByParentID(ctx, groups, parentID)
	} else {
		// Show only root groups
		groups = g.groupService.GetRootGroups(ctx, groups)
	}

	// Apply pagination
	page := ctx.DefaultQuery("page", "1")
	groups = g.groupService.PaginateGroups(ctx, groups, page)

	// Apply archived filter
	archived := ctx.DefaultQuery("archived", "false")
	groups = g.groupService.WithSelectsForList(ctx, groups, archived == "true")

	// Store the groups in the context
	ctx.Set("groups", groups)

	// Respond based on the format
	switch ctx.GetHeader("Accept") {
	case "application/json":
		// Create a serializer
		groupSerializer := serializer.NewGroupChildSerializer(ctx)

		// Apply pagination to the serializer
		groupSerializer.WithPagination(ctx.Request, ctx.Writer)

		// Expand hierarchy if filter is present
		if ctx.Query("filter") != "" {
			groupSerializer.ExpandHierarchy()
		}

		// Render JSON response
		ctx.JSON(http.StatusOK, groupSerializer.Represent(groups))
	default:
		// Render HTML response
		ctx.HTML(http.StatusOK, "groups/index", gin.H{
			"groups": groups,
		})
	}
}

// filteredGroupsWithAncestors filters groups and includes their ancestors
func (g *GroupTree) filteredGroupsWithAncestors(ctx *gin.Context, groups []*model.Group, filter string) []*model.Group {
	// Apply search filter
	filteredGroups := g.groupService.SearchGroups(ctx, groups, filter)

	// Apply pagination
	page := ctx.DefaultQuery("page", "1")
	filteredGroups = g.groupService.PaginateGroups(ctx, filteredGroups, page)

	// Get the IDs of the filtered groups
	filteredGroupIDs := make([]string, len(filteredGroups))
	for i, group := range filteredGroups {
		filteredGroupIDs[i] = group.ID
	}

	// Get the ancestors of the filtered groups
	return g.groupService.GetGroupsWithAncestors(ctx, filteredGroupIDs)
}
