package model

// Diffs represents a collection of diff files
type Diffs struct {
	DiffFiles []*DiffFile
}

// DiffFile represents a file with diffs
type DiffFile struct {
	FileIdentifier string
	FileName       string
	FilePath       string
	Content        string
	RawContent     string
	ParallelView   bool
}

// DiffOptions represents options for diff operations
type DiffOptions struct {
	OffsetIndex int
	View        string
	DiffBlobs   bool
}

// NewDiffOptions creates a new instance of DiffOptions
func NewDiffOptions() *DiffOptions {
	return &DiffOptions{
		OffsetIndex: 0,
		View:        "inline",
		DiffBlobs:   false,
	}
}

// WithOffsetIndex sets the offset index
func (o *DiffOptions) WithOffsetIndex(offsetIndex int) *DiffOptions {
	o.OffsetIndex = offsetIndex
	return o
}

// WithView sets the view
func (o *DiffOptions) WithView(view string) *DiffOptions {
	o.View = view
	return o
}

// WithDiffBlobs sets the diff blobs flag
func (o *DiffOptions) WithDiffBlobs(diffBlobs bool) *DiffOptions {
	o.DiffBlobs = diffBlobs
	return o
}
