package ui

import (
	"fmt"
	"strings"
)

// BroadcastBanner represents a broadcast banner component
type BroadcastBanner struct {
	Message       string
	ID            string
	Theme         string
	Dismissable   bool
	ExpireDate    string
	CookieKey     string
	DismissalPath string
	ButtonTestID  string
	Banner        string
}

// NewBroadcastBanner creates a new broadcast banner with the given parameters
func NewBroadcastBanner(message string, id string, theme string, dismissable bool, expireDate string, cookieKey string, dismissalPath string, buttonTestID string, banner string) *BroadcastBanner {
	return &BroadcastBanner{
		Message:       message,
		ID:            id,
		Theme:         theme,
		Dismissable:   dismissable,
		ExpireDate:    expireDate,
		CookieKey:     cookieKey,
		DismissalPath: dismissalPath,
		ButtonTestID:  buttonTestID,
		Banner:        banner,
	}
}

// Render generates the HTML for the broadcast banner
func (b *BroadcastBanner) Render() string {
	// Build banner attributes
	bannerAttrs := make(map[string]interface{})

	// Add default attributes
	bannerAttrs["role"] = "alert"
	bannerAttrs["class"] = fmt.Sprintf("js-broadcast-notification-%s %s", b.ID, b.Theme)
	bannerAttrs["data-testid"] = "banner-broadcast-message"
	bannerAttrs["data-broadcast-banner"] = b.Banner

	// Build the content
	var contentParts []string

	// Add content
	contentParts = append(contentParts, `<div class="gl-broadcast-message-content">`)
	contentParts = append(contentParts, `<div class="gl-broadcast-message-icon"><span class="gl-icon">bullhorn</span></div>`)
	contentParts = append(contentParts, `<div class="gl-broadcast-message-text">`)
	contentParts = append(contentParts, `<h2 class="gl-sr-only">Admin message</h2>`)
	contentParts = append(contentParts, b.Message)
	contentParts = append(contentParts, `</div></div>`)

	// Add dismiss button if needed
	if b.Dismissable {
		buttonAttrs := make(map[string]interface{})
		buttonAttrs["class"] = "gl-broadcast-message-dismiss js-dismiss-current-broadcast-notification"
		buttonAttrs["aria-label"] = "Close"
		buttonAttrs["data-id"] = b.ID
		buttonAttrs["data-expire-date"] = b.ExpireDate
		buttonAttrs["data-dismissal-path"] = b.DismissalPath
		buttonAttrs["data-cookie-key"] = b.CookieKey
		buttonAttrs["data-testid"] = b.ButtonTestID
		buttonAttrs["category"] = "tertiary"
		buttonAttrs["icon"] = "close"
		buttonAttrs["size"] = "small"

		contentParts = append(contentParts, fmt.Sprintf(`<button type="button" %s><span class="js-close-icon">close</span></button>`, formatAttributes(buttonAttrs)))
	}

	// Render the broadcast banner
	return fmt.Sprintf(`<div %s>%s</div>`,
		formatAttributes(bannerAttrs),
		strings.Join(contentParts, "\n"))
}
