package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"japanese-learning-app/internal/config"
	"japanese-learning-app/internal/database"
	"japanese-learning-app/internal/handlers"
	"japanese-learning-app/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Set JWT secret from config
	middleware.SetJWTSecret(cfg.JWTSecret)

	// Initialize database
	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB(db)

	// Run database migrations
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize Gin router
	r := gin.Default()

	// CORS middleware
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{
		"http://localhost:3000",
		"http://127.0.0.1:3000",
		"http://localhost:8080",
	}
	// Allow local file and null origins during development without listing "file://" (which gin-cors rejects)
	corsConfig.AllowOriginFunc = func(origin string) bool {
		return origin == "null" || strings.HasPrefix(origin, "file://")
	}
	corsConfig.AllowCredentials = true
	corsConfig.AllowMethods = []string{
		"GET",
		"POST",
		"PUT",
		"PATCH",
		"DELETE",
		"OPTIONS",
	}
	corsConfig.AllowHeaders = []string{
		"Origin",
		"Content-Length",
		"Content-Type",
		"Authorization",
		"Accept",
		"X-Requested-With",
	}
	r.Use(cors.New(corsConfig))

	// Serve static files
	r.Static("/static", "../frontend/static")
	r.Static("/uploads", "../data/uploads")

	// Initialize handlers
	h := handlers.New(db)

	// Database middleware - make database available to all routes
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"message": "Japanese Learning API is running",
		})
	})

	// API routes
	api := r.Group("/api")
	{
		// Public routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", h.Register)
			auth.POST("/login", h.Login)
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthRequired())
		{
			// User routes
			protected.GET("/auth/me", h.GetCurrentUser)
			protected.POST("/auth/logout", h.Logout)

			// Book routes
			books := protected.Group("/books")
			{
				books.GET("/", h.GetUserBooks)
				books.POST("/upload", h.UploadBook)
				books.GET("/:id", h.GetBook)
				books.DELETE("/:id", h.DeleteBook)
			}

			// Word routes
			words := protected.Group("/words")
			{
				words.GET("/lookup/:word", h.LookupWord)
				words.POST("/mark-known", h.MarkWordAsKnown)
			}

			// SRS routes
			srs := protected.Group("/srs")
			{
				srs.GET("/due-cards", h.GetDueCards)
				srs.POST("/review", h.ReviewCard)
			}

			// Reading session routes
			reading := protected.Group("/reading")
			{
				reading.POST("/start-session", h.StartReadingSession)
				reading.POST("/end-session", h.EndReadingSession)
				reading.GET("/stats", h.GetReadingStats)
			}
		}
	}

	// Root endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Japanese Learning Web App API",
			"version": "1.0.0",
			"status":  "running",
			"docs":    "/docs",
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("Starting server on port %s", port)
	log.Printf("API documentation available at http://localhost:%s/docs", port)
	
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
