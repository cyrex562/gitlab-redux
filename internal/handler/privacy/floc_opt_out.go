package privacy

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// FlocOptOut handles the Federated Learning of Cohorts (FLoC) opt-out functionality
type FlocOptOut struct {
	settingsService *service.SettingsService
}

// NewFlocOptOut creates a new instance of FlocOptOut
func NewFlocOptOut(settingsService *service.SettingsService) *FlocOptOut {
	return &FlocOptOut{
		settingsService: settingsService,
	}
}

// FlocOptOutMiddleware creates a middleware that sets the FLoC opt-out header
func (f *FlocOptOut) FlocOptOutMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Skip if FLoC is enabled
		if f.IsFlocEnabled() {
			ctx.Next()
			return
		}

		// Set the FLoC opt-out header
		f.SetFlocOptOutHeader(ctx)
		ctx.Next()
	}
}

// IsFlocEnabled checks if FLoC is enabled in the current settings
func (f *FlocOptOut) IsFlocEnabled() bool {
	return f.settingsService.IsFlocEnabled()
}

// SetFlocOptOutHeader sets the FLoC opt-out header in the response
func (f *FlocOptOut) SetFlocOptOutHeader(ctx *gin.Context) {
	ctx.Header("Permissions-Policy", "interest-cohort=()")
}
