package model

import (
	"time"
)

// Group represents a GitLab group
type Group struct {
	ID                int64     `json:"id"`
	Name              string    `json:"name"`
	Path              string    `json:"path"`
	Description       string    `json:"description"`
	Visibility        string    `json:"visibility"`
	LFSEnabled        bool      `json:"lfs_enabled"`
	AvatarURL         string    `json:"avatar_url"`
	WebURL            string    `json:"web_url"`
	RequestAccessEnabled bool   `json:"request_access_enabled"`
	FullName          string    `json:"full_name"`
	FullPath          string    `json:"full_path"`
	ParentID          int64     `json:"parent_id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	ProjectsCount     int       `json:"projects_count"`
	MembersCount      int       `json:"members_count"`
	SharedWithGroups  []*Group  `json:"shared_with_groups"`
	SharedProjects    []*Project `json:"shared_projects"`
	Projects          []*Project `json:"projects"`
}

// RenderedGroup represents a rendered group for display
type RenderedGroup struct {
	*Group
	PathWithNamespace string `json:"path_with_namespace"`
	FullPath          string `json:"full_path"`
	HasParent         bool   `json:"has_parent"`
	ParentName        string `json:"parent_name"`
	ParentPath        string `json:"parent_path"`
}

// GroupDetails represents a group with its associated data
type GroupDetails struct {
	Group     Group         `json:"group"`
	Members   []GroupMember `json:"members"`
	Projects  []Project     `json:"projects"`
}

// GroupMember represents a member of a group
type GroupMember struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	GroupID     uint      `json:"group_id"`
	UserID      uint      `json:"user_id"`
	AccessLevel int       `json:"access_level"`
}

// AdminNote represents an administrative note for a group
type AdminNote struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	GroupID   uint      `json:"group_id"`
	Note      string    `json:"note"`
}

// TableName specifies the table name for Group
func (Group) TableName() string {
	return "groups"
}

// TableName specifies the table name for GroupMember
func (GroupMember) TableName() string {
	return "group_members"
}

// TableName specifies the table name for AdminNote
func (AdminNote) TableName() string {
	return "admin_notes"
}
