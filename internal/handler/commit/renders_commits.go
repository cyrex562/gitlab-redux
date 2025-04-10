package commit

import (
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
)

// RendersCommits handles rendering of commits in GitLab
type RendersCommits struct {
	commitRenderer CommitRenderer
	gitRefValidator GitRefValidator
}

// CommitRenderer defines the interface for rendering commits
type CommitRenderer interface {
	Render(commits []*model.Commit, project *model.Project, currentUser *model.User)
}

// GitRefValidator defines the interface for validating Git references
type GitRefValidator interface {
	Validate(refName string, skipHeadRefCheck bool) bool
}

// NewRendersCommits creates a new instance of RendersCommits
func NewRendersCommits(commitRenderer CommitRenderer, gitRefValidator GitRefValidator) *RendersCommits {
	return &RendersCommits{
		commitRenderer: commitRenderer,
		gitRefValidator: gitRefValidator,
	}
}

// CommitsSafeSize is the maximum number of commits to render safely
const CommitsSafeSize = 1000

// LimitedCommits limits the number of commits based on a safe size
func (r *RendersCommits) LimitedCommits(commits []*model.Commit, commitsCount int) ([]*model.Commit, int) {
	if commitsCount > CommitsSafeSize {
		limitedCommits := make([]*model.Commit, CommitsSafeSize)
		copy(limitedCommits, commits[:CommitsSafeSize])
		return limitedCommits, commitsCount - CommitsSafeSize
	}

	return commits, 0
}

// SetCommitsForRendering sets up commits for rendering with proper counts
func (r *RendersCommits) SetCommitsForRendering(commits []*model.Commit, commitsCount *int) ([]*model.Commit, int, int) {
	var totalCommitCount int
	if commitsCount != nil {
		totalCommitCount = *commitsCount
	} else {
		totalCommitCount = len(commits)
	}

	limited, hiddenCommitCount := r.LimitedCommits(commits, totalCommitCount)
	preparedCommits := r.PrepareCommitsForRendering(limited)

	return preparedCommits, totalCommitCount, hiddenCommitCount
}

// PrepareCommitsForRendering prepares commits for rendering by loading authors and pipelines
func (r *RendersCommits) PrepareCommitsForRendering(commits []*model.Commit) []*model.Commit {
	// Preload commits' authors
	for _, commit := range commits {
		commit.LazyAuthor()
	}

	// Preload commits' latest pipelines
	for _, commit := range commits {
		commit.LazyLatestPipeline()
	}

	// Note: In the Go version, we don't have direct access to @project and current_user
	// These would be passed as parameters in a real implementation
	// For now, we'll assume the commit renderer has access to the necessary context

	return commits
}

// ValidRef checks if a Git reference name is valid
func (r *RendersCommits) ValidRef(refName string) bool {
	if refName == "" {
		return true
	}

	return r.gitRefValidator.Validate(refName, true)
}
