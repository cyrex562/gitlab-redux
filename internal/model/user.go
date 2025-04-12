package model

import "time"

// User represents a GitLab user
type User struct {
	ID                int64     `json:"id"`
	Name              string    `json:"name"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	State             string    `json:"state"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	Bio               string    `json:"bio"`
	Location          string    `json:"location"`
	PublicEmail       string    `json:"public_email"`
	Skype             string    `json:"skype"`
	Linkedin          string    `json:"linkedin"`
	Twitter           string    `json:"twitter"`
	WebsiteURL        string    `json:"website_url"`
	Organization      string    `json:"organization"`
	JobTitle          string    `json:"job_title"`
	WorkInfo          string    `json:"work_info"`
	TwoFactorEnabled  bool      `json:"two_factor_enabled"`
	External          bool      `json:"external"`
	PrivateProfile    bool      `json:"private_profile"`
	CommitEmail       string    `json:"commit_email"`
	ThemeID           int       `json:"theme_id"`
	ColorSchemeID     int       `json:"color_scheme_id"`
	ProjectsLimit     int       `json:"projects_limit"`
	CurrentSignInAt   time.Time `json:"current_sign_in_at"`
	LastSignInAt      time.Time `json:"last_sign_in_at"`
	CurrentSignInIP   string    `json:"current_sign_in_ip"`
	LastSignInIP      string    `json:"last_sign_in_ip"`
	Admin             bool      `json:"admin"`
	CanCreateGroup    bool      `json:"can_create_group"`
	CanCreateProject  bool      `json:"can_create_project"`
	TwoFactorGracePeriod int    `json:"two_factor_grace_period"`
	AvatarURL         string    `json:"avatar_url"`
}

// UserParams represents the parameters for creating or updating a user
type UserParams struct {
	Username         string `json:"username" binding:"required"`
	Email            string `json:"email" binding:"required,email"`
	Name             string `json:"name" binding:"required"`
	Password         string `json:"password"`
	ResetPassword    bool   `json:"reset_password"`
	SkipConfirmation bool   `json:"skip_confirmation"`
	Admin            bool   `json:"admin"`
	// Add other necessary fields
}

// IsAdmin checks if the user has admin privileges
func (u *User) IsAdmin() bool {
	return u.Admin
}

// IsActive checks if the user is active
func (u *User) IsActive() bool {
	return u.State == "active"
}

// IsBlocked checks if the user is blocked
func (u *User) IsBlocked() bool {
	return u.State == "blocked"
}

// IsBanned checks if the user is banned
func (u *User) IsBanned() bool {
	return u.State == "banned"
}

// IsDeactivated checks if the user is deactivated
func (u *User) IsDeactivated() bool {
	return u.State == "deactivated"
}

// IsPendingApproval checks if the user is pending approval
func (u *User) IsPendingApproval() bool {
	return u.State == "pending_approval"
}

// IsLocked checks if the user is locked
func (u *User) IsLocked() bool {
	return u.State == "locked"
}

// IsTrustedWithSpamCheck checks if the user is trusted with spam checks
func (u *User) IsTrustedWithSpamCheck() bool {
	return u.TrustedWithSpamCheck
}

// RenderedUser represents a rendered user for display
type RenderedUser struct {
	*User
	Path              string `json:"path"`
	FullPath          string `json:"full_path"`
	Status            string `json:"status"`
	StatusEmoji       string `json:"status_emoji"`
	StatusMessage     string `json:"status_message"`
	StatusMessageHtml string `json:"status_message_html"`
}
