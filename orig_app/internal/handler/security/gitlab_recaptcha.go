package security

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
	"gitlab.com/gitlab-org/gitlab-redux/internal/util"
)

// GitlabRecaptcha handles reCAPTCHA functionality for GitLab
type GitlabRecaptcha struct {
	recaptchaService *service.RecaptchaService
	flashService     *service.FlashService
	gonService       *service.GonService
	logger           *util.Logger
}

// NewGitlabRecaptcha creates a new instance of GitlabRecaptcha
func NewGitlabRecaptcha(
	recaptchaService *service.RecaptchaService,
	flashService *service.FlashService,
	gonService *service.GonService,
	logger *util.Logger,
) *GitlabRecaptcha {
	return &GitlabRecaptcha{
		recaptchaService: recaptchaService,
		flashService:     flashService,
		gonService:       gonService,
		logger:           logger,
	}
}

// LoadRecaptcha loads the reCAPTCHA configuration if it's enabled
func (g *GitlabRecaptcha) LoadRecaptcha(ctx *gin.Context) bool {
	return g.recaptchaService.IsEnabled() && g.recaptchaService.LoadConfigurations(ctx)
}

// CheckRecaptchaMiddleware creates a middleware that checks reCAPTCHA
func (g *GitlabRecaptcha) CheckRecaptchaMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Skip if reCAPTCHA is not loaded
		if !g.LoadRecaptcha(ctx) {
			ctx.Next()
			return
		}

		// Verify reCAPTCHA
		if g.recaptchaService.VerifyRecaptcha(ctx) {
			ctx.Next()
			return
		}

		// Set flash message
		g.flashService.SetAlert(ctx, "There was an error with the reCAPTCHA. Please solve the reCAPTCHA again.")
		g.flashService.Delete(ctx, "recaptcha_error")

		// Create a new resource
		resource := g.createNewResource(ctx)

		// Add GON variables
		g.addGonVariables(ctx)

		// Render the new action
		ctx.HTML(http.StatusOK, "new", gin.H{
			"resource": resource,
		})
		ctx.Abort()
	}
}

// createNewResource creates a new resource of the appropriate type
func (g *GitlabRecaptcha) createNewResource(ctx *gin.Context) interface{} {
	// Get the resource class from the context
	resourceClass, exists := ctx.Get("resource_class")
	if !exists {
		return nil
	}

	// Create a new resource of the appropriate type
	switch resourceClass.(type) {
	case *model.User:
		return &model.User{}
	case *model.Project:
		return &model.Project{}
	case *model.Group:
		return &model.Group{}
	default:
		return nil
	}
}

// addGonVariables adds GON variables to the context
func (g *GitlabRecaptcha) addGonVariables(ctx *gin.Context) {
	g.gonService.AddVariables(ctx)
}
