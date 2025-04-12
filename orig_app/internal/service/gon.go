package service

import (
	"github.com/gin-gonic/gin"
)

// GonService handles setting up gon variables for JavaScript
type GonService struct {
	config *Config
}

// NewGonService creates a new instance of GonService
func NewGonService(config *Config) *GonService {
	return &GonService{
		config: config,
	}
}

// AddVariables adds gon variables to the context for JavaScript
func (s *GonService) AddVariables(ctx *gin.Context) error {
	// Get the relative URL root from the config
	relativeURLRoot := s.config.RelativeURLRoot

	// Add gon variables to the context
	gon := map[string]interface{}{
		"relative_url_root": relativeURLRoot,
		// Add other gon variables as needed
	}

	// Store gon variables in the context for template rendering
	ctx.Set("gon", gon)

	return nil
}

// Config holds configuration for the GonService
type Config struct {
	RelativeURLRoot string
	// Add other configuration options as needed
}
