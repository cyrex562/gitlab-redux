package rapid_diffs

// DiffFile represents a file in a diff
type DiffFile struct {
	FilePath       string
	OldPath        string
	ContentSha     string
	OldContentSha  string
	Repository     *Repository
	NewFile        bool
	DeletedFile    bool
	RenamedFile    bool
	ContentChanged bool
	ModeChanged    bool
	AMode          string
	BMode          string
	TooLarge       bool
	Collapsed      bool
	Diffable       bool
	DiffableText   bool
}

// Repository represents a git repository
type Repository struct {
	Project *Project
}

// Project represents a GitLab project
type Project struct {
	ID   int64
	Path string
}

// ViewerComponent is the base component for all diff viewers
type ViewerComponent struct {
	DiffFile *DiffFile
}

// NewViewerComponent creates a new viewer component with the given diff file
func NewViewerComponent(diffFile *DiffFile) *ViewerComponent {
	return &ViewerComponent{
		DiffFile: diffFile,
	}
}

// ViewerName returns the name of the viewer
func (v *ViewerComponent) ViewerName() string {
	panic("ViewerName must be implemented by each viewer component")
}
