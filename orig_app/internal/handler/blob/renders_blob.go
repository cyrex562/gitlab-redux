package blob

import (
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/view"
)

// RendersBlob handles rendering of blobs in GitLab
type RendersBlob struct {
	viewRenderer ViewRenderer
}

// ViewRenderer defines the interface for rendering views
type ViewRenderer interface {
	ViewToHTMLString(templatePath string, data map[string]interface{}) string
}

// NewRendersBlob creates a new instance of RendersBlob
func NewRendersBlob(viewRenderer ViewRenderer) *RendersBlob {
	return &RendersBlob{
		viewRenderer: viewRenderer,
	}
}

// BlobViewerJSON determines the appropriate viewer for a blob based on parameters
func (r *RendersBlob) BlobViewerJSON(blob *model.Blob, viewerParam string) map[string]interface{} {
	var viewer interface{}

	switch viewerParam {
	case "rich":
		viewer = blob.RichViewer()
	case "auxiliary":
		viewer = blob.AuxiliaryViewer()
	case "none":
		viewer = nil
	default:
		viewer = blob.SimpleViewer()
	}

	if viewer == nil {
		return map[string]interface{}{}
	}

	html := r.viewRenderer.ViewToHTMLString("projects/blob/_viewer", map[string]interface{}{
		"viewer":     viewer,
		"load_async": false,
	})

	return map[string]interface{}{
		"html": html,
	}
}

// RenderBlobJSON renders the blob as JSON
func (r *RendersBlob) RenderBlobJSON(blob *model.Blob, viewerParam string) (*view.Response, error) {
	json := r.BlobViewerJSON(blob, viewerParam)

	if len(json) == 0 {
		return view.NewNotFoundResponse(), nil
	}

	return view.NewJSONResponse(json), nil
}

// ConditionallyExpandBlob expands a single blob if requested
func (r *RendersBlob) ConditionallyExpandBlob(blob *model.Blob, expandedParam string) {
	r.ConditionallyExpandBlobs([]*model.Blob{blob}, expandedParam)
}

// ConditionallyExpandBlobs expands multiple blobs if requested
func (r *RendersBlob) ConditionallyExpandBlobs(blobs []*model.Blob, expandedParam string) {
	if expandedParam != "true" {
		return
	}

	for _, blob := range blobs {
		blob.Expand()
	}
}
