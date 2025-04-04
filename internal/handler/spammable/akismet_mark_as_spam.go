package spammable

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
	"gitlab.com/gitlab-org/gitlab-redux/internal/util"
)

// AkismetMarkAsSpamAction handles marking content as spam using Akismet
type AkismetMarkAsSpamAction struct {
	akismetService *service.AkismetMarkAsSpamService
	logger         *util.Logger
}

// NewAkismetMarkAsSpamAction creates a new instance of AkismetMarkAsSpamAction
func NewAkismetMarkAsSpamAction(akismetService *service.AkismetMarkAsSpamService, logger *util.Logger) *AkismetMarkAsSpamAction {
	return &AkismetMarkAsSpamAction{
		akismetService: akismetService,
		logger:         logger,
	}
}

// MarkAsSpam marks the spammable content as spam using Akismet
func (a *AkismetMarkAsSpamAction) MarkAsSpam(ctx *gin.Context, spammable model.Spammable, spammablePath string) {
	// Check if the user is authorized to mark content as spam
	if !a.authorizeSubmitSpammable(ctx) {
		ctx.JSON(403, gin.H{"error": "Access denied"})
		return
	}

	// Submit the content to Akismet
	success, err := a.akismetService.Execute(spammable)
	if err != nil {
		a.logger.Error("Error submitting to Akismet", "error", err)
		ctx.Redirect(302, spammablePath)
		ctx.Set("flash_alert", "Error with Akismet. Please check the logs for more info.")
		return
	}

	if success {
		// Format the success message
		entityType := spammable.GetSpammableEntityType()
		titleCaseEntityType := util.TitleCase(entityType)
		successMessage := fmt.Sprintf("%s was submitted to Akismet successfully.", titleCaseEntityType)

		ctx.Redirect(302, spammablePath)
		ctx.Set("flash_notice", successMessage)
	} else {
		ctx.Redirect(302, spammablePath)
		ctx.Set("flash_alert", "Error with Akismet. Please check the logs for more info.")
	}
}

// authorizeSubmitSpammable checks if the current user is authorized to mark content as spam
func (a *AkismetMarkAsSpamAction) authorizeSubmitSpammable(ctx *gin.Context) bool {
	// Get the current user from the context
	currentUser, exists := ctx.Get("current_user")
	if !exists {
		return false
	}

	// Check if the user can admin all resources
	// This assumes there's a method on the user model to check this
	user, ok := currentUser.(model.User)
	if !ok {
		return false
	}

	return user.CanAdminAllResources()
}
