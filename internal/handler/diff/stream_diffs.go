package diff

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
	"gitlab.com/gitlab-org/gitlab-redux/internal/util"
)

// StreamDiffs handles streaming diffs
type StreamDiffs struct {
	featureService *service.FeatureService
	userService    *service.UserService
	diffService    *service.DiffService
	viewService    *service.ViewService
	logger         *util.Logger
}

// NewStreamDiffs creates a new instance of StreamDiffs
func NewStreamDiffs(
	featureService *service.FeatureService,
	userService *service.UserService,
	diffService *service.DiffService,
	viewService *service.ViewService,
	logger *util.Logger,
) *StreamDiffs {
	return &StreamDiffs{
		featureService: featureService,
		userService:    userService,
		diffService:    diffService,
		viewService:    viewService,
		logger:         logger,
	}
}

// Diffs streams diffs
func (s *StreamDiffs) Diffs(c *gin.Context) {
	// Check if rapid diffs is enabled
	enabled, err := s.isRapidDiffsEnabled(c)
	if err != nil {
		s.logger.Error("failed to check if rapid diffs is enabled", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check if rapid diffs is enabled"})
		return
	}

	if !enabled {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rapid diffs is not enabled"})
		return
	}

	// Get streaming start time
	streamingStartTime := time.Now()

	// Set streaming headers
	s.setStreamHeaders(c)

	// Get offset from query parameters
	offsetStr := c.DefaultQuery("offset", "0")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		s.logger.Error("failed to parse offset", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset"})
		return
	}

	// Get diff blobs from query parameters
	diffBlobsStr := c.DefaultQuery("diff_blobs", "false")
	diffBlobs, err := strconv.ParseBool(diffBlobsStr)
	if err != nil {
		s.logger.Error("failed to parse diff_blobs", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid diff_blobs"})
		return
	}

	// Get resource
	resource, err := s.getResource(c)
	if err != nil {
		s.logger.Error("failed to get resource", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get resource"})
		return
	}

	// Create diff options
	options := model.NewDiffOptions().
		WithOffsetIndex(offset).
		WithDiffBlobs(diffBlobs)

	// Stream diff files
	if diffBlobs {
		err = s.streamDiffBlobs(c, resource, options)
	} else {
		err = s.streamDiffFiles(c, resource, options)
	}

	// Calculate streaming time
	streamingTime := time.Since(streamingStartTime).Seconds()

	// Write server timings
	c.Writer.Write([]byte(fmt.Sprintf("<server-timings streaming=\"%.2f\"></server-timings>", streamingTime)))

	// Handle errors
	if err != nil {
		s.logger.Error("error streaming diffs", "error", err)
		errorHTML, err := s.viewService.RenderToString(c, "rapid_diffs/streaming_error", gin.H{
			"message": err.Error(),
		})
		if err != nil {
			s.logger.Error("failed to render error", "error", err)
			return
		}
		c.Writer.Write([]byte(errorHTML))
	}
}

// setStreamHeaders sets the headers for streaming
func (s *StreamDiffs) setStreamHeaders(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Header("X-Accel-Buffering", "no")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
}

// isRapidDiffsEnabled checks if rapid diffs is enabled
func (s *StreamDiffs) isRapidDiffsEnabled(c *gin.Context) (bool, error) {
	// Get current user
	user, err := s.userService.GetCurrentUser(c)
	if err != nil {
		return false, err
	}

	// Check if feature is enabled
	return s.featureService.IsEnabled("rapid_diffs", user, "wip")
}

// getResource gets the resource for streaming diffs
func (s *StreamDiffs) getResource(c *gin.Context) (interface{}, error) {
	// This is a placeholder for actual implementation
	// In a real implementation, this would get the resource from the context
	return nil, nil
}

// streamDiffFiles streams diff files
func (s *StreamDiffs) streamDiffFiles(c *gin.Context, resource interface{}, options *model.DiffOptions) error {
	if resource == nil {
		return nil
	}

	// Get diffs for streaming
	diffs, err := s.diffService.GetDiffsForStreaming(resource)
	if err != nil {
		return err
	}

	// Stream each diff file
	for _, diffFile := range diffs.DiffFiles {
		// Render diff file
		html, err := s.renderDiffFile(c, diffFile)
		if err != nil {
			return err
		}

		// Write diff file to response
		_, err = c.Writer.Write([]byte(html))
		if err != nil {
			return err
		}
	}

	return nil
}

// renderDiffFile renders a diff file
func (s *StreamDiffs) renderDiffFile(c *gin.Context, diffFile *model.DiffFile) (string, error) {
	// Render diff file
	return s.viewService.RenderToString(c, "rapid_diffs/diff_file", gin.H{
		"diff_file": diffFile,
		"parallel_view": diffFile.ParallelView,
	})
}

// streamDiffBlobs streams diff blobs
func (s *StreamDiffs) streamDiffBlobs(c *gin.Context, resource interface{}, options *model.DiffOptions) error {
	// Stream diff files in batches
	return s.diffService.StreamDiffFiles(c, resource, map[string]interface{}{
		"offset_index": options.OffsetIndex,
		"diff_blobs":   options.DiffBlobs,
	}, func(diffFiles []*model.DiffFile) error {
		// Stream each diff file in the batch
		for _, diffFile := range diffFiles {
			// Render diff file
			html, err := s.renderDiffFile(c, diffFile)
			if err != nil {
				return err
			}

			// Write diff file to response
			_, err = c.Writer.Write([]byte(html))
			if err != nil {
				return err
			}
		}

		return nil
	})
}
