package lfs

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

const (
	ContentType = "application/vnd.git-lfs+json"
)

// LfsRequest handles Git LFS request functionality
type LfsRequest struct {
	configService    *service.ConfigService
	projectService   *service.ProjectService
	authService      *service.AuthService
	logger           *service.Logger
	container        interface{}
	project          *service.Project
	user             *service.User
	deployToken      *service.DeployToken
	authenticationResult *service.AuthenticationResult
}

// NewLfsRequest creates a new instance of LfsRequest
func NewLfsRequest(
	configService *service.ConfigService,
	projectService *service.ProjectService,
	authService *service.AuthService,
	logger *service.Logger,
) *LfsRequest {
	return &LfsRequest{
		configService:    configService,
		projectService:   projectService,
		authService:      authService,
		logger:           logger,
	}
}

// SetupMiddleware sets up the LFS request middleware
func (l *LfsRequest) SetupMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := l.requireLfsEnabled(ctx); err != nil {
			return
		}

		if err := l.lfsCheckAccess(ctx); err != nil {
			return
		}

		ctx.Next()
	}
}

// requireLfsEnabled checks if Git LFS is enabled
func (l *LfsRequest) requireLfsEnabled(ctx *gin.Context) error {
	if !l.configService.LfsEnabled() {
		l.renderLfsError(ctx, http.StatusNotImplemented, "Git LFS is not enabled on this GitLab server, contact your admin.")
		return nil
	}
	return nil
}

// lfsCheckAccess checks if the request has proper LFS access
func (l *LfsRequest) lfsCheckAccess(ctx *gin.Context) error {
	// Check if container has LFS enabled
	container, ok := l.container.(interface{ LfsEnabled() bool })
	if !ok || !container.LfsEnabled() {
		l.renderLfsNotFound(ctx)
		return nil
	}

	// Check download access
	if l.isDownloadRequest(ctx) && l.lfsDownloadAccess() {
		return nil
	}

	// Check upload access
	if l.isUploadRequest(ctx) && l.lfsUploadAccess() {
		return nil
	}

	// Return appropriate error based on access level
	if l.lfsDownloadAccess() {
		l.lfsForbidden(ctx)
	} else {
		l.renderLfsNotFound(ctx)
	}

	return nil
}

// lfsForbidden renders a forbidden error
func (l *LfsRequest) lfsForbidden(ctx *gin.Context) {
	l.renderLfsError(ctx, http.StatusForbidden, "Access forbidden. Check your access level.")
}

// renderLfsNotFound renders a not found error
func (l *LfsRequest) renderLfsNotFound(ctx *gin.Context) {
	l.renderLfsError(ctx, http.StatusNotFound, "Not found.")
}

// renderLfsError renders an LFS error response
func (l *LfsRequest) renderLfsError(ctx *gin.Context, status int, message string) {
	ctx.JSON(status, gin.H{
		"message":           message,
		"documentation_url": l.helpURL(),
	})
	ctx.Header("Content-Type", ContentType)
	ctx.Abort()
}

// lfsDownloadAccess checks if the request has download access
func (l *LfsRequest) lfsDownloadAccess() bool {
	return l.isCI() ||
		l.lfsDeployToken() ||
		l.userCanDownloadCode() ||
		l.buildCanDownloadCode() ||
		l.deployTokenCanDownloadCode()
}

// deployTokenCanDownloadCode checks if the deploy token has download access
func (l *LfsRequest) deployTokenCanDownloadCode() bool {
	return l.deployToken != nil &&
		l.deployToken.HasAccessTo(l.project) &&
		l.deployToken.ReadRepository()
}

// lfsUploadAccess checks if the request has upload access
func (l *LfsRequest) lfsUploadAccess() bool {
	if !l.hasAuthenticationAbility("push_code") {
		return false
	}

	if l.limitExceeded() {
		return false
	}

	return l.lfsDeployToken() ||
		l.can(l.user, "push_code", l.project) ||
		l.can(l.deployToken, "push_code", l.project) ||
		l.anyBranchAllowsCollaboration()
}

// anyBranchAllowsCollaboration checks if any branch allows collaboration
func (l *LfsRequest) anyBranchAllowsCollaboration() bool {
	return l.project.MergeRequestsAllowingPushToUser(l.user).Any()
}

// lfsDeployToken checks if the request uses a deploy token
func (l *LfsRequest) lfsDeployToken() bool {
	return l.authenticationResult.LfsDeployToken(l.project)
}

// userCanDownloadCode checks if the user can download code
func (l *LfsRequest) userCanDownloadCode() bool {
	return l.hasAuthenticationAbility("download_code") &&
		l.can(l.user, "download_code", l.project)
}

// buildCanDownloadCode checks if the build can download code
func (l *LfsRequest) buildCanDownloadCode() bool {
	return l.hasAuthenticationAbility("build_download_code") &&
		l.can(l.user, "build_download_code", l.project)
}

// GetObjects gets the LFS objects from the request
func (l *LfsRequest) GetObjects(ctx *gin.Context) []map[string]interface{} {
	objects, _ := ctx.Get("objects")
	if objects == nil {
		return []map[string]interface{}{}
	}
	return objects.([]map[string]interface{})
}

// GetObjectsOids gets the LFS object OIDs from the request
func (l *LfsRequest) GetObjectsOids(ctx *gin.Context) []string {
	objects := l.GetObjects(ctx)
	oids := make([]string, 0, len(objects))
	for _, obj := range objects {
		if oid, ok := obj["oid"].(string); ok {
			oids = append(oids, oid)
		}
	}
	return oids
}

// hasAuthenticationAbility checks if the request has a specific authentication ability
func (l *LfsRequest) hasAuthenticationAbility(capability string) bool {
	abilities := l.authenticationResult.Abilities()
	for _, ability := range abilities {
		if ability == capability {
			return true
		}
	}
	return false
}

// can checks if a subject can perform an action on an object
func (l *LfsRequest) can(subject interface{}, action string, object interface{}) bool {
	return l.authService.Can(subject, action, object)
}

// isCI checks if the request is from CI
func (l *LfsRequest) isCI() bool {
	return l.authService.IsCI()
}

// isDownloadRequest checks if the request is a download request
func (l *LfsRequest) isDownloadRequest(ctx *gin.Context) bool {
	return ctx.Request.Method == http.MethodGet
}

// isUploadRequest checks if the request is an upload request
func (l *LfsRequest) isUploadRequest(ctx *gin.Context) bool {
	return ctx.Request.Method == http.MethodPost
}

// limitExceeded checks if the LFS limit is exceeded
func (l *LfsRequest) limitExceeded() bool {
	return false // Overridden in EE
}

// helpURL returns the help URL
func (l *LfsRequest) helpURL() string {
	return "https://docs.gitlab.com/ee/topics/git/lfs/"
}
