package webhooks

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
	"gitlab.com/gitlab-org/gitlab-redux/internal/util"
)

// HookActions handles webhook management actions
type HookActions struct {
	createService *service.WebHooksCreateService
	destroyService *service.WebHooksDestroyService
	logger        *util.Logger
}

// NewHookActions creates a new instance of HookActions
func NewHookActions(createService *service.WebHooksCreateService, destroyService *service.WebHooksDestroyService, logger *util.Logger) *HookActions {
	return &HookActions{
		createService: createService,
		destroyService: destroyService,
		logger:        logger,
	}
}

// Index displays a list of webhooks
func (h *HookActions) Index(ctx *gin.Context, relation model.WebHookRelation) {
	hooks := relation.SelectPersisted()
	hook := relation.New()

	ctx.HTML(200, "webhooks/index.html", gin.H{
		"hooks": hooks,
		"hook":  hook,
	})
}

// Create creates a new webhook
func (h *HookActions) Create(ctx *gin.Context, relation model.WebHookRelation) {
	currentUser, _ := ctx.Get("current_user")
	user, _ := currentUser.(model.User)

	hookParams := h.getHookParams(ctx)
	result := h.createService.Execute(user, hookParams, relation)

	if result.Success {
		ctx.Set("flash_notice", "Webhook created")
	} else {
		hooks := relation.SelectPersisted()
		ctx.Set("flash_alert", result.Message)
		ctx.HTML(200, "webhooks/index.html", gin.H{
			"hooks": hooks,
			"hook":  relation.New(),
		})
		return
	}

	ctx.Redirect(302, "/webhooks")
}

// Update updates an existing webhook
func (h *HookActions) Update(ctx *gin.Context, hook model.WebHook) {
	hookParams := h.getHookParams(ctx)

	if hook.Update(hookParams) {
		ctx.Set("flash_notice", "Webhook updated")
		ctx.Redirect(302, "/webhooks/"+hook.ID()+"/edit")
	} else {
		ctx.HTML(200, "webhooks/edit.html", gin.H{
			"hook": hook,
		})
	}
}

// Destroy deletes a webhook
func (h *HookActions) Destroy(ctx *gin.Context, hook model.WebHook) {
	h.destroyHook(ctx, hook)
	ctx.Redirect(302, "/webhooks")
}

// Edit displays the edit form for a webhook
func (h *HookActions) Edit(ctx *gin.Context, hook model.WebHook) {
	if hook == nil {
		ctx.Redirect(302, "/webhooks")
		return
	}

	hookLogs := h.getHookLogs(ctx, hook)

	ctx.HTML(200, "webhooks/edit.html", gin.H{
		"hook":     hook,
		"hookLogs": hookLogs,
	})
}

// getHookParams processes and returns the webhook parameters from the request
func (h *HookActions) getHookParams(ctx *gin.Context) map[string]interface{} {
	hook := ctx.PostForm("hook")
	if hook == "" {
		return nil
	}

	// Get the base hook parameters
	params := make(map[string]interface{})

	// Add the basic hook parameter names
	hookParamNames := []string{
		"enable_ssl_verification", "name", "description", "token", "url",
		"push_events_branch_filter", "branch_filter_strategy", "custom_webhook_template",
	}

	for _, name := range hookParamNames {
		if value := ctx.PostForm("hook[" + name + "]"); value != "" {
			params[name] = value
		}
	}

	// Handle token masking
	if ctx.Request.Method == "PUT" && params["token"] == model.WebHookSecretMask {
		delete(params, "token")
	}

	// Handle URL variables
	if urlVariables := ctx.PostFormArray("hook[url_variables][][key]"); len(urlVariables) > 0 {
		urlValues := ctx.PostFormArray("hook[url_variables][][value]")
		urlVars := make(map[string]string)

		for i, key := range urlVariables {
			if i < len(urlValues) && urlValues[i] != "" {
				urlVars[key] = urlValues[i]
			}
		}

		// For updates, merge with existing variables
		if ctx.Request.Method == "PUT" {
			existingHook, _ := ctx.Get("hook")
			if webHook, ok := existingHook.(model.WebHook); ok {
				existingVars := webHook.GetURLVariables()
				for k, v := range existingVars {
					if _, exists := urlVars[k]; !exists {
						urlVars[k] = v
					}
				}
			}
		}

		params["url_variables"] = urlVars
	}

	// Handle custom headers
	if customHeaders := ctx.PostFormArray("hook[custom_headers][][key]"); len(customHeaders) > 0 {
		customValues := ctx.PostFormArray("hook[custom_headers][][value]")
		customHeadersMap := make(map[string]string)

		for i, key := range customHeaders {
			if i < len(customValues) {
				value := customValues[i]

				// Check if we need to use the existing value from the database
				if value == model.WebHookSecretMask {
					existingHook, _ := ctx.Get("hook")
					if webHook, ok := existingHook.(model.WebHook); ok {
						existingHeaders := webHook.GetCustomHeaders()
						if existingValue, exists := existingHeaders[key]; exists {
							value = existingValue
						}
					}
				}

				customHeadersMap[key] = value
			}
		}

		params["custom_headers"] = customHeadersMap
	}

	return params
}

// destroyHook destroys a webhook
func (h *HookActions) destroyHook(ctx *gin.Context, hook model.WebHook) {
	currentUser, _ := ctx.Get("current_user")
	user, _ := currentUser.(model.User)

	result := h.destroyService.Execute(user, hook)

	if result.Status == "success" {
		if result.Async {
			ctx.Set("flash_notice", "Webhook scheduled for deletion")
		} else {
			ctx.Set("flash_notice", "Webhook deleted")
		}
	} else {
		ctx.Set("flash_alert", result.Message)
	}
}

// getHookLogs retrieves the logs for a webhook
func (h *HookActions) getHookLogs(ctx *gin.Context, hook model.WebHook) []model.WebHookLog {
	page := ctx.DefaultQuery("page", "1")
	return hook.GetWebHookLogs().Recent().Page(page).WithoutCount()
}
