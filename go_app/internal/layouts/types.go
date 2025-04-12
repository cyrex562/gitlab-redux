package layouts

// BaseComponent represents the base layout component
type BaseComponent struct {
	Content   string
	IsVisible bool
}

// ComponentOptions contains common options for all components
type ComponentOptions struct {
	Content   string
	IsVisible bool
}

// CRUDComponent represents a CRUD (Create, Read, Update, Delete) component
type CRUDComponent struct {
	BaseComponent
	Title       string
	Description string
	Actions     []Action
	Form        interface{}
}

// Action represents a CRUD action
type Action struct {
	Label     string
	URL       string
	Method    string
	Icon      string
	Class     string
	Data      map[string]string
}

// EmptyResultComponent represents a component for displaying empty states
type EmptyResultComponent struct {
	BaseComponent
	Title       string
	Description string
	Icon        string
	Action      *Action
}

// HorizontalSectionComponent represents a horizontal section layout
type HorizontalSectionComponent struct {
	BaseComponent
	Title       string
	Description string
	Content     string
}

// PageHeadingComponent represents a page heading layout
type PageHeadingComponent struct {
	BaseComponent
	Title       string
	Description string
	Actions     []Action
}

// SettingsBlockComponent represents a settings block layout
type SettingsBlockComponent struct {
	BaseComponent
	Title       string
	Description string
	Content     string
	HelpText    string
}

// SettingsSectionComponent represents a settings section layout
type SettingsSectionComponent struct {
	BaseComponent
	Title       string
	Description string
	Content     string
}
