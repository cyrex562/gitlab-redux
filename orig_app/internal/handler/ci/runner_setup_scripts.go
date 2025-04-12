// Package ci provides handlers for CI/CD related functionality
package ci

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RunnerSetupScripts handles the generation of runner setup scripts
type RunnerSetupScripts struct {
	// Dependencies can be injected here if needed
}

// NewRunnerSetupScripts creates a new RunnerSetupScripts handler
func NewRunnerSetupScripts() *RunnerSetupScripts {
	return &RunnerSetupScripts{}
}

// PrivateRunnerSetupScripts generates installation and registration scripts for GitLab runners
func (h *RunnerSetupScripts) PrivateRunnerSetupScripts(c *gin.Context) {
	// Extract parameters from the request
	os := c.Query("os")
	arch := c.Query("arch")

	// Create a new runner instructions instance
	instructions := NewRunnerInstructions(os, arch)

	// Generate the output
	output := map[string]string{
		"install":  instructions.InstallScript(),
		"register": instructions.RegisterCommand(),
	}

	// Check for errors
	if len(instructions.Errors()) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": instructions.Errors(),
		})
		return
	}

	// Return the successful response
	c.JSON(http.StatusOK, output)
}

// RunnerInstructions represents the instructions for setting up a GitLab runner
type RunnerInstructions struct {
	os      string
	arch    string
	errors  []string
}

// NewRunnerInstructions creates a new RunnerInstructions instance
func NewRunnerInstructions(os, arch string) *RunnerInstructions {
	return &RunnerInstructions{
		os:     os,
		arch:   arch,
		errors: []string{},
	}
}

// InstallScript returns the installation script for the runner
func (ri *RunnerInstructions) InstallScript() string {
	// This is a placeholder implementation
	// In a real implementation, this would generate the appropriate installation script
	// based on the OS and architecture
	return "#!/bin/bash\n# Installation script for " + ri.os + " on " + ri.arch
}

// RegisterCommand returns the registration command for the runner
func (ri *RunnerInstructions) RegisterCommand() string {
	// This is a placeholder implementation
	// In a real implementation, this would generate the appropriate registration command
	return "gitlab-runner register"
}

// Errors returns any errors that occurred during instruction generation
func (ri *RunnerInstructions) Errors() []string {
	return ri.errors
}
