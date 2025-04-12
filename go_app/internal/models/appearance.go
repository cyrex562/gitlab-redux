package models

import (
	"time"
)

// Appearance represents the application's appearance settings
type Appearance struct {
	ID                           uint      `json:"id" gorm:"primaryKey"`
	Title                        string    `json:"title"`
	Description                  string    `json:"description"`
	PWAName                      string    `json:"pwa_name"`
	PWAShortName                 string    `json:"pwa_short_name"`
	PWADescription               string    `json:"pwa_description"`
	Logo                         string    `json:"logo"`
	LogoCache                    string    `json:"logo_cache"`
	HeaderLogo                   string    `json:"header_logo"`
	HeaderLogoCache              string    `json:"header_logo_cache"`
	PWAIcon                      string    `json:"pwa_icon"`
	PWAIconCache                 string    `json:"pwa_icon_cache"`
	Favicon                      string    `json:"favicon"`
	FaviconCache                 string    `json:"favicon_cache"`
	MemberGuidelines            string    `json:"member_guidelines"`
	NewProjectGuidelines        string    `json:"new_project_guidelines"`
	ProfileImageGuidelines      string    `json:"profile_image_guidelines"`
	UpdatedBy                    uint      `json:"updated_by"`
	HeaderMessage               string    `json:"header_message"`
	FooterMessage               string    `json:"footer_message"`
	MessageBackgroundColor      string    `json:"message_background_color"`
	MessageFontColor            string    `json:"message_font_color"`
	EmailHeaderAndFooterEnabled bool      `json:"email_header_and_footer_enabled"`
	CreatedAt                   time.Time `json:"created_at"`
	UpdatedAt                   time.Time `json:"updated_at"`
}

// GetCurrentAppearance returns the current appearance settings
func GetCurrentAppearance() *Appearance {
	var appearance Appearance
	// TODO: Implement database query to get current appearance
	return &appearance
}

// Save saves the appearance settings to the database
func (a *Appearance) Save() error {
	// TODO: Implement database save
	return nil
}

// RemoveLogo removes the logo from appearance settings
func (a *Appearance) RemoveLogo() {
	a.Logo = ""
	a.LogoCache = ""
}

// RemoveHeaderLogo removes the header logo from appearance settings
func (a *Appearance) RemoveHeaderLogo() {
	a.HeaderLogo = ""
	a.HeaderLogoCache = ""
}

// RemovePWAIcon removes the PWA icon from appearance settings
func (a *Appearance) RemovePWAIcon() {
	a.PWAIcon = ""
	a.PWAIconCache = ""
}

// RemoveFavicon removes the favicon from appearance settings
func (a *Appearance) RemoveFavicon() {
	a.Favicon = ""
	a.FaviconCache = ""
}
