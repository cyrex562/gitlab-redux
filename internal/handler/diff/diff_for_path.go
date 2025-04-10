package diff

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
	"gitlab.com/gitlab-org/gitlab-redux/internal/util"
)

// DiffForPath handles rendering diffs for specific file paths
type DiffForPath struct {
	diffService *service.DiffService
	viewService *service.ViewService
	logger      *util.Logger
}

// NewDiffForPath creates a new instance of DiffForPath
func NewDiffForPath(
	diffService *service.DiffService,
	viewService *service.ViewService,
	logger *util.Logger,
) *DiffForPath {
	return &DiffForPath{
		diffService: diffService,
		viewService: viewService,
		logger:      logger,
	}
}

// RenderDiffForPath renders the diff for a specific file path
func (d *DiffForPath) RenderDiffForPath(ctx *gin.Context, diffs *model.Diffs) error {
	// Get the file identifier from the request
	fileIdentifier, exists := ctx.GetQuery("file_identifier")
	if !exists {
		return util.NewBadRequestError("file_identifier is required")
	}

	// Find the diff file with the matching identifier
	var diffFile *model.DiffFile
	for _, diff := range diffs.DiffFiles {
		if diff.FileIdentifier == fileIdentifier {
			diffFile = diff
			break
		}
	}

	// Return 404 if the diff file is not found
	if diffFile == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Diff file not found",
		})
		return util.NewNotFoundError("diff file not found")
	}

	// Render the diff content
	html, err := d.viewService.RenderToString(ctx, "projects/diffs/_content", gin.H{
		"diff_file": diffFile,
	})
	if err != nil {
		d.logger.Error("failed to render diff content", "error", err)
		return err
	}

	// Return the rendered HTML
	ctx.JSON(http.StatusOK, gin.H{
		"html": html,
	})

	return nil
}
