package explore

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/jmadden/gitlab-redux/internal/handler/feature_flags"
	"github.com/jmadden/gitlab-redux/internal/model"
	"github.com/jmadden/gitlab-redux/internal/service"
)

const (
	// MinSearchLength is the minimum length for search queries
	MinSearchLength = 3
	// PageLimit is the maximum number of pages
	PageLimit = 50
	// RssEntriesLimit is the maximum number of entries for RSS feeds
	RssEntriesLimit = 20
)

// ProjectsController handles project-related actions
type ProjectsController struct {
	*ApplicationController
	projectService service.ProjectService
	featureFlagService service.FeatureFlagService
	organizationService service.OrganizationService
}

// NewProjectsController creates a new ProjectsController
func NewProjectsController(
	appController *ApplicationController,
	projectService service.ProjectService,
	featureFlagService service.FeatureFlagService,
	organizationService service.OrganizationService,
) *ProjectsController {
	return &ProjectsController{
		ApplicationController: appController,
		projectService:       projectService,
		featureFlagService:   featureFlagService,
		organizationService:  organizationService,
	}
}

// RegisterRoutes registers the routes for the ProjectsController
func (c *ProjectsController) RegisterRoutes(router *gin.Engine) {
	explore := router.Group("/explore")
	{
		projects := explore.Group("/projects")
		{
			projects.GET("", c.Index)
			projects.GET("/trending", c.Trending)
			projects.GET("/starred", c.Starred)
			projects.GET("/topics", c.Topics)
			projects.GET("/topics/:topic_name", c.Topic)
		}
	}
}

// Index handles the index action
func (c *ProjectsController) Index(ctx *gin.Context) {
	// Show alert if search is disabled
	c.showAlertIfSearchIsDisabled(ctx)

	// Load projects
	projects, err := c.loadProjects(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set sorting
	c.setSorting(ctx)

	// Respond based on format
	format := ctx.DefaultQuery("format", "html")
	if format == "json" {
		// Render projects partial as HTML
		html, err := c.renderProjectsPartial(ctx, projects)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"html": html,
		})
	} else {
		// Render the index template
		ctx.HTML(http.StatusOK, "explore/projects/index", gin.H{
			"layout":   c.GetLayout(),
			"projects": projects,
			"sort":     ctx.Query("sort"),
		})
	}
}

// Trending handles the trending action
func (c *ProjectsController) Trending(ctx *gin.Context) {
	// Set trending parameter
	ctx.Set("trending", true)

	// Load projects
	projects, err := c.loadProjects(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set sorting
	c.setSorting(ctx)

	// Respond based on format
	format := ctx.DefaultQuery("format", "html")
	if format == "json" {
		// Render projects partial as HTML
		html, err := c.renderProjectsPartial(ctx, projects)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"html": html,
		})
	} else {
		// Render the trending template
		ctx.HTML(http.StatusOK, "explore/projects/trending", gin.H{
			"layout":   c.GetLayout(),
			"projects": projects,
			"sort":     ctx.Query("sort"),
		})
	}
}

// Starred handles the starred action
func (c *ProjectsController) Starred(ctx *gin.Context) {
	// Load projects
	projects, err := c.loadProjects(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reorder by star count
	projects = c.projectService.ReorderByStarCount(projects)

	// Set sorting
	c.setSorting(ctx)

	// Respond based on format
	format := ctx.DefaultQuery("format", "html")
	if format == "json" {
		// Render projects partial as HTML
		html, err := c.renderProjectsPartial(ctx, projects)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"html": html,
		})
	} else {
		// Render the starred template
		ctx.HTML(http.StatusOK, "explore/projects/starred", gin.H{
			"layout":   c.GetLayout(),
			"projects": projects,
			"sort":     ctx.Query("sort"),
		})
	}
}

// Topics handles the topics action
func (c *ProjectsController) Topics(ctx *gin.Context) {
	// Load project counts
	allUserProjects, allStarredProjects, err := c.loadProjectCounts(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Load topics
	topics, err := c.loadTopics(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Render the topics template
	ctx.HTML(http.StatusOK, "explore/projects/topics", gin.H{
		"layout":            c.GetLayout(),
		"topics":            topics,
		"allUserProjects":   allUserProjects,
		"allStarredProjects": allStarredProjects,
	})
}

// Topic handles the topic action
func (c *ProjectsController) Topic(ctx *gin.Context) {
	// Load topic
	topic, err := c.loadTopic(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if topic == nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	// Set topic parameter
	ctx.Set("topic", topic.Name)

	// Load projects
	projects, err := c.loadProjects(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond based on format
	format := ctx.DefaultQuery("format", "html")
	if format == "atom" {
		// Limit projects for RSS
		projects = c.projectService.LimitProjects(projects, RssEntriesLimit)

		// Render the atom template
		ctx.HTML(http.StatusOK, "explore/projects/topic.atom", gin.H{
			"layout":   "xml",
			"projects": projects,
			"topic":    topic,
		})
	} else {
		// Render the topic template
		ctx.HTML(http.StatusOK, "explore/projects/topic", gin.H{
			"layout":   c.GetLayout(),
			"projects": projects,
			"topic":    topic,
		})
	}
}

// loadProjectCounts loads the project counts
func (c *ProjectsController) loadProjectCounts(ctx *gin.Context) ([]*model.Project, []*model.Project, error) {
	// Get current user
	user := ctx.MustGet("current_user").(*model.User)

	// Find all user projects
	allUserProjects, err := c.projectService.FindProjects(ctx, user, map[string]interface{}{
		"non_public": true,
	})
	if err != nil {
		return nil, nil, err
	}

	// Find all starred projects
	allStarredProjects, err := c.projectService.FindProjects(ctx, user, map[string]interface{}{
		"starred": true,
	})
	if err != nil {
		return nil, nil, err
	}

	return allUserProjects, allStarredProjects, nil
}

// loadProjects loads the projects
func (c *ProjectsController) loadProjects(ctx *gin.Context) ([]*model.Project, error) {
	// Load project counts
	_, _, err := c.loadProjectCounts(ctx)
	if err != nil {
		return nil, err
	}

	// Get current user
	user := ctx.MustGet("current_user").(*model.User)

	// Get pagination parameters
	page := ctx.DefaultQuery("page", "1")
	perPage := ctx.DefaultQuery("per_page", "20")

	// Get finder parameters
	finderParams := map[string]interface{}{
		"minimum_search_length": MinSearchLength,
		"not_aimed_for_deletion": true,
		"current_organization": c.currentOrganization(ctx),
	}

	// Merge with request parameters
	for k, v := range ctx.Request.URL.Query() {
		if len(v) > 0 {
			finderParams[k] = v[0]
		}
	}

	// Find projects
	projects, err := c.projectService.FindProjects(ctx, user, finderParams)
	if err != nil {
		return nil, err
	}

	// Preload associations
	projects = c.preloadAssociations(projects)

	// Paginate projects
	projects = c.projectService.PaginateProjects(projects, page, perPage)

	// Prepare projects for rendering
	projects = c.projectService.PrepareProjectsForRendering(projects)

	return projects, nil
}

// loadTopics loads the topics
func (c *ProjectsController) loadTopics(ctx *gin.Context) ([]*model.Topic, error) {
	// Get pagination parameters
	page := ctx.DefaultQuery("page", "1")
	perPage := ctx.DefaultQuery("per_page", "20")

	// Get search parameter
	search := ctx.DefaultQuery("search", "")

	// Get organization ID
	organizationID := c.currentOrganization(ctx).ID

	// Find topics
	topics, err := c.projectService.FindTopics(ctx, map[string]interface{}{
		"search":          search,
		"organization_id": organizationID,
	})
	if err != nil {
		return nil, err
	}

	// Paginate topics
	topics = c.projectService.PaginateTopics(topics, page, perPage)

	return topics, nil
}

// loadTopic loads the topic
func (c *ProjectsController) loadTopic(ctx *gin.Context) (*model.Topic, error) {
	// Get topic name
	topicName := ctx.Param("topic_name")

	// Check if feature flag is enabled
	if feature_flags.PushFrontendFeatureFlag(ctx, "explore_topics_cleaned_path", nil) {
		// Decode topic name
		decodedTopicName, err := url.QueryUnescape(topicName)
		if err != nil {
			return nil, err
		}
		topicName = decodedTopicName
	}

	// Get organization
	organization := c.currentOrganization(ctx)
	if organization == nil {
		return nil, nil
	}

	// Find topic
	topic, err := c.projectService.FindTopicByName(ctx, organization.ID, topicName)
	if err != nil {
		return nil, err
	}

	return topic, nil
}

// preloadAssociations preloads the associations for the projects
func (c *ProjectsController) preloadAssociations(projects []*model.Project) []*model.Project {
	return c.projectService.PreloadAssociations(projects)
}

// setSorting sets the sorting
func (c *ProjectsController) setSorting(ctx *gin.Context) {
	// Get sort parameter
	sort := ctx.DefaultQuery("sort", c.defaultSortOrder())

	// Set sort parameter
	ctx.Set("sort", sort)
}

// defaultSortOrder returns the default sort order
func (c *ProjectsController) defaultSortOrder() string {
	return "latest_activity"
}

// sortingField returns the sorting field
func (c *ProjectsController) sortingField() string {
	return "project_sorting_preference"
}

// pageOutOfBounds handles the page out of bounds error
func (c *ProjectsController) pageOutOfBounds(ctx *gin.Context, maxPage int) {
	// Load project counts
	allUserProjects, allStarredProjects, err := c.loadProjectCounts(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond based on format
	format := ctx.DefaultQuery("format", "html")
	if format == "json" {
		// Render the page out of bounds partial as HTML
		html, err := c.renderPageOutOfBoundsPartial(ctx, maxPage)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{
			"html": html,
		})
	} else {
		// Render the page out of bounds template
		ctx.HTML(http.StatusBadRequest, "explore/projects/page_out_of_bounds", gin.H{
			"layout":            c.GetLayout(),
			"allUserProjects":   allUserProjects,
			"allStarredProjects": allStarredProjects,
			"maxPageNumber":     maxPage,
		})
	}
}

// showAlertIfSearchIsDisabled shows an alert if search is disabled
func (c *ProjectsController) showAlertIfSearchIsDisabled(ctx *gin.Context) {
	// Get current user
	user := ctx.MustGet("current_user").(*model.User)

	// Get name and search parameters
	name := ctx.DefaultQuery("name", "")
	search := ctx.DefaultQuery("search", "")

	// Check if HTML request
	isHTMLRequest := ctx.DefaultQuery("format", "html") == "html"

	// Check if feature flag is disabled
	featureDisabled := !feature_flags.PushFrontendFeatureFlag(ctx, "disable_anonymous_project_search", nil)

	// Show alert if conditions are met
	if user == nil && (name != "" || search != "") && isHTMLRequest && featureDisabled {
		ctx.Set("flash_notice", "You must sign in to search for specific projects.")
	}
}

// currentOrganization returns the current organization
func (c *ProjectsController) currentOrganization(ctx *gin.Context) *model.Organization {
	return c.organizationService.GetCurrentOrganization(ctx)
}

// renderProjectsPartial renders the projects partial as HTML
func (c *ProjectsController) renderProjectsPartial(ctx *gin.Context, projects []*model.Project) (string, error) {
	// TODO: Implement rendering of projects partial
	return "", nil
}

// renderPageOutOfBoundsPartial renders the page out of bounds partial as HTML
func (c *ProjectsController) renderPageOutOfBoundsPartial(ctx *gin.Context, maxPage int) (string, error) {
	// TODO: Implement rendering of page out of bounds partial
	return "", nil
}
