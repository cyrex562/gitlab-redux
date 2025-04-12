package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

type DeployKeysController struct {
	deployKeyService *service.DeployKeyService
}

func NewDeployKeysController(deployKeyService *service.DeployKeyService) *DeployKeysController {
	return &DeployKeysController{
		deployKeyService: deployKeyService,
	}
}

// Index displays the list of deploy keys
func (c *DeployKeysController) Index(ctx *gin.Context) {
	deployKeys, err := c.deployKeyService.GetPublicKeys(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"deploy_keys": deployKeys,
	})
}

// New displays the form to create a new deploy key
func (c *DeployKeysController) New(ctx *gin.Context) {
	// In a real application, this would render a template
	// For now, we'll just return a success response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "New deploy key form",
	})
}

// Create creates a new deploy key
func (c *DeployKeysController) Create(ctx *gin.Context) {
	var params service.CreateDeployKeyParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set public to true as per the Ruby implementation
	params.Public = true

	deployKey, err := c.deployKeyService.Create(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"deploy_key": deployKey,
	})
}

// Edit displays the form to edit a deploy key
func (c *DeployKeysController) Edit(ctx *gin.Context) {
	id := ctx.Param("id")
	deployKey, err := c.deployKeyService.GetByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// In a real application, this would render a template
	// For now, we'll just return the deploy key
	ctx.JSON(http.StatusOK, gin.H{
		"deploy_key": deployKey,
	})
}

// Update updates a deploy key
func (c *DeployKeysController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var params service.UpdateDeployKeyParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	deployKey, err := c.deployKeyService.Update(ctx, id, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"deploy_key": deployKey,
	})
}

// Destroy deletes a deploy key
func (c *DeployKeysController) Destroy(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.deployKeyService.Delete(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
