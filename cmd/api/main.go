package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"suitemedia/config"
	"suitemedia/internal/database"
	"suitemedia/internal/handlers"
	"suitemedia/internal/middleware"
	"suitemedia/internal/repository"
	"suitemedia/internal/service"
	"suitemedia/pkg/logger"
	"suitemedia/pkg/redis"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger := logger.NewLogger(cfg.App.LogLevel)
	logger.Info("Starting SuiteMedia API Server")

	// Initialize database connection
	db, err := database.NewConnection(cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", "error", err)
	}
	defer db.Close()

	// Run database migrations
	if err := database.RunMigrations(db); err != nil {
		logger.Fatal("Failed to run migrations", "error", err)
	}

	// Initialize Redis client
	redisClient, err := redis.NewClient(cfg.Redis)
	if err != nil {
		logger.Fatal("Failed to connect to Redis", "error", err)
	}
	defer redisClient.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo, redisClient)
	productService := service.NewProductService(productRepo, redisClient)
	authService := service.NewAuthService(userRepo, cfg.JWT)

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler(db, redisClient)
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	productHandler := handlers.NewProductHandler(productService)

	// Setup Gin router
	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Global middleware
	router.Use(middleware.Logger(logger))
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.CORS(cfg.CORS))
	router.Use(middleware.RequestID())
	router.Use(middleware.Metrics())

	// Health check endpoints
	router.GET("/health", healthHandler.Health)
	router.GET("/ready", healthHandler.Ready)
	router.GET("/metrics", handlers.PrometheusHandler())

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthRequired(cfg.JWT))
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("", userHandler.List)
				users.GET("/:id", userHandler.GetByID)
				users.POST("", userHandler.Create)
				users.PUT("/:id", userHandler.Update)
				users.DELETE("/:id", userHandler.Delete)
				users.GET("/me", userHandler.GetProfile)
				users.PUT("/me", userHandler.UpdateProfile)
			}

			// Product routes
			products := protected.Group("/products")
			{
				products.GET("", productHandler.List)
				products.GET("/:id", productHandler.GetByID)
				products.POST("", middleware.RoleRequired("admin"), productHandler.Create)
				products.PUT("/:id", middleware.RoleRequired("admin"), productHandler.Update)
				products.DELETE("/:id", middleware.RoleRequired("admin"), productHandler.Delete)
			}
		}
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.App.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Server starting", "port", cfg.App.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown with 30 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", "error", err)
	}

	logger.Info("Server exited")
}
