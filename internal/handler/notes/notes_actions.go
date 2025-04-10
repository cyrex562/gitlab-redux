package notes

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// NotesActions handles note-related actions
type NotesActions struct {
	noteService      *service.NoteService
	notesFinder      *service.NotesFinder
	userService      *service.UserService
	projectService   *service.ProjectService
	viewService      *service.ViewService
	rateLimitService *service.RateLimitService
	logger           *service.Logger
}

// NewNotesActions creates a new instance of NotesActions
func NewNotesActions(
	noteService *service.NoteService,
	notesFinder *service.NotesFinder,
	userService *service.UserService,
	projectService *service.ProjectService,
	viewService *service.ViewService,
	rateLimitService *service.RateLimitService,
	logger *service.Logger,
) *NotesActions {
	return &NotesActions{
		noteService:      noteService,
		notesFinder:      notesFinder,
		userService:      userService,
		projectService:   projectService,
		viewService:      viewService,
		rateLimitService: rateLimitService,
		logger:           logger,
	}
}

// SetupMiddleware sets up the middleware for notes actions
func (n *NotesActions) SetupMiddleware(router *gin.Engine) {
	router.Use(func(c *gin.Context) {
		// Set polling interval header for index action
		if c.Request.URL.Path == "/notes" && c.Request.Method == http.MethodGet {
			n.setPollingIntervalHeader(c, 6000)
		}

		// Require last fetched at header for index action
		if c.Request.URL.Path == "/notes" && c.Request.Method == http.MethodGet {
			if err := n.requireLastFetchedAtHeader(c); err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"message": "X-Last-Fetched-At header is required",
				})
				return
			}
		}

		// Require noteable for index and create actions
		if (c.Request.URL.Path == "/notes" && c.Request.Method == http.MethodGet) ||
			(c.Request.URL.Path == "/notes" && c.Request.Method == http.MethodPost) {
			if err := n.requireNoteable(c); err != nil {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
		}

		// Authorize admin note for update and destroy actions
		if (c.Request.URL.Path == "/notes/:id" && c.Request.Method == http.MethodPut) ||
			(c.Request.URL.Path == "/notes/:id" && c.Request.Method == http.MethodDelete) {
			if err := n.authorizeAdminNote(c); err != nil {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
		}

		// Set note project for create action
		if c.Request.URL.Path == "/notes" && c.Request.Method == http.MethodPost {
			if err := n.noteProject(c); err != nil {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
		}

		// Check rate limit for create action
		if c.Request.URL.Path == "/notes" && c.Request.Method == http.MethodPost {
			currentUser, err := c.Get("current_user")
			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			// Get current application settings
			settings, err := n.userService.GetCurrentApplicationSettings()
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			// Check rate limit
			if err := n.rateLimitService.CheckRateLimit(c, "notes_create", currentUser, settings.NotesCreateLimitAllowlist); err != nil {
				c.AbortWithStatus(http.StatusTooManyRequests)
				return
			}
		}

		c.Next()
	})
}

// Index handles the index action
func (n *NotesActions) Index(ctx *gin.Context) error {
	// Get current user from context
	currentUser, err := ctx.Get("current_user")
	if err != nil {
		return err
	}

	// Gather all notes
	notes, meta, err := n.gatherAllNotes(ctx)
	if err != nil {
		return err
	}

	// Prepare notes for rendering
	notes, err = n.prepareNotesForRendering(ctx, notes)
	if err != nil {
		return err
	}

	// Filter notes by readability
	notes = n.filterNotesByReadability(notes, currentUser)

	// Serialize notes
	var serializedNotes interface{}
	if n.useNoteSerializer(ctx) {
		serializedNotes, err = n.noteSerializer(ctx).Represent(notes)
		if err != nil {
			return err
		}
	} else {
		serializedNotes = n.serializeNotes(ctx, notes)
	}

	// Render JSON response
	ctx.JSON(http.StatusOK, gin.H{
		"notes": serializedNotes,
		"meta":  meta,
	})

	return nil
}

// Create handles the create action
func (n *NotesActions) Create(ctx *gin.Context) error {
	// Get current user from context
	currentUser, err := ctx.Get("current_user")
	if err != nil {
		return err
	}

	// Get note project
	noteProject, err := n.noteProject(ctx)
	if err != nil {
		return err
	}

	// Get create note params
	createNoteParams, err := n.createNoteParams(ctx)
	if err != nil {
		return err
	}

	// Create note
	note, err := n.noteService.Create(noteProject, currentUser, createNoteParams)
	if err != nil {
		return err
	}

	// Handle response based on format
	switch ctx.GetHeader("Accept") {
	case "application/json":
		// Initialize JSON response
		json := gin.H{
			"commands_changes": note.CommandsChanges,
		}

		if note.Persisted && n.returnDiscussion(ctx) {
			json["valid"] = true

			// Get discussion
			discussion, err := note.GetDiscussion()
			if err != nil {
				return err
			}

			// Prepare notes for rendering
			discussionNotes, err := n.prepareNotesForRendering(ctx, discussion.Notes)
			if err != nil {
				return err
			}

			// Serialize discussion
			json["discussion"], err = n.discussionSerializer(ctx).Represent(discussion, ctx)
			if err != nil {
				return err
			}
		} else {
			// Prepare notes for rendering
			_, err = n.prepareNotesForRendering(ctx, []*service.Note{note})
			if err != nil {
				return err
			}

			// Merge note JSON
			for k, v := range n.noteJSON(ctx, note) {
				json[k] = v
			}
		}

		// Get quick actions status
		quickActions := note.QuickActionsStatus
		if quickActions != nil {
			json["quick_actions_status"] = quickActions.ToMap()
		}

		// Handle errors
		if note.Errors != nil && len(note.Errors) > 0 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"errors": n.errorsOnCreate(note),
			})
		} else if quickActions != nil && quickActions.HasError() {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"quick_actions_status": quickActions.ToMap(),
			})
		} else {
			ctx.JSON(http.StatusOK, json)
		}
	case "text/html":
		// Redirect back or default
		ctx.Redirect(http.StatusSeeOther, ctx.DefaultQuery("back", "/"))
	default:
		ctx.JSON(http.StatusNotAcceptable, gin.H{
			"error": "Unsupported format",
		})
	}

	return nil
}

// Update handles the update action
func (n *NotesActions) Update(ctx *gin.Context) error {
	// Get current user from context
	currentUser, err := ctx.Get("current_user")
	if err != nil {
		return err
	}

	// Get project from context
	project, err := ctx.Get("project")
	if err != nil {
		return err
	}

	// Get note from context
	note, err := ctx.Get("note")
	if err != nil {
		return err
	}

	// Get update note params
	updateNoteParams, err := n.updateNoteParams(ctx)
	if err != nil {
		return err
	}

	// Update note
	updatedNote, err := n.noteService.Update(project, currentUser, updateNoteParams, note)
	if err != nil {
		return err
	}

	// Check if note was destroyed
	if updatedNote.Destroyed {
		ctx.Status(http.StatusGone)
		return nil
	}

	// Handle response based on format
	switch ctx.GetHeader("Accept") {
	case "application/json":
		// Handle errors
		if updatedNote.Errors != nil && len(updatedNote.Errors) > 0 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"errors": updatedNote.Errors.FullMessages().ToSentence(),
			})
		} else {
			// Prepare notes for rendering
			_, err = n.prepareNotesForRendering(ctx, []*service.Note{updatedNote})
			if err != nil {
				return err
			}

			// Render JSON response
			ctx.JSON(http.StatusOK, n.noteJSON(ctx, updatedNote))
		}
	case "text/html":
		// Redirect back or default
		ctx.Redirect(http.StatusSeeOther, ctx.DefaultQuery("back", "/"))
	default:
		ctx.JSON(http.StatusNotAcceptable, gin.H{
			"error": "Unsupported format",
		})
	}

	return nil
}

// Destroy handles the destroy action
func (n *NotesActions) Destroy(ctx *gin.Context) error {
	// Get project from context
	project, err := ctx.Get("project")
	if err != nil {
		return err
	}

	// Get current user from context
	currentUser, err := ctx.Get("current_user")
	if err != nil {
		return err
	}

	// Get note from context
	note, err := ctx.Get("note")
	if err != nil {
		return err
	}

	// Check if note is editable
	if note.Editable {
		// Destroy note
		_, err = n.noteService.Destroy(project, currentUser, note)
		if err != nil {
			return err
		}
	}

	// Handle response based on format
	switch ctx.GetHeader("Accept") {
	case "application/javascript":
		ctx.Status(http.StatusOK)
	default:
		ctx.JSON(http.StatusNotAcceptable, gin.H{
			"error": "Unsupported format",
		})
	}

	return nil
}

// Private methods

// GatherAllNotes gathers all notes
func (n *NotesActions) gatherAllNotes(ctx *gin.Context) ([]*service.Note, map[string]interface{}, error) {
	// Get current time
	now := time.Now()

	// Get notes finder
	notesFinder, err := n.notesFinder(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Execute notes finder
	notes, err := notesFinder.Execute()
	if err != nil {
		return nil, nil, err
	}

	// Get noteable from context
	noteable, err := ctx.Get("noteable")
	if err != nil {
		return nil, nil, err
	}

	// Include relations for view
	notes, err = notes.IncRelationsForView(noteable)
	if err != nil {
		return nil, nil, err
	}

	// Merge resource events
	notes, err = n.mergeResourceEvents(ctx, notes)
	if err != nil {
		return nil, nil, err
	}

	// Calculate last fetched at
	lastFetchedAt := (now.Unix() * 1000000) + int64(now.Nanosecond()/1000)

	// Return notes and meta
	return notes, map[string]interface{}{
		"last_fetched_at": lastFetchedAt,
	}, nil
}

// MergeResourceEvents merges resource events into notes
func (n *NotesActions) mergeResourceEvents(ctx *gin.Context, notes []*service.Note) ([]*service.Note, error) {
	// Get notes filter
	notesFilter, err := n.notesFilter(ctx)
	if err != nil {
		return notes, nil
	}

	// Check if filter is only comments
	if notesFilter == service.UserPreferenceNotesFiltersOnlyComments {
		return notes, nil
	}

	// Get noteable from context
	noteable, err := ctx.Get("noteable")
	if err != nil {
		return notes, nil
	}

	// Get current user from context
	currentUser, err := ctx.Get("current_user")
	if err != nil {
		return notes, nil
	}

	// Get last fetched at
	lastFetchedAt, err := n.lastFetchedAt(ctx)
	if err != nil {
		return notes, nil
	}

	// Merge resource events
	return n.noteService.MergeResourceEventsIntoNotes(noteable, currentUser, lastFetchedAt, notes)
}

// NoteHTML renders a note to HTML
func (n *NotesActions) noteHTML(ctx *gin.Context, note *service.Note) (string, error) {
	// Render note to HTML
	return n.viewService.RenderPartial(ctx, "shared/notes/_note", map[string]interface{}{
		"note": note,
	})
}

// NoteJSON returns the note JSON
func (n *NotesActions) noteJSON(ctx *gin.Context, note *service.Note) map[string]interface{} {
	// Initialize attributes
	attrs := make(map[string]interface{})

	// Check if note is persisted
	if note.Persisted {
		attrs["valid"] = true

		// Check if return discussion is true
		if n.returnDiscussion(ctx) {
			// Get discussion
			discussion, err := note.GetDiscussion()
			if err != nil {
				n.logger.Error("Failed to get discussion", "error", err)
				return attrs
			}

			// Prepare notes for rendering
			_, err = n.prepareNotesForRendering(ctx, discussion.Notes)
			if err != nil {
				n.logger.Error("Failed to prepare notes for rendering", "error", err)
				return attrs
			}

			// Serialize discussion
			attrs["discussion"], err = n.discussionSerializer(ctx).Represent(discussion, ctx)
			if err != nil {
				n.logger.Error("Failed to serialize discussion", "error", err)
				return attrs
			}
		} else if n.useNoteSerializer(ctx) {
			// Serialize note
			serializedNote, err := n.noteSerializer(ctx).Represent(note)
			if err != nil {
				n.logger.Error("Failed to serialize note", "error", err)
				return attrs
			}

			// Merge serialized note
			for k, v := range serializedNote {
				attrs[k] = v
			}
		} else {
			// Get noteable from context
			noteable, err := ctx.Get("noteable")
			if err != nil {
				n.logger.Error("Failed to get noteable", "error", err)
				return attrs
			}

			// Render note to HTML
			html, err := n.noteHTML(ctx, note)
			if err != nil {
				n.logger.Error("Failed to render note to HTML", "error", err)
				return attrs
			}

			// Merge note attributes
			attrs["id"] = note.ID
			attrs["discussion_id"] = note.GetDiscussionID(noteable)
			attrs["html"] = html
			attrs["note"] = note.Note
			attrs["on_image"] = note.OnImage

			// Get discussion
			discussion, err := note.ToDiscussion(noteable)
			if err != nil {
				n.logger.Error("Failed to get discussion", "error", err)
				return attrs
			}

			// Check if discussion is individual note
			if !discussion.IndividualNote {
				// Merge discussion attributes
				attrs["discussion_resolvable"] = discussion.Resolvable

				// Render diff discussion to HTML
				diffDiscussionHTML, err := n.diffDiscussionHTML(ctx, discussion)
				if err != nil {
					n.logger.Error("Failed to render diff discussion to HTML", "error", err)
				} else {
					attrs["diff_discussion_html"] = diffDiscussionHTML
				}

				// Render discussion to HTML
				discussionHTML, err := n.discussionHTML(ctx, discussion)
				if err != nil {
					n.logger.Error("Failed to render discussion to HTML", "error", err)
				} else {
					attrs["discussion_html"] = discussionHTML
				}

				// Check if discussion is diff discussion
				if discussion.DiffDiscussion {
					attrs["discussion_line_code"] = discussion.LineCode
				}
			}
		}
	} else {
		// Merge error attributes
		attrs["valid"] = false
		attrs["errors"] = note.Errors
	}

	return attrs
}

// DiffDiscussionHTML renders a diff discussion to HTML
func (n *NotesActions) diffDiscussionHTML(ctx *gin.Context, discussion *service.Discussion) (string, error) {
	// Check if discussion is diff discussion
	if !discussion.DiffDiscussion {
		return "", nil
	}

	// Get on image
	onImage := discussion.OnImage

	// Get view from query
	view := ctx.DefaultQuery("view", "")

	// Get line type from query
	lineType := ctx.DefaultQuery("line_type", "")

	// Initialize template and locals
	var template string
	var locals map[string]interface{}

	// Check if view is parallel and not on image
	if view == "parallel" && !onImage {
		template = "discussions/_parallel_diff_discussion"

		// Check line type
		if lineType == "old" {
			locals = map[string]interface{}{
				"discussions_left":  []*service.Discussion{discussion},
				"discussions_right": nil,
			}
		} else {
			locals = map[string]interface{}{
				"discussions_left":  nil,
				"discussions_right": []*service.Discussion{discussion},
			}
		}
	} else {
		template = "discussions/_diff_discussion"

		// Set fresh discussion
		ctx.Set("fresh_discussion", true)

		// Set locals
		locals = map[string]interface{}{
			"discussions": []*service.Discussion{discussion},
			"on_image":    onImage,
		}
	}

	// Render template to HTML
	return n.viewService.RenderPartial(ctx, template, locals)
}

// DiscussionHTML renders a discussion to HTML
func (n *NotesActions) discussionHTML(ctx *gin.Context, discussion *service.Discussion) (string, error) {
	// Check if discussion is individual note
	if discussion.IndividualNote {
		return "", nil
	}

	// Render template to HTML
	return n.viewService.RenderPartial(ctx, "discussions/_discussion", map[string]interface{}{
		"discussion": discussion,
	})
}

// AuthorizeAdminNote authorizes admin note
func (n *NotesActions) authorizeAdminNote(ctx *gin.Context) error {
	// Get current user from context
	currentUser, err := ctx.Get("current_user")
	if err != nil {
		return err
	}

	// Get note from context
	note, err := ctx.Get("note")
	if err != nil {
		return err
	}

	// Check if user can admin note
	canAdminNote, err := n.userService.Can(currentUser, "admin_note", note)
	if err != nil {
		return err
	}

	// Check if user can admin note
	if !canAdminNote {
		return service.ErrAccessDenied
	}

	return nil
}

// CreateNoteParams returns the create note params
func (n *NotesActions) createNoteParams(ctx *gin.Context) (map[string]interface{}, error) {
	// Get note from request
	noteParams, err := ctx.GetPostForm("note")
	if err != nil {
		return nil, err
	}

	// Parse note params
	createParams := make(map[string]interface{})
	if err := ctx.ShouldBindJSON(&createParams); err != nil {
		return nil, err
	}

	// Get merge request diff head SHA from request
	mergeRequestDiffHeadSHA, err := ctx.GetPostForm("merge_request_diff_head_sha")
	if err == nil && mergeRequestDiffHeadSHA != "" {
		createParams["merge_request_diff_head_sha"] = mergeRequestDiffHeadSHA
	}

	// Get in reply to discussion ID from request
	inReplyToDiscussionID, err := ctx.GetPostForm("in_reply_to_discussion_id")
	if err == nil && inReplyToDiscussionID != "" {
		createParams["in_reply_to_discussion_id"] = inReplyToDiscussionID
	}

	// Get noteable from context
	noteable, err := ctx.Get("noteable")
	if err != nil {
		return nil, err
	}

	// Set noteable type
	createParams["noteable_type"] = noteable.(interface{ ClassName() string }).ClassName()

	// Set noteable ID based on type
	switch noteable.(type) {
	case *service.Commit:
		createParams["commit_id"] = noteable.(interface{ ID() string }).ID()
	case *service.MergeRequest:
		createParams["noteable_id"] = noteable.(interface{ ID() string }).ID()

		// Get commit ID from request
		commitID, err := ctx.GetPostForm("note.commit_id")
		if err == nil && commitID != "" {
			createParams["commit_id"] = commitID
		}
	default:
		createParams["noteable_id"] = noteable.(interface{ ID() string }).ID()
	}

	return createParams, nil
}

// UpdateNoteParams returns the update note params
func (n *NotesActions) updateNoteParams(ctx *gin.Context) (map[string]interface{}, error) {
	// Get note from request
	noteParams, err := ctx.GetPostForm("note")
	if err != nil {
		return nil, err
	}

	// Parse note params
	updateParams := make(map[string]interface{})
	if err := ctx.ShouldBindJSON(&updateParams); err != nil {
		return nil, err
	}

	return updateParams, nil
}

// SetPollingIntervalHeader sets the polling interval header
func (n *NotesActions) setPollingIntervalHeader(ctx *gin.Context, interval int) {
	// Set polling interval header
	ctx.Header("Poll-Interval", strconv.Itoa(interval))
}

// Noteable returns the noteable
func (n *NotesActions) noteable(ctx *gin.Context) (interface{}, error) {
	// Get notes finder
	notesFinder, err := n.notesFinder(ctx)
	if err != nil {
		return nil, err
	}

	// Get target from notes finder
	target, err := notesFinder.Target()
	if err != nil {
		return nil, err
	}

	// Check if target is not nil
	if target != nil {
		return target, nil
	}

	// Get note from context
	note, err := ctx.Get("note")
	if err != nil {
		return nil, err
	}

	// Check if note is not nil
	if note != nil {
		return note.(interface{ Noteable() interface{} }).Noteable(), nil
	}

	return nil, service.ErrNotFound
}

// RequireNoteable requires the noteable
func (n *NotesActions) requireNoteable(ctx *gin.Context) error {
	// Get noteable
	noteable, err := n.noteable(ctx)
	if err != nil {
		return err
	}

	// Check if noteable is not nil
	if noteable == nil {
		return service.ErrNotFound
	}

	return nil
}

// RequireLastFetchedAtHeader requires the last fetched at header
func (n *NotesActions) requireLastFetchedAtHeader(ctx *gin.Context) error {
	// Get last fetched at header
	lastFetchedAtHeader := ctx.GetHeader("X-Last-Fetched-At")

	// Check if last fetched at header is present
	if lastFetchedAtHeader == "" {
		return service.ErrBadRequest
	}

	return nil
}

// LastFetchedAt returns the last fetched at
func (n *NotesActions) lastFetchedAt(ctx *gin.Context) (time.Time, error) {
	// Get last fetched at header
	lastFetchedAtHeader := ctx.GetHeader("X-Last-Fetched-At")

	// Parse microseconds
	microseconds, err := strconv.ParseInt(lastFetchedAtHeader, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	// Calculate seconds and fraction
	seconds := microseconds / 1000000
	frac := microseconds % 1000000

	// Create time
	return time.Unix(seconds, frac*1000), nil
}

// NotesFilter returns the notes filter
func (n *NotesActions) notesFilter(ctx *gin.Context) (string, error) {
	// Get current user from context
	currentUser, err := ctx.Get("current_user")
	if err != nil {
		return "", err
	}

	// Get target type from query
	targetType := ctx.DefaultQuery("target_type", "")

	// Get notes filter
	return n.userService.NotesFilterFor(currentUser, targetType)
}

// NotesFinder returns the notes finder
func (n *NotesActions) notesFinder(ctx *gin.Context) (*service.NotesFinder, error) {
	// Get finder params
	finderParams, err := n.finderParams(ctx)
	if err != nil {
		return nil, err
	}

	// Get current user from context
	currentUser, err := ctx.Get("current_user")
	if err != nil {
		return nil, err
	}

	// Create notes finder
	return service.NewNotesFinder(currentUser, finderParams), nil
}

// FinderParams returns the finder params
func (n *NotesActions) finderParams(ctx *gin.Context) (map[string]interface{}, error) {
	// Initialize finder params
	finderParams := make(map[string]interface{})

	// Get target type from query
	targetType := ctx.DefaultQuery("target_type", "")
	if targetType != "" {
		finderParams["target_type"] = targetType
	}

	// Get target ID from query
	targetID := ctx.DefaultQuery("target_id", "")
	if targetID != "" {
		finderParams["target_id"] = targetID
	}

	// Get project ID from query
	projectID := ctx.DefaultQuery("project_id", "")
	if projectID != "" {
		finderParams["project_id"] = projectID
	}

	// Get group ID from query
	groupID := ctx.DefaultQuery("group_id", "")
	if groupID != "" {
		finderParams["group_id"] = groupID
	}

	return finderParams, nil
}

// NoteSerializer returns the note serializer
func (n *NotesActions) noteSerializer(ctx *gin.Context) *service.ProjectNoteSerializer {
	// Get project from context
	project, err := ctx.Get("project")
	if err != nil {
		return nil
	}

	// Get noteable from context
	noteable, err := ctx.Get("noteable")
	if err != nil {
		return nil
	}

	// Get current user from context
	currentUser, err := ctx.Get("current_user")
	if err != nil {
		return nil
	}

	// Create note serializer
	return service.NewProjectNoteSerializer(project, noteable, currentUser)
}

// DiscussionSerializer returns the discussion serializer
func (n *NotesActions) discussionSerializer(ctx *gin.Context) *service.DiscussionSerializer {
	// Get project from context
	project, err := ctx.Get("project")
	if err != nil {
		return nil
	}

	// Get noteable from context
	noteable, err := ctx.Get("noteable")
	if err != nil {
		return nil
	}

	// Get current user from context
	currentUser, err := ctx.Get("current_user")
	if err != nil {
		return nil
	}

	// Create discussion serializer
	return service.NewDiscussionSerializer(project, noteable, currentUser, service.ProjectNoteEntity)
}

// NoteProject returns the note project
func (n *NotesActions) noteProject(ctx *gin.Context) (interface{}, error) {
	// Get project from context
	project, err := ctx.Get("project")
	if err != nil {
		return nil, err
	}

	// Check if project is not nil
	if project == nil {
		return nil, service.ErrNotFound
	}

	// Get note project ID from query
	noteProjectID := ctx.DefaultQuery("note_project_id", "")

	// Initialize the project
	var theProject interface{}

	// Check if note project ID is present
	if noteProjectID != "" {
		// Find project by ID
		theProject, err = n.projectService.Find(noteProjectID)
		if err != nil {
			return nil, err
		}
	} else {
		// Use project from context
		theProject = project
	}

	// Get current user from context
	currentUser, err := ctx.Get("current_user")
	if err != nil {
		return nil, err
	}

	// Check if user can create note
	canCreateNote, err := n.userService.Can(currentUser, "create_note", theProject)
	if err != nil {
		return nil, err
	}

	// Check if user can create note
	if !canCreateNote {
		return nil, service.ErrAccessDenied
	}

	return theProject, nil
}

// ReturnDiscussion returns whether to return discussion
func (n *NotesActions) returnDiscussion(ctx *gin.Context) bool {
	// Get return discussion from query
	returnDiscussion := ctx.DefaultQuery("return_discussion", "false")

	// Check if return discussion is true
	return returnDiscussion == "true"
}

// UseNoteSerializer returns whether to use note serializer
func (n *NotesActions) useNoteSerializer(ctx *gin.Context) bool {
	// Get HTML from query
	html := ctx.DefaultQuery("html", "")

	// Check if HTML is present
	if html != "" {
		return false
	}

	// Get noteable from context
	noteable, err := ctx.Get("noteable")
	if err != nil {
		return false
	}

	// Check if noteable discussions are rendered on frontend
	return noteable.(interface{ DiscussionsRenderedOnFrontend() bool }).DiscussionsRenderedOnFrontend()
}

// ErrorsOnCreate returns the errors on create
func (n *NotesActions) errorsOnCreate(note *service.Note) string {
	// Check if errors are present
	if note.Errors == nil || len(note.Errors) == 0 {
		return ""
	}

	// Return errors as sentence
	return note.Errors.FullMessages().ToSentence()
}

// PrepareNotesForRendering prepares notes for rendering
func (n *NotesActions) prepareNotesForRendering(ctx *gin.Context, notes []*service.Note) ([]*service.Note, error) {
	// Get current user from context
	currentUser, err := ctx.Get("current_user")
	if err != nil {
		return notes, err
	}

	// Prepare notes for rendering
	return n.noteService.PrepareNotesForRendering(notes, currentUser)
}

// FilterNotesByReadability filters notes by readability
func (n *NotesActions) filterNotesByReadability(notes []*service.Note, currentUser interface{}) []*service.Note {
	// Initialize filtered notes
	filteredNotes := make([]*service.Note, 0, len(notes))

	// Filter notes by readability
	for _, note := range notes {
		// Check if note is readable by user
		if note.ReadableBy(currentUser) {
			filteredNotes = append(filteredNotes, note)
		}
	}

	return filteredNotes
}

// SerializeNotes serializes notes
func (n *NotesActions) serializeNotes(ctx *gin.Context, notes []*service.Note) []map[string]interface{} {
	// Initialize serialized notes
	serializedNotes := make([]map[string]interface{}, len(notes))

	// Serialize notes
	for i, note := range notes {
		serializedNotes[i] = n.noteJSON(ctx, note)
	}

	return serializedNotes
}
