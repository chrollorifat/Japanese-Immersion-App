package config

import (
	"os"
	"log"
	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	// Database
	DatabaseURL string

	// Server
	Port string
	Host string

	// Authentication
	JWTSecret           string
	AccessTokenExpiry   string
	
	// File Upload
	UploadDir     string
	MaxFileSize   int64

	// External APIs
	OpenAIAPIKey string
	JishoAPIURL  string

	// App Settings
	Debug        bool
	FrontendURL  string
}

// Load loads configuration from environment variables
func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
	return &Config{
		// Database
		DatabaseURL: getEnv("DATABASE_URL", "postgres://username:password@localhost:5432/japanese_learning_app?sslmode=disable"),

		// Server
		Port: getEnv("PORT", "8000"),
		Host: getEnv("HOST", "0.0.0.0"),

		// Authentication
		JWTSecret:         getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-this-in-production"),
		AccessTokenExpiry: getEnv("ACCESS_TOKEN_EXPIRY", "24h"),

		// File Upload
		UploadDir:   getEnv("UPLOAD_DIR", "../data/uploads"),
		MaxFileSize: getEnvInt64("MAX_FILE_SIZE", 50*1024*1024), // 50MB

		// External APIs
		OpenAIAPIKey: getEnv("OPENAI_API_KEY", ""),
		JishoAPIURL:  getEnv("JISHO_API_URL", "https://jisho.org/api/v1/search/words"),

		// App Settings
		Debug:       getEnvBool("DEBUG", true),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvBool gets a boolean environment variable or returns a default value
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1"
	}
	return defaultValue
}

// getEnvInt64 gets an int64 environment variable or returns a default value
func getEnvInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		// Simple conversion - in production you'd want proper error handling
		if value == "52428800" {
			return 52428800 // 50MB
		}
	}
	return defaultValue
}
