package rapid_diffs

import (
	"fmt"
	"strings"
)

// AppComponent represents the main rapid diffs application component
type AppComponent struct {
	DiffsSlice         []*DiffFile
	ReloadStreamURL    string
	StreamURL          string
	ShowWhitespace     bool
	DiffView           string
	UpdateUserEndpoint string
	MetadataEndpoint   string
	Preload            bool
	DiffsList          string
	InitialSidebarWidth string
}

// NewAppComponent creates a new app component with the given parameters
func NewAppComponent(
	diffsSlice []*DiffFile,
	reloadStreamURL string,
	streamURL string,
	showWhitespace bool,
	diffView string,
	updateUserEndpoint string,
	metadataEndpoint string,
	preload bool,
	diffsList string,
	initialSidebarWidth string,
) *AppComponent {
	return &AppComponent{
		DiffsSlice:         diffsSlice,
		ReloadStreamURL:    reloadStreamURL,
		StreamURL:          streamURL,
		ShowWhitespace:     showWhitespace,
		DiffView:           diffView,
		UpdateUserEndpoint: updateUserEndpoint,
		MetadataEndpoint:   metadataEndpoint,
		Preload:            preload,
		DiffsList:          diffsList,
		InitialSidebarWidth: initialSidebarWidth,
	}
}

// Render generates the HTML for the app component
func (a *AppComponent) Render() string {
	var parts []string

	// Preload section
	if a.Preload {
		if a.StreamURL != "" {
			parts = append(parts, fmt.Sprintf(`
				<script nonce="%s">
					var controller = new AbortController();
					window.gl.rapidDiffsPreload = {
						controller: controller,
						streamRequest: fetch('%s', { signal: controller.signal })
					}
				</script>
			`, "CSP_NONCE", a.StreamURL))
		}
	}

	// Main app container
	parts = append(parts, fmt.Sprintf(`
		<div class="rd-app" data-rapid-diffs="true" data-reload-stream-url="%s" data-metadata-endpoint="%s">
	`, a.ReloadStreamURL, a.MetadataEndpoint))

	// App header
	parts = append(parts, `
		<div class="rd-app-header">
			<div class="rd-app-settings">
				<div data-view-settings="true" data-show-whitespace="%s" data-diff-view-type="%s" data-update-user-endpoint="%s"></div>
			</div>
		</div>
	`)

	// App body
	sidebarStyle := ""
	if a.InitialSidebarWidth != "" {
		sidebarStyle = fmt.Sprintf(` style="width: %spx"`, a.InitialSidebarWidth)
	}

	parts = append(parts, fmt.Sprintf(`
		<div class="rd-app-body">
			<div class="rd-app-sidebar" data-file-browser="true"%s>
				<div class="rd-app-sidebar-loading">
					<span class="gl-spinner gl-spinner-sm gl-spinner-dark !gl-align-text-bottom"></span>
				</div>
			</div>
			<div class="rd-app-content" data-sidebar-visible="true">
				<div class="rd-app-content-header" data-hidden-files-warning="true"></div>
				<div class="code">
					<div data-diffs-list="true">
	`, sidebarStyle))

	// Performance mark
	parts = append(parts, `
		<script nonce="CSP_NONCE">
			requestAnimationFrame(() => { window.performance.mark('rapid-diffs-first-diff-file-shown') })
		</script>
	`)

	// Diffs list
	if a.DiffsList != "" {
		parts = append(parts, a.DiffsList)
	} else {
		for _, diff := range a.DiffsSlice {
			diffFile := NewDiffFileComponent(diff, a.DiffView == "parallel")
			parts = append(parts, diffFile.Render())
		}
	}

	// Stream container
	if a.StreamURL != "" {
		parts = append(parts, fmt.Sprintf(`
			<div id="js-stream-container" data-diffs-stream-url="%s"></div>
		`, a.StreamURL))
	} else {
		parts = append(parts, `
			<script nonce="CSP_NONCE">
				requestAnimationFrame(() => {
					window.performance.mark('rapid-diffs-list-loaded');
					window.performance.measure('rapid-diffs-list-loading', 'rapid-diffs-first-diff-file-shown', 'rapid-diffs-list-loaded');
				})
			</script>
		`)
	}

	// Close all divs
	parts = append(parts, "</div></div></div></div></div>")

	return strings.Join(parts, "\n")
}
