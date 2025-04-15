package server

import (
	"fmt"
	"net/http"

	"github.com/cyrex562/gitlab-redux/internal/config"
	"github.com/cyrex562/gitlab-redux/internal/websocket"
	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/app/handlers"
	"gitlab.com/gitlab-org/gitlab-redux/app/routes"
)

type Server struct {
	router *gin.Engine
	config *config.Config
	hub    *websocket.Hub
}

func NewServer(cfg *config.Config) *Server {
	gin.SetMode(cfg.Server.Mode)
	router := gin.Default()

	// Initialize WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Serve static files
	router.Static("/assets", "./static")
	router.StaticFile("/favicon.ico", "./static/images/favicon.png")

	// WebSocket endpoint
	router.GET("/cable", func(c *gin.Context) {
		websocket.NewHandler(hub).ServeHTTP(c.Writer, c.Request)
	})

	return &Server{
		router: router,
		config: cfg,
		hub:    hub,
	}
}

func (s *Server) setupRoutes() {
	// Health check
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// API v1 group
	v1 := s.router.Group("/api/v1")
	{
		v1.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Welcome to GitLab API v1",
			})
		})
	}

	// Register OAuth routes
	baseHandler := handlers.NewBaseHandler()
	routes.RegisterOAuthRoutes(s.router, baseHandler)

	// Serve index.html for the root path
	s.router.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})
}

func (s *Server) Start() error {
	s.setupRoutes()

	addr := fmt.Sprintf("%s:%s", s.config.Server.Host, s.config.Server.Port)
	return s.router.Run(addr)
}
