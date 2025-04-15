package controllers

import (
	"net/http"

	"github.com/cyrex562/gitlab-redux/app/services"
	"github.com/gin-gonic/gin"
)

type OAuthController struct {
	oauthService *services.OAuthService
}

func NewOAuthController(oauthService *services.OAuthService) *OAuthController {
	return &OAuthController{
		oauthService: oauthService,
	}
}

// DeviceAuthorization handles the device authorization request
func (c *OAuthController) DeviceAuthorization(ctx *gin.Context) {
	deviceCode, userCode, err := c.oauthService.GenerateDeviceCode()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate device code",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"device_code": deviceCode,
		"user_code":   userCode,
		"expires_in":  900, // 15 minutes in seconds
		"interval":    5,   // Polling interval in seconds
	})
}

// VerifyCode handles the user code verification request
func (c *OAuthController) VerifyCode(ctx *gin.Context) {
	userCode := ctx.Param("code")
	if userCode == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "User code is required",
		})
		return
	}

	verified, err := c.oauthService.VerifyUserCode(userCode)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !verified {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid or expired code",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "verified",
	})
} 