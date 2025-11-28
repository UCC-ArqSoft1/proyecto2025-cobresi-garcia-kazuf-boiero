package main

import (
	"log"

	"github.com/alesio/gestion-actividades-deportivas/config"
	"github.com/alesio/gestion-actividades-deportivas/database"
	"github.com/alesio/gestion-actividades-deportivas/handlers"
	"github.com/alesio/gestion-actividades-deportivas/middlewares"
	"github.com/alesio/gestion-actividades-deportivas/services"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("database initialization failed: %v", err)
	}

	router := gin.Default()
	router.Use(middlewares.CORSMiddleware())

	// Initialize services.
	authService := services.NewAuthService(db, cfg)
	userService := services.NewUserService(db)
	activityService := services.NewActivityService(db)
	enrollmentService := services.NewEnrollmentService(db)

	// Initialize handlers.
	healthHandler := handlers.NewHealthHandler()
	authHandler := handlers.NewAuthHandler(authService, userService)
	activitiesHandler := handlers.NewActivitiesHandler(activityService)
	enrollmentsHandler := handlers.NewEnrollmentsHandler(enrollmentService)
	adminActivitiesHandler := handlers.NewAdminActivitiesHandler(activityService)

	// Register health route.
	healthHandler.RegisterRoutes(router)

	apiGroup := router.Group("/api")
	authHandler.RegisterRoutes(apiGroup)
	activitiesHandler.RegisterRoutes(apiGroup)

	authMiddleware := middlewares.NewAuthMiddleware(authService)

	protected := apiGroup.Group("")
	protected.Use(authMiddleware.Handle())
	enrollmentsHandler.RegisterRoutes(protected)

	adminGroup := apiGroup.Group("")
	adminGroup.Use(authMiddleware.Handle(), middlewares.AdminMiddleware())
	adminActivitiesHandler.RegisterRoutes(adminGroup)

	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
