package database

import (
	"fmt"
	"log"

	"japanese-learning-app/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Initialize initializes the database connection
func Initialize(databaseURL string) (*gorm.DB, error) {
	// Configure GORM logger
	gormLogger := logger.Default.LogMode(logger.Info)

	// Connect to PostgreSQL
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get the underlying sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	log.Println("Database connection established successfully")
	return db, nil
}

// Migrate runs database migrations
func Migrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// Auto-migrate all models
	err := db.AutoMigrate(
		&models.User{},
		&models.Book{},
		&models.ReadingSession{},
		&models.BookAnnotation{},
		// Add more models here as we create them
		// &models.Word{},
		// &models.WordDefinition{},
		// &models.UserWordKnowledge{},
		// &models.SRSCard{},
		// &models.ReviewHistory{},
	)

	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// CreateIndexes creates additional database indexes for better performance
func CreateIndexes(db *gorm.DB) error {
	log.Println("Creating database indexes...")

	// Add custom indexes if needed
	// Example:
	// db.Exec("CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_books_user_status ON books(user_id, processing_status);")

	log.Println("Database indexes created successfully")
	return nil
}

// CloseDB closes the database connection
func CloseDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
