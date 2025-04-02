package ui

import (
	"fmt"
	"strings"
)

// AvatarSize represents the possible sizes for an avatar
type AvatarSize int

const (
	// AvatarSize16 represents a 16px avatar
	AvatarSize16 AvatarSize = 16
	// AvatarSize24 represents a 24px avatar
	AvatarSize24 AvatarSize = 24
	// AvatarSize32 represents a 32px avatar
	AvatarSize32 AvatarSize = 32
	// AvatarSize48 represents a 48px avatar
	AvatarSize48 AvatarSize = 48
	// AvatarSize64 represents a 64px avatar
	AvatarSize64 AvatarSize = 64
	// AvatarSize96 represents a 96px avatar
	AvatarSize96 AvatarSize = 96
)

// AvatarItem represents an item that can be displayed as an avatar
type AvatarItem interface {
	GetName() string
	GetAvatarURL() string
	GetID() int
}

// AvatarEmail represents an email address that can be displayed as an avatar
type AvatarEmail struct {
	Email string
}

// GetName returns the email address as the name
func (e AvatarEmail) GetName() string {
	return e.Email
}

// GetAvatarURL returns an empty string as there is no avatar URL for an email
func (e AvatarEmail) GetAvatarURL() string {
	return ""
}

// GetID returns 0 as there is no ID for an email
func (e AvatarEmail) GetID() int {
	return 0
}

// Avatar represents an avatar component
type Avatar struct {
	Item          AvatarItem
	Alt           string
	CustomClass   string
	Size          AvatarSize
	AvatarOptions map[string]interface{}
}

// NewAvatar creates a new avatar with the given parameters
func NewAvatar(item AvatarItem, alt string, customClass string, size AvatarSize, avatarOptions map[string]interface{}) *Avatar {
	if avatarOptions == nil {
		avatarOptions = make(map[string]interface{})
	}

	return &Avatar{
		Item:          item,
		Alt:           alt,
		CustomClass:   customClass,
		Size:          size,
		AvatarOptions: avatarOptions,
	}
}

// AvatarClasses returns the CSS classes for the avatar
func (a *Avatar) AvatarClasses() string {
	var classes []string

	// Add base classes
	classes = append(classes, "gl-avatar", fmt.Sprintf("gl-avatar-s%d", a.Size), a.CustomClass)

	// Add shape class based on item type
	switch a.Item.(type) {
	case AvatarEmail:
		classes = append(classes, "gl-avatar-circle")
	default:
		classes = append(classes, "!gl-rounded-base")
	}

	// Add identicon classes if no source
	if a.Src() == "" {
		classes = append(classes, "gl-avatar-identicon")
		classes = append(classes, fmt.Sprintf("gl-avatar-identicon-bg%d", (a.Item.GetID()%7)+1))
	}

	return strings.Join(classes, " ")
}

// Src returns the source URL for the avatar
func (a *Avatar) Src() string {
	if a.Item == nil {
		return ""
	}

	// If the item has an avatar URL, use it
	if url := a.Item.GetAvatarURL(); url != "" {
		return fmt.Sprintf("%s?width=%d", url, a.Size)
	}

	return ""
}

// Srcset returns the srcset attribute for the avatar
func (a *Avatar) Srcset() string {
	src := a.Src()
	if src == "" {
		return ""
	}

	retinaSrc := strings.Replace(src, fmt.Sprintf("width=%d", a.Size), fmt.Sprintf("width=%d", a.Size*2), 1)
	return fmt.Sprintf("%s 1x, %s 2x", src, retinaSrc)
}

// Alt returns the alt text for the avatar
func (a *Avatar) AltText() string {
	if a.Alt != "" {
		return a.Alt
	}
	return a.Item.GetName()
}

// Initial returns the first letter of the item's name
func (a *Avatar) Initial() string {
	name := a.Item.GetName()
	if name == "" {
		return ""
	}
	return strings.ToUpper(string(name[0]))
}

// Render generates the HTML for the avatar
func (a *Avatar) Render() string {
	// Build avatar attributes
	avatarAttrs := make(map[string]interface{})
	for k, v := range a.AvatarOptions {
		avatarAttrs[k] = v
	}
	avatarAttrs["class"] = a.AvatarClasses()
	avatarAttrs["alt"] = a.AltText()

	// If we have a source, render an image
	if src := a.Src(); src != "" {
		avatarAttrs["src"] = src
		avatarAttrs["srcset"] = a.Srcset()
		avatarAttrs["height"] = a.Size
		avatarAttrs["width"] = a.Size
		avatarAttrs["loading"] = "lazy"

		return fmt.Sprintf(`<img %s>`, formatAttributes(avatarAttrs))
	}

	// Otherwise, render a div with the initial
	return fmt.Sprintf(`<div %s>%s</div>`, formatAttributes(avatarAttrs), a.Initial())
}
