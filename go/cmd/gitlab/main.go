package main

import (
	"log"

	"gitlab.com/gitlab-org/gitlab/internal/config"
	"gitlab.com/gitlab-org/gitlab/internal/server"
)

func main() {
	// Initialize configuration
	cfg := config.NewConfig()

	// Create and start server
	srv := server.NewServer(cfg)
	log.Printf("Starting GitLab server on %s:%s in %s mode...",
		cfg.Server.Host, cfg.Server.Port, cfg.Server.Mode)

	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
