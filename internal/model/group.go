package model

import (
	"time"
)

// Group represents a GitLab group
type Group struct {
	ID                          uint      `gorm:"primarykey" json:"id"`
	CreatedAt                   time.Time `json:"created_at"`
	UpdatedAt                   time.Time `json:"updated_at"`
	Name                        string    `json:"name"`
	Path                        string    `json:"path"`
	Description                 string    `json:"description"`
	VisibilityLevel            int       `json:"visibility_level"`
	LFSEnabled                 bool      `json:"lfs_enabled"`
	RequestAccessEnabled       bool      `json:"request_access_enabled"`
	RunnerRegistrationEnabled  bool      `json:"runner_registration_enabled"`
	RequireTwoFactorAuth       bool      `json:"require_two_factor_authentication"`
	TwoFactorGracePeriod       int       `json:"two_factor_grace_period"`
	EnabledGitAccessProtocol   string    `json:"enabled_git_access_protocol"`
	ProjectCreationLevel       int       `json:"project_creation_level"`
	SubgroupCreationLevel      int       `json:"subgroup_creation_level"`
	OrganizationID            uint      `json:"organization_id"`
	AdminNote                 *AdminNote `json:"admin_note,omitempty" gorm:"foreignKey:GroupID"`
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
