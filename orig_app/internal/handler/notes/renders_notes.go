package notes

import (
	"gitlab.com/gitlab-org/gitlab-redux/internal/access"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
)

// RendersNotes handles rendering of notes in GitLab
type RendersNotes struct {
	preloader Preloader
	teamService TeamService
	renderService RenderService
}

// Preloader defines the interface for preloading associations
type Preloader interface {
	Preload(records interface{}, associations interface{})
}

// TeamService defines the interface for team-related operations
type TeamService interface {
	MaxMemberAccessForUserIDs(userIDs []int64) map[int64]int
	ContributionCheckForUserIDs(userIDs []int64) map[int64]bool
}

// RenderService defines the interface for rendering notes
type RenderService interface {
	Execute(notes []*model.Note, currentUser *model.User)
}

// NewRendersNotes creates a new instance of RendersNotes
func NewRendersNotes(preloader Preloader, teamService TeamService, renderService RenderService) *RendersNotes {
	return &RendersNotes{
		preloader: preloader,
		teamService: teamService,
		renderService: renderService,
	}
}

// PrepareNotesForRendering prepares notes for rendering by preloading various associations
func (r *RendersNotes) PrepareNotesForRendering(notes []*model.Note, project *model.Project, currentUser *model.User) []*model.Note {
	r.preloadNoteableForRegularNotes(notes)
	r.preloadNoteNamespace(notes)
	r.preloadMaxAccessForAuthors(notes, project)
	r.preloadAuthorStatus(notes)
	r.renderService.Execute(notes, currentUser)

	return notes
}

// preloadNoteNamespace preloads the namespace for notes
func (r *RendersNotes) preloadNoteNamespace(notes []*model.Note) {
	r.preloader.Preload(notes, "namespace")
}

// preloadMaxAccessForAuthors preloads the maximum access level for note authors
func (r *RendersNotes) preloadMaxAccessForAuthors(notes []*model.Note, project *model.Project) {
	if project == nil {
		return
	}

	// Extract author IDs from notes
	authorIDs := make([]int64, 0, len(notes))
	for _, note := range notes {
		authorIDs = append(authorIDs, note.AuthorID)
	}

	// Get max member access for authors
	accessMap := r.teamService.MaxMemberAccessForUserIDs(authorIDs)

	// Filter authors with NO_ACCESS
	noAccessUserIDs := make([]int64, 0)
	for userID, accessLevel := range accessMap {
		if accessLevel == access.NoAccess {
			noAccessUserIDs = append(noAccessUserIDs, userID)
		}
	}

	// Check contribution for users with no access
	if len(noAccessUserIDs) > 0 {
		r.teamService.ContributionCheckForUserIDs(noAccessUserIDs)
	}
}

// preloadNoteableForRegularNotes preloads the noteable object for regular notes
func (r *RendersNotes) preloadNoteableForRegularNotes(notes []*model.Note) {
	// Filter notes that are not for commits
	regularNotes := make([]*model.Note, 0, len(notes))
	for _, note := range notes {
		if !note.ForCommit() {
			regularNotes = append(regularNotes, note)
		}
	}

	if len(regularNotes) > 0 {
		r.preloader.Preload(regularNotes, "noteable")
	}
}

// preloadAuthorStatus preloads the author status for notes
func (r *RendersNotes) preloadAuthorStatus(notes []*model.Note) {
	r.preloader.Preload(notes, map[string]string{"author": "status"})
}
