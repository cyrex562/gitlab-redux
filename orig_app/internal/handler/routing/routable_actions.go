package routing

import (
	"net/http"
	"strings"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
)

// RoutableActions handles finding and handling routable objects
type RoutableActions struct {
	authorizer Authorizer
	router Router
	flash Flash
	projectUnauthorized ProjectUnauthorized
}

// Authorizer defines the interface for authorization operations
type Authorizer interface {
	Can(user *model.User, action string, subject interface{}) bool
}

// Router defines the interface for routing operations
type Router interface {
	BuildCanonicalPath(routable interface{}) (string, error)
}

// Flash defines the interface for flash message operations
type Flash interface {
	SetNotice(message string)
}

// ProjectUnauthorized defines the interface for project unauthorized operations
type ProjectUnauthorized interface {
	OnRoutableNotFound() func(routable interface{}, fullPath string)
}

// NewRoutableActions creates a new instance of RoutableActions
func NewRoutableActions(
	authorizer Authorizer,
	router Router,
	flash Flash,
	projectUnauthorized ProjectUnauthorized,
) *RoutableActions {
	return &RoutableActions{
		authorizer: authorizer,
		router: router,
		flash: flash,
		projectUnauthorized: projectUnauthorized,
	}
}

// FindRoutable finds a routable object by its full path
func (r *RoutableActions) FindRoutable(
	w http.ResponseWriter,
	req *http.Request,
	routableKlass RoutableKlass,
	routableFullPath string,
	fullPath string,
	currentUser *model.User,
	extraAuthorizationProc func(interface{}) bool,
) interface{} {
	followRedirects := req.Method == http.MethodGet
	routable := routableKlass.FindByFullPath(routableFullPath, followRedirects)

	if r.RoutableAuthorized(routable, currentUser, extraAuthorizationProc) {
		r.EnsureCanonicalPath(w, req, routable, routableFullPath)
		return routable
	}

	r.PerformNotFoundActions(routable, r.NotFoundActions(), fullPath)

	return nil
}

// NotFoundActions returns a list of actions to perform when a routable is not found
func (r *RoutableActions) NotFoundActions() []func(interface{}, string) {
	return []func(interface{}, string){
		r.projectUnauthorized.OnRoutableNotFound(),
	}
}

// PerformNotFoundActions performs the not found actions
func (r *RoutableActions) PerformNotFoundActions(routable interface{}, actions []func(interface{}, string), fullPath string) {
	for _, action := range actions {
		action(routable, fullPath)
	}
}

// RoutableAuthorized checks if a routable is authorized
func (r *RoutableActions) RoutableAuthorized(routable interface{}, currentUser *model.User, extraAuthorizationProc func(interface{}) bool) bool {
	if routable == nil {
		return false
	}

	routableType := GetRoutableType(routable)
	action := "read_" + ToSnakeCase(routableType)

	if !r.authorizer.Can(currentUser, action, routable) {
		return false
	}

	if extraAuthorizationProc != nil {
		return extraAuthorizationProc(routable)
	}

	return true
}

// EnsureCanonicalPath ensures the canonical path is used
func (r *RoutableActions) EnsureCanonicalPath(w http.ResponseWriter, req *http.Request, routable interface{}, routableFullPath string) {
	if req.Method != http.MethodGet {
		return
	}

	canonicalPath := GetRoutableFullPath(routable)
	if canonicalPath == routableFullPath {
		return
	}

	if !IsXHR(req) && IsHTMLFormat(req) && !strings.EqualFold(canonicalPath, routableFullPath) {
		routableType := GetRoutableType(routable)
		routableTypeTitle := ToTitleCase(routableType)

		message := routableTypeTitle + " '" + routableFullPath + "' was moved to '" + canonicalPath + "'. " +
			"Please update any links and bookmarks that may still have the old path."

		r.flash.SetNotice(message)
	}

	path, err := r.router.BuildCanonicalPath(routable)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	http.Redirect(w, req, path, http.StatusMovedPermanently)
}

// Helper functions

// GetRoutableType returns the type of a routable object
func GetRoutableType(routable interface{}) string {
	// This would be implemented based on the actual type of the routable
	// For example, it could use reflection to get the type name
	return "unknown"
}

// GetRoutableFullPath returns the full path of a routable object
func GetRoutableFullPath(routable interface{}) string {
	// This would be implemented based on the actual type of the routable
	// For example, it could call a method on the routable to get its full path
	return ""
}

// ToSnakeCase converts a string from CamelCase to snake_case
func ToSnakeCase(s string) string {
	// This would be implemented to convert CamelCase to snake_case
	return strings.ToLower(s)
}

// ToTitleCase converts a string to Title Case
func ToTitleCase(s string) string {
	// This would be implemented to convert a string to Title Case
	return s
}

// IsXHR checks if a request is an XHR request
func IsXHR(req *http.Request) bool {
	return req.Header.Get("X-Requested-With") == "XMLHttpRequest"
}

// IsHTMLFormat checks if a request is for HTML format
func IsHTMLFormat(req *http.Request) bool {
	accept := req.Header.Get("Accept")
	return strings.Contains(accept, "text/html")
}

// RoutableKlass defines the interface for routable class operations
type RoutableKlass interface {
	FindByFullPath(fullPath string, followRedirects bool) interface{}
}
