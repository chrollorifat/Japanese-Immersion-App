package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Basic user information
	Username string `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Email    string `json:"email" gorm:"uniqueIndex;size:100;not null"`
	Password string `json:"-" gorm:"size:255;not null"`

	// User status
	IsActive          bool   `json:"is_active" gorm:"default:true"`
	PreferredLanguage string `json:"preferred_language" gorm:"size:10;default:en"`

	// Learning statistics
	TotalWordsLearned int       `json:"total_words_learned" gorm:"default:0"`
	TotalReadingTime  int       `json:"total_reading_time" gorm:"default:0"` // in minutes
	StreakDays        int       `json:"streak_days" gorm:"default:0"`
	LastActivity      *time.Time `json:"last_activity"`

	// JSON field for additional preferences
	LearningPreferences map[string]interface{} `json:"learning_preferences" gorm:"serializer:json"`

	// Relationships
	Books           []Book           `json:"books,omitempty" gorm:"foreignKey:UserID"`
	ReadingSessions []ReadingSession `json:"reading_sessions,omitempty" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for GORM
func (User) TableName() string {
	return "users"
}

// HashPassword hashes the user's password
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword compares the provided password with the user's hashed password
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// ToResponse converts the user to a response format (excluding sensitive data)
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:                  u.ID,
		Username:            u.Username,
		Email:               u.Email,
		CreatedAt:           u.CreatedAt,
		IsActive:            u.IsActive,
		PreferredLanguage:   u.PreferredLanguage,
		TotalWordsLearned:   u.TotalWordsLearned,
		TotalReadingTime:    u.TotalReadingTime,
		StreakDays:          u.StreakDays,
		LastActivity:        u.LastActivity,
		LearningPreferences: u.LearningPreferences,
	}
}

// UserResponse is the response format for user data
type UserResponse struct {
	ID                  uint                   `json:"id"`
	Username            string                 `json:"username"`
	Email               string                 `json:"email"`
	CreatedAt           time.Time              `json:"created_at"`
	IsActive            bool                   `json:"is_active"`
	PreferredLanguage   string                 `json:"preferred_language"`
	TotalWordsLearned   int                    `json:"total_words_learned"`
	TotalReadingTime    int                    `json:"total_reading_time"`
	StreakDays          int                    `json:"streak_days"`
	LastActivity        *time.Time             `json:"last_activity"`
	LearningPreferences map[string]interface{} `json:"learning_preferences"`
}

// CreateUserRequest is the request format for creating a user
type CreateUserRequest struct {
	Username          string `json:"username" binding:"required,min=3,max=50"`
	Email             string `json:"email" binding:"required,email"`
	Password          string `json:"password" binding:"required,min=6"`
	PreferredLanguage string `json:"preferred_language"`
}

// LoginRequest is the request format for user login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse is the response format for successful login
type LoginResponse struct {
	AccessToken string       `json:"access_token"`
	TokenType   string       `json:"token_type"`
	User        UserResponse `json:"user"`
}
