package model

import (
	"time"
)

// ProjectSortingPreferenceField is the field used for project sorting preference
const ProjectSortingPreferenceField = "project_sorting_preference"

// Project represents a GitLab project
type Project struct {
	ID                int64     `json:"id"`
	Name              string    `json:"name"`
	Path              string    `json:"path"`
	Description       string    `json:"description"`
	Visibility        string    `json:"visibility"`
	Archived          bool      `json:"archived"`
	StarCount         int       `json:"star_count"`
	ForksCount        int       `json:"forks_count"`
	LastActivityAt    time.Time `json:"last_activity_at"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	CreatorID         int64     `json:"creator_id"`
	NamespaceID       int64     `json:"namespace_id"`
	DefaultBranch     string    `json:"default_branch"`
	WebURL            string    `json:"web_url"`
	SSHURLToRepo      string    `json:"ssh_url_to_repo"`
	HTTPURLToRepo     string    `json:"http_url_to_repo"`
	ReadmeURL         string    `json:"readme_url"`
	TagList           []string  `json:"tag_list"`
	Topics            []string  `json:"topics"`
	EmptyRepo         bool      `json:"empty_repo"`
	Public            bool      `json:"public"`
	IssuesEnabled     bool      `json:"issues_enabled"`
	WikiEnabled       bool      `json:"wiki_enabled"`
	BuildsEnabled     bool      `json:"builds_enabled"`
	SnippetsEnabled   bool      `json:"snippets_enabled"`
	ContainerRegistryEnabled bool `json:"container_registry_enabled"`
	SharedRunnersEnabled     bool `json:"shared_runners_enabled"`
	LFSEnabled              bool `json:"lfs_enabled"`
	RequestAccessEnabled    bool `json:"request_access_enabled"`
	MergeRequestsEnabled    bool `json:"merge_requests_enabled"`
	OnlyAllowMergeIfPipelineSucceeds bool `json:"only_allow_merge_if_pipeline_succeeds"`
	OnlyAllowMergeIfAllDiscussionsAreResolved bool `json:"only_allow_merge_if_all_discussions_are_resolved"`
	RemoveSourceBranchAfterMerge bool `json:"remove_source_branch_after_merge"`
	ApprovalsBeforeMerge         int  `json:"approvals_before_merge"`
	ResolveOutdatedDiffDiscussions bool `json:"resolve_outdated_diff_discussions"`
	ImportError                   string `json:"import_error"`
	ImportStatus                  string `json:"import_status"`
	OpenIssuesCount               int    `json:"open_issues_count"`
	AvatarURL                     string `json:"avatar_url"`
	Starred                       bool   `json:"starred"`
}

// RenderedProject represents a rendered project for display
type RenderedProject struct {
	*Project
	PathWithNamespace string `json:"path_with_namespace"`
	FullPath          string `json:"full_path"`
	StarredByUser     bool   `json:"starred_by_user"`
}
