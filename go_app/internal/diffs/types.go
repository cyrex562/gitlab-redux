package diffs

// DiffStats represents statistics about a diff
type DiffStats struct {
	Additions    int
	Deletions    int
	TotalChanges int
}

// OverflowWarning represents a warning about diff overflow
type OverflowWarning struct {
	Message     string
	IsVisible   bool
	MaxLines    int
	CurrentLines int
}

// BaseComponent represents the base diff component
type BaseComponent struct {
	Stats     *DiffStats
	Warning   *OverflowWarning
	Content   string
	IsVisible bool
}

// ComponentOptions contains options for creating diff components
type ComponentOptions struct {
	MaxLines     int
	CurrentLines int
	Additions    int
	Deletions    int
	Content      string
	IsVisible    bool
}
