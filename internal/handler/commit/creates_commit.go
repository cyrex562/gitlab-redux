package commit

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
	"gitlab.com/gitlab-org/gitlab-redux/internal/util"
)

// CreatesCommit handles creating commits
type CreatesCommit struct {
	authService *service.AuthService
	userService *service.UserService
	repoService *service.RepositoryService
	mrService   *service.MergeRequestService
	logger      *util.Logger
}

// NewCreatesCommit creates a new instance of CreatesCommit
func NewCreatesCommit(
	authService *service.AuthService,
	userService *service.UserService,
	repoService *service.RepositoryService,
	mrService *service.MergeRequestService,
	logger *util.Logger,
) *CreatesCommit {
	return &CreatesCommit{
		authService: authService,
		userService: userService,
		repoService: repoService,
		mrService:   mrService,
		logger:      logger,
	}
}

// CreateCommit creates a commit
func (c *CreatesCommit) CreateCommit(
	ctx *gin.Context,
	service *service.CommitService,
	successPath string,
	failurePath interface{},
	failureView string,
	successNotice string,
	targetProject *model.Project,
) error {
	// Get the current user from the context
	currentUser, exists := ctx.Get("current_user")
	if !exists {
		return util.NewUnauthorizedError("user not authenticated")
	}
	user := currentUser.(*model.User)

	// Get the project from the context
	project, exists := ctx.Get("project")
	if !exists {
		return util.NewNotFoundError("project not found")
	}
	projectObj := project.(*model.Project)

	// Use the target project if provided, otherwise use the current project
	if targetProject == nil {
		targetProject = projectObj
	}

	// Get the branch name or ref
	branchNameOrRef := c.getBranchNameOrRef(ctx)

	// Check if the user can push to the branch
	userAccess := c.userService.GetUserAccess(ctx, user, targetProject)
	if userAccess.CanPushToBranch(branchNameOrRef) {
		// Set the project to commit into
		projectToCommitInto := targetProject
		differentProject := false
		branchName := branchNameOrRef
		ctx.Set("project_to_commit_into", projectToCommitInto)
		ctx.Set("different_project", differentProject)
		ctx.Set("branch_name", branchName)
	} else {
		// Get the user's fork of the target project
		fork := c.userService.GetForkOf(ctx, user, targetProject)
		if fork == nil {
			return util.NewForbiddenError("you don't have a fork of this project")
		}

		// Set the project to commit into
		projectToCommitInto := fork
		differentProject := true
		branchName := c.repoService.NextBranch(ctx, fork, "patch")
		ctx.Set("project_to_commit_into", projectToCommitInto)
		ctx.Set("different_project", differentProject)
		ctx.Set("branch_name", branchName)
	}

	// Get the start branch
	startBranch := c.getStartBranch(ctx)
	ctx.Set("start_branch", startBranch)

	// Get the commit parameters
	commitParams := c.getCommitParams(ctx)
	commitParams["start_project"] = ctx.Get("project_to_commit_into")
	commitParams["start_branch"] = startBranch
	commitParams["source_project"] = projectObj
	commitParams["target_project"] = targetProject
	commitParams["branch_name"] = ctx.Get("branch_name")

	// Execute the commit service
	result, err := service.New(ctx.Get("project_to_commit_into").(*model.Project), user, commitParams).Execute(ctx)
	if err != nil {
		return err
	}

	// Handle the result
	if result.Status == "success" {
		// Get the final success path
		finalSuccessPath := c.getFinalSuccessPath(ctx, successPath, targetProject)

		// Update the flash notice
		c.updateFlashNotice(ctx, successNotice, finalSuccessPath)

		// Handle the response format
		switch ctx.GetHeader("Accept") {
		case "application/json":
			ctx.JSON(http.StatusOK, gin.H{
				"message":   "success",
				"filePath": finalSuccessPath,
			})
		default:
			ctx.Redirect(http.StatusFound, finalSuccessPath)
		}
	} else {
		// Set the flash alert
		ctx.Set("flash_alert", c.formatFlashNotice(result.Message))

		// Get the failure path
		var finalFailurePath string
		switch p := failurePath.(type) {
		case string:
			finalFailurePath = p
		case func() string:
			finalFailurePath = p()
		default:
			finalFailurePath = "/"
		}

		// Handle the response format
		switch ctx.GetHeader("Accept") {
		case "application/json":
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":    result.Message,
				"filePath": finalFailurePath,
			})
		default:
			if failureView != "" {
				ctx.HTML(http.StatusUnprocessableEntity, failureView, nil)
			} else {
				ctx.Redirect(http.StatusFound, finalFailurePath)
			}
		}
	}

	return nil
}

// AuthorizeEditTree authorizes editing the tree
func (c *CreatesCommit) AuthorizeEditTree(ctx *gin.Context) error {
	// Get the project from the context
	project, exists := ctx.Get("project")
	if !exists {
		return util.NewNotFoundError("project not found")
	}
	projectObj := project.(*model.Project)

	// Get the branch name or ref
	branchNameOrRef := c.getBranchNameOrRef(ctx)

	// Check if the user can collaborate with the project
	if c.userService.CanCollaborateWithProject(ctx, projectObj, branchNameOrRef) {
		return nil
	}

	// Return an access denied error
	return util.NewForbiddenError("you don't have permission to edit this tree")
}

// formatFlashNotice formats a flash notice
func (c *CreatesCommit) formatFlashNotice(message string) string {
	// Replace newlines with <br> tags
	formattedMessage := strings.ReplaceAll(message, "\n", "<br>")

	// Sanitize the message
	return template.HTMLEscapeString(formattedMessage)
}

// updateFlashNotice updates the flash notice
func (c *CreatesCommit) updateFlashNotice(ctx *gin.Context, successNotice, successPath string) {
	// Create the changes link
	changesLink := fmt.Sprintf(`<a href="%s" class="gl-link">changes</a>`, successPath)

	// Create the default message
	defaultMessage := fmt.Sprintf("Your %s have been committed successfully.", changesLink)

	// Set the flash notice
	if successNotice != "" {
		ctx.Set("flash_notice", successNotice)
	} else {
		ctx.Set("flash_notice", defaultMessage)
	}

	// Check if we should create a merge request
	if c.shouldCreateMergeRequest(ctx) {
		// Check if the merge request exists
		if c.mergeRequestExists(ctx) {
			ctx.Set("flash_notice", nil)
		} else {
			// Get the different project flag
			differentProject, _ := ctx.Get("different_project")
			isDifferentProject := differentProject.(bool)

			// Create the merge request message
			var mrMessage string
			if isDifferentProject {
				mrMessage = "You can now submit a merge request to get this change into the original project."
			} else {
				mrMessage = "You can now submit a merge request to get this change into the original branch."
			}

			// Append the merge request message to the flash notice
			ctx.Set("flash_notice", ctx.Get("flash_notice").(string)+" "+mrMessage)
		}
	}
}

// getFinalSuccessPath gets the final success path
func (c *CreatesCommit) getFinalSuccessPath(ctx *gin.Context, successPath string, targetProject *model.Project) string {
	// Check if we should create a merge request
	if c.shouldCreateMergeRequest(ctx) {
		// Check if the merge request exists
		if c.mergeRequestExists(ctx) {
			return c.getExistingMergeRequestPath(ctx)
		} else {
			return c.getNewMergeRequestPath(ctx, targetProject)
		}
	}

	// Handle the success path
	switch p := successPath.(type) {
	case string:
		return p
	case func() string:
		return p()
	default:
		return "/"
	}
}

// getNewMergeRequestPath gets the path for a new merge request
func (c *CreatesCommit) getNewMergeRequestPath(ctx *gin.Context, targetProject *model.Project) string {
	// Get the project to commit into
	projectToCommitInto, _ := ctx.Get("project_to_commit_into")
	projectToCommitIntoObj := projectToCommitInto.(*model.Project)

	// Get the branch name
	branchName, _ := ctx.Get("branch_name")
	branchNameStr := branchName.(string)

	// Get the start branch
	startBranch, _ := ctx.Get("start_branch")
	startBranchStr := startBranch.(string)

	// Get the default merge request target
	defaultTarget := projectToCommitIntoObj.DefaultMergeRequestTarget

	// Create the merge request path
	return fmt.Sprintf(
		"/projects/%d/merge_requests/new?merge_request[target_project_id]=%d&merge_request[source_branch]=%s&merge_request[target_branch]=%s",
		projectToCommitIntoObj.ID,
		defaultTarget.ID,
		branchNameStr,
		startBranchStr,
	)
}

// getExistingMergeRequestPath gets the path for an existing merge request
func (c *CreatesCommit) getExistingMergeRequestPath(ctx *gin.Context) string {
	// Get the project from the context
	project, _ := ctx.Get("project")
	projectObj := project.(*model.Project)

	// Get the merge request from the context
	mergeRequest, _ := ctx.Get("merge_request")
	mergeRequestObj := mergeRequest.(*model.MergeRequest)

	// Create the merge request path
	return fmt.Sprintf("/projects/%d/merge_requests/%d", projectObj.ID, mergeRequestObj.ID)
}

// mergeRequestExists checks if a merge request exists
func (c *CreatesCommit) mergeRequestExists(ctx *gin.Context) bool {
	// Get the current user from the context
	currentUser, _ := ctx.Get("current_user")
	user := currentUser.(*model.User)

	// Get the project from the context
	project, _ := ctx.Get("project")
	projectObj := project.(*model.Project)

	// Get the project to commit into
	projectToCommitInto, _ := ctx.Get("project_to_commit_into")
	projectToCommitIntoObj := projectToCommitInto.(*model.Project)

	// Get the branch name
	branchName, _ := ctx.Get("branch_name")
	branchNameStr := branchName.(string)

	// Get the start branch
	startBranch, _ := ctx.Get("start_branch")
	startBranchStr := startBranch.(string)

	// Find the merge request
	mr := c.mrService.FindBySourceAndTarget(
		ctx,
		user,
		projectObj.ID,
		projectToCommitIntoObj.ID,
		branchNameStr,
		startBranchStr,
	)

	// Set the merge request in the context
	if mr != nil {
		ctx.Set("merge_request", mr)
	}

	return mr != nil
}

// shouldCreateMergeRequest checks if we should create a merge request
func (c *CreatesCommit) shouldCreateMergeRequest(ctx *gin.Context) bool {
	// Get the create merge request parameter
	createMR, exists := ctx.GetQuery("create_merge_request")
	if !exists {
		return false
	}

	// Parse the boolean value
	createMRBool, err := strconv.ParseBool(createMR)
	if err != nil {
		return false
	}

	// Get the different project flag
	differentProject, _ := ctx.Get("different_project")
	isDifferentProject := differentProject.(bool)

	// Get the start branch
	startBranch, _ := ctx.Get("start_branch")
	startBranchStr := startBranch.(string)

	// Get the branch name
	branchName, _ := ctx.Get("branch_name")
	branchNameStr := branchName.(string)

	// Check if we should create a merge request
	return createMRBool && (isDifferentProject || startBranchStr != branchNameStr)
}

// getBranchNameOrRef gets the branch name or ref
func (c *CreatesCommit) getBranchNameOrRef(ctx *gin.Context) string {
	// Get the branch name
	branchName, exists := ctx.Get("branch_name")
	if exists {
		return branchName.(string)
	}

	// Get the ref
	ref, exists := ctx.Get("ref")
	if exists {
		return ref.(string)
	}

	return ""
}

// getStartBranch gets the start branch
func (c *CreatesCommit) getStartBranch(ctx *gin.Context) string {
	// Get the ref
	ref, exists := ctx.Get("ref")
	if exists {
		return ref.(string)
	}

	// Get the branch name
	branchName, exists := ctx.Get("branch_name")
	if exists {
		return branchName.(string)
	}

	return ""
}

// getCommitParams gets the commit parameters
func (c *CreatesCommit) getCommitParams(ctx *gin.Context) map[string]interface{} {
	// Get the commit parameters from the context
	commitParams, exists := ctx.Get("commit_params")
	if exists {
		return commitParams.(map[string]interface{})
	}

	return make(map[string]interface{})
}
