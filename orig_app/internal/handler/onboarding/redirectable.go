package onboarding

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// Redirectable provides functionality for handling redirects after sign-up
type Redirectable struct {
	onboardingService *service.OnboardingService
}

// NewRedirectable creates a new instance of Redirectable
func NewRedirectable(onboardingService *service.OnboardingService) *Redirectable {
	return &Redirectable{
		onboardingService: onboardingService,
	}
}

// AfterSignUpPath returns the path to redirect to after sign-up
func (r *Redirectable) AfterSignUpPath(ctx *gin.Context) string {
	// Get the onboarding status presenter
	presenter := r.onboardingService.GetOnboardingStatusPresenter(ctx)

	// Check if there's a single invite
	if presenter.IsSingleInvite() {
		// Get the last invited member
		member := presenter.GetLastInvitedMember()

		// Set flash notice
		ctx.SetCookie("notice", r.getInviteAcceptedNotice(ctx, member), 3600, "/", "", false, true)

		// Return the polymorphic path for the last invited member source
		return r.getPolymorphicPath(ctx, presenter.GetLastInvitedMemberSource())
	}

	// If there are multiple invites, return the path for signed-in user
	return r.PathForSignedInUser(ctx)
}

// PathForSignedInUser returns the path for a signed-in user
func (r *Redirectable) PathForSignedInUser(ctx *gin.Context) string {
	// Get the stored location for the user
	storedLocation := r.getStoredLocation(ctx, "user")
	if storedLocation != "" {
		return storedLocation
	}

	// If no stored location, return the last member source path
	return r.LastMemberSourcePath(ctx)
}

// LastMemberSourcePath returns the path for the last member source
func (r *Redirectable) LastMemberSourcePath(ctx *gin.Context) string {
	// Get the onboarding status presenter
	presenter := r.onboardingService.GetOnboardingStatusPresenter(ctx)

	// Get the last invited member source
	source := presenter.GetLastInvitedMemberSource()
	if source == nil {
		return "/dashboard/projects"
	}

	// Return the polymorphic path for the source
	return r.getPolymorphicPath(ctx, source)
}

// getStoredLocation gets the stored location for a key
func (r *Redirectable) getStoredLocation(ctx *gin.Context, key string) string {
	// Get the stored location from the session
	location, exists := ctx.Get("stored_location_" + key)
	if !exists {
		return ""
	}
	return location.(string)
}

// getPolymorphicPath gets the polymorphic path for a source
func (r *Redirectable) getPolymorphicPath(ctx *gin.Context, source interface{}) string {
	// This is a simplified version of Rails' polymorphic_path
	// In a real implementation, this would be more complex

	switch s := source.(type) {
	case *model.Project:
		return "/projects/" + s.ID
	case *model.Group:
		return "/groups/" + s.ID
	case *model.Namespace:
		return "/namespaces/" + s.ID
	default:
		return "/"
	}
}

// getInviteAcceptedNotice gets the invite accepted notice
func (r *Redirectable) getInviteAcceptedNotice(ctx *gin.Context, member *model.Member) string {
	// This is a simplified version of Rails' invite_accepted_notice
	// In a real implementation, this would be more complex

	return "You have been invited to join " + member.SourceType + " " + member.SourceName
}
