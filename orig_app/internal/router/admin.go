package router

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/your-project/model"
	"github.com/your-project/service"
)

// RegisterAdminRoutes registers all admin routes
func RegisterAdminRoutes(r *gin.Engine, db *sql.DB) {
	admin := r.Group("/admin")
	{
		// ... existing routes ...

		// Background Migrations routes
		migrationService := service.NewBackgroundMigrationService(db)
		migrationController := admin.NewBackgroundMigrationsController(migrationService)

		admin.GET("/background_migrations", migrationController.Index)
		admin.GET("/background_migrations/:id", migrationController.Show)
		admin.POST("/background_migrations/:id/pause", migrationController.Pause)
		admin.POST("/background_migrations/:id/resume", migrationController.Resume)
		admin.POST("/background_migrations/:id/retry", migrationController.Retry)

		// Batched Jobs routes
		jobService := service.NewBatchedJobService(db)
		jobController := admin.NewBatchedJobsController(jobService)

		admin.GET("/batched_jobs/:id", jobController.Show)

		// Broadcast Messages routes
		broadcastService := service.NewBroadcastMessageService(db)
		broadcastController := admin.NewBroadcastMessagesController(broadcastService)

		admin.GET("/broadcast_messages", broadcastController.Index)
		admin.GET("/broadcast_messages/:id", broadcastController.Show)
		admin.POST("/broadcast_messages", broadcastController.Create)
		admin.PUT("/broadcast_messages/:id", broadcastController.Update)
		admin.DELETE("/broadcast_messages/:id", broadcastController.Delete)
		admin.POST("/broadcast_messages/preview", broadcastController.Preview)

		// Clusters routes
		clusterService := service.NewClusterService(db)
		featureService := service.NewFeatureService(db)
		user := &model.User{} // This should be properly initialized with the current user
		clusterController := admin.NewAdminClustersController(clusterService, featureService, user)

		// Apply feature check middleware
		admin.Use(clusterController.EnsureFeatureEnabled())

		admin.GET("/clusters", clusterController.Index)
		admin.GET("/clusters/new", clusterController.New)
		admin.GET("/clusters/:id", clusterController.Show)
		admin.POST("/clusters", clusterController.Create)
		admin.PUT("/clusters/:id", clusterController.Update)
		admin.DELETE("/clusters/:id", clusterController.Delete)

		// Cohorts routes
		cohortsService := service.NewCohortsService(db, redisClient)
		analyticsService := service.NewAnalyticsService(redisClient)
		cohortsController := admin.NewCohortsController(cohortsService, analyticsService)

		admin.GET("/cohorts", cohortsController.Index)
	}
}
